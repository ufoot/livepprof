// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Location struct {
	Addr     uint64
	Function string
	File     string
	Line     int
	Stack    string
}

var _ fmt.Stringer = &Location{}

func (loc *Location) String() string {
	if loc == nil {
		return "{}"
	}
	js, err := json.Marshal(loc)
	if err != nil {
		return "{}"
	}
	return string(js)
}

func funcOnly(f string) string {
	li := strings.LastIndex(f, "/")
	if li < 0 {
		return f
	}
	return f[li+1:]
}
