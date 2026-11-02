package main

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	limiter := &Limiter{
		clock: clock,
		rate: ratePerSec,
		burst: burst,
	}

	if clock == nil || burst <= 0 {
		return limiter
	}

	limiter.tokens = float64(burst)
	limiter.last = clock.Now()

	return limiter
}
func (l *Limiter) AllowN(n int) bool {
	if l.clock == nil || l.burst <= 0 || n <= 0 || n > l.burst {
		return false
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	
	now := l.clock.Now()

	if now.After(l.last) {
		elapsed := now.Sub(l.last).Seconds()

		if l.rate > 0 {
			refill := elapsed * l.rate

			if refill >= float64(l.burst) - l.tokens {
				l.tokens = float64(l.burst)
			} else {
				l.tokens += refill
			}
		}

		l.last = now
	}

	if l.tokens < float64(n) {
		return false
	}

	l.tokens -= float64(n)
	return true
}
