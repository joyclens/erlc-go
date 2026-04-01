package erlc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		serverKey string
		wantErr   bool
	}{
		{
			name:      "valid server key",
			serverKey: "test-key-12345",
			wantErr:   false,
		},
		{
			name:      "empty server key",
			serverKey: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.serverKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClientWithOptions(t *testing.T) {
	client, err := NewClient(
		"test-key",
		WithTimeout(15*time.Second),
		WithMaxRetries(5),
		WithRateLimiting(true, 50),
	)

	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	if client == nil {
		t.Error("NewClient() returned nil client")
	}

	if client.config.Timeout != 15*time.Second {
		t.Errorf("Expected timeout 15s, got %v", client.config.Timeout)
	}

	if client.config.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries 5, got %d", client.config.MaxRetries)
	}
}

func TestCachingEnabled(t *testing.T) {
	client, err := NewClient(
		"test-key",
		WithCaching(true, 5*time.Minute, 1024*1024),
	)

	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	if client.cache == nil {
		t.Error("Expected cache to be enabled")
	}

	if client.config.CacheTTL != 5*time.Minute {
		t.Errorf("Expected cache TTL 5m, got %v", client.config.CacheTTL)
	}
}

func TestAPIErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code":401,"message":"invalid server key"}`))
	}))
	defer server.Close()

	client, _ := NewClient("invalid-key", WithBaseURL(server.URL))
	_, err := client.GetServer(context.Background())

	if err == nil {
		t.Error("Expected error for unauthorized request")
	}

	if apiErr, ok := err.(*APIError); ok {
		if !apiErr.IsAuthError() {
			t.Errorf("Expected auth error, got %v", apiErr)
		}
	} else {
		t.Errorf("Expected APIError, got %T", err)
	}
}

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(2)

	if !rl.Allow() {
		t.Error("First request should be allowed")
	}
	if !rl.Allow() {
		t.Error("Second request should be allowed")
	}

	if rl.Allow() {
		t.Error("Third request should be rate limited")
	}

	time.Sleep(600 * time.Millisecond)
	if !rl.Allow() {
		t.Error("Request should be allowed after token replenishment")
	}
}

func TestRateLimiterWait(t *testing.T) {
	rl := NewRateLimiter(10)
	ctx := context.Background()

	start := time.Now()
	if err := rl.Wait(ctx); err != nil {
		t.Errorf("Wait failed: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Unexpected delay: %v", elapsed)
	}
}

func TestCache(t *testing.T) {
	cache := NewCache(1024, time.Second)

	cache.Set("key1", []byte("value1"))
	data, ok := cache.Get("key1")
	if !ok {
		t.Error("Expected cache hit for key1")
	}
	if string(data) != "value1" {
		t.Errorf("Expected value1, got %s", string(data))
	}

	_, ok = cache.Get("nonexistent")
	if ok {
		t.Error("Expected cache miss for nonexistent key")
	}

	cache.SetWithTTL("key2", []byte("value2"), 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, ok = cache.Get("key2")
	if ok {
		t.Error("Expected cache miss after expiration")
	}
}

func TestCacheStats(t *testing.T) {
	cache := NewCache(1024, time.Minute)

	cache.Set("key1", []byte("value1"))
	cache.Get("key1")
	cache.Get("key2")

	stats := cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
	if stats.HitRate < 0.49 || stats.HitRate > 0.51 {
		t.Errorf("Expected hit rate ~0.5, got %f", stats.HitRate)
	}
}

func TestCacheEviction(t *testing.T) {
	cache := NewCache(100, time.Minute)

	cache.Set("key1", make([]byte, 50))
	cache.Set("key2", make([]byte, 40))
	cache.Set("key3", make([]byte, 50))

	stats := cache.Stats()
	if stats.Entries < 1 {
		t.Error("Expected at least 1 entry after eviction")
	}
}

func TestErrorTypes(t *testing.T) {
	apiErr := &APIError{
		Code:       429,
		Message:    "Rate limited",
		StatusCode: 429,
	}

	if !apiErr.IsRateLimit() {
		t.Error("Expected IsRateLimit() to return true")
	}

	valErr := NewValidationError("field", "required", nil)
	if valErr.Field != "field" {
		t.Errorf("Expected field 'field', got %s", valErr.Field)
	}

	err := NewNetworkError(nil, true, false)
	if !err.IsTimeout {
		t.Error("Expected IsTimeout to be true")
	}
}

func TestParseErrorResponse(t *testing.T) {
	body := []byte(`{"code":422,"message":"validation failed","details":"missing required field"}`)
	err := ParseErrorResponse(422, body)

	if err.Code != 422 {
		t.Errorf("Expected code 422, got %d", err.Code)
	}

	if err.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got %s", err.Message)
	}

	if err.Details != "missing required field" {
		t.Errorf("Expected details, got %s", err.Details)
	}
}

func TestValidation(t *testing.T) {
	client, _ := NewClient("test-key")

	_, err := client.ExecuteCommand(context.Background(), "", "args")
	if err == nil {
		t.Error("Expected validation error for empty command")
	}

	if valErr, ok := err.(*ValidationError); ok {
		if valErr.Field != "command" {
			t.Errorf("Expected field 'command', got %s", valErr.Field)
		}
	} else {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	rl := NewRateLimiter(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow()
	}
}

func BenchmarkCache(b *testing.B) {
	cache := NewCache(1024*1024, time.Minute)
	data := []byte("test data")

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Set("key", data)
		}
	})

	b.Run("Get", func(b *testing.B) {
		cache.Set("key", data)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get("key")
		}
	})
}
