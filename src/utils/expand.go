package utils

import (
	"strings"
)

// ExpandVars 展开 $VAR / ${VAR} / $$ 转义。
// lookup 负责提供变量值，不存在时保留原文便于排查。
func ExpandVars(s string, lookup func(name string) (string, bool)) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		c := s[i]
		if c != '$' {
			b.WriteByte(c)
			i++
			continue
		}
		// $$ -> $
		if i+1 < len(s) && s[i+1] == '$' {
			b.WriteByte('$')
			i += 2
			continue
		}
		// ${NAME}
		if i+1 < len(s) && s[i+1] == '{' {
			end := strings.IndexByte(s[i+2:], '}')
			if end < 0 {
				b.WriteByte(c)
				i++
				continue
			}
			name := s[i+2 : i+2+end]
			if v, ok := lookup(name); ok {
				b.WriteString(v)
			} else {
				b.WriteString(s[i : i+2+end+1])
			}
			i += 2 + end + 1
			continue
		}
		// $NAME：字母数字下划线
		j := i + 1
		for j < len(s) && isVarChar(s[j]) {
			j++
		}
		if j == i+1 {
			b.WriteByte(c)
			i++
			continue
		}
		name := s[i+1 : j]
		if v, ok := lookup(name); ok {
			b.WriteString(v)
		} else {
			b.WriteString(s[i:j])
		}
		i = j
	}
	return b.String()
}

func isVarChar(c byte) bool {
	return c == '_' ||
		(c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9')
}
