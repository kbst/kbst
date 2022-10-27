package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cwd, _ = os.Getwd()
var fixturesPath = path.Join(cwd, "../", "../", "test_fixtures", "generator")

func TestStackUnmarshal(t *testing.T) {
	p := filepath.Join(fixturesPath, "single_eks.json")
	f, err := ioutil.ReadFile(p)
	assert.Equal(t, nil, err, nil)

	s := LegacyStack{}
	s.Unmarshal(f)

	assert.IsType(t, []LegacyEnvironment{}, s.Environments, nil)
	assert.IsType(t, []LegacyModule{}, s.Modules, nil)
}

func TestStackTerraformSingleAKS(t *testing.T) {
	p := filepath.Join(fixturesPath, "single_aks.json")
	f, err := ioutil.ReadFile(p)
	assert.Equal(t, nil, err, nil)

	ls := LegacyStack{}
	s, err := ls.Unmarshal(f)
	assert.Equal(t, nil, err, nil)

	files, err := s.Files()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"aks_kbst_westeurope_cluster.tf",
		"aks_kbst_westeurope_providers.tf",
		"aks_kbst_westeurope_service_nginx.tf"}

	assert.Equal(t, len(expected), len(files), "list of files does not have expected length")

	for _, k := range expected {
		_, ok := files[k]
		assert.Equal(t, true, ok, fmt.Sprintf("%q not in list of files %+v", k, files))
	}
}

func TestStackTerraformSingleEKS(t *testing.T) {
	p := filepath.Join(fixturesPath, "single_eks.json")
	f, err := ioutil.ReadFile(p)
	assert.Equal(t, nil, err, nil)

	ls := LegacyStack{}
	s, err := ls.Unmarshal(f)
	assert.Equal(t, nil, err, nil)

	files, err := s.Files()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"eks_kbst_eu-west-1_cluster.tf",
		"eks_kbst_eu-west-1_providers.tf",
		"eks_kbst_eu-west-1_service_nginx.tf"}

	assert.Equal(t, len(expected), len(files), "list of files does not have expected length")

	for _, k := range expected {
		_, ok := files[k]
		assert.Equal(t, true, ok, fmt.Sprintf("%q not in list of files %+v", k, files))
	}
}

func TestStackTerraformSingleGKE(t *testing.T) {
	p := filepath.Join(fixturesPath, "single_gke.json")
	f, err := ioutil.ReadFile(p)
	assert.Equal(t, nil, err, nil)

	ls := LegacyStack{}
	s, err := ls.Unmarshal(f)
	assert.Equal(t, nil, err, nil)

	files, err := s.Files()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"gke_kbst_europe-west1_cluster.tf",
		"gke_kbst_europe-west1_providers.tf",
		"gke_kbst_europe-west1_service_nginx.tf"}

	assert.Equal(t, len(expected), len(files), "list of files does not have expected length")

	for _, k := range expected {
		_, ok := files[k]
		assert.Equal(t, true, ok, fmt.Sprintf("%q not in list of files %+v", k, files))
	}
}

func TestStackTerraformMultiCloud(t *testing.T) {
	p := filepath.Join(fixturesPath, "multi_cloud.json")
	f, err := ioutil.ReadFile(p)
	assert.Equal(t, nil, err, nil)

	ls := LegacyStack{}
	s, err := ls.Unmarshal(f)
	assert.Equal(t, nil, err, nil)

	files, err := s.Files()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"aks_kbst_westeurope_cluster.tf",
		"aks_kbst_westeurope_providers.tf",
		"aks_kbst_westeurope_node_pool_extra.tf",
		"aks_kbst_westeurope_service_nginx.tf",
		"aks_kbst_westeurope_service_cert-manager.tf",
		"eks_kbst_eu-west-1_cluster.tf",
		"eks_kbst_eu-west-1_providers.tf",
		"eks_kbst_eu-west-1_node_pool_extra.tf",
		"eks_kbst_eu-west-1_service_nginx.tf",
		"eks_kbst_eu-west-1_service_cert-manager.tf",
		"gke_kbst_europe-west1_cluster.tf",
		"gke_kbst_europe-west1_providers.tf",
		"gke_kbst_europe-west1_node_pool_extra.tf",
		"gke_kbst_europe-west1_service_nginx.tf",
		"gke_kbst_europe-west1_service_cert-manager.tf"}

	assert.Equal(t, len(expected), len(files), "list of files does not have expected length")

	for _, k := range expected {
		_, ok := files[k]
		assert.Equal(t, true, ok, fmt.Sprintf("%q not in list of files", k))
	}
}
