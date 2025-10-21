// Package utils 提供了 zrun 的工具函数。
package utils

import (
	"runtime"
	"strings"
)

var currentOS string

func init() {
	// 初始化操作系统类型
	switch runtime.GOOS {
	case "windows":
		currentOS = "windows"
	case "linux":
		currentOS = "linux"
	case "darwin":
		currentOS = "macos"
	default:
		currentOS = runtime.GOOS
	}
}

// GetOS 获取当前操作系统
// 返回值"windows", "linux", "macos"
func GetOS() string {
	return currentOS
}

// 检查平台是否匹配
// 参数:
//   - platform: 脚本中声明的平台标识符
//   - currentOS: 当前操作系统标识
//
// 返回值:
//   - bool: 平台是否匹配
func MatchPlatform(platform, currentOS string) bool {
	// default匹配所有平台
	if platform == "default" {
		return true
	}



	if !strings.Contains(platform, ",") {
		if platform == currentOS {
			return true
		}
		


		// Unix特殊处理
		if platform == "unix" && (currentOS == "linux" || currentOS == "macos") {
			return true
		}
		
		return false
	}


	


	platforms := strings.Split(platform, ",")
	for _, p := range platforms {
		// 去除空白字符
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