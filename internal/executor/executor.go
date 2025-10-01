package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"zrun/internal/types"
	"zrun/internal/utils"
)

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

func processEchoCommand(cmd types.ScriptCommand, script *types.ZRunScript) error {
	switch cmd.Param {
	case "off":
		script.EchoOn = false
	case "on":
		script.EchoOn = true
	}
	return nil
}

func ExecuteCommand(command string, echoOn bool) error {
	var cmd *exec.Cmd

	// 根据操作系统选择hell执行命令
	// fack microsoft
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 根据echoOn决定是否显示命令
	if echoOn {
		fmt.Printf("$ %s\n", command)
	}

	// 执行
	return cmd.Run()
}
