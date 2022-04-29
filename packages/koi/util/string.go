package util

import "strings"

func Trim(s string) string {
	return strings.Trim(s, " 　\f\n\r\t\v\a\b")
}
