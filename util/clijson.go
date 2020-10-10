package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Version struct {
	Name     string            `json:"name"`
	Archive  string            `json:"archive,omitempty"`
	Archives map[string]string `json:"archives,omitempty"`
}

type Entry struct {
	Name     string    `json:"name"`
	Versions []Version `json:"versions"`
}

func (e Entry) GetReleaseOrLatest(r string) (v Version, err error) {
	if len(e.Versions) < 1 {
		return v, fmt.Errorf("No versions for '%s'", e.Name)
	}

	v = e.Versions[0]
	if r != "latest" {
		var cv Version
		for i := range e.Versions {
			cv = e.Versions[i]
			if cv.Name == r {
				v = cv
				break
			}
		}

		if cv.Name != r {
			return Version{}, fmt.Errorf(
				"'%s' is not a valid version, try the latest version '%s'",
				r,
				v.Name,
			)
		}
	}

	return v, nil
}

type CliJSON struct {
	Catalog   map[string]Entry `json:"catalog"`
	Framework Entry            `json:"framework"`
	Cli       Entry            `json:"cli"`
}

func GetCatalog(d Downloader) (catalog map[string]Entry, err error) {
	cliJson, err := getCliJson(d)
	if err != nil {
		return catalog, err
	}
	return cliJson.Catalog, nil
}

func GetFramework(d Downloader) (framework Entry, err error) {
	cliJson, err := getCliJson(d)
	if err != nil {
		return framework, err
	}
	return cliJson.Framework, nil
}

func GetCli(d Downloader) (cli Entry, err error) {
	cliJson, err := getCliJson(d)
	if err != nil {
		return cli, err
	}
	return cliJson.Cli, nil
}

func getCliJson(d Downloader) (cliJson CliJSON, err error) {
	resp, err := d.Download("https://www.kubestack.com/cli.json")
	if err != nil {
		return cliJson, err
	}
	defer resp.Body.Close()

	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cliJson, err
	}

	json.Unmarshal([]byte(respJson), &cliJson)

	return cliJson, nil
}
