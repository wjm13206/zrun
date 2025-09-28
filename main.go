package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ScriptCommand 表示脚本中的一个命令
type ScriptCommand struct {
	Command string
	Type    string // "echo" 或 "platform"
	Param   string // platform名称或其他参数
}

// ZRunScript 表示整个zrun脚本
type ZRunScript struct {
	Commands []ScriptCommand
	EchoOn   bool // 控制是否显示命令
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: zrun <脚本文件.zr>")
		os.Exit(1)
	}

	filename := os.Args[1]
	script, err := parseScript(filename)
	if err != nil {
		fmt.Printf("解析脚本错误: %v\n", err)
		os.Exit(1)
	}

	err = executeScript(script)
	if err != nil {
		fmt.Printf("执行脚本错误: %v\n", err)
		os.Exit(1)
	}
}

// parseScript 解析.zr脚本文件
func parseScript(filename string) (*ZRunScript, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 默认开启echo
	script := &ZRunScript{
		Commands: make([]ScriptCommand, 0),
		EchoOn:   true,
	}

	scanner := bufio.NewScanner(file)
	var currentPlatform string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 检查@echo指令
		if strings.HasPrefix(line, "@echo ") {
			echoParam := strings.TrimSpace(strings.TrimPrefix(line, "@echo"))
			if echoParam == "off" {
				script.EchoOn = false
			} else if echoParam == "on" {
				script.EchoOn = true
			}
			// 添加到命令列表中，以便在执行时处理
			script.Commands = append(script.Commands, ScriptCommand{
				Command: line,
				Type:    "echo",
				Param:   echoParam,
			})
			continue
		}

		// 检查是否是平台块开始
		if strings.HasPrefix(line, "@") && strings.HasSuffix(line, "{") {
			platform := strings.TrimPrefix(line, "@")
			platform = strings.TrimSuffix(platform, " {")
			currentPlatform = platform
			continue
		}

		// 检查是否是块结束
		if line == "}" {
			currentPlatform = ""
			continue
		}

		// 添加命令到当前平台块
		if currentPlatform != "" {
			script.Commands = append(script.Commands, ScriptCommand{
				Command: line,
				Type:    "platform",
				Param:   currentPlatform,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return script, nil
}

// executeScript 执行解析后的脚本
func executeScript(script *ZRunScript) error {
	currentOS := getOS()
	
	for _, cmd := range script.Commands {
		// 处理echo指令
		if cmd.Type == "echo" {
			if cmd.Param == "off" {
				script.EchoOn = false
			} else if cmd.Param == "on" {
				script.EchoOn = true
			}
			continue
		}
		
		// 处理平台命令
		if cmd.Type == "platform" && matchPlatform(cmd.Param, currentOS) {
			err := executeCommand(cmd.Command, script.EchoOn)
			if err != nil {
				return fmt.Errorf("执行命令 '%s' 失败: %v", cmd.Command, err)
			}
		}
	}
	
	return nil
}

// getOS 获取当前操作系统
func getOS() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	default:
		return runtime.GOOS
	}
}

// matchPlatform 检查平台是否匹配
func matchPlatform(platform, currentOS string) bool {
	// default块总是匹配
	if platform == "default" {
		return true
	}
	
	// 支持多个平台名称
	platforms := strings.Split(platform, ",")
	for _, p := range platforms {
		p = strings.TrimSpace(p)
		if p == currentOS {
			return true
		}
		
		// 支持通用Unix平台
		if p == "unix" && (currentOS == "linux" || currentOS == "macos") {
			return true
		}
	}
	return false
}

// executeCommand 执行单个系统命令
func executeCommand(command string, echoOn bool) error {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	// 设置输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// 根据echoOn参数决定是否显示命令
	if echoOn {
		fmt.Printf("$ %s\n", command)
	}
	return cmd.Run()
}