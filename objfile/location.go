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

// Location identifies a place in the code. It is used to aggregate data.
// It does not have the uint64 address or the file line because this would
// lead to high cardinality and, for instance, different points of the same
// function would be counted in different entries. OTOH the stack trace is
// considered a key field, to know where the call comes from.
type Location struct {
	Function string
	File     string
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
