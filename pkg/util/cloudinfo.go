package util

import (
	_ "embed"
	"encoding/json"
	"sort"

	"golang.org/x/exp/maps"
)

//go:embed cloudinfo.json
var cloudinfo string

type CloudInfo struct {
	data map[string]Provider
}

type Provider map[string]Region

type Region map[string]Instance

type Instance struct {
	Name   string   `json:"name"`
	Family string   `json:"family"`
	Zones  []string `json:"zones"`
}

func (ci *CloudInfo) Load(d Downloader) (err error) {
	json.Unmarshal([]byte(cloudinfo), &ci.data)

	return nil
}

func (ci *CloudInfo) Providers() []string {
	keys := maps.Keys(ci.data)
	sort.Strings(keys)

	return keys
}

func (ci *CloudInfo) Regions(p string) []string {
	keys := maps.Keys(ci.data[p])
	sort.Strings(keys)

	return keys
}

func (ci *CloudInfo) Instances(p, r string) []string {
	keys := maps.Keys(ci.data[p][r])
	sort.Strings(keys)

	return keys
}

func (ci *CloudInfo) Zones(p, r, i string) []string {
	zones := ci.data[p][r][i].Zones
	sort.Strings(zones)

	return zones
}
