package internal

import (
	"strings"
	"time"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func StringInSlice(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

func StringInSliceIgnoreCase(s string, list []string) bool {
	lower := strings.ToLower(s)
	for _, item := range list {
		if strings.ToLower(item) == lower {
			return true
		}
	}
	return false
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func FormatDuration(d time.Duration) string {
	switch {
	case d < time.Second:
		return d.String()
	case d < time.Minute:
		return d.Round(time.Millisecond).String()
	case d < time.Hour:
		return d.Round(time.Second).String()
	default:
		return d.Round(time.Minute).String()
	}
}

func ParseBool(s string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "yes", "1", "on":
		return true, nil
	case "false", "no", "0", "off":
		return false, nil
	default:
		return false, ErrInvalidBool
	}
}

var ErrInvalidBool = &parseError{"invalid boolean value"}

type parseError struct {
	msg string
}

func (e *parseError) Error() string {
	return e.msg
}

func FilterString(list []string, predicate func(string) bool) []string {
	var result []string
	for _, item := range list {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func MapString(list []string, fn func(string) string) []string {
	result := make([]string, len(list))
	for i, item := range list {
		result[i] = fn(item)
	}
	return result
}

func Unique(list []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range list {
		if !seen[item] {
			result = append(result, item)
			seen[item] = true
		}
	}
	return result
}

func PointerTo[T any](v T) *T {
	return &v
}

func ValueOf[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func Retry(maxAttempts int, delay time.Duration, fn func() error) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < maxAttempts {
			time.Sleep(delay * time.Duration(1<<uint(attempt-1)))
		}
	}
	return lastErr
}
