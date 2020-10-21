package cli

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEntryVariants(t *testing.T) {
	p := filepath.Join(fixturesPath)
	e := "catalog-entry"
	v := "variant1"
	variants, err := getEntryVariants(p, e, v)

	expVariants := []string{"variant1", "variant2"}
	assert.Equal(t, expVariants, variants, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetEntryVariantsPathError(t *testing.T) {
	p := filepath.Join(fixturesPath, "_may_not_exist")
	e := "catalog-entry"
	v := "variant1"
	variants, err := getEntryVariants(p, e, v)

	assert.Equal(t, []string(nil), variants, nil)
	assert.Error(t, err, nil)
}

func TestGetEntryVariantsNotFound(t *testing.T) {
	p := filepath.Join(fixturesPath)
	e := "catalog-entry"
	v := "does_not_exist"
	variants, err := getEntryVariants(p, e, v)

	assert.Equal(t, []string(nil), variants, nil)
	assert.Error(t, err, nil)
}

func TestGetOverlays(t *testing.T) {
	p := filepath.Join(fixturesPath, "overlays")
	overlays, err := getOverlays(p)

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
	p := filepath.Join(fixturesPath, "_may_not_exist")
	overlays, err := getOverlays(p)

	assert.Equal(t, []string(nil), overlays, nil)
	assert.Error(t, err, nil)
}

func TestGetManifestDownloadUrlGitRef(t *testing.T) {
	entry := "test"
	gitRef := "master-4ce4e92"
	url, err := getManifestDownloadUrl(entry, "", gitRef)

	expUrl := fmt.Sprintf(
		"https://storage.googleapis.com/dev.catalog.kubestack.com/%s-%s.zip",
		entry,
		gitRef,
	)

	assert.Equal(t, expUrl, url, nil)
	assert.Equal(t, nil, err, nil)
}
