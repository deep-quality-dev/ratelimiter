package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	// the current number of tokens in the bucket
	tokens int
	// the maximum number of tokens the bucket can hold
	maxTokens int
	// the number of tokens that are added to the bucket during a refill
	refillAmount int
	// the amount of time between refills
	refillInterval time.Duration
	// the last time the bucket was updated
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewTokenBucket creates a new TokenBucket with the specific parameters
func NewTokenBucket(maxTokens, refillAmount int, refillInterval time.Duration) (*TokenBucket, error) {
	if maxTokens < 0 {
		return nil, fmt.Errorf("max tokens not valid")
	}
	if refillAmount < 0 {
		return nil, fmt.Errorf("refill amount not valid")
	}
	return &TokenBucket{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillAmount:   refillAmount,
		refillInterval: refillInterval,
		lastUpdate:     time.Now(),
	}, nil
}

// ConsumeTokens tries to consume the specific amount of tokens and
// return true if consumed successfully
func (tb *TokenBucket) ConsumeTokens(tokens int) bool {
	if tokens <= 0 {
		return false
	}
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens
	now := time.Now()
	refillTokens := int(now.Sub(tb.lastUpdate)/tb.refillInterval) * tb.refillAmount
	if refillTokens > 0 {
		tb.tokens += refillTokens
		tb.lastUpdate = now
	}

	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}

	// Check if enough tokens are available to consume
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}
