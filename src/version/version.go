package version

import (
	"fmt"
	"strconv"
	"strings"
)

// 版本格式：YY.大版本.小版本，例如 26.3.1 表示 2026 年的大版本 3、小版本 1。
const Version = "26.1.0"

// SyntaxVersion 为脚本语法版本，与应用版本独立演进。
const SyntaxVersion = "2.0"

// VersionInfo 解析后的版本号。
type VersionInfo struct {
	Year  int
	Major int
	Minor int
}

// Parse 解析 YY.大版本.小版本，允许带 v/V 前缀，两侧空格会被忽略。
func Parse(s string) (VersionInfo, error) {
	t := strings.TrimSpace(s)
	t = strings.TrimPrefix(t, "v")
	t = strings.TrimPrefix(t, "V")
	parts := strings.Split(t, ".")
	if len(parts) != 3 {
		return VersionInfo{}, fmt.Errorf("版本号 %q 非法，应为 YY.大版本.小版本，例如 26.3.1", s)
	}
	nums := make([]int, 3)
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return VersionInfo{}, fmt.Errorf("版本号 %q 非法，第 %d 段为空", s, i+1)
		}
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return VersionInfo{}, fmt.Errorf("版本号 %q 非法，第 %d 段 %q 不是非负整数", s, i+1, p)
		}
		nums[i] = n
	}
	if nums[0] > 99 {
		return VersionInfo{}, fmt.Errorf("版本号 %q 非法，年份段应为两位数 00-99", s)
	}
	return VersionInfo{Year: nums[0], Major: nums[1], Minor: nums[2]}, nil
}

// String 还原为 26.3.1 形式。
func (v VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Year, v.Major, v.Minor)
}

// FullYear 展开为四位数年份，例如 26 -> 2026。
func (v VersionInfo) FullYear() int {
	return 2000 + v.Year
}

// Compare 按年、大版本、小版本逐段数字比较。
// 返回 -1/0/1，分别表示 a<b、a==b、a>b。
func Compare(a, b VersionInfo) int {
	if a.Year != b.Year {
		if a.Year < b.Year {
			return -1
		}
		return 1
	}
	if a.Major != b.Major {
		if a.Major < b.Major {
			return -1
		}
		return 1
	}
	if a.Minor != b.Minor {
		if a.Minor < b.Minor {
			return -1
		}
		return 1
	}
	return 0
}

// CompareStrings 解析后比较，任一解析失败时回退为字符串比较，保证不误判为相等。
func CompareStrings(a, b string) int {
	va, errA := Parse(a)
	vb, errB := Parse(b)
	if errA != nil || errB != nil {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	}
	return Compare(va, vb)
}
