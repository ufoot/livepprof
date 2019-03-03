// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptsEnabled(t *testing.T) {
	assert := assert.New(t)

	var o opts

	assert.True(o.enabled())
	o.disabled = true
	assert.False(o.enabled())
	o.enabledFunc = func() bool { return true }
	assert.True(o.enabled())
	o.enabledFunc = func() bool { return false }
	assert.False(o.enabled())
}

func TestOptsJitteredDelay(t *testing.T) {
	assert := assert.New(t)

	const n = 1000

	var o opts

	o.delay = defaultDelay

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < n; i++ {
		// no jitter at all, we get the exact value
		assert.Equal(time.Minute, o.jitteredDelay(r))
	}

	o.jitter = defaultJitter
	var sum float64
	var sumSq float64
	for i := 0; i < n; i++ {
		delay := float64(o.jitteredDelay(r))
		sum += delay
		sumSq += (float64(o.delay) - delay) * (float64(o.delay) - delay)
	}
	avg := sum / n
	assert.InEpsilon(float64(o.delay), avg, 0.1)
	dev := math.Sqrt(sumSq/n) / avg
	assert.InEpsilon(o.jitter/math.Sqrt(12), dev, 0.1)
}

func TestWithFilter(t *testing.T) {
	assert := assert.New(t)

	o := defaultOpts
	assert.Equal("", o.filter)
	assert.Nil(WithFilter("toto")(&o))
	assert.Equal("toto", o.filter)
}

func TestWithErrorHandler(t *testing.T) {
	assert := assert.New(t)

	o := defaultOpts
	assert.Nil(o.errHandler)
	assert.Nil(WithErrorHandler(func(err error) {})(&o))
	assert.NotNil(o.errHandler)
}

func TestWithDelay(t *testing.T) {
	assert := assert.New(t)

	o := defaultOpts
	assert.Equal(time.Minute, o.delay)
	assert.Nil(WithDelay(time.Second)(&o))
	assert.Equal(time.Second, o.delay)
	assert.NotNil(WithDelay(0)(&o))
	assert.Equal(time.Second, o.delay)
	assert.NotNil(WithDelay(-1)(&o))
	assert.Equal(time.Second, o.delay)
}

func TestWithJitter(t *testing.T) {
	assert := assert.New(t)

	o := defaultOpts
	assert.Equal(0.1, o.jitter)
	assert.Nil(WithJitter(0.0)(&o))
	assert.Equal(0.0, o.jitter)
	assert.Nil(WithJitter(1.0)(&o))
	assert.Equal(1.0, o.jitter)
	assert.Nil(WithJitter(0.5)(&o))
	assert.Equal(0.5, o.jitter)
	assert.NotNil(WithJitter(1.1)(&o))
	assert.Equal(0.5, o.jitter)
	assert.NotNil(WithJitter(-0.1)(&o))
	assert.Equal(0.5, o.jitter)
}

func TestWithLimit(t *testing.T) {
	assert := assert.New(t)

	o := defaultOpts
	assert.Equal(20, o.limit)
	assert.Nil(WithLimit(100)(&o))
	assert.Equal(100, o.limit)
	assert.NotNil(WithLimit(0)(&o))
	assert.Equal(100, o.limit)
	assert.NotNil(WithLimit(-1)(&o))
	assert.Equal(100, o.limit)
}

func TestWithEnabled(t *testing.T) {
	assert := assert.New(t)

	o := defaultOpts
	assert.False(o.disabled)
	assert.True(o.enabled())
	assert.Nil(WithEnabled(false)(&o))
	assert.True(o.disabled)
	assert.False(o.enabled())
	assert.Nil(WithEnabled(true)(&o))
	assert.False(o.disabled)
	assert.True(o.enabled())
}

func TestWithEnabledFunc(t *testing.T) {
	assert := assert.New(t)

	o := defaultOpts
	assert.False(o.disabled)
	assert.True(o.enabled())
	assert.Nil(WithEnabledFunc(func() bool { return false })(&o))
	assert.False(o.disabled)
	assert.False(o.enabled())
	assert.Nil(WithEnabledFunc(func() bool { return true })(&o))
	assert.False(o.disabled)
	assert.True(o.enabled())
	o.disabled = true
	assert.True(o.enabled())
}
