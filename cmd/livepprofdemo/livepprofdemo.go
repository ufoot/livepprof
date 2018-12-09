// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package main

import (
	"log"
	"math"
	"time"

	"github.com/ufoot/livepprof"
)

// something really stupid, but triggers activity...
func something(exit <-chan struct{}) float64 {
	var f float64
	var i uint64

	const n = 1000

	for {
		bufs := make([][]float64, n)
		for i = 0; i < 1000; i++ {
			f += math.Sqrt(float64(i))
			bufs[i] = make([]float64, i+1)
			bufs[i][i] = f
		}

		select {
		case <-exit:
			return f
		default: // non-blocking
		}
	}
}

func main() {
	lp := livepprof.New(func(err error) { log.Printf("%v", err) }, 3*time.Second, 100)
	defer lp.Close()

	exit := make(chan struct{})
	go func() {
		log.Printf("f=%0.1f", something(exit))
	}()

	go func() {
		log.Printf("ready to log cpu")
		for cpu := range lp.Cpu() {
			log.Printf("cpu timestamp=%v", cpu.Timestamp)
			for i, entry := range cpu.Entries {
				log.Printf("cpu %d/%d: %s -> %0.1f",
					i+1, len(cpu.Entries),
					entry.Key.String(),
					entry.Value,
				)
			}
		}
	}()

	go func() {
		log.Printf("ready to log heap")
		for heap := range lp.Heap() {
			log.Printf("heap timestamp=%v", heap.Timestamp)
			for i, entry := range heap.Entries {
				log.Printf("heap %d/%d: %s -> %0.1f",
					i+1, len(heap.Entries),
					entry.Key.String(),
					entry.Value,
				)
			}
		}
	}()

	time.Sleep(5 * time.Minute)
	close(exit)
}
