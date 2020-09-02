package cli

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/docopt/docopt-go"
	"github.com/kbst/kbst/util"
)

func Repository(argv []string) (err error) {
	usage := `
Usage:
  kbst repository init <starter> [--release=release] [--path=path]

Options:
  -r. --release=release  Release version to use [default: latest].
  -p, --path=path  		 Path to initialize the repository in [default: .].
  -h, --help	   		 Show this help.
`

	args, _ := docopt.ParseDoc(usage)
	fmt.Println(args)

	if args["init"] == true {
		starter := args["<starter>"].(string)
		release := args["--release"].(string)
		path := args["--path"].(string)
		return repoInit(starter, release, path)
	}

	return
}

func repoInit(starter string, release string, path string) (err error) {
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
