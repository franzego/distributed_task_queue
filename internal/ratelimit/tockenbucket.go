package ratelimit

import (
	"sync"
	"time"
)

type TokenBucket struct {
	Tokens         float64
	Capacity       int
	LastRefillTime time.Time
	RefillRate     float64
	mu             sync.RWMutex
}

func NewTokenBucketService(capacity int, refillRate int) *TokenBucket {
	return &TokenBucket{
		Tokens:         float64(capacity),
		Capacity:       capacity,
		LastRefillTime: time.Now(),
		RefillRate:     float64(refillRate),
	}
}

// function that checks if the tokens are enough to make a request
func (t *TokenBucket) TryConsume() bool {

	t.mu.Lock()
	defer t.mu.Unlock()
	t.Refill()
	if t.Tokens > 0 {
		t.Tokens--
		return true
	}
	return false

}
func (t *TokenBucket) Refill() {
	now := time.Now().UTC()
	elapsedTime := now.Sub(t.LastRefillTime).Seconds()
	tokenstoAdd := elapsedTime * t.RefillRate
	if tokenstoAdd > 0 {
		t.Tokens = min(tokenstoAdd+t.Tokens, float64(t.Capacity))
		t.LastRefillTime = now
	}
}
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b

}
