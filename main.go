package main

import (
	"fmt"
	"os"
	"time"
	"zrun/src/executor"
	"zrun/src/parser"
	"zrun/src/utils"
)

const version = "2025.10.26"
const SyntaxVersion = "1.1"

// 是否启用测量
var enablePerfMeasurement = false

// 解析命令行参数，加载并解析脚本文件，然后执行
func main() {
	// 检查参数
	if len(os.Args) < 2 {
		fmt.Println("用法: ./zrun <文件.zr>")
		os.Exit(1)
	}

	// 检查版本参数
	if os.Args[1] == "-version" || os.Args[1] == "-v" {
		fmt.Printf("zrun: %s\n", version)
		os.Exit(0)
	}

	// 检查语法版本更新参数
	if os.Args[1] == "-u" || os.Args[1] == "--update" {
		utils.CheckSyntaxUpdates(version, SyntaxVersion)
		os.Exit(0)
	}

	filename := os.Args[1]

	// 记录开始时间
	var start time.Time
	if enablePerfMeasurement {
		start = time.Now()
	}

	// 解析
	script, err := parser.ParseScript(filename)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		os.Exit(1)
	}

	// 执行
	err = executor.ExecuteScript(script)
	if err != nil {
		fmt.Printf("执行错误: %v\n", err)
		os.Exit(1)
	}

	// 输出执行时间
	if enablePerfMeasurement {
		elapsed := time.Since(start)
		fmt.Printf("\n执行完成，总耗时: %v\n", elapsed)
	}
}
