// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

type NoArgs0Error struct{}

func (e NoArgs0Error) Error() string {
	return "no args[0]"
}

type NoFrame0Error struct{}

func (e NoFrame0Error) Error() string {
	return "no frame[0]"
}

type NoAddrError struct{}

func (e NoAddrError) Error() string {
	return "no addr"
}

type NotEnoughAddrsError struct{}

func (e NotEnoughAddrsError) Error() string {
	return "not enough addrs"
}

type NilObjFileError struct{}

func (e NilObjFileError) Error() string {
	return "nil obj file"
}
