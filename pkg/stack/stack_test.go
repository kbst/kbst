package stack

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestStackFromPathAKS2Envs(t *testing.T) {
	_, p, err := newTestRepoFromFixture("kubestack-starter-aks-2envs")
	assert.Equal(t, nil, err, nil)

	os.RemoveAll(p)
}

func TestStackFromPathEKS3Envs(t *testing.T) {
	_, p, err := newTestRepoFromFixture("kubestack-starter-eks-3envs")
	assert.Equal(t, nil, err, nil)

	os.RemoveAll(p)
}

func TestStackFromPathEKSELBDNS(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-eks-3envs")
	assert.Equal(t, nil, err, nil)

	assert.Len(t, s.Clusters(), 1, "incorrect number of clusters")
	assert.Len(t, s.NodePools(), 0, "incorrect number of node pools")
	assert.Len(t, s.Modules(), 0, "incorrect number of custom modules")

	services := s.Services()
	assert.Len(t, services, 1, "incorrect number of services")
	assert.Equal(t, "eks_gc0_eu-west-1_nginx", services[0].Name())

	os.RemoveAll(p)
}

func TestStackFromPathGKE4Envs(t *testing.T) {
	_, p, err := newTestRepoFromFixture("kubestack-starter-gke-4envs")
	assert.Equal(t, nil, err, nil)

	os.RemoveAll(p)
}

func TestStackFromPathMulti4Envs(t *testing.T) {
	_, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	os.RemoveAll(p)
}

func TestFilesMulti4Envs(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	f := s.root.Parser.Files()
	assert.Equal(t, 25, len(f), nil)

	os.RemoveAll(p)
}

func TestWriteChanges(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-eks-3envs")
	assert.Equal(t, nil, err, nil)

	namePrefix := "rm"
	region := "europe-west4"

	baseCfg := map[string]cty.Value{
		"name_prefix":                cty.StringVal(namePrefix),
		"project_id":                 cty.StringVal("terraform-kubestack-testing"),
		"region":                     cty.StringVal(region),
		"cluster_min_node_count":     cty.NumberIntVal(1),
		"cluster_initial_node_count": cty.NumberIntVal(1),
		"cluster_max_node_count":     cty.NumberIntVal(3),
		"cluster_node_locations":     cty.StringVal("europe-west4-a,europe-west4-b,europe-west4-c"),
		"cluster_machine_type":       cty.StringVal("e2-standard-8"),
		"cluster_min_master_version": cty.StringVal("1.22"),
	}

	c, err := s.AddCluster(namePrefix, "google", region, "", GenerateConfigurations(s.Environments, baseCfg))
	assert.Equal(t, err, nil, nil)

	clusterName := c.Name()
	poolName := "rm"

	npBaseCfg := map[string]cty.Value{
		"name":               cty.StringVal(poolName),
		"min_node_count":     cty.NumberIntVal(1),
		"initial_node_count": cty.NumberIntVal(1),
		"max_node_count":     cty.NumberIntVal(3),
		"machine_type":       cty.StringVal("e2-standard-8"),
	}

	_, err = s.AddNodePool(clusterName, poolName, GenerateConfigurations(s.Environments, npBaseCfg))
	assert.Equal(t, err, nil, nil)

	_, err = s.AddService(clusterName, "prometheus", "")
	assert.Equal(t, err, nil, nil)

	diffs, err := getGitDiffs(p)
	if err != nil {
		log.Println(err)
	}

	// after adding cluster, node pool and service
	// there must these exact changes
	assert.Equal(t, []string{
		" M Dockerfile",
		"?? gke_rm_europe-west4_cluster.tf",
		"?? gke_rm_europe-west4_node_pool_rm.tf",
		"?? gke_rm_europe-west4_providers.tf",
		"?? gke_rm_europe-west4_service_prometheus.tf",
	}, diffs, nil)

	// refresh Stack in memory
	s.FromPath()

	// removing the cluster also removes
	// it's node pools and services
	err = s.Remove("gke_rm_europe-west4")
	if err != nil {
		log.Println(err)
	}

	hasDiff, err := hasGitDiff(p)
	if err != nil {
		log.Println(err)
	}

	// after removing the cluster again
	// there must not be any changes
	diffs, err = getGitDiffs(p)
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, false, hasDiff, diffs)

	//os.RemoveAll(p)
}

func TestAddClusterRejectDuplicate(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Clusters() {
		_, err = s.AddCluster(ex.NamePrefix, ex.Provider, ex.Region, ex.Version, ex.Configurations)
		assert.EqualError(t, err, fmt.Sprintf("error: cluster %q already exists", ex.Name()), nil)
	}

	os.RemoveAll(p)
}

func TestAddClusterDiffNamePrefix(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Clusters() {
		// different namePrefix, everything else identical
		_, err = s.AddCluster("gc1", ex.Provider, ex.Region, ex.Version, ex.Configurations)
		assert.Equal(t, err, nil, nil)

	}

	os.RemoveAll(p)
}

func TestAddClusterDiffRegion(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	region := map[string]string{
		"aws":     "us-east-1",
		"azurerm": "centralus",
		"google":  "northamerica-northeast1",
	}

	for _, ex := range s.Clusters() {
		// different region, everything else identical
		cfg := ex.Configurations
		if ex.Provider == "aws" {
			cfg[0].Attributes["cluster_availability_zones"] = cty.StringVal("us-east-1a,us-east-1b,us-east-1c")
		}
		if ex.Provider == "azurerm" {
			cfg[0].Attributes["availability_zones"] = cty.StringVal("1,2,3")
		}
		if ex.Provider == "google" {
			cfg[0].Attributes["cluster_node_locations"] = cty.StringVal("northamerica-northeast1-a,northamerica-northeast1-b,northamerica-northeast1-c")
		}
		_, err = s.AddCluster(ex.NamePrefix, ex.Provider, region[ex.Provider], ex.Version, cfg)
		assert.Equal(t, err, nil, nil)
	}

	os.RemoveAll(p)
}

func TestAddNodePool(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.NodePools() {
		_, err = s.AddNodePool(ex.ClusterName, "test", ex.Configurations)
		assert.Equal(t, err, nil, nil)
	}

	os.RemoveAll(p)
}

func TestAddNodePoolRejectNoCluster(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	_, err = s.AddNodePool("no_such_cluster", "test", s.NodePools()[0].Configurations)
	assert.EqualError(t, err, "no cluster named \"no_such_cluster\" found", nil)

	os.RemoveAll(p)
}

func TestAddNodePoolRejectDuplicate(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.NodePools() {
		_, err = s.AddNodePool(ex.ClusterName, ex.PoolName, ex.Configurations)
		assert.EqualError(t, err, fmt.Sprintf("error: node pool %q already exists", ex.Name()), nil)
	}

	os.RemoveAll(p)
}

func TestAddService(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Clusters() {
		_, err = s.AddService(ex.Name(), "sealed-secrets", "")
		assert.Equal(t, err, nil, nil)
	}

	os.RemoveAll(p)
}

func TestAddServiceRejectNoCluster(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	_, err = s.AddService("no_such_cluster", "", "")
	assert.EqualError(t, err, "no cluster named \"no_such_cluster\" found", nil)

	os.RemoveAll(p)
}

func TestAddServiceRejectDuplicate(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Services() {
		_, err = s.AddService(ex.ClusterName, ex.EntryName, "")
		assert.EqualError(t, err, fmt.Sprintf("error: service %q already exists", ex.Name()), nil)
	}

	os.RemoveAll(p)
}

func TestRemoveCluster(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	clusters := s.Clusters()
	nodePools := s.NodePools()
	services := s.Services()

	expLenC := len(clusters) - 1
	expLenNP := len(nodePools) - 1
	expLenSVC := len(services) - 3

	err = s.Remove(clusters[0].Name())
	assert.Equal(t, nil, err, nil)

	assert.Equal(t, expLenC, len(s.Clusters()), nil)
	assert.Equal(t, expLenNP, len(s.NodePools()), nil)
	assert.Equal(t, expLenSVC, len(s.Services()), nil)

	os.RemoveAll(p)
}

func TestRemoveClusterRejectLast(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-eks-3envs")
	assert.Equal(t, nil, err, nil)

	clusters := s.Clusters()

	assert.Len(t, clusters, 1, nil)

	err = s.Remove(clusters[0].Name())
	assert.EqualError(t, err, "stacks require one cluster, not removing \"eks_gc0_eu-west-1\"", nil)

	os.RemoveAll(p)
}

func TestRemoveNodePool(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	nodePools := s.NodePools()

	expLenNP := len(nodePools) - 1

	err = s.Remove(nodePools[0].Name())
	assert.Equal(t, nil, err, nil)

	assert.Equal(t, expLenNP, len(s.NodePools()), nil)

	os.RemoveAll(p)
}

func TestRemoveService(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	services := s.Services()

	expLenSVC := len(services) - 1

	err = s.Remove(services[0].Name())
	assert.Equal(t, nil, err, nil)

	assert.Equal(t, expLenSVC, len(s.Services()), nil)

	os.RemoveAll(p)
}

func getGitDiffs(p string) (diffs []string, err error) {
	out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && git status -s", p)).CombinedOutput()
	if err != nil {
		return []string{}, fmt.Errorf("%s: %s", err, out)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))

	for scanner.Scan() {
		diffs = append(diffs, scanner.Text())
	}

	return diffs, nil
}

func hasGitDiff(p string) (bool, error) {
	diffs, err := getGitDiffs(p)
	if err != nil {
		return false, err
	}

	return len(diffs) > 0, nil
}

func newTestRepoFromFixture(n string) (*Stack, string, error) {
	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")

	// copy fixture into temp test directory
	out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s/* %s/", filepath.Join(fixturesPath, n), p)).CombinedOutput()
	if err != nil {
		return &Stack{}, p, fmt.Errorf("%s: %s", err, out)
	}

	// initialize git repo in temp test directory
	out, err = exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && git init . && git add . && git commit -m initial", p)).CombinedOutput()
	if err != nil {
		return &Stack{}, p, fmt.Errorf("%s: %s", err, out)
	}

	r := tfhcl.NewRoot(p)
	cj := util.CliJSON{}
	err = cj.Load(util.CachedDownloader{})
	if err != nil {
		return &Stack{}, "", err
	}
	s := NewStack(r, cj)

	err = s.FromPath()

	return s, p, err
}
