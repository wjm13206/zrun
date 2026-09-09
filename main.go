package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"zrun/src/executor"
	"zrun/src/parser"
	"zrun/src/utils"
)

const version = "2026.05.28"
const SyntaxVersion = "2.0"

func main() {
	// 切分 -- 后的透传参数为 $ARGS
	rawArgs := os.Args[1:]
	var passThrough string
	if i := indexOf(rawArgs, "--"); i >= 0 {
		passThrough = strings.Join(rawArgs[i+1:], " ")
		rawArgs = rawArgs[:i]
	}

	fs := flag.NewFlagSet("zrun", flag.ContinueOnError)
	showList := fs.Bool("list", false, "列出任务")
	dryRun := fs.Bool("dry-run", false, "只打印不执行")
	showTime := fs.Bool("time", true, "显示总耗时")
	shortV := fs.Bool("v", false, "显示版本")
	longV := fs.Bool("version", false, "显示版本")
	shortU := fs.Bool("u", false, "检查更新")
	longU := fs.Bool("update", false, "检查更新")
	showHelp := fs.Bool("h", false, "显示帮助")
	longHelp := fs.Bool("help", false, "显示帮助")
	// 兼容 --no-time 写法
	noTime := fs.Bool("no-time", false, "不显示总耗时")
	_ = fs.Parse(rawArgs)

	if *showHelp || *longHelp {
		printUsage()
		return
	}
	if *shortV || *longV {
		fmt.Printf("zrun version: %s (syntax %s)\n", version, SyntaxVersion)
		return
	}
	if *shortU || *longU {
		utils.CheckSyntaxUpdates(version, SyntaxVersion)
		return
	}

	pos := fs.Args()
	if len(pos) < 1 {
		printUsage()
		os.Exit(1)
	}
	filename := pos[0]
	tasks := pos[1:]
	if *noTime {
		*showTime = false
	}

	if *showList {
		script, err := parser.ParseScript(filename)
		if err != nil {
			fmt.Printf("解析错误: %v\n", err)
			os.Exit(1)
		}
		for _, t := range script.Tasks {
			on := strings.Join(t.On, ",")
			if on == "" {
				on = "all"
			}
			deps := strings.Join(t.Deps, ",")
			if deps == "" {
				deps = "-"
			}
			fmt.Printf("%-16s on=[%s] deps=[%s] %s\n", t.Name, on, deps, t.Desc)
		}
		return
	}

	var start time.Time
	if *showTime {
		start = time.Now()
	}

	script, err := parser.ParseScript(filename)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		os.Exit(1)
	}

	err = executor.ExecuteScript(script, executor.Options{
		Tasks:  tasks,
		Args:   passThrough,
		DryRun: *dryRun,
	})
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		os.Exit(1)
	}

	if *showTime {
		fmt.Printf("\n执行完成，总耗时: %v\n", time.Since(start))
	}
}

func indexOf(args []string, target string) int {
	for i, a := range args {
		if a == target {
			return i
		}
	}
	return -1
}

func printUsage() {
	fmt.Printf(`zrun %s (syntax %s) - 跨平台任务脚本

用法:
  zrun [选项] <文件.zr> [任务...] [-- 透传参数]

选项:
  --list            列出任务
  --dry-run         只打印不执行
  --time / --no-time 是否显示总耗时（默认显示）
  -v, --version     显示版本
  -u, --update      检查更新
  -h, --help        显示帮助

示例:
  zrun script.zr
  zrun script.zr build
  zrun --list script.zr
  zrun --dry-run script.zr build
  zrun script.zr run -- --port 8080
`, version, SyntaxVersion)
}
