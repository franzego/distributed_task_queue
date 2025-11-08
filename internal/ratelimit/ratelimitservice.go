package ratelimit

import "sync"

type RateLimiter struct {
	capacity   int
	refillRate int
	Buckets    map[string]*TokenBucket
	mu         *sync.Mutex
}

func NewRateLimiterService(capacity, refillRate int) *RateLimiter {
	return &RateLimiter{
		Buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
		mu:         &sync.Mutex{},
	}
}
func (r *RateLimiter) Allow(keyID string) bool {
	r.mu.Lock()
	// defer r.mu.Unlock()
	bucket, exists := r.Buckets[keyID]
	if !exists {
		bucket = NewTokenBucketService(r.capacity, r.refillRate)
		r.Buckets[keyID] = bucket
	}
	r.mu.Unlock()
	return bucket.TryConsume()
}
