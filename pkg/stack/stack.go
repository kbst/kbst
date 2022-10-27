package stack

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Stack struct {
	root         *tfhcl.Root
	BaseDomain   string
	Environments []Environment
	Clusters     []Cluster
	NodePools    []NodePool
	Services     []Service
}

type Environment struct {
	Key       string
	IsBaseKey bool
}

type Configuration struct {
	EnvironmentKey string
	Attributes     map[string]cty.Value
}

func NewStack(r *tfhcl.Root) *Stack {
	s := &Stack{
		root: r,
	}

	return s
}

func (s *Stack) FromPath(p string) error {
	err := s.root.Read(p)
	if err != nil {
		return err
	}

	// TODO: parse TF variables and read base domain from vars
	s.BaseDomain = "kubestack.example.com"

	for _, mf := range s.root.Modules {
		for _, m := range mf {
			prefix, region := parsePrefixRegion(m.Name)
			kind, provider, version := parseKindProviderVersion(m.Source, m.Version)
			cbk := m.ConfigurationBaseKey
			if cbk == "" {
				cbk = "apps"
			}

			switch kind {
			case "cluster":
				c := Cluster{
					NamePrefix: prefix,
					Provider:   provider,
					Region:     region,
					Version:    version,
				}

				keys := maps.Keys(m.Configuration)
				sort.Strings(keys)

				for _, ek := range keys {
					isBk := false
					if ek == cbk {
						isBk = true
					}

					env := Environment{
						Key:       ek,
						IsBaseKey: isBk,
					}

					if !slices.Contains(s.Environments, env) {
						if env.IsBaseKey {
							// base environment always comes first
							s.Environments = append([]Environment{env}, s.Environments...)
							continue
						}
						s.Environments = append(s.Environments, env)
					}
				}

				c.Configurations = parseConfiguration(m.ConfigurationBaseKey, m.Configuration)

				s.Clusters = append(s.Clusters, c)
			case "node_pool":
				_, region := parsePrefixRegion(m.Name)
				_, provider, version := parseKindProviderVersion(m.Source, m.Version)
				cbk := m.ConfigurationBaseKey
				if cbk == "" {
					cbk = "apps"
				}

				clusterName, nameSuffix := parseNodePoolClusteNameNameSuffix(m.Name)
				np := NodePool{
					NameSuffix:  nameSuffix,
					ClusterName: clusterName,
					Provider:    provider,
					Region:      region,
					Version:     version,
				}

				np.Configurations = parseConfiguration(m.ConfigurationBaseKey, m.Configuration)

				s.NodePools = append(s.NodePools, np)
			case "service":
				clusterName, entryName := parseServiceClusteNameEntryName(m.Name)
				svc := Service{
					EntryName:   entryName,
					ClusterName: clusterName,
					Provider:    "kustomization",
					Version:     m.Version,
				}

				svc.Configurations = parseConfiguration(m.ConfigurationBaseKey, m.Configuration)

				s.Services = append(s.Services, svc)
			default:
				return fmt.Errorf("error loading stack: %q is not a valid kind", kind)
			}
		}
	}

	return nil
}

func (s *Stack) Files() (map[string]*hclwrite.File, error) {
	files := make(map[string]*hclwrite.File)

	p := make(map[string]map[string]string)
	p["kustomization"] = make(map[string]string)
	p["kustomization"]["source"] = "kbst/kustomization"
	fver := hclwrite.NewEmptyFile()
	tfhcl.BlockTerraform(fver, p)
	files["versions.tf"] = fver

	fvar := hclwrite.NewEmptyFile()
	tfhcl.BlockVariable(fvar, "base_domain", "string", "Used to generate fully qualified domain names for all clusters.")
	files["variables.tf"] = fvar

	ftfvars := hclwrite.NewEmptyFile()
	ftfvars.Body().SetAttributeValue("base_domain", cty.StringVal(s.BaseDomain))
	files["config.auto.tfvars"] = ftfvars

	for _, c := range s.Clusters {
		for k, v := range c.ToHCL() {
			files[k] = v
			//fmt.Printf("%s: %s\n", k, v.Bytes())
		}
	}

	for _, np := range s.NodePools {
		for k, v := range np.ToHCL() {
			files[k] = v
			//fmt.Printf("%s: %s\n", k, v.Bytes())
		}
	}

	for _, svc := range s.Services {
		for k, v := range svc.ToHCL() {
			files[k] = v
			//fmt.Printf("%s: %s\n", k, v.Bytes())
		}
	}

	return files, nil
}

func (s *Stack) AddCluster(namePrefix, provider, region, version string, configurations []Configuration) {
	s.Clusters = append(s.Clusters, Cluster{
		NamePrefix:     namePrefix,
		Provider:       provider,
		Region:         region,
		Version:        version,
		Configurations: configurations,
	})
}
