// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package heap

import (
	"bytes"
	"runtime/pprof"

	"github.com/google/pprof/profile"

	"github.com/ufoot/livepprof/collector"
	"github.com/ufoot/livepprof/objfile"
)

// NoHeapProfileError when there's no heap profile.
type NoHeapProfileError struct{}

// Error string.
func (e NoHeapProfileError) Error() string {
	return "no heap profile"
}

// NoLocationError when no location can be found.
type NoLocationError struct{}

// Error string.
func (e NoLocationError) Error() string {
	return "no location"
}

// Heap collector.
type Heap struct {
	contains string
}

var _ collector.Collector = &Heap{}

// New heap collector.
func New(contains string) *Heap {
	return &Heap{
		contains: contains,
	}
}

// Collect data.
func (h *Heap) Collect(exit <-chan struct{}) (map[objfile.Location]float64, error) {
	rp := pprof.Lookup("heap")
	if rp == nil {
		return nil, NoHeapProfileError{}
	}

	var buf bytes.Buffer

	err := rp.WriteTo(&buf, 2)
	if err != nil {
		return nil, err
	}

	gp, err := profile.Parse(&buf)
	if err != nil {
		return nil, err
	}

	objFile, err := objfile.New()
	if err != nil {
		return nil, err
	}

	ret := make(map[objfile.Location]float64)
	for _, sample := range gp.Sample {
		if len(sample.Location) < 1 {
			return nil, NoLocationError{}
		}
		addresses := make([]uint64, 0, len(sample.Location))
		for _, loc := range sample.Location {
			addresses = append(addresses, loc.Address)
		}
		loc, err := objFile.Resolve(h.contains, addresses)
		if err != nil {
			return nil, err
		}
		if loc == nil {
			return nil, NoLocationError{}
		}
		bytesValues := sample.NumLabel["bytes"]
		var d float64
		for _, v := range bytesValues {
			d += float64(v)
		}
		if d > 0 {
			ret[*loc] += d
		}
	}

	return ret, nil
}
