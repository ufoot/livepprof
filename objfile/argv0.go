// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"os"
	"path"
	"strings"
)

func canOpen(filename string) bool {
	f, err := os.Open(filename)
	if err != nil || f == nil {
		return false
	}
	if f.Close() != nil {
		return false
	}
	return true
}

// findBestFile does the job for findArgv0, isolated to be easier to test.
func findBestFile(filename string) string {
	if canOpen(filename) {
		// easiest case, we could open it directly
		return filename
	}
	envPath := os.Getenv("PATH")
	dirs := strings.Split(envPath, ":") // [TODO:ufoot] make this cross-platform...
	for _, dir := range dirs {
		f := path.Join(dir, filename)
		if canOpen(f) {
			return f
		}
	}

	// Still return what we have, you never know, binary might show up later, etc.
	return filename
}

// findArgv0 is a small utility to find the current running binary.
// Indeed using os.Args[0] does not work out-of-the-box, it might
// be that it's not accessible from CWD so PATH needs to be inspected.
func findArgv0() (string, error) {
	if len(os.Args) < 1 {
		return "", NoArgs0Error{}
	}
	return findBestFile(os.Args[0]), nil
}
