package utils

import (
	"runtime"
	"strings"
)

// GetOS 获取当前操作系统
func GetOS() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	default:
		return runtime.GOOS
	}
}

// MatchPlatform 检查平台是否匹配
func MatchPlatform(platform, currentOS string) bool {
	// default块总是匹配
	if platform == "default" {
		return true
	}
	
	// 支持多个平台名称
	platforms := strings.Split(platform, ",")
	for _, p := range platforms {
		p = strings.TrimSpace(p)
		if p == currentOS {
			return true
		}
		
		// 支持通用Unix平台
		if p == "unix" && (currentOS == "linux" || currentOS == "macos") {
			return true
		}
	}
	return false
}