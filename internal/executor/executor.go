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
func ExecuteScript(script *types.ZRunScript) error {
	currentOS := utils.GetOS()
	
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
		if cmd.Type == "platform" && utils.MatchPlatform(cmd.Param, currentOS) {
			err := ExecuteCommand(cmd.Command, script.EchoOn)
			if err != nil {
				return fmt.Errorf("执行命令 '%s' 失败: %v", cmd.Command, err)
			}
		}
	}
	
	return nil
}

// ExecuteCommand 执行单个系统命令
func ExecuteCommand(command string, echoOn bool) error {
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