package types

// TaskCommand 是任务中的单条 shell 命令，保留行号便于报错。
type TaskCommand struct {
	Raw  string `json:"raw"`
	Line int    `json:"line"`
}


type Task struct {
	Name     string        `json:"name"`
	Desc     string        `json:"desc,omitempty"`
	On       []string      `json:"on,omitempty"`   // 为空表示全平台
	Deps     []string      `json:"deps,omitempty"` // 依赖的任务名
	Shell    string        `json:"shell,omitempty"`
	Workdir  string        `json:"workdir,omitempty"`
	EchoOn   *bool         `json:"echo_on,omitempty"` // nil 表示继承全局
	Commands []TaskCommand `json:"commands"`
	Line     int           `json:"line"`
}

// EchoOrDefault 返回任务的回显设置，nil 时继承全局。
func (t *Task) EchoOrDefault(global bool) bool {
	if t.EchoOn == nil {
		return global
	}
	return *t.EchoOn
}

// ZRunScript 只承载 v2，不兼容 v1。
type ZRunScript struct {
	File          string            `json:"file,omitempty"`
	SyntaxVersion string            `json:"syntax_version"`
	EchoOn        bool              `json:"echo_on"`
	Vars          map[string]string `json:"vars,omitempty"`
	Env           map[string]string `json:"env,omitempty"`
	Tasks         []Task            `json:"tasks"`
}
