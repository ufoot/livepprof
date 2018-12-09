// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package livepprof

import (
	"sort"
	"strings"
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

type sortEntries struct {
	entries []Entry
}

func (se *sortEntries) Len() int {
	return len(se.entries)
}

func (se *sortEntries) Swap(i, j int) {
	se.entries[i], se.entries[j] = se.entries[j], se.entries[i]
}

func (se *sortEntries) Less(i, j int) bool {
	// using i>j as we want the greatest value 1st, reverse sort
	if se.entries[i].Value > se.entries[j].Value {
		return true
	}
	if se.entries[i].Value < se.entries[j].Value {
		return false
	}

	// unable to sort on value, sorting by key (should be rare)
	keyI := se.entries[i].Key
	keyJ := se.entries[j].Key
	if cmp := strings.Compare(keyI.Function, keyJ.Function); cmp < 0 {
		return true
	} else if cmp > 0 {
		return false
	}
	if cmp := strings.Compare(keyI.File, keyJ.File); cmp < 0 {
		return true
	} else if cmp > 0 {
		return false
	}
	if cmp := strings.Compare(keyI.Stack, keyJ.Stack); cmp < 0 {
		return true
	}
	return false
}

func buildData(ts time.Time, rawData map[objfile.Location]float64, limit int) Data {
	ts = ts.Truncate(time.Millisecond) // makes logs easier to read

	if limit <= 0 {
		return Data{Timestamp: ts}
	}

	ret := Data{
		Timestamp: ts,
		Entries:   make([]Entry, 0, limit),
	}

	for k, v := range rawData {
		ret.Entries = append(ret.Entries, Entry{Key: k, Value: v})
	}

	se := sortEntries{entries: ret.Entries}
	sort.Sort(&se)
	ret.Entries = se.entries

	if len(ret.Entries) > limit {
		ret.Entries = ret.Entries[:limit]
	}

	return ret
}
