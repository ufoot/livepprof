// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

import (
	"sync"
	"time"

	"github.com/ufoot/livepprof/collector"
	"github.com/ufoot/livepprof/collector/heap"
)

type LP struct {
	errHandler    func(err error)
	delay         time.Duration
	heapCollector collector.Collector
	heaps         chan Data
	exit          chan struct{}
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

var _ Profiler = &LP{}

func New(errHandler func(err error), delay time.Duration) *LP {
	lp := &LP{
		errHandler:    errHandler,
		delay:         delay,
		heapCollector: heap.New(),
	}
	lp.Start()
	return lp
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
			data := Data{
				Timestamp: now,
			}
			for k, v := range rawData {
				data.Entries = append(data.Entries, Entry{Key: k, Value: v})
			}
			timeout := time.NewTimer(lp.delay)
			select {
			case heaps <- data:
				if !timeout.Stop() {
					<-timeout.C
				}
			case <-timeout.C:
				lp.handleErr(err)
				continue
			}
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
	lp.heaps = make(chan Data)

	lp.wg.Add(1)
	go lp.runHeaps(lp.heaps)
}

func (lp *LP) Stop() {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.exit == nil {
		return
	}

	close(lp.exit)
	lp.wg.Wait()
	lp.exit = nil
	close(lp.heaps)
	lp.heaps = nil
}

func (lp *LP) Close() {
	lp.Stop()
}
