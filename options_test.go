// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

import (
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
