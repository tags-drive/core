package limiter

import (
	"time"
)

// limit is used for RateLimiter. It is NOT thread safe
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

	if l.limit == 0 {
		return false
	}

	l.lastTime = time.Now()
	l.limit--
	return true
}

// update updates limit.
//
// It must be called right at the beginning of take()
func (l *limit) update() {
	if l.limit == l.max {
		return
	}

	delta := time.Since(l.lastTime)
	n := uint32(delta / l.timeout)
	if l.limit+n >= l.max {
		l.limit = l.max
		return
	}

	l.limit += n
}
