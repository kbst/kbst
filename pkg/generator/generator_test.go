package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

var cwd, _ = os.Getwd()
var fixturesPath = path.Join(cwd, "../", "../", "test_fixtures", "generator")

func TestCfgToHCLNodePool(t *testing.T) {
	m := Module{
		Name: "",
		Provider: "test",
		Type: "node_pool",
		Children: []Module{},
		Configuration: map[string]interface{}{
			"name": "test",
			"instance_type": "test",
		},
	}

	hcl, err := m.cfgToHCL()

	assert.Equal(t, nil, err, nil)

	expected := "{\n    \"instance_type\" = \"test\"\n    \"name\" = \"test\"\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLService(t *testing.T) {
	m := Module{
		Name: "test",
		Provider: "kustomization",
		Type: "service",
		Children: []Module{},
		Configuration: map[string]interface{}{
			"variant": nil,
		},
	}

	hcl, err := m.cfgToHCL()

	assert.Equal(t, nil, err, nil)

	expected := "{\n    \"variant\" = null\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLCluster(t *testing.T) {
	m := Module{
		Name: "",
		Provider: "test",
		Type: "cluster",
		Children: []Module{},
		Configuration: map[string]interface{}{
			"name_prefix": "test",
			"region": "test",
		},
	}

	hcl, err := m.cfgToHCL()

	assert.Equal(t, nil, err, nil)

	expected := "{\n    \"base_domain\" = var.base_domain\n    \"name_prefix\" = \"test\"\n    \"region\" = \"test\"\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLClusterAWS(t *testing.T) {
	m := Module{
		Name: "",
		Provider: "aws",
		Type: "cluster",
		Children: []Module{},
		Configuration: map[string]interface{}{
			"name_prefix": "test",
			"region": "test",
		},
	}

	hcl, err := m.cfgToHCL()

	assert.Equal(t, nil, err, nil)

	expected := "{\n    \"base_domain\" = var.base_domain\n    \"name_prefix\" = \"test\"\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestRender(t *testing.T) {
	tpl := template.Must(template.New("tpl").Parse("{{.Test}}"))
	d := make(map[string]string)
	d["Test"] = "test"
	s, err := render(tpl, d)

	assert.Equal(t, nil, err, nil)
	assert.Equal(t, d["Test"], s, nil)
}

func TestStackUnmarshal(t *testing.T) {
	p := filepath.Join(fixturesPath, "single_eks.json")
	f, err := ioutil.ReadFile(p)
	assert.Equal(t, nil, err, nil)

	s := Stack{}
	s.Unmarshal(f)

	assert.IsType(t, []Environment{}, s.Environments, nil)
	assert.IsType(t, []Module{}, s.Modules, nil)
}

func TestStackTerraformSingleAKS(t *testing.T) {
	p := filepath.Join(fixturesPath, "single_aks.json")
	f, err := ioutil.ReadFile(p)
	assert.Equal(t, nil, err, nil)

	s := Stack{}
	s.Unmarshal(f)

	files, err := s.Terraform()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"aks_kbst_westeurope_cluster.tf",
		"aks_kbst_westeurope_providers.tf",
		"aks_kbst_westeurope_node_pool_default.tf",
		"aks_kbst_westeurope_cluster_service_nginx.tf"}

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

	s := Stack{}
	s.Unmarshal(f)

	files, err := s.Terraform()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"eks_kbst_eu-west-1_cluster.tf",
		"eks_kbst_eu-west-1_providers.tf",
		"eks_kbst_eu-west-1_node_pool_default.tf",
		"eks_kbst_eu-west-1_cluster_service_nginx.tf"}

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

	s := Stack{}
	s.Unmarshal(f)

	files, err := s.Terraform()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"gke_kbst_europe-west1_cluster.tf",
		"gke_kbst_europe-west1_providers.tf",
		"gke_kbst_europe-west1_node_pool_default.tf",
		"gke_kbst_europe-west1_cluster_service_nginx.tf"}

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

	s := Stack{}
	s.Unmarshal(f)

	files, err := s.Terraform()
	assert.Equal(t, nil, err, nil)

	expected := []string{
		"versions.tf",
		"variables.tf",
		"config.auto.tfvars",
		"aks_kbst_westeurope_cluster.tf",
		"aks_kbst_westeurope_providers.tf",
		"aks_kbst_westeurope_node_pool_default.tf",
		"aks_kbst_westeurope_cluster_service_nginx.tf",
		"eks_kbst_eu-west-1_cluster.tf",
		"eks_kbst_eu-west-1_providers.tf",
		"eks_kbst_eu-west-1_node_pool_default.tf",
		"eks_kbst_eu-west-1_cluster_service_nginx.tf",
		"gke_kbst_europe-west1_cluster.tf",
		"gke_kbst_europe-west1_providers.tf",
		"gke_kbst_europe-west1_node_pool_default.tf",
		"gke_kbst_europe-west1_cluster_service_nginx.tf"}

	assert.Equal(t, len(expected), len(files), "list of files does not have expected length")

	for _, k := range expected {
		_, ok := files[k]
		assert.Equal(t, true, ok, fmt.Sprintf("%q not in list of files", k))
	}
}