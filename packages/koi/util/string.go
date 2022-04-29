package util

import (
	"strings"
)

var (
	trimSet      = string([]rune{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 28, 29, 30, 31, 32, 12288})
	ResetCtrlStr = string([]rune{27, 91, 48, 109})
)

func Trim(s string) string {
	s = strings.Trim(s, trimSet)
	lenS := len(s)
	if lenS >= 4 && s[lenS-4:] == ResetCtrlStr {
		return Trim(s[:lenS-4])
	}
	return s
}
