package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/kbst/kbst/util"
)

func RepoInit(starter string, release string, path string) (err error) {
	framework, err := util.GetFramework()
	if err != nil {
		return err
	}

	// determine version
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

	// download archive
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

	filenames, err := util.Unzip(archive, path)
	if err != nil {
		return err
	}

	for _, name := range filenames {
		log.Println(name)
	}

	// initialize git repository
	repoPath, err := filepath.Abs(filenames[0])
	if err != nil {
		return err
	}

	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		return err
	}

	// make initial commit
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	copts := &git.CommitOptions{All: true}
	msg := fmt.Sprintf("Initialized from %s starter", strings.ToUpper(starter))

	_, err = wt.Add(".")
	if err != nil {
		return err
	}

	_, err = wt.Commit(msg, copts)
	if err != nil {
		return err
	}

	return
}
