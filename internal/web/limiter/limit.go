package limiter

import (
	"sync/atomic"
	"time"
)

// limit is used for RateLimiter. It is thread safe
//
// It inceases limit by 1 after "timeout" since LAST successful take() call
type limit struct {
	limit    uint32
	lastTime time.Time

	timeout time.Duration
	max     uint32
}

func newLimit(maxRates uint32, timeout time.Duration) *limit {
	return &limit{
		limit:    maxRates,
		lastTime: time.Now(),
		timeout:  timeout,
		max:      maxRates,
	}
}

// take reduces l.limit by 1 when l.limit > 0. If it isn't possible, function returns false
func (l *limit) take() bool {
	l.update()

	curr := atomic.LoadUint32(&l.limit)
	if curr == 0 {
		return false
	}

	l.lastTime = time.Now()
	atomic.StoreUint32(&l.limit, curr-1)
	return true
}

// update updates limit.
//
// It must be called right at the beginning of take()
func (l *limit) update() {
	curr := atomic.LoadUint32(&l.limit)
	if curr == l.max {
		return
	}

	delta := time.Since(l.lastTime)
	n := uint32(delta / l.timeout)
	if curr+n >= l.max {
		atomic.StoreUint32(&l.limit, l.max)
		return
	}

	atomic.AddUint32(&l.limit, n)
}
