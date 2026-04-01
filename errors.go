package erlc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type APIError struct {
	Code       int           `json:"code,omitempty"`
	Message    string        `json:"message,omitempty"`
	Details    string        `json:"details,omitempty"`
	StatusCode int           `json:"-"`
	RetryAfter time.Duration `json:"-"`
	URL        string        `json:"-"`
	Method     string        `json:"-"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API error %d: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API error %d: %s", e.Code, e.Message)
}

func (e *APIError) IsRateLimit() bool {
	return e.StatusCode == StatusTooManyRequests
}

func (e *APIError) IsAuthError() bool {
	return e.StatusCode == StatusUnauthorized || e.StatusCode == StatusForbidden
}

func (e *APIError) IsNotFound() bool {
	return e.StatusCode == StatusNotFound
}

func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == StatusBadRequest || e.StatusCode == StatusUnprocessableEntity
}

func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}

func (e *APIError) IsTransient() bool {
	switch e.StatusCode {
	case StatusTooManyRequests, StatusInternalServerError, StatusServiceUnavailable:
		return true
	default:
		return false
	}
}

type NetworkError struct {
	Err           error
	IsTimeout     bool
	IsUnavailable bool
	IsTemporary   bool
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Field   string
	Message string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s (got: %v)", e.Field, e.Message, e.Value)
}

type RateLimitError struct {
	APIError
	RetryAtTime time.Time
	Limit       int
	Remaining   int
	ResetAt     time.Time
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited: %s (retry after %v)", e.APIError.Error(), e.RetryAfter)
}

func ParseErrorResponse(statusCode int, body []byte) *APIError {
	err := &APIError{
		StatusCode: statusCode,
		Code:       statusCode,
	}

	if len(body) > 0 {
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err != nil {
			err.Message = string(body)
		} else {
			if msg, ok := errResp["message"].(string); ok {
				err.Message = msg
			}
			if details, ok := errResp["details"].(string); ok {
				err.Details = details
			}
			if code, ok := errResp["code"].(float64); ok {
				err.Code = int(code)
			}
		}
	}

	if err.Message == "" {
		err.Message = getHTTPStatusMessage(statusCode)
	}

	return err
}

func ParseRateLimitHeaders(headers http.Header) (retryAfter time.Duration, resetTime time.Time) {
	if retryAfterStr := headers.Get("Retry-After"); retryAfterStr != "" {
		if seconds, err := strconv.Atoi(retryAfterStr); err == nil {
			retryAfter = time.Duration(seconds) * time.Second
		} else {
			if t, err := time.Parse(time.RFC1123, retryAfterStr); err == nil {
				retryAfter = time.Until(t)
				resetTime = t
			}
		}
	}

	if resetStr := headers.Get("X-RateLimit-Reset"); resetStr != "" {
		if unix, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			resetTime = time.Unix(unix, 0)
		}
	}

	return
}

func getHTTPStatusMessage(statusCode int) string {
	switch statusCode {
	case StatusOK:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusBadRequest:
		return "Bad request: invalid parameters"
	case StatusUnauthorized:
		return "Unauthorized: invalid or missing server key"
	case StatusForbidden:
		return "Forbidden: insufficient permissions"
	case StatusNotFound:
		return "Not found: resource does not exist"
	case StatusUnprocessableEntity:
		return "Unprocessable entity: validation failed"
	case StatusTooManyRequests:
		return "Too many requests: rate limited"
	case StatusInternalServerError:
		return "Internal server error"
	case StatusServiceUnavailable:
		return "Service unavailable"
	default:
		return http.StatusText(statusCode)
	}
}

func NewValidationError(field, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

func NewNetworkError(err error, isTimeout, isUnavailable bool) *NetworkError {
	return &NetworkError{
		Err:           err,
		IsTimeout:     isTimeout,
		IsUnavailable: isUnavailable,
		IsTemporary:   isTimeout || isUnavailable,
	}
}

func NewRateLimitError(retryAfter time.Duration) *RateLimitError {
	return &RateLimitError{
		APIError: APIError{
			StatusCode: StatusTooManyRequests,
			Code:       StatusTooManyRequests,
			Message:    "Rate limited",
			RetryAfter: retryAfter,
		},
		RetryAfter: retryAfter,
		RetryAtTime: time.Now().Add(retryAfter),
	}
}
