package cli

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/kbst/kbst/util"
)

type Repo struct {
	Framework  util.Entry
	Downloader util.Downloader
}

func (r Repo) Init(starter string, release string, gitRef string, path string) (err error) {
	// download archive
	url, err := r.downloadUrl(starter, release, gitRef)
	if err != nil {
		return err
	}

	resp, err := r.Downloader.Download(url)
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

func (r Repo) downloadUrl(starter string, release string, gitRef string) (url string, err error) {
	if gitRef != "" {
		return fmt.Sprintf(
			"https://storage.googleapis.com/dev.quickstart.kubestack.com/kubestack-starter-%s-%s.zip",
			starter,
			gitRef,
		), nil
	}

	// determine version
	version, err := r.Framework.GetReleaseOrLatest(release)
	if err != nil {
		return url, err
	}

	url, ok := version.Archives[starter]
	if !ok {
		options := []string{}
		for k := range version.Archives {
			options = append(options, k)
		}
		sort.Strings(options)

		return url, fmt.Errorf(
			"'%s' is not a valid starter name, choose one of %v",
			starter,
			options,
		)
	}

	return url, nil
}
