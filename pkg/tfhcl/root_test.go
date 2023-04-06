package tfhcl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	r := NewRoot("fixtures")
	err := r.Read()

	assert.Equal(t, nil, err)

	// r.Variables
	for fn, vs := range r.Variables {
		expL := 0
		if fn == "fixtures/test_root_read_variable.tf" {
			expL = 1
		}
		assert.Equal(t, expL, len(vs), fmt.Sprintf("%s: found %d unexpected variables", fn, len(vs)))
	}

	// r.Modules
	mCount := 0
	for _, ms := range r.Modules {
		mCount += len(ms)
	}
	assert.Equal(t, 14, mCount, nil)

	// r.Providers
	for fn, ps := range r.Providers {
		expL := 0
		if fn == "fixtures/test_root_read_providers.tf" {
			expL = 2
		}
		assert.Equal(t, expL, len(ps), fmt.Sprintf("%s: found %d unexpected providers", fn, len(ps)))
	}
}

func TestReadTwoModules(t *testing.T) {
	r := NewRoot("fixtures")
	err := r.Read()

	assert.Equal(t, nil, err)

	mods := r.Modules["fixtures/test_root_read_two_modules.tf"]

	assert.Len(t, mods, 2, nil)
	assert.Equal(t, "test_mod1", mods[0].Name, nil)
	assert.Equal(t, "test_mod2", mods[1].Name, nil)
}
