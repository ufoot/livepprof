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

func TestLPBasic(t *testing.T) {
	assert := assert.New(t)

	buf1 := allocator1(1e5)
	buf2a := allocator2(1e6)
	buf2b := allocator2(1e6)

	errHandler := func(err error) {
		assert.Nil(err)
	}

	a := New("livepprof", errHandler, time.Second/10, 100)
	defer a.Close()
	assert.NotNil(a)
	assert.NotNil(a.Heap())

	f := func() {
		for h := range a.Heap() {
			t.Logf("%v", h)
		}
	}
	go f()

	time.Sleep(time.Second)

	assert.Equal(byte(0), buf1[1])
	assert.Equal(byte(0), buf2a[1])
	assert.Equal(byte(0), buf2b[1])

}
