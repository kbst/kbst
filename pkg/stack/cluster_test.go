package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestClusterNameEKS(t *testing.T) {
	c := Cluster{
		NamePrefix:     "kbsteks",
		Provider:       "aws",
		Region:         "test-region-1",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	n := c.Name()

	assert.Equal(t, "eks_kbsteks_test-region-1", n)
}

func TestClusterNameAKS(t *testing.T) {
	c := Cluster{
		NamePrefix:     "kbstaks",
		Provider:       "azurerm",
		Region:         "test-continent",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	n := c.Name()

	assert.Equal(t, "aks_kbstaks_test-continent", n)
}

func TestClusterNameGKE(t *testing.T) {
	c := Cluster{
		NamePrefix:     "kbstgke",
		Provider:       "google",
		Region:         "test-region1",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	n := c.Name()

	assert.Equal(t, "gke_kbstgke_test-region1", n)
}

func TestClusterToHCL(t *testing.T) {
	c := Cluster{
		NamePrefix:     "kbstgke",
		Provider:       "google",
		Region:         "test-region1",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	files := c.ToHCL()

	assert.ElementsMatch(t, maps.Keys(files), []string{"gke_kbstgke_test-region1_cluster.tf", "gke_kbstgke_test-region1_providers.tf"})

	for _, d := range files {
		assert.NotEqual(t, 0, len(d.Bytes()))
	}
}
