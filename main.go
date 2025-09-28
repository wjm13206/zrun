package main

import (
	"fmt"
	"os"
	"zrun/internal/executor"
	"zrun/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: zrun <脚本文件.zr>")
		os.Exit(1)
	}

	filename := os.Args[1]
	script, err := parser.ParseScript(filename)
	if err != nil {
		fmt.Printf("解析脚本错误: %v\n", err)
		os.Exit(1)
	}

	err = executor.ExecuteScript(script)
	if err != nil {
		fmt.Printf("执行脚本错误: %v\n", err)
		os.Exit(1)
	}
}