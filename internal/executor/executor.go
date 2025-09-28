// Package executor 提供了执行解析后脚本的功能。
package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"zrun/internal/types"
	"zrun/internal/utils"
)

// ExecuteScript 执行解析后的脚本
// 参数:
//   - script: 解析后的脚本结构体
//
// 返回值:
//   - error: 执行过程中可能发生的错误
func ExecuteScript(script *types.ZRunScript) error {
	currentOS := utils.GetOS()

	for _, cmd := range script.Commands {
		// 根据命令类型处理
		err := processCommand(cmd, script, currentOS)
		if err != nil {
			return err
		}
	}

	return nil
}

// processCommand 处理单个命令
// 参数:
//   - cmd: 要处理的命令
//   - script: 脚本结构体引用
//   - currentOS: 当前操作系统
//
// 返回值:
//   - error: 处理过程中可能发生的错误
func processCommand(cmd types.ScriptCommand, script *types.ZRunScript, currentOS string) error {
	// 处理echo指令
	if cmd.Type == "echo" {
		return processEchoCommand(cmd, script)
	}

	// 处理平台命令
	if cmd.Type == "platform" && utils.MatchPlatform(cmd.Param, currentOS) {
		return ExecuteCommand(cmd.Command, script.EchoOn)
	}

	return nil
}

// processEchoCommand 处理 echo 控制命令
// 参数:
//   - cmd: echo 命令
//   - script: 脚本结构体引用
//
// 返回值:
//   - error: 处理过程中可能发生的错误
func processEchoCommand(cmd types.ScriptCommand, script *types.ZRunScript) error {
	if cmd.Param == "off" {
		script.EchoOn = false
	} else if cmd.Param == "on" {
		script.EchoOn = true
	}
	return nil
}

// ExecuteCommand 执行单个系统命令
// 参数:
//   - command: 要执行的命令字符串
//   - echoOn: 是否回显命令
//
// 返回值:
//   - error: 命令执行过程中可能发生的错误
func ExecuteCommand(command string, echoOn bool) error {
	var cmd *exec.Cmd

	// 根据操作系统选择合适的shell执行命令
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	// 设置命令的输出流
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 根据echoOn参数决定是否显示命令
	if echoOn {
		fmt.Printf("$ %s\n", command)
	}

	// 执行命令并返回结果
	return cmd.Run()
}
