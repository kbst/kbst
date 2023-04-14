package export

import (
	"github.com/kbst/kbst/pkg/stack"
	"github.com/zclconf/go-cty/cty"
	ctyJson "github.com/zclconf/go-cty/cty/json"
)

type Stack struct {
	BaseDomain   string        `json:"base_domain"`
	Environments []Environment `json:"environments"`
	Clusters     []Cluster     `json:"clusters"`
	NodePools    []NodePool    `json:"node_pools"`
	Services     []Service     `json:"services"`
}

type Environment struct {
	Key       string `json:"key"`
	IsBaseKey bool   `json:"is_base_key"`
}

type Configuration struct {
	EnvironmentKey string                             `json:"environment_key"`
	Attributes     map[string]ctyJson.SimpleJSONValue `json:"attributes"`
}

type Cluster struct {
	NamePrefix     string          `json:"name_prefix"`
	Provider       string          `json:"provider"`
	Region         string          `json:"region"`
	Version        string          `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

type NodePool struct {
	PoolName       string          `json:"pool_name"`
	ClusterName    string          `json:"cluster_name"`
	Provider       string          `json:"provider"`
	Region         string          `json:"region"`
	Version        string          `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

type Service struct {
	EntryName      string          `json:"entry_name"`
	ClusterName    string          `json:"cluster_name"`
	Provider       string          `json:"provider"`
	Version        string          `json:"version"`
	Configurations []Configuration `json:"configurations"`
}

func (s *Stack) StackEnvironments() []stack.Environment {
	envs := []stack.Environment{}
	for _, e := range s.Environments {
		envs = append(envs, stack.Environment{
			Key:       e.Key,
			IsBaseKey: e.IsBaseKey,
		})
	}

	return envs
}

func (s *Stack) StackConfigurations(in []Configuration) []stack.Configuration {
	cfgs := []stack.Configuration{}
	for _, cfg := range in {
		sAttrs := map[string]cty.Value{}

		for k, v := range cfg.Attributes {
			sAttrs[k] = v.Value
		}

		cfgs = append(cfgs, stack.Configuration{
			EnvironmentKey: cfg.EnvironmentKey,
			Attributes:     sAttrs,
		})
	}

	return cfgs
}
