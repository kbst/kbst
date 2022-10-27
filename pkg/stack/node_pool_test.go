package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestNodePoolNameEKS(t *testing.T) {
	np := NodePool{
		NameSuffix:     "test-extra",
		ClusterName:    "eks_kbsteks_test-region-1",
		Provider:       "aws",
		Region:         "test-region",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	n := np.Name()

	assert.Equal(t, "eks_kbsteks_test-region-1_node_pool_test-extra", n)
}

func TestNodePoolNameAKS(t *testing.T) {
	np := NodePool{
		NameSuffix:     "test-extra",
		ClusterName:    "aks_kbstaks_test-continent",
		Provider:       "azurerm",
		Region:         "test-region",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	n := np.Name()

	assert.Equal(t, "aks_kbstaks_test-continent_node_pool_test-extra", n)
}

func TestNodePoolNameGKE(t *testing.T) {
	np := NodePool{
		NameSuffix:     "test-extra",
		ClusterName:    "gke_kbstgke_test-region1",
		Provider:       "azurerm",
		Region:         "test-region",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	n := np.Name()

	assert.Equal(t, "gke_kbstgke_test-region1_node_pool_test-extra", n)
}

func TestNodePoolToHCL(t *testing.T) {
	np := NodePool{
		NameSuffix:     "test-extra",
		ClusterName:    "gke_kbstgke_test-region1",
		Provider:       "azurerm",
		Region:         "test-region",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	files := np.ToHCL()

	assert.ElementsMatch(t, maps.Keys(files), []string{"gke_kbstgke_test-region1_node_pool_test-extra.tf"})

	for _, d := range files {
		assert.NotEqual(t, 0, len(d.Bytes()))
	}
}
