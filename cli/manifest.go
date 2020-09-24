package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/kbst/kbst/util"
)

func ManifestInstall(entry string, variant string, overlay string, release string, gitRef string, path string) (err error) {
	url, err := getManifestDownloadUrl(entry, release, gitRef)
	if err != nil {
		return err
	}

	// download entry archive
	resp, err := util.CachedDownload(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	archive, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// extract into tempdir first
	tempEntry, err := ioutil.TempDir(os.TempDir(), "kbst-manifest-tmp-")
	defer os.RemoveAll(tempEntry)

	_, err = util.Unzip(archive, tempEntry)
	if err != nil {
		return err
	}

	// check variant is in entry
	_, err = getEntryVariants(tempEntry, entry, variant)
	if err != nil {
		return err
	}

	// now that we know entry and variant are ok
	// we extract into basesPath
	basesPath := filepath.Join(path, "manifests", "bases")

	_, err = util.Unzip(archive, basesPath)
	if err != nil {
		return err
	}

	// add kustomization resources
	overlayPath := filepath.Join(path, "manifests", "overlays", overlay)

	fSys := util.MakeRelFsOnDisk(overlayPath)
	mf, err := util.NewKustomizationFile(fSys)
	if err != nil {
		return err
	}

	m, err := mf.Read()
	if err != nil {
		return err
	}

	variantPath := filepath.Join(basesPath, entry, variant)
	resource, err := filepath.Rel(overlayPath, variantPath)
	if err != nil {
		return err
	}

	if util.StringInSlice(resource, m.Resources) {
		// what we're trying to add is already in the list
		return nil
	}
	m.Resources = append(m.Resources, resource)

	err = mf.Write(m)
	if err != nil {
		return err
	}

	return nil
}

func ManifestRemove(entry string, overlay string, path string) (err error) {
	basesPath := filepath.Join(path, "manifests", "bases")
	entryPath := filepath.Join(basesPath, entry)

	// remove kustomization resources
	overlaysPath := filepath.Join(path, "manifests", "overlays")
	overlayPath := filepath.Join(overlaysPath, overlay)

	fSys := util.MakeRelFsOnDisk(overlayPath)
	mf, err := util.NewKustomizationFile(fSys)
	if err != nil {
		return err
	}

	m, err := mf.Read()
	if err != nil {
		return err
	}

	relEntryPath, err := filepath.Rel(overlayPath, entryPath)
	if err != nil {
		return err
	}

	// just remove all references to the base we're deleting
	newResources := []string{}
	for _, resource := range m.Resources {
		if strings.HasPrefix(resource, relEntryPath) {
			continue
		}

		newResources = append(newResources, resource)
	}

	m.Resources = newResources

	err = mf.Write(m)
	if err != nil {
		return err
	}

	// delete the base
	err = os.RemoveAll(entryPath)
	if err != nil {
		return err
	}

	return nil
}

func getEntryVariants(path string, entry string, variant string) (variants []string, err error) {
	variants = []string{}
	found := false
	entryPath := filepath.Join(path, entry)
	entryPathEntries, err := ioutil.ReadDir(entryPath)
	if err != nil {
		return nil, err
	}
	for _, e := range entryPathEntries {
		if e.IsDir() {
			variants = append(variants, e.Name())

			if e.Name() == variant {
				found = true
			}
		}
	}

	if !found {
		return nil, fmt.Errorf(
			"'%s' is not a valid variant for '%s', choose one of %v",
			variant,
			entry,
			variants,
		)
	}

	return variants, nil
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
