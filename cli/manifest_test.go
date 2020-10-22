package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/kbst/kbst/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type MockDownloaderCatalogArchive struct{}

func (c MockDownloaderCatalogArchive) Download(url string) (resp *http.Response, err error) {
	p := filepath.Join(fixturesPath, "prometheus-v0.42.1-kbst.0.zip")
	f, err := ioutil.ReadFile(p)
	if err != nil {
		return resp, err
	}

	r := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(f)),
	}
	return r, nil
}

func TestManifestInstall(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog:    cj.Catalog,
		Downloader: MockDownloaderCatalogArchive{},
	}
	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	r.Init("multi-cloud", "latest", "", p)
	mp := filepath.Join(p, "kubestack-starter-multi-cloud")

	err := m.Install("prometheus", "clusterwide", "apps", "latest", "", mp, false)
	assert.Equal(t, nil, err, nil)
	assert.DirExists(t, filepath.Join(mp, "manifests", "bases", "prometheus"), nil)

	kF, _ := ioutil.ReadFile(filepath.Join(mp, "manifests", "overlays", "apps", "kustomization.yaml"))
	k := make(map[string]interface{})
	yaml.Unmarshal(kF, k)

	assert.Contains(t, k["resources"], "../../bases/prometheus/clusterwide", nil)

	os.RemoveAll(p)
}

func TestManifestInstallSkipEditKustomization(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog:    cj.Catalog,
		Downloader: MockDownloaderCatalogArchive{},
	}
	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	r.Init("multi-cloud", "latest", "", p)
	mp := filepath.Join(p, "kubestack-starter-multi-cloud")

	err := m.Install("prometheus", "clusterwide", "apps", "latest", "", mp, true)
	assert.Equal(t, nil, err, nil)
	assert.DirExists(t, filepath.Join(mp, "manifests", "bases", "prometheus"), nil)

	kF, _ := ioutil.ReadFile(filepath.Join(mp, "manifests", "overlays", "apps", "kustomization.yaml"))
	k := make(map[string]interface{})
	yaml.Unmarshal(kF, k)

	assert.NotContains(t, k["resources"], "../../bases/prometheus/clusterwide", nil)

	os.RemoveAll(p)
}

func TestManifestRemove(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog:    cj.Catalog,
		Downloader: MockDownloaderCatalogArchive{},
	}
	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	r.Init("multi-cloud", "latest", "", p)
	mp := filepath.Join(p, "kubestack-starter-multi-cloud")
	m.Install("prometheus", "clusterwide", "apps", "latest", "", mp, false)

	err := m.Remove("prometheus", mp, false)

	assert.Equal(t, nil, err, nil)
	bases, _ := ioutil.ReadDir(filepath.Join(mp, "manifests", "bases"))
	for _, b := range bases {
		assert.NotEqual(t, "prometheus", b.Name(), nil)
	}

	kF, _ := ioutil.ReadFile(filepath.Join(mp, "manifests", "overlays", "apps", "kustomization.yaml"))
	k := make(map[string]interface{})
	yaml.Unmarshal(kF, k)

	assert.NotContains(t, k["resources"], "../../bases/prometheus/clusterwide", nil)

	os.RemoveAll(p)
}

func TestManifestRemoveSkipEditKustomization(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog:    cj.Catalog,
		Downloader: MockDownloaderCatalogArchive{},
	}
	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	r.Init("multi-cloud", "latest", "", p)
	mp := filepath.Join(p, "kubestack-starter-multi-cloud")
	m.Install("prometheus", "clusterwide", "apps", "latest", "", mp, false)

	err := m.Remove("prometheus", mp, true)

	assert.Equal(t, nil, err, nil)
	bases, _ := ioutil.ReadDir(filepath.Join(mp, "manifests", "bases"))
	for _, b := range bases {
		assert.NotEqual(t, "prometheus", b.Name(), nil)
	}

	kF, _ := ioutil.ReadFile(filepath.Join(mp, "manifests", "overlays", "apps", "kustomization.yaml"))
	k := make(map[string]interface{})
	yaml.Unmarshal(kF, k)

	assert.Contains(t, k["resources"], "../../bases/prometheus/clusterwide", nil)

	os.RemoveAll(p)
}

func TestManifestUpdate(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog:    cj.Catalog,
		Downloader: MockDownloaderCatalogArchive{},
	}
	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	r.Init("multi-cloud", "latest", "", p)
	mp := filepath.Join(p, "kubestack-starter-multi-cloud")

	m.Install("prometheus", "clusterwide", "apps", "latest", "", mp, false)

	err := m.Update("prometheus", "apps", "latest", "", mp)

	assert.Equal(t, nil, err, nil)
	assert.DirExists(t, filepath.Join(mp, "manifests", "bases", "prometheus"), nil)

	kF, _ := ioutil.ReadFile(filepath.Join(mp, "manifests", "overlays", "apps", "kustomization.yaml"))
	k := make(map[string]interface{})
	yaml.Unmarshal(kF, k)

	assert.Contains(t, k["resources"], "../../bases/prometheus/clusterwide", nil)

	os.RemoveAll(p)
}

func TestGetEntryVariants(t *testing.T) {
	m := Manifest{}
	p := filepath.Join(fixturesPath)
	e := "catalog-entry"
	v := "variant1"
	variants, err := m.getEntryVariants(p, e, v)

	expVariants := []string{"variant1", "variant2"}
	assert.Equal(t, expVariants, variants, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetEntryVariantsPathError(t *testing.T) {
	m := Manifest{}
	p := filepath.Join(fixturesPath, "_may_not_exist")
	e := "catalog-entry"
	v := "variant1"
	variants, err := m.getEntryVariants(p, e, v)

	assert.Equal(t, []string(nil), variants, nil)
	assert.Error(t, err, nil)
}

func TestGetEntryVariantsNotFound(t *testing.T) {
	m := Manifest{}
	p := filepath.Join(fixturesPath)
	e := "catalog-entry"
	v := "does_not_exist"
	variants, err := m.getEntryVariants(p, e, v)

	assert.Equal(t, []string(nil), variants, nil)
	assert.Error(t, err, nil)
}

func TestGetOverlays(t *testing.T) {
	m := Manifest{}
	p := filepath.Join(fixturesPath, "overlays")
	overlays, err := m.getOverlays(p)

	expOverlays := []string{
		filepath.Join(p, "overlay1"),
		filepath.Join(p, "overlay2"),
	}
	notExpContains := filepath.Join(p, "missing-kustomization")
	assert.NotContains(t, overlays, notExpContains, nil)
	assert.Equal(t, expOverlays, overlays, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetOverlaysPathError(t *testing.T) {
	m := Manifest{}
	p := filepath.Join(fixturesPath, "_may_not_exist")
	overlays, err := m.getOverlays(p)

	assert.Equal(t, []string(nil), overlays, nil)
	assert.Error(t, err, nil)
}

func TestGetManifestDownloadUrl(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog: cj.Catalog,
	}

	url, err := m.getManifestDownloadUrl("prometheus", "latest", "")

	assert.Equal(t, "https://storage.googleapis.com/catalog.kubestack.com/prometheus-v0.42.1-kbst.0.zip", url, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetManifestDownloadUrlRelease(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog: cj.Catalog,
	}

	url, err := m.getManifestDownloadUrl("prometheus", "v0.29.0-kbst.0", "")

	assert.Equal(t, "https://storage.googleapis.com/catalog.kubestack.com/prometheus-v0.29.0-kbst.0.zip", url, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetManifestDownloadUrlGitRef(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog: cj.Catalog,
	}
	entry := "test"
	gitRef := "master-4ce4e92"
	url, err := m.getManifestDownloadUrl(entry, "", gitRef)

	expUrl := fmt.Sprintf(
		"https://storage.googleapis.com/dev.catalog.kubestack.com/%s-%s.zip",
		entry,
		gitRef,
	)

	assert.Equal(t, expUrl, url, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetManifestDownloadNoSuchEntry(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog: cj.Catalog,
	}
	url, err := m.getManifestDownloadUrl("no-such-entry", "", "")

	assert.Equal(t, "", url, nil)
	assert.EqualError(t, err, "'no-such-entry' is not a valid entry name, choose one of [argo-cd cert-manager flux nginx postgresql prometheus sealed-secrets tektoncd]", nil)
}

func TestGetManifestDownloadNoSuchRelease(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	m := Manifest{
		Catalog: cj.Catalog,
	}
	url, err := m.getManifestDownloadUrl("prometheus", "no-such-release", "")

	assert.Equal(t, "", url, nil)
	assert.EqualError(t, err, "'no-such-release' is not a valid version, try the latest version 'v0.42.1-kbst.0'", nil)
}
