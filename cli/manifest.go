package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/kbst/kbst/util"
)

func ManifestInstall(entry string, variant string, overlay string, release string, gitRef string, path string) (err error) {
	url, err := getManifestDownloadUrl(entry, release, gitRef)
	if err != nil {
		return err
	}

	resp, err := util.CachedDownload(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// extract archive
	archive, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	basesPath := filepath.Join(path, "manifests", "bases")
	absBasesPath, err := filepath.Abs(basesPath)
	if err != nil {
		return err
	}

	// extract into tempdir first
	tempEntry, err := ioutil.TempDir(os.TempDir(), "kbst-manifest-tmp-")
	defer os.RemoveAll(tempEntry)

	absTempEntry, err := filepath.Abs(tempEntry)
	if err != nil {
		return err
	}

	_, err = util.Unzip(archive, absTempEntry)
	if err != nil {
		return err
	}

	// determine available variants
	tempEntryPath := filepath.Join(absTempEntry, entry)
	entryVariants := []string{}
	entryVariantFound := false
	entryPathEntries, err := ioutil.ReadDir(tempEntryPath)
	if err != nil {
		return err
	}
	for _, e := range entryPathEntries {
		if e.IsDir() {
			entryVariants = append(entryVariants, e.Name())

			if e.Name() == variant {
				entryVariantFound = true
			}
		}
	}

	if !entryVariantFound {
		return fmt.Errorf(
			"'%s' is not a valid variant for '%s', choose one of %v",
			variant,
			entry,
			entryVariants,
		)
	}

	// now that we know entry and variant are ok
	// we extract into basesPath
	_, err = util.Unzip(archive, absBasesPath)
	if err != nil {
		return err
	}

	// add kustomization resources
	overlaysPath := filepath.Join(path, "manifests", "overlays")

	if overlay == "" {
		overlay = "apps"
	}
	overlayPath := filepath.Join(overlaysPath, overlay)
	absOverlayPath, err := filepath.Abs(overlayPath)
	if err != nil {
		return err
	}

	absVariantPath := filepath.Join(absBasesPath, entry, variant)
	relOverlayPath, err := filepath.Rel(absOverlayPath, absVariantPath)
	if err != nil {
		return err
	}

	resources := []string{relOverlayPath}

	fSys := util.MakeRelFsOnDisk(absOverlayPath)
	mf, err := util.NewKustomizationFile(fSys)
	if err != nil {
		return err
	}

	m, err := mf.Read()
	if err != nil {
		return err
	}

	for _, resource := range resources {
		if util.StringInSlice(resource, m.Resources) {
			log.Printf("resource %s already in kustomization file", resource)
			continue
		}
		m.Resources = append(m.Resources, resource)
	}

	mf.Write(m)

	return nil
}

func getManifestDownloadUrl(entry string, release string, gitRef string) (url string, err error) {
	if gitRef != "" {
		return fmt.Sprintf(
			"https://storage.googleapis.com/dev.catalog.kubestack.com/%s-%s.zip",
			entry,
			release,
		), nil
	}

	// determine version
	catalog, err := util.GetCatalog()
	if err != nil {
		return url, err
	}

	current, ok := catalog[entry]
	if !ok {
		return url, fmt.Errorf(
			"'%s' is not a valid entry name, choose one of %v",
			entry,
			reflect.ValueOf(catalog).MapKeys(),
		)
	}

	version, err := current.GetReleaseOrLatest(release)
	if err != nil {
		return url, err
	}

	return version.Archive, nil
}
