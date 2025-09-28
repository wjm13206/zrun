package parser

import (
	"bufio"
	"os"
	"strings"
	"zrun/internal/types"
)

// ParseScript 解析.zr脚本文件
func ParseScript(filename string) (*types.ZRunScript, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 默认开启echo
	script := &types.ZRunScript{
		Commands: make([]types.ScriptCommand, 0),
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
			script.Commands = append(script.Commands, types.ScriptCommand{
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
			script.Commands = append(script.Commands, types.ScriptCommand{
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