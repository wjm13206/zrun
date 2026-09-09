package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"zrun/src/types"
)

func ParseScript(filename string) (*types.ZRunScript, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开脚本 %s 失败: %w", filename, err)
	}
	defer file.Close()

	script := &types.ZRunScript{
		File:          filename,
		EchoOn:        true,
		Vars:          make(map[string]string),
		Env:           make(map[string]string),
		Tasks:         nil,
		SyntaxVersion: "",
	}

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	type rawLine struct {
		no   int
		text string
	}
	var lines []rawLine
	for scanner.Scan() {
		lines = append(lines, rawLine{no: len(lines) + 1, text: scanner.Text()})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取脚本 %s 失败: %w", filename, err)
	}

	taskIndex := make(map[string]int) // 任务名 -> 下标，用于判重
	curIdx := -1                      // -1 表示当前不在 task 块内
	seenSyntax := false
	firstMeaningfulSeen := false

	fail := func(no int, format string, args ...interface{}) error {
		return fmt.Errorf("%s:%d: %s", filename, no, fmt.Sprintf(format, args...))
	}

	for _, rl := range lines {
		no := rl.no
		line := strings.TrimSpace(rl.text)
		if line == "" {
			continue
		}
		// 全行注释：# 或 //
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// 老语法拦截
		if isLegacyBlock(line) {
			return nil, fail(no, "v1 语法已废弃（%s），请改写为 task 形式，例如 task main --on [%s] { ... }", line, legacyHint(line))
		}

		// 块结束
		if line == "}" {
			if curIdx < 0 {
				return nil, fail(no, "多余的 }，当前不在任何 task 内")
			}
			curIdx = -1
			continue
		}

		// 块内：全部视为 shell 命令
		if curIdx >= 0 {
			script.Tasks[curIdx].Commands = append(
				script.Tasks[curIdx].Commands,
				types.TaskCommand{Raw: line, Line: no},
			)
			continue
		}

		// 顶层指令
		lower := strings.ToLower(line)
		switch {
		case strings.HasPrefix(lower, "@syntax"):
			if seenSyntax {
				return nil, fail(no, "@syntax 重复声明")
			}
			ver, err := parseSyntaxVersion(line)
			if err != nil {
				return nil, fail(no, "%v", err)
			}
			if !strings.HasPrefix(ver, "2.") {
				return nil, fail(no, "@syntax %s 不受支持，仅支持 2.x", ver)
			}
			script.SyntaxVersion = ver
			seenSyntax = true
			firstMeaningfulSeen = true
		case strings.HasPrefix(lower, "@echo"):
			if !firstMeaningfulSeen && !seenSyntax {
				return nil, fail(no, "缺少 @syntax 2.x 声明，v2 文件必须以 @syntax 2.x 开头")
			}
			on, err := parseEcho(line)
			if err != nil {
				return nil, fail(no, "%v", err)
			}
			script.EchoOn = on
			firstMeaningfulSeen = true
		case strings.HasPrefix(line, "var ") || strings.HasPrefix(line, "var\t"):
			if !seenSyntax {
				return nil, fail(no, "缺少 @syntax 2.x 声明，v2 文件必须以 @syntax 2.x 开头")
			}
			name, value, err := parseAssignment(line, "var")
			if err != nil {
				return nil, fail(no, "%v", err)
			}
			if _, exists := script.Vars[name]; exists {
				return nil, fail(no, "var %s 重复定义", name)
			}
			script.Vars[name] = value
			firstMeaningfulSeen = true
		case strings.HasPrefix(line, "env ") || strings.HasPrefix(line, "env\t"):
			if !seenSyntax {
				return nil, fail(no, "缺少 @syntax 2.x 声明，v2 文件必须以 @syntax 2.x 开头")
			}
			name, value, err := parseAssignment(line, "env")
			if err != nil {
				return nil, fail(no, "%v", err)
			}
			script.Env[name] = value
			firstMeaningfulSeen = true
		case strings.HasPrefix(line, "task ") || strings.HasPrefix(line, "task\t"):
			if !seenSyntax {
				return nil, fail(no, "缺少 @syntax 2.x 声明，v2 文件必须以 @syntax 2.x 开头")
			}
			t, err := parseTaskHeader(line, no)
			if err != nil {
				return nil, fail(no, "%v", err)
			}
			if _, exists := taskIndex[t.Name]; exists {
				return nil, fail(no, "task %s 重复定义", t.Name)
			}
			taskIndex[t.Name] = len(script.Tasks)
			script.Tasks = append(script.Tasks, *t)
			curIdx = len(script.Tasks) - 1
			// curIdx 仅用于判空，命令追加直接操作 slice 元素
			firstMeaningfulSeen = true
		default:
			return nil, fail(no, "无法识别的顶层指令 %q，仅允许 @syntax/@echo/var/env/task", line)
		}
	}

	if curIdx >= 0 {
		t := script.Tasks[curIdx]
		return nil, fmt.Errorf("%s:%d: task %s 缺少闭合 }", filename, t.Line, t.Name)
	}
	if !seenSyntax {
		return nil, fmt.Errorf("%s: 缺少 @syntax 2.x 声明，v2 文件必须以 @syntax 2.x 开头", filename)
	}
	if len(script.Tasks) == 0 {
		return nil, fmt.Errorf("%s: 未定义任何 task", filename)
	}
	// 校验 deps 指向存在
	for _, t := range script.Tasks {
		for _, d := range t.Deps {
			if _, ok := taskIndex[d]; !ok {
				return nil, fmt.Errorf("%s:%d: task %s 依赖了不存在的 task %s", filename, t.Line, t.Name, d)
			}
			if d == t.Name {
				return nil, fmt.Errorf("%s:%d: task %s 不能依赖自身", filename, t.Line, t.Name)
			}
		}
	}
	return script, nil
}

// isLegacyBlock 识别 v1 平台块开头。
func isLegacyBlock(line string) bool {
	lower := strings.ToLower(strings.TrimSpace(line))
	if !strings.HasPrefix(lower, "@") {
		return false
	}
	// 去掉 @ 后取第一个词
	rest := strings.TrimSpace(lower[1:])
	end := len(rest)
	for i, c := range rest {
		if c == ' ' || c == '\t' || c == '{' || c == ',' || c == '/' {
			end = i
			break
		}
	}
	word := rest[:end]
	switch word {
	case "windows", "linux", "macos", "darwin", "unix", "default":
		return true
	default:
		return false
	}
}

func legacyHint(line string) string {
	lower := strings.ToLower(strings.TrimSpace(line))
	rest := strings.TrimSpace(lower[1:])
	end := len(rest)
	for i, c := range rest {
		if c == ' ' || c == '\t' || c == '{' {
			end = i
			break
		}
	}
	return strings.TrimSpace(rest[:end])
}

func parseSyntaxVersion(line string) (string, error) {
	rest := strings.TrimSpace(line[len("@syntax"):])
	if rest == "" {
		return "", fmt.Errorf("@syntax 缺少版本号，例如 @syntax 2.0")
	}
	fields := strings.Fields(rest)
	return fields[0], nil
}

func parseEcho(line string) (bool, error) {
	rest := strings.TrimSpace(line[len("@echo"):])
	switch strings.ToLower(rest) {
	case "on":
		return true, nil
	case "off":
		return false, nil
	default:
		return false, fmt.Errorf("@echo 后仅允许 on/off，例如 @echo off")
	}
}

func parseAssignment(line, keyword string) (string, string, error) {
	rest := strings.TrimSpace(line[len(keyword):])
	eq := strings.Index(rest, "=")
	if eq < 0 {
		return "", "", fmt.Errorf("%s 缺少 =，例如 %s NAME = value", keyword, keyword)
	}
	name := strings.TrimSpace(rest[:eq])
	value := strings.TrimSpace(rest[eq+1:])
	if !isValidVarName(name) {
		return "", "", fmt.Errorf("%s 名称 %q 非法，应为 [A-Za-z_][A-Za-z0-9_]*", keyword, name)
	}
	if value == "" {
		return "", "", fmt.Errorf("%s %s 的值不能为空", keyword, name)
	}
	return name, stripQuotes(value), nil
}

func parseTaskHeader(line string, no int) (*types.Task, error) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasSuffix(trimmed, "{") {
		return nil, fmt.Errorf("task 缺少结尾 {，例如 task build {")
	}
	inner := strings.TrimSpace(trimmed[:len(trimmed)-1])
	rest := strings.TrimSpace(inner[len("task"):])
	if rest == "" {
		return nil, fmt.Errorf("task 缺少名称，例如 task build {")
	}
	// 任务名
	nameEnd := len(rest)
	for i, c := range rest {
		if c == ' ' || c == '\t' {
			nameEnd = i
			break
		}
	}
	name := rest[:nameEnd]
	if !isValidTaskName(name) {
		return nil, fmt.Errorf("task 名称 %q 非法，应为 [A-Za-z_][A-Za-z0-9_-]*", name)
	}
	optStr := strings.TrimSpace(rest[nameEnd:])
	t := &types.Task{Name: name, Line: no}
	if optStr == "" {
		return t, nil
	}
	opts, err := parseTaskOptions(optStr)
	if err != nil {
		return nil, err
	}
	for k, v := range opts {
		switch k {
		case "desc":
			t.Desc = v.str
		case "shell":
			t.Shell = v.str
		case "workdir":
			t.Workdir = v.str
		case "on":
			t.On = v.list
		case "deps":
			t.Deps = v.list
		case "echo":
			switch strings.ToLower(v.str) {
			case "on":
				b := true
				t.EchoOn = &b
			case "off":
				b := false
				t.EchoOn = &b
			default:
				return nil, fmt.Errorf("--echo 仅允许 on/off")
			}
		}
	}
	return t, nil
}

type optValue struct {
	str  string
	list []string
}

func parseTaskOptions(s string) (map[string]optValue, error) {
	out := make(map[string]optValue)
	i := 0
	for i < len(s) {
		// 跳空格
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
		if i >= len(s) {
			break
		}
		if !(i+1 < len(s) && s[i] == '-' && s[i+1] == '-') {
			return nil, fmt.Errorf("task 选项必须以 -- 开头，位置 %d 附近：%q", i, s)
		}
		i += 2
		ks := i
		for i < len(s) && s[i] != ' ' && s[i] != '\t' && s[i] != '=' {
			i++
		}
		key := s[ks:i]
		if key == "" {
			return nil, fmt.Errorf("task 选项名为空")
		}
		if _, dup := out[key]; dup {
			return nil, fmt.Errorf("--%s 重复出现", key)
		}
		switch key {
		case "desc", "on", "deps", "shell", "workdir", "echo":
		default:
			return nil, fmt.Errorf("未知选项 --%s，仅允许 --desc/--on/--deps/--shell/--workdir/--echo", key)
		}
		for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '=') {
			i++
		}
		if i >= len(s) {
			return nil, fmt.Errorf("--%s 缺少值", key)
		}
		var str string
		var list []string
		switch s[i] {
		case '[':
			end := strings.IndexByte(s[i:], ']')
			if end < 0 {
				return nil, fmt.Errorf("--%s 的 [ 缺少闭合 ]", key)
			}
			raw := s[i+1 : i+end]
			list = splitList(raw)
			str = strings.Join(list, ",")
			i += end + 1
		case '"', '\'':
			q := s[i]
			end := -1
			for j := i + 1; j < len(s); j++ {
				if s[j] == q && s[j-1] != '\\' {
					end = j
					break
				}
			}
			if end < 0 {
				return nil, fmt.Errorf("--%s 的引号未闭合", key)
			}
			str = s[i+1 : end]
			if key == "on" || key == "deps" {
				list = splitList(str)
			}
			i = end + 1
		default:
			js := i
			for i < len(s) && s[i] != ' ' && s[i] != '\t' {
				i++
			}
			str = s[js:i]
			if key == "on" || key == "deps" {
				list = splitList(str)
			}
		}
		if (key == "on" || key == "deps") && len(list) == 0 && str != "" {
			list = []string{str}
		}
		out[key] = optValue{str: str, list: list}
	}
	return out, nil
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = stripQuotes(strings.TrimSpace(p))
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func isValidVarName(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i == 0 {
			if !(c == '_' || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
				return false
			}
		} else if !(c == '_' || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

func isValidTaskName(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i == 0 {
			if !(c == '_' || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
				return false
			}
		} else if !(c == '_' || c == '-' || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}
