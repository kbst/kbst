package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/flosch/pongo2/v4"
	"github.com/stretchr/testify/assert"
)

var cwd, _ = os.Getwd()
var fixturesPath = path.Join(cwd, "../", "../", "test_fixtures", "generator")

func TestCfgToHCLNodePoolGoogle(t *testing.T) {
	m := Module{
		Name:     "",
		Provider: "google",
		Type:     "node_pool",
		Children: []Module{},
		Configuration: map[string]map[string]interface{}{
			"apps": {
				"name":           "test",
				"machine_type":   "test",
				"node_locations": "test1,test2,test3",
			},

			"ops": {},
		},
	}

	hcl, err := m.cfgToHCL("apps", "gke_test_test")

	assert.Equal(t, nil, err, nil)

	expected := "{\n    apps = {\n      machine_type = \"test\"\n      name = \"test\"\n      node_locations = [\n        \"test1\",\n        \"test2\",\n        \"test3\",\n      ]\n      project_id = module.gke_test_test.current_config[\"project_id\"]\n    }\n    ops = {}\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLNodePoolAzure(t *testing.T) {
	m := Module{
		Name:     "",
		Provider: "azurerm",
		Type:     "node_pool",
		Children: []Module{},
		Configuration: map[string]map[string]interface{}{
			"apps": {
				"node_pool_name": "test",
				"vm_size":        "test",
			},

			"ops": {},
		},
	}

	hcl, err := m.cfgToHCL("apps", "")

	assert.Equal(t, nil, err, nil)

	expected := "{\n    apps = {\n      node_pool_name = \"test\"\n      vm_size = \"test\"\n    }\n    ops = {}\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLNodePoolAWS(t *testing.T) {
	m := Module{
		Name:     "",
		Provider: "aws",
		Type:     "node_pool",
		Children: []Module{},
		Configuration: map[string]map[string]interface{}{
			"apps": {
				"name":          "test",
				"instance_type": "test",
			},

			"ops": {},
		},
	}

	hcl, err := m.cfgToHCL("apps", "")

	assert.Equal(t, nil, err, nil)

	expected := "{\n    apps = {\n      instance_type = \"test\"\n      name = \"test\"\n    }\n    ops = {}\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLService(t *testing.T) {
	m := Module{
		Name:     "test",
		Provider: "kustomization",
		Type:     "service",
		Children: []Module{},
		Configuration: map[string]map[string]interface{}{
			"apps": {
				"variant": nil,
			},

			"ops": {},
		},
	}

	hcl, err := m.cfgToHCL("apps", "")

	assert.Equal(t, nil, err, nil)

	expected := "{\n    apps = {}\n    ops = {}\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLClusterGoogle(t *testing.T) {
	m := Module{
		Name:     "",
		Provider: "google",
		Type:     "cluster",
		Children: []Module{},
		Configuration: map[string]map[string]interface{}{
			"apps": {
				"project_id":  "test",
				"name_prefix": "test",
				"region":      "test",
			},

			"ops": {},
		},
	}

	hcl, err := m.cfgToHCL("apps", "")

	assert.Equal(t, nil, err, nil)

	expected := "{\n    apps = {\n      base_domain = var.base_domain\n      name_prefix = \"test\"\n      project_id = \"test\"\n      region = \"test\"\n    }\n    ops = {}\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLClusterAWS(t *testing.T) {
	m := Module{
		Name:     "",
		Provider: "aws",
		Type:     "cluster",
		Children: []Module{},
		Configuration: map[string]map[string]interface{}{
			"apps": {
				"name_prefix": "test",
				"region":      "test",
			},

			"ops": {},
		},
	}

	hcl, err := m.cfgToHCL("apps", "")

	assert.Equal(t, nil, err, nil)

	expected := "{\n    apps = {\n      base_domain = var.base_domain\n      name_prefix = \"test\"\n    }\n    ops = {}\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestCfgToHCLClusterAzurerm(t *testing.T) {
	m := Module{
		Name:     "",
		Provider: "azurerm",
		Type:     "cluster",
		Children: []Module{},
		Configuration: map[string]map[string]interface{}{
			"apps": {
				"name_prefix":    "test",
				"region":         "test",
				"resource_group": "test",
			},

			"ops": {},
		},
	}

	hcl, err := m.cfgToHCL("apps", "")

	assert.Equal(t, nil, err, nil)

	expected := "{\n    apps = {\n      base_domain = var.base_domain\n      name_prefix = \"test\"\n      resource_group = \"test\"\n    }\n    ops = {}\n  }"
	assert.Equal(t, expected, hcl, nil)
}

func TestRender(t *testing.T) {
	tpl := pongo2.Must(pongo2.FromString("{{test}}"))
	d := pongo2.Context{"test": "test"}
	s, err := render(tpl, d)

	assert.Equal(t, nil, err, nil)
	assert.Equal(t, "test\n", s, nil)
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
