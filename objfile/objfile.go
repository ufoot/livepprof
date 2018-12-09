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

type Resolver interface {
	Name() string
	Resolve(addrs []uint64) (*Location, error)
}

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

func (bof *ObjFile) Name() string {
	if bof == nil {
		return ""
	}
	return bof.objFile.Name()
}

// Resolve returns the leaf source line for a location.
func (bof *ObjFile) Resolve(addrs []uint64) (*Location, error) {
	if bof == nil {
		return nil, NilObjFileError{}
	}
	if len(addrs) < 1 {
		return nil, NoAddrError{}
	}
	frames, err := bof.objFile.SourceLine(addrs[0])
	if err != nil {
		return nil, err
	}
	if len(frames) < 1 {
		return nil, NoFrame0Error{}
	}
	funcs := make([]string, 0, len(addrs))
	// Starting at len(addrs)-2, len(addrs)-1 is usually runtime.goexit, not interesting
	for i := len(addrs) - 2; i >= 1; i-- {
		frames, err := bof.objFile.SourceLine(addrs[i])
		if err != nil {
			return nil, err
		}
		if len(frames) < 1 {
			return nil, NoFrame0Error{}
		}
		funcs = append(funcs, funcOnly(frames[0].Func))
	}
	funcs = append(funcs, funcOnly(frames[0].Func))
	return &Location{
		Addr:     addrs[0],
		Function: frames[0].Func,
		File:     frames[0].File,
		Line:     frames[0].Line,
		Stack:    strings.Join(funcs, "/"),
	}, nil
}
