package limiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate    uint32
	timeout time.Duration

	limits map[string]*limit
	mutex  *sync.Mutex
}

func NewRateLimiter(rate uint32, timeout time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:    rate,
		timeout: timeout,

		limits: make(map[string]*limit),
		mutex:  new(sync.Mutex),
	}
}

func (l *RateLimiter) Take(id string) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	lim, ok := l.limits[id]
	if !ok {
		l.limits[id] = newLimit(l.rate, l.timeout)
		lim = l.limits[id]
	}

	return lim.take()
}
