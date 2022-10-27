package cli

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/kbst/kbst/pkg/util"
	"github.com/stretchr/testify/assert"
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
	p := filepath.Join(fixturesPath, "kubestack-starter-multi-cloud-v0.11.0-beta.0.zip")
	f, err := ioutil.ReadFile(p)
	if err != nil {
		return resp, err
	}

	r := &http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(f)),
	}
	return r, nil
}

func TestRepoInit(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	err := r.Init("multi-cloud", "kubestack.example.com", "latest", "", p)

	assert.Equal(t, nil, err, nil)
	assert.DirExists(t, filepath.Join(p, "kubestack-starter-multi-cloud"), nil)
	assert.DirExists(t, filepath.Join(p, "kubestack-starter-multi-cloud", ".git"), nil)

	os.RemoveAll(p)
}

func TestRepoInitGitRef(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	p, _ := ioutil.TempDir(os.TempDir(), "kbst-unit-test-*")
	err := r.Init("multi-cloud", "kubestack.example.com", "", "test", p)

	assert.Equal(t, nil, err, nil)
	assert.DirExists(t, filepath.Join(p, "kubestack-starter-multi-cloud"), nil)
	assert.DirExists(t, filepath.Join(p, "kubestack-starter-multi-cloud", ".git"), nil)

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

	err := r.Init("multi-cloud", "kubestack.example.com", "latest", "", "")

	assert.Error(t, err, nil)
}

func TestRepoInitNoSuchRelease(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}
	err := r.Init("no-such-starter", "kubestack.example.com", "no-such-release", "", "")

	assert.EqualError(t, err, "'no-such-release' is not a valid version, try the latest version 'v0.11.0-beta.0'", nil)
}

func TestRepoInitNoSuchStarter(t *testing.T) {
	cj := util.CliJSON{}
	cj.Load(MockDownloaderCliJson{})
	r := Repo{
		Framework:  cj.Framework,
		Downloader: MockDownloaderFrameworkArchive{},
	}

	err := r.Init("no-such-starter", "kubestack.example.com", "latest", "", "")

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
