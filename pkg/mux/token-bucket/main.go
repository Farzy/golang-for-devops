// Package token_bucket implements a simple rate limiter using the Token Bucket algorithm.
// The Token Bucket algorithm is commonly used for traffic shaping and rate limiting.
// This implementation is adapted from https://himanshu-007.medium.com/simple-rate-limiter-in-golang-using-token-bucket-algorithm-388d0596d1e4.
//
// Here's a way to convert requests/sec into MaxTokens and Rate, provided by Google Bard:
// To convert "requests per second" into MaxToken and Rate in the Token Bucket algorithm, you can use the following formulas:
//
// ```
// MaxToken = requests per second * burst interval
// Rate = requests per second / burst interval
// ```
//
// Where:
//
// * **burst interval** is the number of seconds that the system allows to burst at the maximum rate.
//
// For example, if you want to allow a maximum of 100 requests per second with a burst interval of 5 seconds, then your MaxToken would be 500 and your Rate would be 20.
//
// ```
// MaxToken = 100 requests per second * 5 seconds = 500
// Rate = 100 requests per second / 5 seconds = 20
// ```
//
// This means that the system would allow a maximum of 500 requests in a 5-second period, and then it would start limiting requests to 20 per second.
//
// Note that the burst interval is a trade-off between performance and fairness. A longer burst interval will allow the system to handle sudden spikes in traffic better, but it will also allow users to send more requests in a short period of time. A shorter burst interval will be fairer to all users, but it may also impact performance during sudden spikes in traffic.
//
// You can adjust the MaxToken and Rate parameters to achieve the desired rate limiting behavior for your web API.

package token_bucket

import (
	"math"
	"sync"
	"time"
)

// Store MaxToken and Rate multiplied by a fixed coefficient, so that for very small fractions of time
// the token increment does not compute as a number smaller than 1.
const coefficient = 1000

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
		rate:                Rate * coefficient,
		maxTokens:           MaxTokens * coefficient,
		currentTokens:       MaxTokens * coefficient,
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
	if tokensToBeAdded == 0 {
		panic("TokenBucket: tokensToBeAdded is 0!")
	}
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
	tokens *= coefficient
	if tb.currentTokens >= tokens {
		tb.currentTokens -= tokens
		return true
	}
	return false
}
