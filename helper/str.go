package helper

import (
	"strings"
)

// UnderscoreCamelCase 下划线转驼峰
func UnderscoreCamelCase(str string) string {
	b := make([]byte, 0)
	nextToUp := false
	for _, v := range str {
		vs := string(v)
		if vs == "_" {
			nextToUp = true
			continue
		}
		if nextToUp {
			vs = strings.ToUpper(vs)
		}
		b = append(b, []byte(vs)...)
		nextToUp = false
	}
	return string(b)
}
