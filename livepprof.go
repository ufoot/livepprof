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

type LP struct {
	errHandler    func(err error)
	delay         time.Duration
	limit         int
	cpuCollector  collector.Collector
	heapCollector collector.Collector
	cpus          chan Data
	heaps         chan Data
	exit          chan struct{}
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

var _ Profiler = &LP{}

func New(contains string, errHandler func(err error), delay time.Duration, limit int) *LP {
	lp := &LP{
		errHandler:    errHandler,
		delay:         delay,
		limit:         limit,
		cpuCollector:  cpu.New(contains, delay),
		heapCollector: heap.New(contains),
	}
	lp.Start()
	return lp
}

func (lp *LP) Cpu() <-chan Data {
	lp.mu.RLock()
	defer lp.mu.RUnlock()

	return lp.cpus
}

func (lp *LP) Heap() <-chan Data {
	lp.mu.RLock()
	defer lp.mu.RUnlock()

	return lp.heaps
}

func (lp *LP) handleErr(err error) {
	if lp.errHandler != nil {
		lp.errHandler(err)
	}
}

func (lp *LP) runCpus(cpus chan<- Data) {
	defer lp.wg.Done()

	ticker := time.NewTicker(lp.delay)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			rawData, err := lp.cpuCollector.Collect()
			if err != nil {
				lp.handleErr(err)
				continue
			}
			data := buildData(now, rawData, lp.limit)
			cpus <- data
		case <-lp.exit:
			return
		}
	}
}

func (lp *LP) runHeaps(heaps chan<- Data) {
	defer lp.wg.Done()

	ticker := time.NewTicker(lp.delay)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			rawData, err := lp.heapCollector.Collect()
			if err != nil {
				lp.handleErr(err)
				continue
			}
			data := buildData(now, rawData, lp.limit)
			heaps <- data
		case <-lp.exit:
			return
		}
	}
}

func (lp *LP) Start() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.exit != nil {
		return
	}

	lp.exit = make(chan struct{})
	lp.cpus = make(chan Data)
	lp.heaps = make(chan Data)

	lp.wg.Add(2)
	go lp.runCpus(lp.cpus)
	go lp.runHeaps(lp.heaps)
}

func (lp *LP) Stop() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.exit == nil {
		return
	}

	close(lp.exit)

	// Drain chan to avoid it blocking.
	go func() {
		for _ = range lp.heaps {
		}
	}()
	go func() {
		for _ = range lp.cpus {
		}
	}()

	lp.wg.Wait()
	lp.exit = nil
	close(lp.heaps)
	lp.heaps = nil
	close(lp.cpus)
	lp.cpus = nil
}

func (lp *LP) Close() {
	lp.Stop()
}
