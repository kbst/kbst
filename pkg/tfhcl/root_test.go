package tfhcl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestRead(t *testing.T) {
	r := NewRoot()
	err := r.Read("fixtures")

	assert.Equal(t, nil, err)

	// r.Variables
	for fn, vs := range r.Variables {
		expL := 0
		if fn == "fixtures/test_root_read_variable.tf" {
			expL = 1
		}
		assert.Equal(t, expL, len(vs), fmt.Sprintf("%s: found %d unexpected variables", fn, len(vs)))
	}

	// r.VariableValues
	assert.Equal(t, 1, len(r.VariableValues), nil)
	assert.Equal(t, cty.Value(cty.StringVal("value")), r.VariableValues["testvar"], nil)

	// r.Modules
	mCount := 0
	for _, ms := range r.Modules {
		mCount += len(ms)
	}
	assert.Equal(t, 10, mCount, nil)

	// r.Providers
	for fn, ps := range r.Providers {
		expL := 0
		if fn == "fixtures/test_root_read_providers.tf" {
			expL = 2
		}
		assert.Equal(t, expL, len(ps), fmt.Sprintf("%s: found %d unexpected providers", fn, len(ps)))
	}

	// r.Dockerfiles
	assert.Equal(t, 1, len(r.Dockerfiles), nil)
}
