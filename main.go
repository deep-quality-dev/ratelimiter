package main

import (
	"fmt"
	"time"
)

func main() {
	// Example usage for a user sending one message per second
	userRateLimiter := NewTokenBucket(1, 1, time.Second)

	// Example usage for a user with three allowed failed credit card transactions per day
	userCreditCardLimiter := NewTokenBucket(3, 3, 24*time.Hour)

	// Example usage for a single IP creating twenty accounts per day
	ipAccountLimiter := NewTokenBucket(20, 20, 24*time.Hour)

	for i := 0; i < 21; i++ {
		// User sending one message per second
		if userRateLimiter.ConsumeTokens(1) {
			fmt.Println("[True] User allowed to send a message")
		} else {
			fmt.Println("[False] User rate limit exceeded")
		}

		// User with three allowed failed credit card transactions per day
		if userCreditCardLimiter.ConsumeTokens(1) {
			fmt.Println("[True] User allowed for credit card transaction")
		} else {
			fmt.Println("[False] User exceeded credit card transactions limit")
		}

		// Single IP creating twenty accounts per day
		if ipAccountLimiter.ConsumeTokens(1) {
			fmt.Println("[True] IP allowed to create an account")
		} else {
			fmt.Println("[False] IP exceeded account creation limit")
		}

		time.Sleep(time.Second)
	}
}
