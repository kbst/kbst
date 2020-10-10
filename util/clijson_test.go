package util

import (
	"errors"
	"fmt"
	"net/http"
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
	r := &http.Response{Body: http.NoBody}
	return r, nil
}

func TestGetCatalog(t *testing.T) {
	c, err := GetCatalog(MockDownloader{})
	assert.IsType(t, map[string]Entry{}, c, nil)
	assert.Equal(t, nil, err, nil)
}

func TestFramework(t *testing.T) {
	c, err := GetFramework(MockDownloader{})
	assert.IsType(t, Entry{}, c, nil)
	assert.Equal(t, nil, err, nil)
}

func TestGetCli(t *testing.T) {
	c, err := GetCli(MockDownloader{})
	assert.IsType(t, Entry{}, c, nil)
	assert.Equal(t, nil, err, nil)
}

type MockDownloaderError struct{}

func (c MockDownloaderError) Download(url string) (resp *http.Response, err error) {
	return resp, errors.New("Mock HTTP error")
}

func TestGetCatalogError(t *testing.T) {
	c, err := GetCatalog(MockDownloaderError{})
	assert.IsType(t, map[string]Entry{}, c, nil)
	assert.Equal(t, fmt.Errorf("Mock HTTP error"), err, nil)
}

func TestGetFrameworkError(t *testing.T) {
	c, err := GetFramework(MockDownloaderError{})
	assert.IsType(t, Entry{}, c, nil)
	assert.Equal(t, fmt.Errorf("Mock HTTP error"), err, nil)
}

func TestGetCliError(t *testing.T) {
	c, err := GetCli(MockDownloaderError{})
	assert.IsType(t, Entry{}, c, nil)
	assert.Equal(t, fmt.Errorf("Mock HTTP error"), err, nil)
}
