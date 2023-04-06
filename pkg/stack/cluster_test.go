package stack

import (
	"log"
	"testing"

	"github.com/kbst/kbst/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/maps"
)

var testCfgs []Configuration = []Configuration{
	{
		EnvironmentKey: "apps",
		Attributes:     make(map[string]cty.Value),
	},
	{
		EnvironmentKey: "ops",
		Attributes:     make(map[string]cty.Value),
	},
}

func TestClusterNameEKS(t *testing.T) {
	c := Cluster{
		NamePrefix:     "kbsteks",
		Provider:       "aws",
		Region:         "test-region-1",
		Version:        "test-version",
		Configurations: testCfgs,
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
		Configurations: testCfgs,
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
		Configurations: testCfgs,
	}

	n := c.Name()

	assert.Equal(t, "gke_kbstgke_test-region1", n)
}

func TestClusterRejectEmptyConfiguration(t *testing.T) {
	c := Cluster{
		NamePrefix:     "kbstgke",
		Provider:       "google",
		Region:         "test-region1",
		Version:        "test-version",
		Configurations: []Configuration{},
	}

	cj := util.CliJSON{}
	err := cj.Load(util.CachedDownloader{})
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	err = c.Validate(cj)
	assert.EqualError(t, err, "invalid empty configuration []", nil)
}

func TestClusterToHCL(t *testing.T) {
	c := Cluster{
		NamePrefix:     "kbstgke",
		Provider:       "google",
		Region:         "test-region1",
		Version:        "test-version",
		Configurations: testCfgs,
	}

	files := c.ToHCL()

	assert.ElementsMatch(t, maps.Keys(files), []string{"gke_kbstgke_test-region1_cluster.tf", "gke_kbstgke_test-region1_providers.tf"})

	for _, d := range files {
		assert.NotEqual(t, 0, len(d))
	}
}
