package utils

import "testing"

func TestNormalizeOS(t *testing.T) {
	cases := map[string]string{
		"windows": "windows",
		"linux":   "linux",
		"darwin":  "macos",
		"macos":   "macos",
		"Windows": "windows",
	}
	for in, want := range cases {
		if got := NormalizeOS(in); got != want {
			t.Errorf("NormalizeOS(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestNormalizeArch(t *testing.T) {
	cases := map[string]string{
		"amd64": "amd64",
		"64":    "amd64",
		"x64":   "amd64",
		"386":   "386",
		"32":    "386",
		"arm":   "arm",
		"arm64": "arm64",
	}
	for in, want := range cases {
		if got := NormalizeArch(in); got != want {
			t.Errorf("NormalizeArch(%q)=%q, want %q", in, got, want)
		}
	}
}

func TestMatchPlatformWith(t *testing.T) {
	if !MatchPlatformWith("linux", "amd64", "linux") {
		t.Error("linux 应匹配")
	}
	if !MatchPlatformWith("linux", "amd64", "unix") {
		t.Error("linux 应匹配 unix")
	}
	if MatchPlatformWith("windows", "amd64", "unix") {
		t.Error("windows 不应匹配 unix")
	}
	if !MatchPlatformWith("windows", "amd64", "windows/amd64") {
		t.Error("windows/amd64 应匹配")
	}
	if MatchPlatformWith("windows", "386", "windows/amd64") {
		t.Error("arch 不一致不应匹配")
	}
	if !MatchPlatformWith("linux", "amd64", "default") {
		t.Error("default 应常真")
	}
	if !MatchPlatformWith("windows", "amd64", "WINDOWS") {
		t.Error("应大小写不敏感")
	}
	if !MatchPlatformWith("windows", "amd64", "64") {
		t.Error("64 别名应匹配 amd64")
	}
}

func TestMatchAnyWith(t *testing.T) {
	if !MatchAnyWith("linux", "amd64", nil) {
		t.Error("空 on 应全平台匹配")
	}
	if !MatchAnyWith("linux", "amd64", []string{"windows", "linux"}) {
		t.Error("多选一命中应匹配")
	}
}

func TestExpandVars(t *testing.T) {
	got := ExpandVars("hi $A/${B} $$C", func(n string) (string, bool) {
		if n == "A" {
			return "1", true
		}
		if n == "B" {
			return "2", true
		}
		return "", false
	})
	if got != "hi 1/2 $C" {
		t.Errorf("展开结果=%q", got)
	}
}

func TestCheckUpdate(t *testing.T) {
	has, _ := CheckUpdate("26.1.0", "2.0", &RemoteVersionInfo{LatestVersion: "26.1.0", LatestSyntaxVersion: "2.0"})
	if has {
		t.Error("同版本不应提示更新")
	}
	has, inc := CheckUpdate("26.1.0", "1.2", &RemoteVersionInfo{LatestVersion: "26.2.0", LatestSyntaxVersion: "2.0"})
	if !has || !inc {
		t.Error("旧版本应提示更新且语法不兼容")
	}
	// 数字逐段比较：26.10.0 > 26.9.9，字符串比较会误判
	has, _ = CheckUpdate("26.10.0", "2.0", &RemoteVersionInfo{LatestVersion: "26.9.9", LatestSyntaxVersion: "2.0"})
	if has {
		t.Error("26.10.0 不应比 26.9.9 旧")
	}
	has, _ = CheckUpdate("26.9.9", "2.0", &RemoteVersionInfo{LatestVersion: "26.10.0", LatestSyntaxVersion: "2.0"})
	if !has {
		t.Error("26.9.9 应比 26.10.0 旧")
	}
}
