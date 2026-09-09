// Package utils 提供了 zrun 的工具函数。
package utils

import (
	"runtime"
	"strings"
)

var currentOS string
var currentArch string

func init() {
	currentOS = NormalizeOS(runtime.GOOS)
	currentArch = NormalizeArch(runtime.GOARCH)
}

// GetOS 返回归一化后的操作系统名：windows / linux / macos。
func GetOS() string {
	return currentOS
}

// GetArch 返回归一化后的架构名：amd64 / 386 / arm / arm64。
func GetArch() string {
	return currentArch
}

// SetOSArch 仅供测试注入，生产代码不要调用。
func SetOSArch(os, arch string) {
	currentOS = os
	currentArch = arch
}

// NormalizeOS 将 runtime.GOOS 归一化为 windows / linux / macos。
// 未知值原样小写返回，便于透传和报错。
func NormalizeOS(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin", "macos", "osx":
		return "macos"
	default:
		return strings.ToLower(strings.TrimSpace(s))
	}
}

// NormalizeArch 将多种架构别名归一化为 amd64 / 386 / arm / arm64。
// 兼容老语法的 64 / 32 写法：64 -> amd64，32 -> 386。
func NormalizeArch(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "amd64", "x64", "x86_64", "64":
		return "amd64"
	case "386", "x86", "x32", "32":
		return "386"
	case "arm":
		return "arm"
	case "arm64", "aarch64":
		return "arm64"
	default:
		return strings.ToLower(strings.TrimSpace(s))
	}
}

// MatchPlatform 兼容 v1 的单字符串匹配，大小写不敏感。
// 支持逗号分隔、os/arch 组合、unix 展开、default/any/* 常真。
func MatchPlatform(platform string) bool {
	return MatchPlatformWith(currentOS, currentArch, platform)
}

// MatchPlatformWith 纯函数版本，便于测试。
func MatchPlatformWith(os, arch, platform string) bool {
	os = NormalizeOS(os)
	arch = NormalizeArch(arch)
	p := strings.TrimSpace(platform)
	if p == "" {
		return false
	}
	if strings.EqualFold(p, "default") || strings.EqualFold(p, "any") || p == "*" {
		return true
	}
	for _, part := range strings.Split(p, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if matchSingleSelector(os, arch, part) {
			return true
		}
	}
	return false
}

// MatchAny 供 v2 任务 --on 使用：空列表表示全平台。
func MatchAny(selectors []string) bool {
	if len(selectors) == 0 {
		return true
	}
	for _, s := range selectors {
		if MatchPlatform(s) {
			return true
		}
	}
	return false
}

// MatchAnyWith 纯函数版本，便于测试。
func MatchAnyWith(os, arch string, selectors []string) bool {
	if len(selectors) == 0 {
		return true
	}
	for _, s := range selectors {
		if MatchPlatformWith(os, arch, s) {
			return true
		}
	}
	return false
}

func matchSingleSelector(os, arch, sel string) bool {
	sel = strings.TrimSpace(sel)
	if sel == "" {
		return false
	}
	if strings.EqualFold(sel, "default") || strings.EqualFold(sel, "any") || sel == "*" {
		return true
	}
	// os/arch 组合，例如 windows/amd64、unix/arm64。
	if strings.Contains(sel, "/") {
		parts := strings.Split(sel, "/")
		if len(parts) != 2 {
			return false
		}
		osPart := strings.ToLower(strings.TrimSpace(parts[0]))
		archPart := NormalizeArch(parts[1])
		return matchOS(os, osPart) && archPart == arch
	}
	// 裸架构选择器，例如 64 / amd64。
	if isArchToken(sel) {
		return NormalizeArch(sel) == arch
	}
	// 裸 OS 选择器。
	return matchOS(os, strings.ToLower(sel))
}

func matchOS(current, want string) bool {
	want = strings.ToLower(strings.TrimSpace(want))
	if want == current {
		return true
	}
	// darwin 别名容错。
	if want == "darwin" && current == "macos" {
		return true
	}
	// unix 展开为 linux 或 macos。
	if want == "unix" && (current == "linux" || current == "macos") {
		return true
	}
	return false
}

func isArchToken(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "amd64", "x64", "x86_64", "64",
		"386", "x86", "x32", "32",
		"arm", "arm64", "aarch64":
		return true
	default:
		return false
	}
}
