// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package cpu

import (
	"bytes"
	"os"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/google/pprof/profile"

	"github.com/ufoot/livepprof/collector"
	"github.com/ufoot/livepprof/objfile"
)

// NoLocationError when no location can be found.
type NoLocationError struct{}

// Error string.
func (e NoLocationError) Error() string {
	return "no location"
}

// UnexpectedValueLenError when the value array does not have expected size.
type UnexpectedValueLenError struct{}

// Error string.
func (e UnexpectedValueLenError) Error() string {
	return "unexpected value len"
}

// CPU collector.
type CPU struct {
	contains string
	delay    time.Duration
}

var _ collector.Collector = &CPU{}

// New CPU collector.
func New(contains string, delay time.Duration) *CPU {
	return &CPU{
		contains: contains,
		delay:    delay,
	}
}

func sigProfile() error {
	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGPROF)
}

// Collect data.
func (c *CPU) Collect() (map[objfile.Location]float64, error) {
	var buf bytes.Buffer

	err := pprof.StartCPUProfile(&buf)
	if err := sigProfile(); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	time.Sleep(c.delay)
	pprof.StopCPUProfile()

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
		loc, err := objFile.Resolve(c.contains, addresses)
		if err != nil {
			return nil, err
		}
		if loc == nil {
			return nil, NoLocationError{}
		}
		if len(sample.Value) != 2 {
			return nil, UnexpectedValueLenError{}
		}
		// [TODO:ufoot], really figure out what those numbers are...
		d := float64(sample.Value[0])
		if d > 0 {
			ret[*loc] += d
		}
	}

	return ret, nil
}
