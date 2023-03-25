package generator

import (
	"encoding/json"
	"fmt"

	"github.com/kbst/kbst/pkg/stack"
	"github.com/zclconf/go-cty/cty"
)

type LegacyEnvironment struct {
	IsConfigurationBaseKey bool   `json:"is_configuration_base_key"`
	Name                   string `json:"name"`
	Key                    string `json:"key"`
}

type LegacyStack struct {
	BaseDomain      string              `json:"base_domain"`
	BaseEnvironment string              `json:"base_environment"`
	Environments    []LegacyEnvironment `json:"environments"`
	Modules         []LegacyModule      `json:"modules"`
	stack           *stack.Stack
}

func (ls *LegacyStack) Unmarshal(d []byte) (s *stack.Stack, err error) {
	err = json.Unmarshal(d, &ls)
	if err != nil {
		return s, err
	}

	s = &stack.Stack{}
	s.BaseDomain = ls.BaseDomain

	cbk := ls.BaseEnvironment

	for _, le := range ls.Environments {
		s.Environments = append(s.Environments, stack.Environment{
			Key:       le.Key,
			IsBaseKey: le.IsConfigurationBaseKey,
		})
	}

	for _, lm := range ls.Modules {
		name_prefix, region := getNamePrefixRegion(cbk, lm.Configurations)

		c := stack.Cluster{
			NamePrefix:     name_prefix,
			Provider:       lm.Provider,
			Region:         region,
			Version:        lm.Version,
			Configurations: convertLegacyConfigurations(lm.Configurations),
		}

		s.Clusters = append(s.Clusters, c)

		for _, lcm := range lm.Children {
			switch lcm.Type {
			case "node_pool":
				np := stack.NodePool{
					PoolName:       getNodePoolSuffix(cbk, lcm.Configurations),
					ClusterName:    c.Name(),
					Provider:       lcm.Provider,
					Region:         region,
					Version:        lcm.Version,
					Configurations: convertLegacyConfigurations(lcm.Configurations),
				}

				s.NodePools = append(s.NodePools, np)
			case "service":
				svc := stack.Service{
					EntryName:      lcm.Name,
					ClusterName:    c.Name(),
					Provider:       lcm.Provider,
					Version:        lcm.Version,
					Configurations: convertLegacyConfigurations(lcm.Configurations),
				}

				s.Services = append(s.Services, svc)
			default:
				return s, fmt.Errorf("invalid module type: %s: %q", lcm.Name, lcm.Type)
			}
		}
	}

	return s, err
}

func getNamePrefixRegion(cbk string, in []LegacyConfiguration) (name_prefix, region string) {
	for _, cfg := range in {
		if cfg.EnvKey != cbk {
			continue
		}

		for k, v := range cfg.Data {
			if k == "name_prefix" {
				name_prefix = v.(string)
			}

			if k == "region" {
				region = v.(string)
			}
		}
	}

	return name_prefix, region
}

func getNodePoolSuffix(cbk string, in []LegacyConfiguration) (suffix string) {
	for _, cfg := range in {
		if cfg.EnvKey != cbk {
			continue
		}

		for k, v := range cfg.Data {
			if k == "name" {
				suffix = v.(string)
			}

			if k == "node_pool_name" {
				suffix = v.(string)
			}
		}
	}

	return suffix
}

func convertLegacyConfigurations(in []LegacyConfiguration) (out []stack.Configuration) {
	for _, cfg := range in {
		attrs := make(map[string]cty.Value)
		for k, v := range cfg.Data {
			switch t := v.(type) {
			case nil:
				continue
			case string:
				attrs[k] = cty.StringVal(v.(string))
			case float64:
				attrs[k] = cty.NumberFloatVal(v.(float64))
			default:
				fmt.Printf("%s: %s: %+v\n", t, k, v)
			}
		}
		out = append(out, stack.Configuration{
			EnvironmentKey: cfg.EnvKey,
			Attributes:     attrs,
		})
	}

	return out
}

func (lm *LegacyModule) GetK8sServiceName() string {
	k8sServiceName := map[string]string{
		"aws":     "eks",
		"azurerm": "aks",
		"google":  "gke",
	}

	return k8sServiceName[lm.Provider]
}

type LegacyModule struct {
	Name           string                `json:"name"`
	Provider       string                `json:"provider"`
	Type           string                `json:"type"`
	Version        string                `json:"version"`
	Children       []LegacyModule        `json:"children"`
	Configurations []LegacyConfiguration `json:"configurations"`
	Configuration  map[string]map[string]interface{}
}

type LegacyConfiguration struct {
	Data   map[string]interface{} `json:"data"`
	EnvKey string                 `json:"env_key"`
}
