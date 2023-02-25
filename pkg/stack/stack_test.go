package stack

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
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

	f, err := s.Files()
	assert.Equal(t, nil, err, nil)

	assert.Equal(t, 21, len(f), nil)

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
		"cluster_node_locations":     cty.StringVal(""),
		"cluster_machine_type":       cty.StringVal("e2-standard-8"),
		"cluster_min_master_version": cty.StringVal("1.22"),
	}

	err = s.AddCluster(namePrefix, "google", region, "", GenerateConfigurations(s.Environments, baseCfg))
	assert.Equal(t, err, nil, nil)

	clusterName := s.Clusters[len(s.Clusters)-1].Name()
	poolName := "rm"

	npBaseCfg := map[string]cty.Value{
		"name":               cty.StringVal(poolName),
		"min_node_count":     cty.NumberIntVal(1),
		"initial_node_count": cty.NumberIntVal(1),
		"max_node_count":     cty.NumberIntVal(3),
		"machine_type":       cty.StringVal("e2-standard-8"),
	}

	err = s.AddNodePool(clusterName, poolName, GenerateConfigurations(s.Environments, npBaseCfg))
	assert.Equal(t, err, nil, nil)

	err = s.AddService(clusterName, "prometheus", "")
	assert.Equal(t, err, nil, nil)

	s.WriteChanges()

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
	s.FromPath(p)

	// removing the cluster also removes
	// it's node pools and services
	err = s.Remove("gke_rm_europe-west4")
	if err != nil {
		log.Println(err)
	}

	s.WriteChanges()

	hasDiff, err := hasGitDiff(p)
	if err != nil {
		log.Println(err)
	}

	// after removing the cluster again
	// there must not be any changes
	assert.Equal(t, false, hasDiff, nil)

	os.RemoveAll(p)
}

func TestAddClusterRejectDuplicate(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Clusters {
		err = s.AddCluster(ex.NamePrefix, ex.Provider, ex.Region, ex.Version, ex.Configurations)
		assert.EqualError(t, err, fmt.Sprintf("error: cluster %q already exists", ex.Name()), nil)
	}

	os.RemoveAll(p)
}

func TestAddClusterDiffNamePrefix(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Clusters {
		// different namePrefix, everything else identical
		err = s.AddCluster("gc1", ex.Provider, ex.Region, ex.Version, ex.Configurations)
		assert.Equal(t, err, nil, nil)

	}

	os.RemoveAll(p)
}

func TestAddClusterDiffRegion(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	region := map[string]string{
		"aws":     "us-east-1",
		"azurerm": "central-us",
		"google":  "northamerica-northeast1",
	}

	for _, ex := range s.Clusters {
		// different region, everything else identical
		err = s.AddCluster(ex.NamePrefix, ex.Provider, region[ex.Provider], ex.Version, ex.Configurations)
		assert.Equal(t, err, nil, nil)
	}

	os.RemoveAll(p)
}

func TestAddNodePool(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.NodePools {
		err = s.AddNodePool(ex.ClusterName, "test", ex.Configurations)
		assert.Equal(t, err, nil, nil)
	}

	os.RemoveAll(p)
}

func TestAddNodePoolRejectNoCluster(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	err = s.AddNodePool("no_such_cluster", "test", s.NodePools[0].Configurations)
	assert.EqualError(t, err, "no cluster named \"no_such_cluster\" found", nil)

	os.RemoveAll(p)
}

func TestAddNodePoolRejectDuplicate(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.NodePools {
		err = s.AddNodePool(ex.ClusterName, ex.PoolName, ex.Configurations)
		assert.EqualError(t, err, fmt.Sprintf("error: node pool %q already exists", ex.Name()), nil)
	}

	os.RemoveAll(p)
}

func TestAddService(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Clusters {
		err = s.AddService(ex.Name(), "sealed-secrets", "")
		assert.Equal(t, err, nil, nil)
	}

	os.RemoveAll(p)
}

func TestAddServiceRejectNoCluster(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	err = s.AddService("no_such_cluster", "", "")
	assert.EqualError(t, err, "no cluster named \"no_such_cluster\" found", nil)

	os.RemoveAll(p)
}

func TestAddServiceRejectDuplicate(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	for _, ex := range s.Services {
		err = s.AddService(ex.ClusterName, ex.EntryName, "")
		assert.EqualError(t, err, fmt.Sprintf("error: service %q already exists", ex.Name()), nil)
	}

	os.RemoveAll(p)
}

func TestRemoveCluster(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	expLenC := len(s.Clusters) - 1
	expLenNP := len(s.NodePools) - 1
	expLenSVC := len(s.Services) - 3

	err = s.Remove(s.Clusters[0].Name())
	assert.Equal(t, nil, err, nil)

	assert.Equal(t, expLenC, len(s.Clusters), nil)
	assert.Equal(t, expLenNP, len(s.NodePools), nil)
	assert.Equal(t, expLenSVC, len(s.Services), nil)

	os.RemoveAll(p)
}

func TestRemoveClusterRejectLast(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-eks-3envs")
	assert.Equal(t, nil, err, nil)

	assert.Len(t, s.Clusters, 1, nil)

	err = s.Remove(s.Clusters[0].Name())
	assert.EqualError(t, err, "stacks require one cluster, not removing \"eks_gc0_eu-west-1\"", nil)

	os.RemoveAll(p)
}

func TestRemoveNodePool(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	expLenNP := len(s.NodePools) - 1

	err = s.Remove(s.NodePools[0].Name())
	assert.Equal(t, nil, err, nil)

	assert.Equal(t, expLenNP, len(s.NodePools), nil)

	os.RemoveAll(p)
}

func TestRemoveService(t *testing.T) {
	s, p, err := newTestRepoFromFixture("kubestack-starter-multi-4envs")
	assert.Equal(t, nil, err, nil)

	expLenSVC := len(s.Services) - 1

	err = s.Remove(s.Services[0].Name())
	assert.Equal(t, nil, err, nil)

	assert.Equal(t, expLenSVC, len(s.Services), nil)

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

	fmt.Printf("diffs: %v\n", diffs)

	return len(diffs) > 0, nil
}

func newTestRepoFromFixture(n string) (*Stack, string, error) {
	r := tfhcl.NewRoot()
	cj := util.CliJSON{}
	err := cj.Load(util.CachedDownloader{})
	if err != nil {
		return &Stack{}, "", err
	}
	s := NewStack(r, cj)

	// copy fixture into temp test directory
	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	out, err := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s/* %s/", path.Join(fixturesPath, n), p)).CombinedOutput()
	if err != nil {
		return &Stack{}, p, fmt.Errorf("%s: %s", err, out)
	}

	// initialize git repo in temp test directory
	out, err = exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && git init . && git add . && git commit -m initial", p)).CombinedOutput()
	if err != nil {
		return &Stack{}, p, fmt.Errorf("%s: %s", err, out)
	}

	err = s.FromPath(p)

	return s, p, err
}
