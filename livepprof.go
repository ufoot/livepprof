// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

import (
	"sync"
	"time"

	"github.com/ufoot/livepprof/collector"
	"github.com/ufoot/livepprof/collector/cpu"
	"github.com/ufoot/livepprof/collector/heap"
)

// LP is an implementation of a live profiler.
type LP struct {
	opts          opts
	cpuCollector  collector.Collector
	heapCollector collector.Collector
	cpus          chan Data
	heaps         chan Data
	exit          chan struct{}
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

// Profiler is a generic profiler interface.
var _ Profiler = &LP{}

// New live profiler.
// The contains parameter is used to choose the leaf on which to aggregate data.
// Just choose something that is in your source files path, typically a top-level
// package name, namespace, whatever identifies your code.
func New(options ...Option) (*LP, error) {
	opts := defaultOpts
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}
	lp := &LP{
		opts:          opts,
		cpuCollector:  cpu.New(opts.filter, opts.delay),
		heapCollector: heap.New(opts.filter),
		cpus:          make(chan Data),
		heaps:         make(chan Data),
	}

	lp.Start()
	return lp, nil
}

// CPU channel on which cpu data is sent.
func (lp *LP) CPU() <-chan Data {
	lp.mu.RLock()
	defer lp.mu.RUnlock()

	return lp.cpus
}

// Heap channel on which heap data is sent.
func (lp *LP) Heap() <-chan Data {
	lp.mu.RLock()
	defer lp.mu.RUnlock()

	return lp.heaps
}

func (lp *LP) handleErr(err error) {
	if lp.opts.errHandler != nil {
		lp.opts.errHandler(err)
	}
}

func (lp *LP) runCPUs(cpus chan<- Data) {
	defer lp.wg.Done()

	ticker := time.NewTicker(lp.opts.delay)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			if !lp.opts.enabled() {
				continue
			}
			rawData, err := lp.cpuCollector.Collect(lp.exit)
			if err != nil {
				lp.handleErr(err)
				continue
			}
			data := buildData(now, rawData, lp.opts.limit)
			cpus <- data
		case <-lp.exit:
			return
		}
	}
}

func (lp *LP) runHeaps(heaps chan<- Data) {
	defer lp.wg.Done()

	ticker := time.NewTicker(lp.opts.delay)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			if !lp.opts.enabled() {
				continue
			}
			rawData, err := lp.heapCollector.Collect(lp.exit)
			if err != nil {
				lp.handleErr(err)
				continue
			}
			data := buildData(now, rawData, lp.opts.limit)
			heaps <- data
		case <-lp.exit:
			return
		}
	}
}

// Start the profiler.
func (lp *LP) Start() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.exit != nil ||
		lp.cpuCollector == nil || lp.heapCollector == nil ||
		lp.cpus == nil || lp.heaps == nil {
		return
	}

	lp.exit = make(chan struct{})

	lp.wg.Add(2)
	go lp.runCPUs(lp.cpus)
	go lp.runHeaps(lp.heaps)
}

func (lp *LP) stop() {
	if lp.exit == nil {
		return
	}

	close(lp.exit)

	// Drain chan to avoid it blocking.
	go func() {
		for range lp.heaps {
		}
	}()
	go func() {
		for range lp.cpus {
		}
	}()

	lp.wg.Wait()
	lp.exit = nil
}

// Stop the profiler.
func (lp *LP) Stop() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	lp.stop()
}

// Close the profiler. It can't be started again.
func (lp *LP) Close() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	lp.stop()

	lp.cpuCollector = nil
	lp.heapCollector = nil

	close(lp.heaps)
	lp.heaps = nil
	close(lp.cpus)
	lp.cpus = nil
}
