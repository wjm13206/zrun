package main

import (
	"fmt"
	"os"
	"zrun/src/executor"
	"zrun/src/parser"
)

const version = "1.0"

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

	filename := os.Args[1]

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
}
