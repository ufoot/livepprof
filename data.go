// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

import (
	"time"

	"github.com/ufoot/livepprof/objfile"
)

type Entry struct {
	Key   objfile.Location
	Value float64
}

type Data struct {
	Timestamp time.Time
	Entries   []Entry
}
