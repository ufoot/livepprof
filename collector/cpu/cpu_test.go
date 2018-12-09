// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package cpu

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func busy1(exit <-chan struct{}) float64 {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	var j float64
	for {
		select {
		case <-ticker.C:
			for i := 0; i < 1e6; i++ {
				j += math.Sqrt(float64(i))
			}
		case <-exit:
			return j
		}
	}
}

func busy2(exit <-chan struct{}) float64 {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	var j float64
	for {
		select {
		case <-ticker.C:
			for i := 0; i < 1e5; i++ {
				j += math.Sqrt(float64(i))
			}
		case <-exit:
			return j
		}
	}
}

func TestCollect(t *testing.T) {
	assert := assert.New(t)

	exit := make(chan struct{})
	go func() {
		t.Logf("busy1: %0.1f", busy1(exit))
	}()
	go func() {
		t.Logf("busy2: %0.1f", busy2(exit))
	}()

	h := New(3 * time.Second)
	assert.NotNil(h)

	data, err := h.Collect()
	assert.Nil(err)
	assert.NotNil(data)

	assert.True(len(data) > 0, fmt.Sprintf("len(data): %d should be >0", len(data)))
	for k, v := range data {
		t.Logf("%s: %0.1f", k.String(), v)
	}

	close(exit)
}
