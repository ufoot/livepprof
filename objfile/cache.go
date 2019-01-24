// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"encoding/binary"
	"hash/fnv"
	"sync"
)

// cache for locations, avoids resolving the same things over and over.
type cache struct {
	mu   sync.RWMutex
	data map[uint64]Location
}

func cacheKey(addrs []uint64) uint64 {
	h := fnv.New64()
	buf := make([]byte, 8)
	for _, addr := range addrs {
		binary.LittleEndian.PutUint64(buf, addr)
		_, err := h.Write(buf)
		if err != nil {
			return addr // really, should never happen
		}
	}
	return h.Sum64()
}

func newCache() *cache {
	return &cache{}
}

func (c *cache) set(addrs []uint64, l *Location) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.data == nil {
		c.data = make(map[uint64]Location)
	}
	key := cacheKey(addrs)
	c.data[key] = *l
}

func (c *cache) get(addrs []uint64) *Location {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.data == nil {
		return nil
	}
	key := cacheKey(addrs)
	v, ok := c.data[key]
	if !ok {
		return nil
	}
	return &v
}

func (c *cache) len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}
