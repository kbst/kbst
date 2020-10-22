package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testEntry = Entry{
	Name: "test",
	Versions: []Version{
		{Name: "0.0.2"},
		{Name: "0.0.1"},
		{Name: "0.0.0"},
	},
}

func TestGetReleaseOrLatestEmptyVersions(t *testing.T) {
	emptyEntry := Entry{Name: "testEmpty", Versions: []Version{}}
	v, err := emptyEntry.GetReleaseOrLatest("latest")
	assert.Equal(t, Version{}, v, nil)
	assert.Equal(t, fmt.Errorf("No versions for '%s'", emptyEntry.Name), err, nil)
}

func TestGetReleaseOrLatestVersionLatest(t *testing.T) {
	v, err := testEntry.GetReleaseOrLatest("latest")
	assert.Equal(t, testEntry.Versions[0], v, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetReleaseOrLatestVersionSpecific(t *testing.T) {
	v, err := testEntry.GetReleaseOrLatest("0.0.1")
	assert.Equal(t, testEntry.Versions[1], v, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetReleaseOrLatestVersionSpecificMissing(t *testing.T) {
	r := "0.0.4"
	v, err := testEntry.GetReleaseOrLatest(r)
	assert.Equal(t, Version{}, v, nil)
	assert.Equal(t, fmt.Errorf(
		"'%s' is not a valid version, try the latest version '%s'",
		r,
		testEntry.Versions[0].Name,
	), err, nil)
}

type MockDownloader struct{}

func (c MockDownloader) Download(url string) (resp *http.Response, err error) {
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

func TestCliJSON(t *testing.T) {
	cj := CliJSON{}
	err := cj.Load(MockDownloader{})

	assert.IsType(t, map[string]Entry{}, cj.Catalog, nil)
	assert.IsType(t, Entry{}, cj.Framework, nil)
	assert.IsType(t, Entry{}, cj.Cli, nil)
	assert.Equal(t, nil, err, nil)
}

type MockDownloaderError struct{}

func (c MockDownloaderError) Download(url string) (resp *http.Response, err error) {
	return resp, errors.New("Mock HTTP error")
}

func TestCliJSONError(t *testing.T) {
	cj := CliJSON{}
	err := cj.Load(MockDownloaderError{})

	assert.IsType(t, map[string]Entry{}, cj.Catalog, nil)
	assert.IsType(t, Entry{}, cj.Framework, nil)
	assert.IsType(t, Entry{}, cj.Cli, nil)
	assert.Error(t, err, nil)
}
