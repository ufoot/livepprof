// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"os"
	"strings"
	"sync"

	"github.com/ufoot/livepprof/internal/google/binutils"
	"github.com/ufoot/livepprof/internal/google/plugin"
)

var globalMu sync.Mutex
var globalBinutils binutils.Binutils
var globalObjFile *ObjFile

// Resolver resolves addresses to locations.
type Resolver interface {
	// Name of the resolver. Typically, binary file name.
	Name() string
	// Resolve some addresses to a location.
	// The contains string is used to find the leaf on which to aggregate data.
	// The addrs should be ordered with the leaf in first positions, and callers after.
	Resolve(contains string, addrs []uint64) (*Location, error)
}

// ObjFile is an object file representation, used to resolve addresses.
type ObjFile struct {
	objFile plugin.ObjFile
}

var _ Resolver = &ObjFile{}

// New returns a global object allowing to resolve addresses to locations.
// This is a slingleton, and reports data only for self, the current program
// identified by os.Args[0]. When testing, go test tools do not embed symbols
// by default, you need to explicitly use `go test -o filename` else this
// will fail and not be able to get the info.
func New() (*ObjFile, error) {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalObjFile != nil {
		return globalObjFile, nil
	}
	if len(os.Args) < 1 {
		return nil, NoArgs0Error{}
	}
	f, err := globalBinutils.Open(os.Args[0], 0, ^uint64(0), 0)
	if err != nil {
		return nil, err
	}
	globalObjFile = &ObjFile{objFile: f}
	return globalObjFile, nil
}

// Name of the binary file.
func (bof *ObjFile) Name() string {
	if bof == nil {
		return ""
	}
	return bof.objFile.Name()
}

// Resolve returns the leaf source line for a location.
func (bof *ObjFile) Resolve(contains string, addrs []uint64) (*Location, error) {
	if bof == nil {
		return nil, NilObjFileError{}
	}
	if len(addrs) < 1 {
		return nil, NoAddrError{}
	}
	var leaf int
	for i, addr := range addrs {
		frames, err := bof.objFile.SourceLine(addr)
		if err != nil {
			return nil, err
		}
		if len(frames) < 1 {
			return nil, NoFrame0Error{}
		}
		if strings.Contains(frames[0].File, contains) {
			leaf = i
			break
		}
	}

	n := len(addrs) - leaf
	funcs := make([]string, 0, n)
	// Starting at len(addrs)-2, len(addrs)-1 is usually runtime.goexit, not interesting
	i0 := len(addrs) - 2
	if i0 < 0 {
		i0 = 0
	}

	loc := Location{}
	for i := i0; i >= leaf; i-- {
		frames, err := bof.objFile.SourceLine(addrs[i])
		if err != nil {
			return nil, err
		}
		if len(frames) < 1 {
			return nil, NoFrame0Error{}
		}
		funcs = append(funcs, funcOnly(frames[0].Func))
		if i == leaf {
			loc.Function = frames[0].Func
			loc.File = frames[0].File
		}
	}

	loc.Stack = strings.Join(funcs, "/")
	return &loc, nil
}
