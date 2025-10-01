// Package types 定义了 zrun 使用的核心数据结构。
package types

// ScriptCommand 表示 zrun 脚本中的单个命令。
// 它包含命令文本、类型标识符和附加参数。
type ScriptCommand struct {
	// Command 是要执行的实际命令
	Command string

	// Type 标识命令的类型:
	// - "echo": 回显控制命令 (@echo on/off)
	// - "platform": 平台特定命令
	Type string

	// Param 包含附加参数:
	// - 对于 "echo" 类型: "on" 或 "off"
	// - 对于 "platform" 类型: 平台标识符 (例如 "windows", "linux")
	Param string
}

// ZRunScript 表示一个完整的 zrun 脚本。
// 它包含命令列表和当前的回显状态。
type ZRunScript struct {
	// Commands 是按出现顺序排列的脚本中所有命令的列表
	Commands []ScriptCommand

	// EchoOn 控制命令在执行前是否回显到标准输出
	// true = 回显命令（默认），false = 不回显命令
	EchoOn bool
}
