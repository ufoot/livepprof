// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package collector

import (
	"github.com/ufoot/livepprof/objfile"
)

type Collector interface {
	Collect() (map[objfile.Location]float64, error)
}
