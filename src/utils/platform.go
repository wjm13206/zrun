// Package utils 提供了 zrun 的工具函数。
package utils

import (
	"runtime"
	"strings"
)

var currentOS string
var currentArch string

func init() {
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

	switch runtime.GOARCH {
	case "amd64":
		currentArch = "64"
	case "386", "arm":
		currentArch = "32"
	case "arm64":
		currentArch = "arm64"
	default:
		currentArch = runtime.GOARCH
	}
}

func GetOS() string {
	return currentOS
}

func GetArch() string {
	return currentArch
}

func MatchPlatform(platform string) bool {
	if platform == "default" {
		return true
	}

	platforms := strings.Split(platform, ",")
	for _, p := range platforms {
		p = strings.TrimSpace(p)

		if strings.Contains(p, "/") {
			parts := strings.Split(p, "/")
			if len(parts) != 2 {
				continue
			}
			osPart := parts[0]
			archPart := parts[1]

			osMatch := osPart == currentOS || (osPart == "unix" && (currentOS == "linux" || currentOS == "macos"))
			archMatch := archPart == currentArch

			if osMatch && archMatch {
				return true
			}
			continue
		}

		if p == currentOS {
			return true
		}
		if p == "unix" && (currentOS == "linux" || currentOS == "macos") {
			return true
		}
		if p == currentArch {
			return true
		}
	}

	return false
}