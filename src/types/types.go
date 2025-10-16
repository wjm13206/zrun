package types




type ScriptCommand struct {

	Command string



	Type string

	Param string
}


type ZRunScript struct {

	
	Commands []ScriptCommand

	//控制命令在执行前是否回显到标准输出
	// true = 回显命令（默认），false = 不回显命令
	EchoOn bool
}
