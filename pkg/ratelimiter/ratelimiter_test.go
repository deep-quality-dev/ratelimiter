package ratelimiter

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_Refill(t *testing.T) {
	testcases := []struct {
		name           string
		maxTokens      int
		refillAmount   int
		refillInterval time.Duration

		sleep            time.Duration
		remainingTokens  int
		sleep2           time.Duration
		remainingTokens2 int
	}{
		{"1 cycle, same refill amount with max amount", 3, 3, time.Second, time.Second, 2, time.Second, 2},
		{"1 cycle, refill amount less than max amount", 10, 3, time.Second, time.Second, 2, time.Second, 4},
		{"n cycles, same refill amount with max amount", 3, 3, time.Second, time.Second * 3, 2, time.Second * 3, 2},
		{"n cycles, refill amount less than max amount", 20, 3, time.Second, time.Second * 3, 8, time.Second * 3, 16},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tb, _ := NewTokenBucket(tt.maxTokens, tt.refillAmount, tt.refillInterval)

			// Should consume tokens within the limit
			for i := 0; i < tb.maxTokens; i++ {
				assert.True(t, tb.ConsumeTokens(1))
			}

			// Should not consume tokens beyond the limit
			assert.False(t, tb.ConsumeTokens(1))

			// Sleep for refill interval to refill tokens
			time.Sleep(tt.sleep)

			// Should continue consuming tokens after refill
			assert.True(t, tb.ConsumeTokens(1))

			assert.Equal(t, tt.remainingTokens, tb.tokens)

			// Sleep for refill interval to refill tokens
			time.Sleep(tt.sleep2)

			// Should continue consuming tokens after refill
			assert.True(t, tb.ConsumeTokens(1))

			assert.Equal(t, tt.remainingTokens2, tb.tokens)
		})
	}
}

func TestTokenBucket_BeforeRefill(t *testing.T) {
	tb, _ := NewTokenBucket(3, 3, time.Second*3)

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
	tb, _ := NewTokenBucket(5, 0, time.Second)

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
	tb, _ := NewTokenBucket(100, 100, time.Second*3)

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
