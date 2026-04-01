package erlc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	config      *Config
	rateLimiter *RateLimiter
	cache       *Cache
	httpClient  *http.Client
}

func NewClient(serverKey string, opts ...Option) (*Client, error) {
	if serverKey == "" {
		return nil, NewValidationError("serverKey", "server key is required", nil)
	}

	config := DefaultConfig()
	config.ServerKey = serverKey

	for _, opt := range opts {
		opt(config)
	}

	var rateLimiter *RateLimiter
	if config.RateLimitEnabled {
		rateLimiter = NewRateLimiter(config.RateLimitPerSecond)
	}

	var cache *Cache
	if config.CacheEnabled {
		cache = NewCache(config.MaxCacheSize, config.CacheTTL)
	}

	return &Client{
		config:      config,
		rateLimiter: rateLimiter,
		cache:       cache,
		httpClient:  config.HTTPClient,
	}, nil
}

func (c *Client) Close() error {
	if c.cache != nil {
		c.cache.Clear()
	}
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return err
		}
	}

	url := c.config.BaseURL + path
	var bodyReader io.Reader
	var contentType string

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		contentType = ContentTypeJSON
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setRequestHeaders(req, contentType)

	if c.config.LoggingEnabled && c.config.Logger != nil {
		c.config.Logger.Debugf("%s %s", method, path)
	}

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = c.handleNetworkError(err)
			if attempt < c.config.MaxRetries {
				waitTime := c.getRetryWait(attempt)
				select {
				case <-time.After(waitTime):
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			if attempt < c.config.MaxRetries {
				waitTime := c.getRetryWait(attempt)
				select {
				case <-time.After(waitTime):
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			continue
		}

		if resp.StatusCode == StatusTooManyRequests {
			retryAfter, _ := ParseRateLimitHeaders(resp.Header)
			if retryAfter == 0 {
				retryAfter = time.Second * time.Duration((attempt+1)*2)
			}

			if c.config.LoggingEnabled && c.config.Logger != nil {
				c.config.Logger.Warnf("Rate limited, retrying after %v", retryAfter)
			}

			if attempt < c.config.MaxRetries {
				select {
				case <-time.After(retryAfter):
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			return NewRateLimitError(retryAfter)
		}

		if resp.StatusCode >= 400 {
			apiErr := ParseErrorResponse(resp.StatusCode, respBody)
			apiErr.URL = url
			apiErr.Method = method

			if c.config.LoggingEnabled && c.config.Logger != nil {
				c.config.Logger.Errorf("%s %s: %v", method, path, apiErr)
			}

			if apiErr.IsTransient() && attempt < c.config.MaxRetries {
				waitTime := c.getRetryWait(attempt)
				select {
				case <-time.After(waitTime):
					continue
				case <-ctx.Done():
					return ctx.Err()
				}
			}

			return apiErr
		}

		if c.config.LoggingEnabled && c.config.Logger != nil {
			c.config.Logger.Debugf("Response %d", resp.StatusCode)
		}

		if result != nil {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}
		}

		return nil
	}

	return lastErr
}

func (c *Client) doWithCache(ctx context.Context, path string, result interface{}, cacheOpt *CacheOptions) error {
	if c.cache != nil && (cacheOpt == nil || cacheOpt.Enabled) {
		if cached, ok := c.cache.Get(path); ok {
			if c.config.LoggingEnabled && c.config.Logger != nil {
				c.config.Logger.Debugf("Cache hit for %s", path)
			}
			return json.Unmarshal(cached, result)
		}
	}

	if err := c.do(ctx, http.MethodGet, path, nil, result); err != nil {
		return err
	}

	if c.cache != nil && (cacheOpt == nil || cacheOpt.Enabled) {
		if data, err := json.Marshal(result); err == nil {
			ttl := c.config.CacheTTL
			if cacheOpt != nil && cacheOpt.TTL > 0 {
				ttl = cacheOpt.TTL
			}
			c.cache.SetWithTTL(path, data, ttl)
		}
	}

	return nil
}

func (c *Client) setRequestHeaders(req *http.Request, contentType string) {
	req.Header.Set(ServerKeyHeader, c.config.ServerKey)
	req.Header.Set(UserAgentHeader, UserAgent)

	if contentType != "" {
		req.Header.Set(ContentTypeHeader, contentType)
	}

	for k, v := range c.config.CustomHeaders {
		req.Header.Set(k, v)
	}
}

func (c *Client) handleNetworkError(err error) error {
	if err == context.Canceled {
		return fmt.Errorf("request canceled: %w", err)
	}

	if err == context.DeadlineExceeded {
		return NewNetworkError(err, true, false)
	}

	type timeoutError interface {
		Timeout() bool
	}
	if te, ok := err.(timeoutError); ok && te.Timeout() {
		return NewNetworkError(err, true, false)
	}

	type temporaryError interface {
		Temporary() bool
	}
	if te, ok := err.(temporaryError); ok && te.Temporary() {
		return NewNetworkError(err, false, true)
	}

	return NewNetworkError(err, false, false)
}

func (c *Client) getRetryWait(attempt int) time.Duration {
	base := c.config.RetryInterval
	multiplier := 1 << uint(attempt)
	return base * time.Duration(multiplier)
}

func (c *Client) ClearCache() {
	if c.cache != nil {
		c.cache.Clear()
	}
}

func (c *Client) CacheStats() *CacheStats {
	if c.cache != nil {
		return c.cache.Stats()
	}
	return nil
}
