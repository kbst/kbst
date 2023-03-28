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

// EKS
func TestParseKindProviderVersionEKSCluster(t *testing.T) {
	k, p, v, err := parseKindProviderVersion("github.com/kbst/terraform-kubestack//aws/cluster?ref=test-version", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, "cluster", k)
	assert.Equal(t, "aws", p)
	assert.Equal(t, "test-version", v)
}

func TestParseKindProviderVersionEKSNodePool(t *testing.T) {
	k, p, v, err := parseKindProviderVersion("github.com/kbst/terraform-kubestack//aws/cluster/node-pool?ref=test-version", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, "node_pool", k)
	assert.Equal(t, "aws", p)
	assert.Equal(t, "test-version", v)
}

// GKE
func TestParseKindProviderVersionGKECluster(t *testing.T) {
	k, p, v, err := parseKindProviderVersion("github.com/kbst/terraform-kubestack//google/cluster?ref=test-version", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, "cluster", k)
	assert.Equal(t, "google", p)
	assert.Equal(t, "test-version", v)
}

func TestParseKindProviderVersionGKENodePool(t *testing.T) {
	k, p, v, err := parseKindProviderVersion("github.com/kbst/terraform-kubestack//google/cluster/node-pool?ref=test-version", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, "node_pool", k)
	assert.Equal(t, "google", p)
	assert.Equal(t, "test-version", v)
}

// AKS
func TestParseKindProviderVersionAKSCluster(t *testing.T) {
	k, p, v, err := parseKindProviderVersion("github.com/kbst/terraform-kubestack//azurerm/cluster?ref=test-version", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, "cluster", k)
	assert.Equal(t, "azurerm", p)
	assert.Equal(t, "test-version", v)
}

func TestParseKindProviderVersionAKSNodePool(t *testing.T) {
	k, p, v, err := parseKindProviderVersion("github.com/kbst/terraform-kubestack//azurerm/cluster/node-pool?ref=test-version", "")
	assert.Equal(t, nil, err)

	assert.Equal(t, "node_pool", k)
	assert.Equal(t, "azurerm", p)
	assert.Equal(t, "test-version", v)
}

// Service
func TestParseKindProviderVersionService(t *testing.T) {
	k, p, v, err := parseKindProviderVersion("kbst.xyz/catalog/test/kustomization", "test-version")
	assert.Equal(t, nil, err)

	assert.Equal(t, "service", k)
	assert.Equal(t, "kustomization", p)
	assert.Equal(t, "test-version", v)
}
