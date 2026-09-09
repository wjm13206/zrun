package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "test.zr")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestParseBasic(t *testing.T) {
	p := writeTemp(t, "@syntax 2.0\nvar APP = \"demo\"\ntask build --desc \"构建\" --on [linux] --deps [clean] {\n echo hi\n}\ntask clean {\n echo clean\n}\n")
	s, err := ParseScript(p)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	if len(s.Tasks) != 2 || s.Vars["APP"] != "demo" {
		t.Fatalf("解析结果不对: %+v", s)
	}
	if len(s.Tasks[0].On) != 1 || s.Tasks[0].On[0] != "linux" {
		t.Errorf("on 解析不对: %+v", s.Tasks[0].On)
	}
}

func TestParseMissingSyntax(t *testing.T) {
	p := writeTemp(t, "task build {\n echo hi\n}\n")
	if _, err := ParseScript(p); err == nil {
		t.Error("缺少 @syntax 应报错")
	}
}

func TestParseLegacyRejected(t *testing.T) {
	p := writeTemp(t, "@syntax 2.0\n@windows {\n echo hi\n}\n")
	if _, err := ParseScript(p); err == nil {
		t.Error("v1 语法应被拒绝")
	}
}

func TestParseBadDep(t *testing.T) {
	p := writeTemp(t, "@syntax 2.0\ntask a --deps [missing] {\n echo hi\n}\n")
	if _, err := ParseScript(p); err == nil {
		t.Error("依赖不存在应报错")
	}
}

func TestParseDupTask(t *testing.T) {
	p := writeTemp(t, "@syntax 2.0\ntask a {\n echo hi\n}\ntask a {\n echo hi\n}\n")
	if _, err := ParseScript(p); err == nil {
		t.Error("重复 task 应报错")
	}
}
