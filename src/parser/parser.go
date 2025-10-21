package parser

import (
	"bufio"
	"os"
	"strings"
	"zrun/src/types"
)

// 解析
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

	// 预估命令数量以减少切片重新分配
	script := &types.ZRunScript{
		Commands: make([]types.ScriptCommand, 0, 32), // 预分配容量
		EchoOn:   true,
	}

	scanner := bufio.NewScanner(file)
	currentPlatform := "" // 使用值

	// 复用字符串
	echoPrefix := "@echo "
	platformPrefix := "@"
	blockSuffix := " {"
	blockEnd := "}"

	for scanner.Scan() {
		line := scanner.Text() // 直接处理原始行
		// 手动去除行首尾空格
		start, end := 0, len(line)-1
		for start <= end && (line[start] == ' ' || line[start] == '\t') {
			start++
		}
		for end >= start && (line[end] == ' ' || line[end] == '\t') {
			end--
		}
		if start > end {
			continue // 空行
		}
		line = line[start : end+1]

		// 跳过注释
		if line[0] == '#' {
			continue
		}

		// 检查@echo指令
		if strings.HasPrefix(line, echoPrefix) {
			echoParam := strings.TrimSpace(line[len(echoPrefix):])
			// 添加到命令列表中
			script.Commands = append(script.Commands, types.ScriptCommand{
				Command: line,
				Type:    "echo",
				Param:   echoParam,
			})
			continue
		}

		// 检查是否是平台块开始
		if strings.HasPrefix(line, platformPrefix) && strings.HasSuffix(line, blockSuffix) {
			currentPlatform = line[1 : len(line)-2] // 去掉@和 {
			continue
		}

		// 检查是否是块结束
		if line == blockEnd {
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