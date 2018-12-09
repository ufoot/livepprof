// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

// NoArgs0Error when program name can't be found.
type NoArgs0Error struct{}

// Error string.
func (e NoArgs0Error) Error() string {
	return "no args[0]"
}

// NoFrame0Error when there's not even one frame in sample.
type NoFrame0Error struct{}

// Error string.
func (e NoFrame0Error) Error() string {
	return "no frame[0]"
}

// NoAddrError when there's no address.
type NoAddrError struct{}

// Error string.
func (e NoAddrError) Error() string {
	return "no addr"
}

// NilObjFileError when obj file is not initialized.
type NilObjFileError struct{}

// Error string.
func (e NilObjFileError) Error() string {
	return "nil obj file"
}
