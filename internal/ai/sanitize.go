package ai

import "strings"

// Sanitize remove unwanted characters (example: control characters)
func Sanitize(input string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 { // Remove control characters
			return -1
		}
		return r
	}, input)
}
