package parser

import (
	"bufio"
	"os"
	"strings"
	"zrun/internal/types"
)

// ParseScript 解析 .zr 脚本文件
// 参数:
//   - filename: 要解析的脚本文件名
//
// 返回值:
//   - *types.ZRunScript: 解析后的脚本结构体
//   - error: 解析过程中可能发生的错误
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
		// 处理扫描到的每一行
		err := processLine(line, script, &currentPlatform)
		if err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return script, nil
}

// processLine 处理脚本中的一行
// 参数:
//   - line: 当前处理的行内容
//   - script: 脚本结构体引用
//   - currentPlatform: 当前平台块引用
//
// 返回值:
//   - error: 处理过程中可能发生的错误
func processLine(line string, script *types.ZRunScript, currentPlatform *string) error {
	// 跳过空行和注释
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}

	// 检查@echo指令
	if strings.HasPrefix(line, "@echo ") {
		return processEchoCommand(line, script)
	}

	// 检查是否是平台块开始
	if strings.HasPrefix(line, "@") && strings.HasSuffix(line, "{") {
		platform := strings.TrimPrefix(line, "@")
		platform = strings.TrimSuffix(platform, " {")
		*currentPlatform = platform
		return nil
	}

	// 检查是否是块结束
	if line == "}" {
		*currentPlatform = ""
		return nil
	}

	// 添加命令到当前平台块
	if *currentPlatform != "" {
		script.Commands = append(script.Commands, types.ScriptCommand{
			Command: line,
			Type:    "platform",
			Param:   *currentPlatform,
		})
	}

	return nil
}

// processEchoCommand 处理 @echo 指令
// 参数:
//   - line: 包含 @echo 指令的行
//   - script: 脚本结构体引用
//
// 返回值:
//   - error: 处理过程中可能发生的错误
func processEchoCommand(line string, script *types.ZRunScript) error {
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

	return nil
}
