package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// PlatformBlock 表示一个平台特定的代码块
type PlatformBlock struct {
	Platform string
	Commands []string
}

// ZRunScript 表示整个zrun脚本
type ZRunScript struct {
	Blocks []PlatformBlock
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

	script := &ZRunScript{
		Blocks: make([]PlatformBlock, 0),
	}

	scanner := bufio.NewScanner(file)
	var currentBlock *PlatformBlock

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 检查是否是平台块开始
		if strings.HasPrefix(line, "@") && strings.HasSuffix(line, "{") {
			platform := strings.TrimPrefix(line, "@")
			platform = strings.TrimSuffix(platform, " {")
			
			block := PlatformBlock{
				Platform: platform,
				Commands: make([]string, 0),
			}
			
			currentBlock = &block
			continue
		}

		// 检查是否是块结束
		if line == "}" && currentBlock != nil {
			script.Blocks = append(script.Blocks, *currentBlock)
			currentBlock = nil
			continue
		}

		// 添加命令到当前块
		if currentBlock != nil {
			currentBlock.Commands = append(currentBlock.Commands, line)
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
	var defaultBlock *PlatformBlock
	
	// 首先查找匹配当前平台的块
	for _, block := range script.Blocks {
		if block.Platform == "default" {
			defaultBlock = &block
			continue
		}
		
		if matchPlatform(block.Platform, currentOS) {
			fmt.Printf("# 执行 %s 平台命令:\n", block.Platform)
			return executeCommands(block.Commands)
		}
	}
	
	// 如果没有找到匹配的平台块，尝试执行默认块
	if defaultBlock != nil {
		fmt.Printf("# 执行默认命令块:\n")
		return executeCommands(defaultBlock.Commands)
	}
	
	fmt.Println("未找到适用于当前平台的命令块")
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

// executeCommands 执行命令列表
func executeCommands(commands []string) error {
	for _, command := range commands {
		err := executeCommand(command)
		if err != nil {
			return fmt.Errorf("执行命令 '%s' 失败: %v", command, err)
		}
	}
	return nil
}

// executeCommand 执行单个系统命令
func executeCommand(command string) error {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	
	// 设置输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// 执行命令
	fmt.Printf("$ %s\n", command)
	return cmd.Run()
}