package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"zrun/src/types"
	"zrun/src/utils"
)

func ExecuteScript(script *types.ZRunScript) error {
	currentOS := utils.GetOS()
	
	// 预分配命令切片，减少内存重新分配
	matchingCommands := make([]types.ScriptCommand, 0, len(script.Commands))
	
	// 第一次遍历：筛选出需要执行的命令
	for _, cmd := range script.Commands {
		// 处理echo指令
		if cmd.Type == "echo" {
			matchingCommands = append(matchingCommands, cmd)
			continue
		}
		
		// 处理平台命令
		if cmd.Type == "platform" && utils.MatchPlatform(cmd.Param, currentOS) {
			matchingCommands = append(matchingCommands, cmd)
		}
	}
	
	// 执行筛选后的命令
	for _, cmd := range matchingCommands {
		// 处理echo指令
		if cmd.Type == "echo" {
			switch cmd.Param {
			case "off":
				script.EchoOn = false
			case "on":
				script.EchoOn = true
			}
			continue
		}
		
		// 执行平台命令
		if cmd.Type == "platform" {
			err := ExecuteCommand(cmd.Command, script.EchoOn)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}

func ExecuteCommand(command string, echoOn bool) error {
	var cmd *exec.Cmd

	// 根据操作系统选择shell执行命令
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