package types

type ScriptCommand struct {
	Command string `json:"command"`
	Type    string `json:"type"`
	Param   string `json:"param"`
}



	
type ZRunScript struct {
	Commands []ScriptCommand `json:"commands"`
		// true = 回显命令（默认），false = 不回显命令
	EchoOn bool `json:"echo_on"`
}
