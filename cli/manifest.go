package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kbst/kbst/util"
)

type Manifest struct {
	Catalog    map[string]util.Entry
	Downloader util.Downloader
}

func (m Manifest) Install(entry string, variant string, overlay string, release string, gitRef string, path string, skipEditKustomization bool) (err error) {
	url, err := m.getManifestDownloadUrl(entry, release, gitRef)
	if err != nil {
		return err
	}

	// download entry archive
	resp, err := m.Downloader.Download(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	archive, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if skipEditKustomization == false {
		// extract into tempdir first
		tempEntry, err := ioutil.TempDir(os.TempDir(), "kbst-manifest-tmp-")
		defer os.RemoveAll(tempEntry)

		_, err = util.Unzip(archive, tempEntry)
		if err != nil {
			return err
		}

		// check variant is in entry
		_, err = m.getEntryVariants(tempEntry, entry, variant)
		if err != nil {
			return err
		}
	}

	// now that we know entry and variant are ok
	// we extract into basesPath
	basesPath := filepath.Join(path, "manifests", "bases")

	_, err = util.Unzip(archive, basesPath)
	if err != nil {
		return err
	}

	if skipEditKustomization == false {
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
	}

	return nil
}

func (m Manifest) Remove(entry string, path string, skipEditKustomization bool) (err error) {
	entryPath := filepath.Join(path, "manifests", "bases", entry)

	if skipEditKustomization == false {
		overlaysPath := filepath.Join(path, "manifests", "overlays")
		overlays, err := m.getOverlays(overlaysPath)
		if err != nil {
			return err
		}

		// loop through all defined overlays
		for _, overlayPath := range overlays {
			// remove kustomization resources
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
			changed := false
			newResources := []string{}
			for _, resource := range m.Resources {
				if strings.HasPrefix(resource, relEntryPath) {
					changed = true
					continue
				}

				newResources = append(newResources, resource)
			}

			// only write the file if we changed something
			if changed {
				m.Resources = newResources

				err = mf.Write(m)
				if err != nil {
					return err
				}
			}
		}
	}

	// delete the base
	err = os.RemoveAll(entryPath)
	if err != nil {
		return err
	}

	return nil
}

func (m Manifest) Update(entry string, overlay string, release string, gitRef string, path string) (err error) {
	err = m.Remove(entry, path, true)
	if err != nil {
		return err
	}

	variant := ""
	err = m.Install(entry, variant, overlay, release, gitRef, path, true)
	if err != nil {
		return err
	}

	return nil
}

func (m Manifest) getEntryVariants(path string, entry string, variant string) (variants []string, err error) {
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

func (m Manifest) getOverlays(path string) (overlays []string, err error) {
	overlayPathEntries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	potentialOverlays := []string{}
	for _, e := range overlayPathEntries {
		if e.IsDir() {
			potentialOverlays = append(potentialOverlays, e.Name())
		}
	}

	for _, p := range potentialOverlays {
		op := filepath.Join(path, p)
		fSys := util.MakeRelFsOnDisk(op)
		_, err := util.NewKustomizationFile(fSys)
		if err != nil {
			continue
		}
		overlays = append(overlays, op)
	}

	return overlays, nil
}

func (m Manifest) getManifestDownloadUrl(entry string, release string, gitRef string) (url string, err error) {
	if gitRef != "" {
		return fmt.Sprintf(
			"https://storage.googleapis.com/dev.catalog.kubestack.com/%s-%s.zip",
			entry,
			gitRef,
		), nil
	}

	// determine version

	current, ok := m.Catalog[entry]
	if !ok {
		options := []string{}
		for k := range m.Catalog {
			options = append(options, k)
		}
		sort.Strings(options)

		return url, fmt.Errorf(
			"'%s' is not a valid entry name, choose one of %v",
			entry,
			options,
		)
	}

	version, err := current.GetReleaseOrLatest(release)
	if err != nil {
		return url, err
	}

	return version.Archive, nil
}
