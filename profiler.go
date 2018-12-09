// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

// Profiler can profile code dynamically.
type Profiler interface {
	// Cpu channel on which cpu data is sent.
	Cpu() <-chan Data
	// Heap channel on which heap data is sent.
	Heap() <-chan Data
	// Close the profiler.
	Close()
}
