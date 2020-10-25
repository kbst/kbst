package runner

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathHashChars(t *testing.T) {
	p, _ := pathHash(fixturesPath)

	assert.Len(t, p, 7, nil)
}

func TestPathHashDestinct(t *testing.T) {
	path1 := path.Join(fixturesPath, "eks")
	p1, _ := pathHash(path1)

	path2 := path.Join(fixturesPath, "multi-cloud")
	p2, _ := pathHash(path2)

	assert.NotEqual(t, p1, p2, nil)
}
