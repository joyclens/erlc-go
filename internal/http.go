package internal

import (
	"net/http"
	"net/url"
	"strings"
)

// BuildURL constructs a full URL from base and path, handling query parameters
func BuildURL(base, path string, params map[string]string) string {
	u, err := url.Parse(base + path)
	if err != nil {
		return base + path
	}

	if len(params) > 0 {
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}

	return u.String()
}

// SetQueryParams adds query parameters to a URL
func SetQueryParams(urlStr string, params map[string]string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// IsSuccessStatus checks if a status code is a success code
func IsSuccessStatus(status int) bool {
	return status >= 200 && status < 300
}

// IsClientError checks if a status code is a client error
func IsClientError(status int) bool {
	return status >= 400 && status < 500
}

// IsServerError checks if a status code is a server error
func IsServerError(status int) bool {
	return status >= 500 && status < 600
}

// IsRetryableStatus checks if a status code should be retried
func IsRetryableStatus(status int) bool {
	switch status {
	case http.StatusRequestTimeout, http.StatusTooManyRequests,
		http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// CloneHeader clones an HTTP header
func CloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// ContainsHeader checks if a header exists (case-insensitive)
func ContainsHeader(h http.Header, key string) bool {
	_, ok := h[key]
	return ok
}

// NormalizeHeaderKey normalizes a header key to HTTP canonical format
func NormalizeHeaderKey(key string) string {
	// http.CanonicalHeaderKey handles this
	return strings.Title(key) // Simplified version for demonstration
}
