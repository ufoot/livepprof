// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"fmt"
)

type Location struct {
	Addr     uint64
	Function string
	File     string
	Line     int
}

var _ fmt.Stringer = &Location{}

func (loc *Location) String() string {
	if loc == nil {
		return "{}"
	}
	return fmt.Sprintf(`{"addr":"%x","function":"%s","file":"%s","line":%d}`,
		loc.Addr, loc.Function, loc.File, loc.Line)
}
