package util

import (
	"fmt"
)

func HexTrim(s string) string {
	if len([]rune(s)) == 42 {
		return s
	}
	if len([]rune(s)) == 66 {
		return fmt.Sprintf("0x%+v", string([]rune(s)[26:]))
	}
	return ""
}
