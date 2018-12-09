// Live pprof is a Golang library to generate and use live profiles.
// Copyright (C)  2018  Christian Mauduit <ufoot@ufoot.org>
// Live pprof homepage: https://github.com/ufoot/livepprof
// Contact author: ufoot@ufoot.org

package objfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// [WARNING] you need to run these tests with `-o` else symbols won't be found

func TestObjFile(t *testing.T) {
	assert := assert.New(t)

	of, err := New()
	assert.Nil(err)
	assert.NotNil(of)

	l, err := of.Resolve("livepprof", nil)
	assert.NotNil(err, "this fails, but does not crash")
	assert.Nil(l)
}
