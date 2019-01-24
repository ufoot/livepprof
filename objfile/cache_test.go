// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheKey(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(uint64(0xcbf29ce484222325), cacheKey(nil))
	assert.Equal(uint64(0xcbf29ce484222325), cacheKey([]uint64{}))
	assert.Equal(uint64(0xe3757ca7d64666ea), cacheKey([]uint64{1}))
	assert.Equal(uint64(0x0c4c250d07de76e8), cacheKey([]uint64{1, 3, 5}))
}

func TestCache(t *testing.T) {
	assert := assert.New(t)

	key1 := []uint64{1}
	key2 := []uint64{2, 3}
	l1 := Location{
		Function: "function1",
		File:     "file1",
		Stack:    "stack1",
	}
	l2 := Location{
		Function: "function2",
		File:     "file2",
		Stack:    "stack2",
	}

	c := newCache()

	assert.Nil(c.get(key1))
	assert.Nil(c.get(key2))
	assert.Equal(0, c.len())

	c.set(key1, &l1)
	l3 := c.get(key1)
	assert.NotNil(l3)
	assert.Equal(l1, *l3)
	assert.Nil(c.get(key2))
	assert.Equal(1, c.len())

	c.set(key2, &l2)
	l4 := c.get(key2)
	assert.NotNil(l4)
	assert.Equal(l2, *l4)
	assert.Equal(2, c.len())

	c.set(key2, &l1)
	l5 := c.get(key2)
	assert.NotNil(l5)
	assert.Equal(l1, *l5)
	assert.Equal(2, c.len())
}
