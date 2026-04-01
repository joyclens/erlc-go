package internal

import (
	"fmt"
	"net/url"
	"regexp"
	"unicode"
)

func ValidateServerKey(key string) error {
	if key == "" {
		return fmt.Errorf("server key cannot be empty")
	}

	if len(key) < 8 {
		return fmt.Errorf("server key too short (minimum 8 characters)")
	}

	if len(key) > 256 {
		return fmt.Errorf("server key too long (maximum 256 characters)")
	}

	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if len(username) < 3 {
		return fmt.Errorf("username too short (minimum 3 characters)")
	}

	if len(username) > 20 {
		return fmt.Errorf("username too long (maximum 20 characters)")
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return fmt.Errorf("username contains invalid characters")
	}

	return nil
}

func ValidateRobloxID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("roblox ID must be positive")
	}

	if id > 9223372036854775807 {
		return fmt.Errorf("roblox ID exceeds maximum value")
	}

	return nil
}

func ValidateURL(urlStr string) error {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	return nil
}

func IsValidCommand(command string) error {
	if command == "" {
		return fmt.Errorf("command cannot be empty")
	}

	if len(command) > 100 {
		return fmt.Errorf("command too long (maximum 100 characters)")
	}

	if len(command) > 0 && !unicode.IsLetter(rune(command[0])) {
		return fmt.Errorf("command must start with a letter")
	}

	return nil
}

func IsSafeString(s string) bool {
	for _, r := range s {
		if r < 32 || r == 127 {
			return false
		}
	}
	return true
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

func SanitizeInput(s string) string {
	return regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(s, "")
}
