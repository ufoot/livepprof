// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"os"
	"sync"

	"github.com/ufoot/livepprof/internal/google/binutils"
	"github.com/ufoot/livepprof/internal/google/plugin"
)

var globalMu sync.Mutex
var globalBinutils binutils.Binutils
var globalObjFile *ObjFile

type Resolver interface {
	Name() string
	Resolve(uint64) (*Location, error)
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
func (bof *ObjFile) Resolve(addr uint64) (*Location, error) {
	if bof == nil {
		return nil, NilObjFileError{}
	}
	frames, err := bof.objFile.SourceLine(addr)
	if err != nil {
		return nil, err
	}
	if len(frames) < 1 {
		return nil, NoFrame0Error{}
	}
	return &Location{
		Addr:     addr,
		Function: frames[0].Func,
		File:     frames[0].File,
		Line:     frames[0].Line,
	}, nil
}
