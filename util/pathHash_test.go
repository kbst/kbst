package util

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cwd, _ = os.Getwd()
var fixturesPath = path.Join(cwd, "test_fixtures")

func TestPathHashChars(t *testing.T) {
	p, _ := PathHash(fixturesPath)

	assert.Len(t, p, 7, nil)
}

func TestPathHashDestinct(t *testing.T) {
	path1 := path.Join(fixturesPath, "eks")
	p1, _ := PathHash(path1)

	path2 := path.Join(fixturesPath, "multi-cloud")
	p2, _ := PathHash(path2)

	assert.NotEqual(t, p1, p2, nil)
}
