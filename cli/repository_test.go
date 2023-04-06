package cli

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kbst/kbst/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

type MockDownloaderCliJson struct{}

func (c MockDownloaderCliJson) Download(url string) (resp *http.Response, err error) {
	p := filepath.Join(fixturesPath, "cli.json")
	f, err := ioutil.ReadFile(p)
	if err != nil {
		return resp, err
	}

	r := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(f)),
	}
	return r, nil
}

type MockDownloaderFrameworkArchive struct{}

func (c MockDownloaderFrameworkArchive) Download(url string) (resp *http.Response, err error) {
	fn := strings.Split(url, "/")[4]
	p := filepath.Join(fixturesPath, fn)
	f, err := ioutil.ReadFile(p)
	if err != nil {
		return resp, err
	}

	r := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(f)),
	}
	return r, nil
}

func TestRepoInitEKS(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	baseCfg := map[string]cty.Value{
		"name_prefix":                cty.StringVal("test"),
		"cluster_availability_zones": cty.StringVal("eu-west-1a,eu-west-1b,eu-west-1c"),
		"cluster_instance_type":      cty.StringVal("m5a.2xlarge"),
		"cluster_min_size":           cty.NumberIntVal(3),
		"cluster_desired_capacity":   cty.NumberIntVal(3),
		"cluster_max_size":           cty.NumberIntVal(9),
	}

	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	err := r.Init("eks", "kubestack.example.com", "test", "eu-west-1", []string{"apps", "ops"}, baseCfg, "latest", "", p)
	assert.Equal(t, nil, err, nil)

	fp := filepath.Join(p, "kubestack-starter-eks")
	assert.DirExists(t, fp, nil)
	assert.DirExists(t, filepath.Join(fp, ".git"), nil)
	assert.FileExists(t, filepath.Join(fp, "eks_test_eu-west-1_cluster.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "eks_test_eu-west-1_providers.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "versions.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "variables.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "config.auto.tfvars"), nil)
	assert.FileExists(t, filepath.Join(fp, "Dockerfile"), nil)
	assert.FileExists(t, filepath.Join(fp, "README.md"), nil)

	os.RemoveAll(p)
}

func TestRepoInitGKE(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	baseCfg := map[string]cty.Value{
		"name_prefix":                cty.StringVal("test"),
		"project_id":                 cty.StringVal("kubestack-testing"),
		"region":                     cty.StringVal("europe-west4"),
		"cluster_min_node_count":     cty.NumberIntVal(1),
		"cluster_initial_node_count": cty.NumberIntVal(1),
		"cluster_max_node_count":     cty.NumberIntVal(3),
		"cluster_node_locations":     cty.StringVal("europe-west4-a,europe-west4-b,europe-west4-c"),
		"cluster_machine_type":       cty.StringVal("e2-standard-8"),
		"cluster_min_master_version": cty.StringVal("1.20"),
	}

	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	err := r.Init("gke", "kubestack.example.com", "test", "europe-west4", []string{"apps", "ops"}, baseCfg, "latest", "", p)
	assert.Equal(t, nil, err, nil)

	fp := filepath.Join(p, "kubestack-starter-gke")
	assert.DirExists(t, fp, nil)
	assert.DirExists(t, filepath.Join(fp, ".git"), nil)
	assert.FileExists(t, filepath.Join(fp, "gke_test_europe-west4_cluster.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "gke_test_europe-west4_providers.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "versions.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "variables.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "config.auto.tfvars"), nil)
	assert.FileExists(t, filepath.Join(fp, "Dockerfile"), nil)
	assert.FileExists(t, filepath.Join(fp, "README.md"), nil)

	os.RemoveAll(p)
}

func TestRepoInitAKS(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	baseCfg := map[string]cty.Value{
		"name_prefix":                  cty.StringVal("test"),
		"resource_group":               cty.StringVal("kubestack-testing"),
		"default_node_pool_vm_size":    cty.StringVal("Standard_D4_v4"),
		"default_node_pool_min_count":  cty.NumberIntVal(3),
		"default_node_pool_node_count": cty.NumberIntVal(3),
		"default_node_pool_max_count":  cty.NumberIntVal(9),
		"availability_zones":           cty.StringVal("1,2,3"),
	}

	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	err := r.Init("aks", "kubestack.example.com", "test", "westeurope", []string{"apps", "ops"}, baseCfg, "latest", "", p)
	assert.Equal(t, nil, err, nil)

	fp := filepath.Join(p, "kubestack-starter-aks")
	assert.DirExists(t, fp, nil)
	assert.DirExists(t, filepath.Join(fp, ".git"), nil)
	assert.FileExists(t, filepath.Join(fp, "aks_test_westeurope_cluster.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "aks_test_westeurope_providers.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "versions.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "variables.tf"), nil)
	assert.FileExists(t, filepath.Join(fp, "config.auto.tfvars"), nil)
	assert.FileExists(t, filepath.Join(fp, "Dockerfile"), nil)
	assert.FileExists(t, filepath.Join(fp, "README.md"), nil)

	os.RemoveAll(p)
}

type MockDownloaderArchiveError struct{}

func (c MockDownloaderArchiveError) Download(url string) (resp *http.Response, err error) {
	return resp, errors.New("Mock HTTP error")
}

func TestRepoInitDownloadError(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderArchiveError{},
	}

	err := r.Init("aks", "kubestack.example.com", "test", "europe-west4", []string{"apps", "ops"}, map[string]cty.Value{}, "latest", "", "")

	assert.Error(t, err, nil)
}

func TestRepoInitNoSuchRelease(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}
	err := r.Init("no-such-starter", "kubestack.example.com", "test", "europe-west4", []string{"apps", "ops"}, map[string]cty.Value{}, "no-such-release", "", "")

	assert.EqualError(t, err, "'no-such-release' is not a valid version, try the latest version 'v0.18.0-beta.0'", nil)
}

func TestRepoInitNoSuchStarter(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	err := r.Init("no-such-starter", "kubestack.example.com", "test", "europe-west4", []string{"apps", "ops"}, map[string]cty.Value{}, "latest", "", "")

	assert.EqualError(t, err, "'no-such-starter' is not a valid starter name, choose one of [aks eks gke kind multi-cloud]", nil)
}

func TestRepoDownloadUrlGitRef(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	url, err := r.downloadUrl("test", "", "test")

	assert.Equal(t, "https://storage.googleapis.com/dev.quickstart.kubestack.com/kubestack-starter-test-test.zip", url, nil)
	assert.Equal(t, nil, err, nil)
}
