package erlc

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	rate       int
	tokens     int
	maxTokens  int
	lastUpdate time.Time
	mu         sync.Mutex
}

func NewRateLimiter(ratePerSecond int) *RateLimiter {
	if ratePerSecond < 1 {
		ratePerSecond = DefaultRateLimitPerSecond
	}
	return &RateLimiter{
		rate:       ratePerSecond,
		tokens:     ratePerSecond,
		maxTokens:  ratePerSecond,
		lastUpdate: time.Now(),
	}
}

func (r *RateLimiter) Wait(ctx context.Context) error {
	for {
		r.mu.Lock()
		r.replenish()

		if r.tokens > 0 {
			r.tokens--
			r.mu.Unlock()
			return nil
		}

		r.mu.Unlock()

		select {
		case <-time.After(time.Second / time.Duration(r.rate)):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.replenish()
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

func (r *RateLimiter) SetRate(ratePerSecond int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rate = ratePerSecond
	r.maxTokens = ratePerSecond
}

func (r *RateLimiter) GetTokens() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.replenish()
	return r.tokens
}

func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens = r.maxTokens
	r.lastUpdate = time.Now()
}

func (r *RateLimiter) replenish() {
	now := time.Now()
	elapsed := now.Sub(r.lastUpdate)
	tokensToAdd := int(elapsed.Seconds() * float64(r.rate))

	if tokensToAdd > 0 {
		r.tokens = min(r.tokens+tokensToAdd, r.maxTokens)
		r.lastUpdate = now
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type RequestQueue struct {
	maxRetries int
	baseDelay  time.Duration
	mu         sync.Mutex
	queue      chan func()
	stop       chan struct{}
	wg         sync.WaitGroup
}

func NewRequestQueue(maxRetries int, baseDelay time.Duration) *RequestQueue {
	if maxRetries < 1 {
		maxRetries = 3
	}
	if baseDelay < 1 {
		baseDelay = time.Second
	}

	rq := &RequestQueue{
		maxRetries: maxRetries,
		baseDelay:  baseDelay,
		queue:      make(chan func(), 100),
		stop:       make(chan struct{}),
	}

	rq.wg.Add(1)
	go rq.process()

	return rq
}

func (rq *RequestQueue) Enqueue(fn func()) error {
	select {
	case rq.queue <- fn:
		return nil
	case <-rq.stop:
		return fmt.Errorf("request queue is stopped")
	}
}

func (rq *RequestQueue) process() {
	defer rq.wg.Done()

	for {
		select {
		case fn := <-rq.queue:
			fn()
		case <-rq.stop:
			return
		}
	}
}

func (rq *RequestQueue) Stop() {
	close(rq.stop)
	rq.wg.Wait()
}

func (rq *RequestQueue) Len() int {
	return len(rq.queue)
}
