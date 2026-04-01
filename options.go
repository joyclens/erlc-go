package erlc

import (
	"net"
	"net/http"
	"time"
)

type Option func(*Config)

type Config struct {
	ServerKey string

	HTTPClient *http.Client
	Timeout    time.Duration

	MaxRetries    int
	RetryInterval time.Duration

	RateLimitPerSecond int
	RateLimitEnabled   bool

	CacheEnabled bool
	CacheTTL     time.Duration
	MaxCacheSize int

	LoggingEnabled bool
	Logger         Logger

	CustomHeaders map[string]string

	BaseURL string
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func DefaultConfig() *Config {
	return &Config{
		Timeout:            DefaultTimeout,
		MaxRetries:         MaxRetries,
		RetryInterval:      time.Second,
		RateLimitPerSecond: DefaultRateLimitPerSecond,
		RateLimitEnabled:   true,
		CacheEnabled:       false,
		CacheTTL:           DefaultCacheTTL,
		MaxCacheSize:       DefaultMaxCacheSize,
		LoggingEnabled:     false,
		CustomHeaders:      make(map[string]string),
		BaseURL:            V1BaseURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
			},
		},
	}
}

func WithServerKey(key string) Option {
	return func(c *Config) {
		c.ServerKey = key
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
		if c.HTTPClient != nil {
			c.HTTPClient.Timeout = timeout
		}
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(c *Config) {
		c.MaxRetries = maxRetries
	}
}

func WithRetryInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.RetryInterval = interval
	}
}

func WithRateLimiting(enabled bool, perSecond int) Option {
	return func(c *Config) {
		c.RateLimitEnabled = enabled
		c.RateLimitPerSecond = perSecond
	}
}

func WithCaching(enabled bool, ttl time.Duration, maxSize int) Option {
	return func(c *Config) {
		c.CacheEnabled = enabled
		c.CacheTTL = ttl
		c.MaxCacheSize = maxSize
	}
}

func WithCacheTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.CacheTTL = ttl
	}
}

func WithMaxCacheSize(maxSize int) Option {
	return func(c *Config) {
		c.MaxCacheSize = maxSize
	}
}

func WithLogging(enabled bool, logger Logger) Option {
	return func(c *Config) {
		c.LoggingEnabled = enabled
		c.Logger = logger
	}
}

func WithCustomHeader(key, value string) Option {
	return func(c *Config) {
		c.CustomHeaders[key] = value
	}
}

func WithCustomHeaders(headers map[string]string) Option {
	return func(c *Config) {
		for k, v := range headers {
			c.CustomHeaders[k] = v
		}
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

func WithV2API() Option {
	return func(c *Config) {
		c.BaseURL = V2BaseURL
	}
}
