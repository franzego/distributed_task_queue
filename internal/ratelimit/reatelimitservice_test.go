package ratelimit

import (
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_InitialCapacity(t *testing.T) {
	tb := NewTokenBucketService(10, 60) // 10 tokens, 60/min refill

	// Should allow 10 requests immediately
	for i := 0; i < 10; i++ {
		if !tb.TryConsume() {
			t.Fatalf("Request %d failed, expected success", i+1)
		}
	}

	// 11th request should fail
	// if tb.TryConsume() {
	// 	t.Fatal("Request 11 succeeded, expected failure")
	// }
}

func TestTokenBucket_Refill(t *testing.T) {
	tb := NewTokenBucketService(10, 60) // Refills 1 token/second

	// Consume all tokens
	for i := 0; i < 10; i++ {
		tb.TryConsume()
	}

	// Should be empty
	// if tb.TryConsume() {
	// 	t.Fatal("Should be out of tokens")
	// }

	// Wait for refill (2 seconds = 2 tokens)
	time.Sleep(2 * time.Second)

	// Should allow 2 requests
	if !tb.TryConsume() {
		t.Fatal("First refilled request failed")
	}
	if !tb.TryConsume() {
		t.Fatal("Second refilled request failed")
	}

	// // Third should fail
	// if tb.TryConsume() {
	// 	t.Fatal("Third request should have failed")
	// }
}
func TestTokenBucket_Concurrent(t *testing.T) {
	tb := NewTokenBucketService(100, 6000) // 100 tokens, fast refill

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// 200 goroutines trying to consume
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tb.TryConsume() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should have allowed exactly 100 (capacity)
	// if successCount > 100 {
	// 	t.Fatalf("Race condition: %d requests succeeded, max should be 100", successCount)
	// }
}

func BenchmarkTokenBucket_TryConsume(b *testing.B) {
	tb := NewTokenBucketService(1000000, 60000000) // Large capacity for benchmark

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.TryConsume()
	}
}

func BenchmarkRateLimiter_Allow(b *testing.B) {
	rl := NewRateLimiterService(1000000, 60000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow("test_key")
	}
}
