package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrefixRegion(t *testing.T) {
	p, r, _ := parsePrefixRegion("test_prefix_region")

	assert.Equal(t, "prefix", p)
	assert.Equal(t, "region", r)
}

func TestParsePrefixRegionError(t *testing.T) {
	_, _, err := parsePrefixRegion("test_invalid")

	assert.Error(t, err)
}
