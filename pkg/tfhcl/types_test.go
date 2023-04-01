package tfhcl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeProviderVersion(t *testing.T) {
	r := NewRoot()
	err := r.Read("fixtures")
	assert.Equal(t, nil, err)

	for _, fs := range r.Modules {
		for _, m := range fs {
			k, p, v, err := m.TypeProviderVersion()

			if m.Source == "test_source" {
				assert.NotEqual(t, nil, err)
			}

			// EKS cluster
			if m.Source == "github.com/kbst/terraform-kubestack//aws/cluster?ref=test-version" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "cluster", k)
				assert.Equal(t, "aws", p)
				assert.Equal(t, "test-version", v)
			}

			// EKS node-pool
			if m.Source == "github.com/kbst/terraform-kubestack//aws/cluster/node-pool?ref=test-version" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "node_pool", k)
				assert.Equal(t, "aws", p)
				assert.Equal(t, "test-version", v)
			}

			// EKS elb-dns
			if m.Source == "github.com/kbst/terraform-kubestack//aws/cluster/elb-dns?ref=test-version" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "elb-dns", k)
				assert.Equal(t, "aws", p)
				assert.Equal(t, "test-version", v)
			}

			// GKE cluster
			if m.Source == "github.com/kbst/terraform-kubestack//google/cluster?ref=test-version" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "cluster", k)
				assert.Equal(t, "google", p)
				assert.Equal(t, "test-version", v)
			}

			// GKE node-pool
			if m.Source == "github.com/kbst/terraform-kubestack//google/cluster/node-pool?ref=test-version" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "node_pool", k)
				assert.Equal(t, "google", p)
				assert.Equal(t, "test-version", v)
			}

			// AKS cluster
			if m.Source == "github.com/kbst/terraform-kubestack//azurerm/cluster?ref=test-version" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "cluster", k)
				assert.Equal(t, "azurerm", p)
				assert.Equal(t, "test-version", v)
			}

			// AKS node-pool
			if m.Source == "github.com/kbst/terraform-kubestack//azurerm/cluster/node-pool?ref=test-version" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "node_pool", k)
				assert.Equal(t, "azurerm", p)
				assert.Equal(t, "test-version", v)
			}

			// catalog service
			if m.Source == "kbst.xyz/catalog/test/kustomization" {
				assert.Equal(t, nil, err)
				assert.Equal(t, "service", k)
				assert.Equal(t, "kustomization", p)
				assert.Equal(t, "0.0.0-test.0", v)
			}
		}
	}
}
