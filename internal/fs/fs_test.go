package fs_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/nokia/ntt/internal/fs"
	"github.com/stretchr/testify/assert"
)

// TestBytesFromURL verifies, that files specified by an URL can be read.
func TestBytesFromURL(t *testing.T) {
	path, err := filepath.Abs("fs_test.go")
	if err != nil {
		panic(err)
	}

	expected, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	f := fs.Open("file://" + path)
	b, err := f.Bytes()
	assert.Nil(t, err)
	assert.Equal(t, expected, b)
}