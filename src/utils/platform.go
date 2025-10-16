// Package utils 提供了 zrun 的工具函数。
package utils

import (
	"runtime"
	"strings"
)

// GetOS 获取当前操作系统
// 返回值"windows", "linux", "macos"
func GetOS() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	default:
		return runtime.GOOS // 默认
	}
}

// 检查平台是否匹配
// 参数:
//   - platform: 脚本中声明的平台标识符
//   - currentOS: 当前操作系统标识
//
// 返回值:
//   - bool: 平台是否匹配

func MatchPlatform(platform, currentOS string) bool {
	// default匹配
	if platform == "default" {
		return true
	}


	
	platforms := strings.Split(platform, ",")
	for _, p := range platforms {
		p = strings.TrimSpace(p)
		if p == currentOS {
			return true
		}

		// Unix平台
		if p == "unix" && (currentOS == "linux" || currentOS == "macos") {
			return true
		}
	}
	return false
}
