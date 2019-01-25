// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package collector

import (
	"github.com/ufoot/livepprof/objfile"
)

// Collector is a generic interface to collect data.
type Collector interface {
	// Collect data, and return a map of values by location.
	// Can be interrupted by closing exit chan.
	Collect(exit <-chan struct{}) (map[objfile.Location]float64, error)
}
