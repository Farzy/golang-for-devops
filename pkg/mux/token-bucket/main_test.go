package token_bucket

import (
	"testing"
	"time"
)

func TestNewTokenBucket(t *testing.T) {
	tb := NewTokenBucket(10, 100)
	if tb.rate != 10*coefficient {
		t.Errorf("Expected bucket rate to be 10, but got %d", tb.rate)
	}

	if tb.maxTokens != 100*coefficient {
		t.Errorf("Expected max tokens to be 100, but got %d", tb.maxTokens)
	}
}

func TestTokenBucket_IsRequestAllowed(t *testing.T) {
	tb := NewTokenBucket(10, 100)
	if !tb.IsRequestAllowed(10) {
		t.Errorf("Expected IsRequestAllowed to return true when tokens are sufficient")
	}
	if tb.IsRequestAllowed(200) {
		t.Errorf("Expected IsRequestAllowed to return false when tokens are insufficient")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	tb := NewTokenBucket(10, 100)
	if !tb.IsRequestAllowed(100) {
		t.Errorf("Expected IsRequestAllowed to return true when tokens are sufficient")
	}
	if tb.currentTokens != 0 {
		t.Errorf("Expected tokens count to be 0 after consuming all tokens")
	}

	time.Sleep(300 * time.Millisecond) // just for test purposes
	tb.refill()
	if tb.currentTokens/coefficient != 3 {
		t.Errorf("Expected tokens count to be sufficient after sleep")
	}
	if tb.IsRequestAllowed(5) {
		t.Errorf("Expected IsRequestAllowed to return false when tokens are insufficient")
	}

	time.Sleep(200 * time.Millisecond) // just for test purposes
	if !tb.IsRequestAllowed(5) {
		t.Errorf("Expected IsRequestAllowed to return true when tokens are sufficient")
	}
}
