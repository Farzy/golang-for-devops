// Package token_bucket implements a simple rate limiter using the Token Bucket algorithm.
// The Token Bucket algorithm is commonly used for traffic shaping and rate limiting.
// This implementation is adapted from https://himanshu-007.medium.com/simple-rate-limiter-in-golang-using-token-bucket-algorithm-388d0596d1e4.
package token_bucket

import (
	"math"
	"sync"
	"time"
)

// TokenBucket represents a token bucket, a counter for a type of resource,
// filled at a steady rate. The rate is specified in tokens per second.
// The bucket has a capacity, and if a token arrives when the bucket is full,
// it is discarded.
type TokenBucket struct {
	rate                int64      // Tokens per second
	maxTokens           int64      // Maximum tokens the bucket can hold
	currentTokens       int64      // Current tokens in the bucket
	lastRefillTimestamp time.Time  // The last time the bucket was refilled
	mutex               sync.Mutex // Protects the TokenBucket
}

// NewTokenBucket creates a new TokenBucket with the specified rate and capacity.
func NewTokenBucket(Rate int64, MaxTokens int64) *TokenBucket {
	return &TokenBucket{
		rate:                Rate,
		maxTokens:           MaxTokens,
		currentTokens:       MaxTokens,
		lastRefillTimestamp: time.Now(),
	}
}

// refill refills the tokens in the bucket based on the time elapsed since
// the last refill. If more tokens would be added than the bucket can hold,
// the extra tokens are discarded.
func (tb *TokenBucket) refill() {
	now := time.Now()
	end := time.Since(tb.lastRefillTimestamp)
	tokensToBeAdded := (end.Nanoseconds() * tb.rate) / 1000000000
	tb.currentTokens = int64(math.Min(float64(tb.currentTokens+tokensToBeAdded), float64(tb.maxTokens)))
	tb.lastRefillTimestamp = now
}

// IsRequestAllowed checks to see if the specified number of tokens are
// available in the bucket. If the tokens are available, it decreases
// the token count and returns true. If the tokens are not available,
// it returns false. The function ensures thread safety using a mutex lock.
func (tb *TokenBucket) IsRequestAllowed(tokens int64) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()
	tb.refill()
	if tb.currentTokens >= tokens {
		tb.currentTokens -= tokens
		return true
	}
	return false
}
