// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFile(t *testing.T) {
	assert := assert.New(t)

	f := findBestFile("this-certainly-does-not-exist")
	assert.Equal("this-certainly-does-not-exist", f)

	envPath := os.Getenv("PATH")
	dirs := strings.Split(envPath, ":")
	var foundBin bool
	for _, dir := range dirs {
		if dir == "/bin" {
			foundBin = true
			break
		}
	}
	if !foundBin {
		t.Logf("bin not in path \"%s\", skipping test", envPath)
		t.Skip()
		return
	}
	foundBinLs := canOpen("/bin/ls")
	if !foundBinLs {
		t.Logf("cound not open \"/bin/ls\", skipping test")
		t.Skip()
		return
	}
	foundDotLs := canOpen("ls")
	if foundDotLs {
		t.Logf("cound open \"ls\", skipping test")
		t.Skip()
		return
	}

	f = findBestFile("ls")
	assert.Contains(f, "/ls", "path should contain /ls, be /bin/ls, or /usr/bin/ls")
	assert.True(path.IsAbs(f))
	t.Logf("found ls in \"%s\"", f)
}
