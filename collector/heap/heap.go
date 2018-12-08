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

type NoHeapProfileError struct{}

func (e NoHeapProfileError) Error() string {
	return "no heap profile"
}

type NoLocationError struct{}

func (e NoLocationError) Error() string {
	return "no args[0]"
}

type Heap struct{}

var _ collector.Collector = &Heap{}

func New() *Heap {
	return &Heap{}
}

func (h *Heap) Collect() (map[objfile.Location]float64, error) {
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
		loc, err := objFile.Resolve(sample.Location[0].Address)
		if err != nil {
			return nil, err
		}
		if loc == nil {
			return nil, NoLocationError{}
		}
		bytesValues := sample.NumLabel["bytes"]
		bytesValue := ret[*loc]
		for _, v := range bytesValues {
			bytesValue += float64(v)
		}
		ret[*loc] = bytesValue
	}

	return ret, nil
}
