package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_Refill(t *testing.T) {
	tb := NewTokenBucket(3, 3, time.Second)

	// Should consume tokens within the limit
	for i := 0; i < tb.maxTokens; i++ {
		assert.True(t, tb.ConsumeTokens(1))
	}

	// Should not consume tokens beyond the limit
	assert.False(t, tb.ConsumeTokens(1))

	// Sleep for refill interval to refill tokens
	time.Sleep(tb.refillInterval)

	// Should continue consuming tokens after refill
	assert.True(t, tb.ConsumeTokens(1))
}

func TestTokenBucket_BeforeRefill(t *testing.T) {
	tb := NewTokenBucket(3, 3, time.Second*3)

	// Should consume tokens within the limit
	for i := 0; i < tb.maxTokens; i++ {
		assert.True(t, tb.ConsumeTokens(1))
	}

	// Should not consume tokens beyond the limit
	assert.False(t, tb.ConsumeTokens(1))

	// Sleep half of refill interval not to refill tokens
	time.Sleep(tb.refillInterval / 2)

	// Should not consume tokens then
	assert.False(t, tb.ConsumeTokens(1))
}

func TestTokenBucket_WithZeroRefill(t *testing.T) {
	tb := NewTokenBucket(5, 0, time.Second)

	// Should consume tokens should always return true
	for i := 0; i < 5; i++ {
		assert.True(t, tb.ConsumeTokens(1))
	}

	// Sleep for refill interval to replenish tokens
	time.Sleep(tb.refillInterval)

	// Should not continue consuming tokens after refill
	assert.False(t, tb.ConsumeTokens(1))
}

func TestTokenBucket_Concurrent(t *testing.T) {
	tb := NewTokenBucket(100, 100, time.Second*3)

	// Consume tokens concurrently
	var wg sync.WaitGroup
	concurrentCalls := tb.maxTokens
	wg.Add(concurrentCalls)
	for i := 0; i < concurrentCalls; i++ {
		go func() {
			defer wg.Done()
			assert.True(t, tb.ConsumeTokens(1))
		}()
	}
	wg.Wait()
}
