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
	v = e.Versions[0]
	if r != "latest" {
		for i := range e.Versions {
			cv := e.Versions[i]
			if cv.Name == r {
				v = cv
				break
			}
		}

		if e.Name != r {
			return v, fmt.Errorf(
				"'%s' is not a valid version, try the latest version '%s'",
				r,
				v.Name,
			)
		}
	}

	return v, nil
}

type Catalog map[string]Entry
type Framework Entry
type Cli Entry

type CliJSON struct {
	Catalog   Catalog   `json:"catalog"`
	Framework Framework `json:"framework"`
	Cli       Cli       `json:"cli"`
}

func GetCatalog() (catalog Catalog, err error) {
	cliJson, err := getCliJson()
	if err != nil {
		return catalog, err
	}
	return cliJson.Catalog, nil
}

func GetFramework() (framework Framework, err error) {
	cliJson, err := getCliJson()
	if err != nil {
		return framework, err
	}
	return cliJson.Framework, nil
}

func GetCli() (cli Cli, err error) {
	cliJson, err := getCliJson()
	if err != nil {
		return cli, err
	}
	return cliJson.Cli, nil
}

func getCliJson() (cliJson CliJSON, err error) {
	resp, err := CachedDownload("https://www.kubestack.com/cli.json")
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
