package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimit(t *testing.T) {
	const (
		rate    = 5
		timeout = time.Millisecond * 200
	)

	assert := assert.New(t)
	l := newLimit(rate, timeout)

	// -1
	assert.Equal(true, l.take())
	assert.Equal(uint32(rate-1), l.limit)
	time.Sleep(timeout)
	// limit wasn't updated
	assert.Equal(uint32(rate-1), l.limit)
	// limit was updated
	l.update()
	assert.Equal(uint32(rate), l.limit)

	// -4
	assert.Equal(true, l.take())
	assert.Equal(true, l.take())
	assert.Equal(true, l.take())
	assert.Equal(true, l.take())
	assert.Equal(uint32(rate-4), l.limit)
	//
	time.Sleep(timeout)
	l.update()
	// recover 1 pts
	assert.Equal(uint32(rate-3), l.limit)
	//
	// reset limit.lastTime
	l.lastTime = time.Now()
	time.Sleep(2 * timeout)
	l.update()
	// recover 2 pts
	assert.Equal(uint32(rate-1), l.limit)

	// reset limiter to not wait timeout * 5
	l = newLimit(rate, timeout)
	for i := 0; i < rate; i++ {
		assert.Equal(true, l.take())
	}
	assert.Equal(false, l.take())
	assert.Equal(false, l.take())
	assert.Equal(false, l.take())
	time.Sleep(timeout)
	assert.Equal(true, l.take())
}

func TestRateLimiterWithRateFive(t *testing.T) {
	const (
		rate    = 5
		timeout = time.Second / rate
	)

	assert := assert.New(t)
	l := NewRateLimiter(rate, timeout)

	b := l.Take("500")
	assert.Equal(b, true)
	b = l.Take("400")
	assert.Equal(b, true)
	assert.Equal(l.limits["500"].limit, uint32(rate-1))
	assert.Equal(l.limits["400"].limit, uint32(rate-1))

	time.Sleep(timeout)

	// weren't updated
	assert.Equal(l.limits["500"].limit, uint32(rate-1))
	assert.Equal(l.limits["400"].limit, uint32(rate-1))

	// limits must not change after Take() (+1 after update and -1 after return)
	assert.Equal(l.Take("500"), true)
	assert.Equal(l.Take("400"), true)
	assert.Equal(l.limits["500"].limit, uint32(rate-1))
	assert.Equal(l.limits["400"].limit, uint32(rate-1))

	for i := 0; i < rate; i++ {
		assert.Equal(true, l.Take("300"))
	}
	assert.Equal(false, l.Take("300"))
}

func TestRateLimiterWithRateOne(t *testing.T) {
	const (
		rate    = 1
		timeout = time.Millisecond * 10
		id      = "500"
	)

	assert := assert.New(t)
	l := NewRateLimiter(rate, timeout)

	assert.Equal(true, l.Take(id))
	assert.Equal(false, l.Take(id))

	time.Sleep(timeout)
	assert.Equal(true, l.Take(id))
	assert.Equal(false, l.Take(id))
	assert.Equal(false, l.Take(id))
}
