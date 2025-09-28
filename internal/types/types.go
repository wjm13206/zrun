package types

// ScriptCommand 表示脚本中的一个命令
type ScriptCommand struct {
	Command string
	Type    string // "echo" 或 "platform"
	Param   string // platform名称或其他参数
}

// ZRunScript 表示整个zrun脚本
type ZRunScript struct {
	Commands []ScriptCommand
	EchoOn   bool // 控制是否显示命令
}