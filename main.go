// Package main provides the entry point for the zrun cross-platform script executor.
package main

import (
	"fmt"
	"os"
	"zrun/internal/executor"
	"zrun/internal/parser"
)

// main is the entry point of the zrun application.
// It parses command line arguments, loads and parses the script file,
// and then executes the parsed script.
func main() {
	// Check if script file argument is provided
	if len(os.Args) < 2 {
		fmt.Println("用法: zrun <脚本文件.zr>")
		os.Exit(1)
	}

	// Get script filename from command line argument
	filename := os.Args[1]

	// Parse the script file
	script, err := parser.ParseScript(filename)
	if err != nil {
		fmt.Printf("解析脚本错误: %v\n", err)
		os.Exit(1)
	}

	// Execute the parsed script
	err = executor.ExecuteScript(script)
	if err != nil {
		fmt.Printf("执行脚本错误: %v\n", err)
		os.Exit(1)
	}
}
