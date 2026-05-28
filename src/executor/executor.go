package executor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"zrun/src/types"
	"zrun/src/utils"
)

// 并发执行多个命令的新函数
func ExecuteScriptConcurrent(script *types.ZRunScript) error {
	// 分离echo命令和普通命令，因为echo命令会影响全局状态
	var echoCmds []types.ScriptCommand
	var platformCmds []types.ScriptCommand

	for _, cmd := range script.Commands {
		if cmd.Type == "echo" {
			echoCmds = append(echoCmds, cmd)
		} else if cmd.Type == "platform" && utils.MatchPlatform(cmd.Param) {
			platformCmds = append(platformCmds, cmd)
		}
	}

	// 先执行echo命令
	for _, cmd := range echoCmds {
		switch cmd.Param {
		case "off":
			script.EchoOn = false
		case "on":
			script.EchoOn = true
		}
	}

	// 并发执行命令
	var wg sync.WaitGroup
	errChan := make(chan error, len(platformCmds))
	var outputMu sync.Mutex

	for _, cmd := range platformCmds {
		wg.Add(1)
		go func(c types.ScriptCommand) {
			defer wg.Done()
			if err := ExecuteCommand(c.Command, script.EchoOn, &outputMu); err != nil {
				errChan <- err
			}
		}(cmd)
	}

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func ExecuteCommand(command string, echoOn bool, mu *sync.Mutex) error {
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
		mu.Lock()
		fmt.Printf("$ %s\n", command)
		mu.Unlock()
	}

	// 执行
	return cmd.Run()
}