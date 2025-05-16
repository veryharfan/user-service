package pkg

import (
	"regexp"
	"unicode"
)

// Simple email regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Basic phone check: only digits (allow + at start optionally)
func IsPhone(input string) bool {
	if input == "" {
		return false
	}

	// Allow "+1234567890" or "1234567890"
	for i, r := range input {
		if i == 0 && r == '+' {
			continue
		}
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func GetUsernameType(input string) string {
	switch {
	case emailRegex.MatchString(input):
		return "email"
	case IsPhone(input):
		return "phone"
	default:
		return "unknown"
	}
}
