// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package heap

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func allocator1(n int) []byte {
	if n <= 0 {
		return nil
	}
	buf := make([]byte, n)
	buf[0] = 255
	buf[n-1] = 255
	return buf
}

func allocator2(n int) []byte {
	if n <= 0 {
		return nil
	}
	buf := make([]byte, n)
	buf[0] = 255
	buf[n-1] = 255
	return buf
}

func TestCollect(t *testing.T) {
	assert := assert.New(t)

	buf1 := allocator1(1e5)
	buf2a := allocator2(1e6)
	buf2b := allocator2(1e6)

	time.Sleep(time.Second / 10)

	h := New("livepprof")
	assert.NotNil(h)

	data, err := h.Collect(nil)
	assert.Nil(err)
	assert.NotNil(data)

	assert.True(len(data) > 0, fmt.Sprintf("len(data): %d should be >0", len(data)))
	for k, v := range data {
		t.Logf("%s: %0.1f", k.String(), v)
	}

	assert.Equal(byte(0), buf1[1])
	assert.Equal(byte(0), buf2a[1])
	assert.Equal(byte(0), buf2b[1])
}
