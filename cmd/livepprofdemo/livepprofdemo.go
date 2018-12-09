// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package main

import (
	"log"
	"time"

	"github.com/ufoot/livepprof"
)

func main() {
	lp := livepprof.New(func(err error) { log.Printf("%v", err) }, 3*time.Second, 5)
	defer lp.Close()

	go func() {
		log.Printf("ready to log heap")
		for heap := range lp.Heap() {
			log.Printf("heap timesamp=%v", heap.Timestamp)
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
}
