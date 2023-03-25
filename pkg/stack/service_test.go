package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestServiceName(t *testing.T) {
	s := Service{
		EntryName:      "test",
		ClusterName:    "aks_kbstaks_test-continent",
		Provider:       "kustomization",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	n := s.Name()

	assert.Equal(t, "aks_kbstaks_test-continent_service_test", n)
}

func TestServiceToHCL(t *testing.T) {
	s := Service{
		EntryName:      "test",
		ClusterName:    "aks_kbstaks_test-continent",
		Provider:       "kustomization",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	files := s.ToHCL()

	assert.ElementsMatch(t, maps.Keys(files), []string{"aks_kbstaks_test-continent_service_test.tf"})

	for _, d := range files {
		assert.NotEqual(t, 0, len(d.Bytes()))
	}
}
