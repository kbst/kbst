package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrefixRegion(t *testing.T) {
	p, r := parsePrefixRegion("test_prefix_region")

	assert.Equal(t, "prefix", p)
	assert.Equal(t, "region", r)
}

// EKS
func TestParseKindProviderVersionEKSCluster(t *testing.T) {
	k, p, v := parseKindProviderVersion("github.com/kbst/terraform-kubestack//aws/cluster?ref=test-version", "")

	assert.Equal(t, "cluster", k)
	assert.Equal(t, "aws", p)
	assert.Equal(t, "test-version", v)
}

func TestParseKindProviderVersionEKSNodePool(t *testing.T) {
	k, p, v := parseKindProviderVersion("github.com/kbst/terraform-kubestack//aws/cluster/node-pool?ref=test-version", "")

	assert.Equal(t, "node_pool", k)
	assert.Equal(t, "aws", p)
	assert.Equal(t, "test-version", v)
}

// GKE
func TestParseKindProviderVersionGKECluster(t *testing.T) {
	k, p, v := parseKindProviderVersion("github.com/kbst/terraform-kubestack//google/cluster?ref=test-version", "")

	assert.Equal(t, "cluster", k)
	assert.Equal(t, "google", p)
	assert.Equal(t, "test-version", v)
}

func TestParseKindProviderVersionGKENodePool(t *testing.T) {
	k, p, v := parseKindProviderVersion("github.com/kbst/terraform-kubestack//google/cluster/node-pool?ref=test-version", "")

	assert.Equal(t, "node_pool", k)
	assert.Equal(t, "google", p)
	assert.Equal(t, "test-version", v)
}

// AKS
func TestParseKindProviderVersionAKSCluster(t *testing.T) {
	k, p, v := parseKindProviderVersion("github.com/kbst/terraform-kubestack//azurerm/cluster?ref=test-version", "")

	assert.Equal(t, "cluster", k)
	assert.Equal(t, "azurerm", p)
	assert.Equal(t, "test-version", v)
}

func TestParseKindProviderVersionAKSNodePool(t *testing.T) {
	k, p, v := parseKindProviderVersion("github.com/kbst/terraform-kubestack//azurerm/cluster/node-pool?ref=test-version", "")

	assert.Equal(t, "node_pool", k)
	assert.Equal(t, "azurerm", p)
	assert.Equal(t, "test-version", v)
}

// Service
func TestParseKindProviderVersionService(t *testing.T) {
	k, p, v := parseKindProviderVersion("kbst.xyz/catalog/test/kustomization", "test-version")

	assert.Equal(t, "service", k)
	assert.Equal(t, "kustomization", p)
	assert.Equal(t, "test-version", v)
}

func TestParseNodePoolClusterNameNameSuffix(t *testing.T) {
	cn, ns := parseNodePoolClusteNameNameSuffix("cluster-name_node_pool_name_suffix")

	assert.Equal(t, "cluster-name", cn)
	assert.Equal(t, "name_suffix", ns)
}

func TestParseServiceClusteNameEntryName(t *testing.T) {
	cn, en := parseServiceClusteNameEntryName("cluster_name_service_entry-name")

	assert.Equal(t, "cluster_name", cn)
	assert.Equal(t, "entry-name", en)
}
