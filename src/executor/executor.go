package executor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"zrun/src/types"
	"zrun/src/utils"
)

// Options 控制脚本执行。
type Options struct {
	// Tasks 为空表示运行全部任务，否则只运行指定任务及其依赖。
	Tasks []string
	// Args 透传给 $ARGS / ${ARGS}。
	Args string
	// DryRun 只打印不执行。
	DryRun bool
}

// ExecuteScript 按文件顺序拓扑执行，默认顺序、遇错即停。
// 不匹配 --on 的任务自动跳过；被跳过任务的依赖视为已满足。
func ExecuteScript(script *types.ZRunScript, opts Options) error {
	byName := make(map[string]*types.Task, len(script.Tasks))
	for i := range script.Tasks {
		byName[script.Tasks[i].Name] = &script.Tasks[i]
	}

	var requested []string
	if len(opts.Tasks) == 0 {
		for _, t := range script.Tasks {
			requested = append(requested, t.Name)
		}
	} else {
		for _, name := range opts.Tasks {
			if _, ok := byName[name]; !ok {
				return fmt.Errorf("task %q 不存在", name)
			}
			requested = append(requested, name)
		}
	}

	ordered, err := resolveOrder(byName, requested)
	if err != nil {
		return err
	}

	// 平台过滤：跳过不匹配的任务
	var plan []*types.Task
	var skipped []string
	for _, t := range ordered {
		if !utils.MatchAny(t.On) {
			skipped = append(skipped, t.Name)
			continue
		}
		plan = append(plan, t)
	}
	for _, s := range skipped {
		fmt.Printf("跳过 task %s（平台不匹配）\n", s)
	}
	if len(plan) == 0 {
		fmt.Println("没有可执行的任务")
		return nil
	}

	lookupBase := func(name string) (string, bool) {
		if v, ok := script.Vars[name]; ok {
			return v, true
		}
		if v, ok := script.Env[name]; ok {
			return v, true
		}
		switch name {
		case "OS":
			return utils.GetOS(), true
		case "ARCH":
			return utils.GetArch(), true
		case "ARGS":
			return opts.Args, true
		}
		if v, ok := os.LookupEnv(name); ok {
			return v, true
		}
		return "", false
	}

	for _, t := range plan {
		echoOn := t.EchoOrDefault(script.EchoOn)
		fmt.Printf("==> %s\n", t.Name)

		workdir := ""
		if t.Workdir != "" {
			workdir = utils.ExpandVars(t.Workdir, lookupBase)
		}
		// 任务 env：script.Env 全量继承，值同样做插值
		taskEnv := os.Environ()
		for k, v := range script.Env {
			taskEnv = append(taskEnv, k+"="+utils.ExpandVars(v, lookupBase))
		}

		for _, c := range t.Commands {
			expanded := utils.ExpandVars(c.Raw, lookupBase)
			if echoOn || opts.DryRun {
				fmt.Printf("$ %s\n", expanded)
			}
			if opts.DryRun {
				continue
			}
			if err := runOne(expanded, t.Shell, workdir, taskEnv); err != nil {
				return fmt.Errorf("task %s 第 %d 行失败 %q: %w", t.Name, c.Line, expanded, err)
			}
		}
	}
	return nil
}

// resolveOrder 对 requested 做依赖闭包 + 拓扑排序，保持文件顺序优先。
func resolveOrder(byName map[string]*types.Task, requested []string) ([]*types.Task, error) {
	// 文件顺序索引，用于稳定排序
	orderIdx := make(map[string]int, len(byName))
	idx := 0
	// byName 丢失顺序，调用方保证 requested 已按文件序；这里用 DFS 保持依赖先行
	const (
		white = 0
		gray  = 1
		black = 2
	)
	state := make(map[string]int, len(byName))
	var out []*types.Task

	var visit func(name string, stack []string) error
	visit = func(name string, stack []string) error {
		switch state[name] {
		case black:
			return nil
		case gray:
			cycle := append(append([]string{}, stack...), name)
			return fmt.Errorf("task 依赖成环：%s", strings.Join(cycle, " -> "))
		}
		t, ok := byName[name]
		if !ok {
			return fmt.Errorf("task %q 不存在", name)
		}
		state[name] = gray
		// 依赖按声明顺序访问
		for _, d := range t.Deps {
			if err := visit(d, append(stack, name)); err != nil {
				return err
			}
		}
		state[name] = black
		out = append(out, t)
		_ = orderIdx
		_ = idx
		return nil
	}
	for _, name := range requested {
		if err := visit(name, nil); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func runOne(command, shell, workdir string, env []string) error {
	cmd := buildCmd(command, shell)
	if workdir != "" {
		cmd.Dir = workdir
	}
	if env != nil {
		cmd.Env = env
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildCmd(command, shell string) *exec.Cmd {
	s := strings.ToLower(strings.TrimSpace(shell))
	switch s {
	case "", "auto":
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", command)
		}
		return exec.Command("sh", "-c", command)
	case "cmd":
		return exec.Command("cmd", "/C", command)
	case "powershell", "pwsh":
		return exec.Command(shell, "-Command", command)
	case "sh", "bash", "zsh":
		return exec.Command(s, "-c", command)
	default:
		// 未知 shell 按 sh -c 语义透传首词，避免静默失败
		parts := strings.Fields(shell)
		if len(parts) == 0 {
			return exec.Command("sh", "-c", command)
		}
		return exec.Command(parts[0], append(parts[1:], command)...)
	}
}
