package cli

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kbst/kbst/util"
)

func RepoInit(starter string, release string, path string) (err error) {
	framework, err := util.GetFramework()
	if err != nil {
		return err
	}

	initVersion := framework.Versions[0]
	if release != "latest" {
		for i := range framework.Versions {
			v := framework.Versions[i]
			if v.Name == release {
				initVersion = v
				break
			}
		}

		if initVersion.Name != release {
			return fmt.Errorf("'%s' is not a valid version", release)
		}
	}

	url, ok := initVersion.Archives[starter]
	if !ok {
		return fmt.Errorf("'%s' is not a valid starter", starter)
	}

	resp, err := util.CachedDownload(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	archive, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	filenames, err := util.Unzip(archive, path)
	if err != nil {
		return err
	}

	for _, name := range filenames {
		log.Println(name)
	}

	return
}
