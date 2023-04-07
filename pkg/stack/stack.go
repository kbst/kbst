package stack

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Stack struct {
	toDelete     []string
	toWrite      []string
	root         *tfhcl.Root
	cliJSON      util.CliJSON
	baseDomain   string
	Environments []Environment
}

type Environment struct {
	Key       string
	IsBaseKey bool
}

type Configuration struct {
	EnvironmentKey string
	Attributes     map[string]cty.Value
}

func NewStack(r *tfhcl.Root, cj util.CliJSON) *Stack {
	s := &Stack{
		root:    r,
		cliJSON: cj,
	}

	return s
}

func (s *Stack) FromPath() error {
	err := s.root.Read()
	if err != nil {
		return err
	}

	// read base_domain
	bd, ok := s.root.GetVariableValue("base_domain")
	if !ok {
		return fmt.Errorf("value for required var %q not found", "base_domain")
	}
	s.SetBaseDomain(bd)

	// read environments
	for _, mods := range s.root.Modules {
		for _, m := range mods {
			kind, _, _, err := m.TypeProviderVersion()
			if err != nil {
				continue
			}

			if kind != "cluster" {
				continue
			}

			keys := maps.Keys(m.Configuration)
			sort.Strings(keys)

			for _, ek := range keys {
				isBk := false
				if ek == m.ConfigurationBaseKeyOrDefault() {
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
		}
	}

	return nil
}

func (s *Stack) SetBaseDomain(bd cty.Value) {
	s.baseDomain = bd.AsString()
	s.root.SetVariableValue("base_domain", bd)
}

func (s *Stack) Clusters() (clusters []Cluster) {
	for _, mods := range s.root.Modules {
		for i := range mods {
			m := mods[i]
			kind, provider, version, err := m.TypeProviderVersion()
			if err != nil {
				continue
			}

			if kind != "cluster" {
				continue
			}

			region, err := m.Region()
			if err != nil {
				log.Printf("skipping cluster: %q, could not parse region: source: %q, version: %q", m.Name, m.Source, m.Version)
				continue
			}

			namePrefix, err := m.NamePrefix()
			if err != nil {
				log.Printf("skipping cluster: %q, could not parse name_prefix: source: %q, version: %q", m.Name, m.Source, m.Version)
				continue
			}

			c := Cluster{
				mod:        &m,
				NamePrefix: namePrefix,
				Provider:   provider,
				Region:     region,
				Version:    version,
			}

			keys := maps.Keys(m.Configuration)
			sort.Strings(keys)

			for _, ek := range keys {
				isBk := false
				if ek == m.ConfigurationBaseKeyOrDefault() {
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

			clusters = append(clusters, c)
		}
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Name() < clusters[j].Name()
	})

	return clusters
}

func (s *Stack) NodePools() (nodePools []NodePool) {
	for _, mods := range s.root.Modules {
		for i := range mods {
			m := mods[i]
			kind, provider, version, err := m.TypeProviderVersion()
			if err != nil {
				continue
			}

			if kind != "node_pool" {
				continue
			}

			nameSuffix, err := m.NodePoolName()
			if err != nil {
				log.Printf("skipping node-pool: %q, could not parse name: source: %q, version: %q", m.Name, m.Source, m.Version)
				continue
			}

			var parentCluster, region string
			for _, c := range s.Clusters() {
				cn, err := m.ParentCluster()
				if err != nil {
					log.Printf("skipping node-pool: %q, could not parse cluster name: source: %q, version: %q", m.Name, m.Source, m.Version)
					continue
				}
				if c.Name() == cn {
					parentCluster = c.Name()
					region = c.Region
				}
			}

			if region == "" {
				log.Printf("skipping node-pool: %q, could not parse region: source: %q, version: %q", m.Name, m.Source, m.Version)
				continue
			}

			np := NodePool{
				mod:         &m,
				PoolName:    nameSuffix,
				ClusterName: parentCluster,
				Provider:    provider,
				Region:      region,
				Version:     version,
			}

			np.Configurations = parseConfiguration(m.ConfigurationBaseKey, m.Configuration)

			nodePools = append(nodePools, np)
		}
	}

	sort.Slice(nodePools, func(i, j int) bool {
		return nodePools[i].Name() < nodePools[j].Name()
	})

	return nodePools
}

func (s *Stack) Services() (services []Service) {
	for _, mods := range s.root.Modules {
		for i := range mods {
			m := mods[i]
			kind, _, _, err := m.TypeProviderVersion()
			if err != nil {
				continue
			}

			if kind != "service" {
				continue
			}

			parentCluster, err := m.ParentCluster()
			if err != nil {
				log.Printf("skipping service: %q, could not parse parent cluster name", m.Name)
				continue
			}

			var entryName string
			if strings.HasSuffix(m.Source, "/kustomization") {
				spl := strings.Split(m.Source, "/")
				entryName = spl[len(spl)-2]
			} else if strings.HasPrefix(m.Name, parentCluster) {
				spl := strings.Split(m.Name, "_")
				entryName = spl[len(spl)-1]
			}

			if entryName == "" {
				log.Printf("skipping service: %q, could not parse entry name: source: %q, version: %q", m.Name, m.Source, m.Version)
				continue
			}

			svc := Service{
				mod:         &m,
				EntryName:   entryName,
				ClusterName: parentCluster,
				Provider:    "kustomization",
				Version:     m.Version,
			}

			svc.Configurations = parseConfiguration(m.ConfigurationBaseKey, m.Configuration)

			services = append(services, svc)
		}
	}

	sort.Slice(services, func(i, j int) bool {
		return services[i].Name() < services[j].Name()
	})

	return services
}

func (s *Stack) Modules() (modules []Module) {
	for _, mods := range s.root.Modules {
		for i := range mods {
			m := mods[i]
			_, _, _, err := m.TypeProviderVersion()
			if err == nil {
				continue
			}

			modules = append(modules, Module{
				mod: &m,
			})
		}
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name() < modules[j].Name()
	})

	return modules
}

func (s *Stack) dockerfile() (out []byte, err error) {
	for k, v := range s.root.Parser.Files() {
		if !strings.HasSuffix(k, "Dockerfile") {
			continue
		}
		out = v.Bytes

		nd := dockerfile(out, s.Clusters())
		if !bytes.Equal(out, nd) {
			s.root.WriteFiles(map[string][]byte{"Dockerfile": nd})
			out = nd
		}
	}

	err = s.root.Write()
	if err != nil {
		return out, err
	}

	return out, nil
}

func (s *Stack) InitFiles(baseDomain string) error {
	data := make(map[string][]byte)

	p := make(map[string]map[string]string)
	p["kustomization"] = make(map[string]string)
	p["kustomization"]["source"] = "kbst/kustomization"
	fver := hclwrite.NewEmptyFile()
	tfhcl.BlockTerraform(fver, p)
	data["versions.tf"] = fver.Bytes()

	fvar := hclwrite.NewEmptyFile()
	tfhcl.BlockVariable(fvar, "base_domain", "string", "Used to generate fully qualified domain names for all clusters.")
	data["variables.tf"] = fvar.Bytes()

	ftfvars := hclwrite.NewEmptyFile()
	ftfvars.Body().SetAttributeValue("base_domain", cty.StringVal(baseDomain))
	data["config.auto.tfvars"] = ftfvars.Bytes()

	err := s.root.WriteFiles(data)
	if err != nil {
		return err
	}

	return s.root.Write()
}

func (s *Stack) AddCluster(namePrefix, provider, region, version string, configurations []Configuration) (c Cluster, err error) {
	if version == "" {
		version = "latest"
	}

	frameworkVersion, err := s.cliJSON.Framework.GetReleaseOrLatest(version)
	if err != nil {
		return c, err
	}

	c.NamePrefix = namePrefix
	c.Provider = provider
	c.Region = region
	c.Version = frameworkVersion.Name
	c.Configurations = configurations

	for _, ec := range s.Clusters() {
		if ec.NamePrefix == c.NamePrefix &&
			ec.Provider == c.Provider &&
			ec.Region == c.Region {
			return c, fmt.Errorf("error: cluster %q already exists", ec.Name())
		}
	}

	err = c.Validate(s.cliJSON)
	if err != nil {
		return c, err
	}

	err = s.root.WriteFiles(c.ToHCL())
	if err != nil {
		return c, err
	}

	err = s.root.Write()
	if err != nil {
		return c, err
	}

	_, err = s.dockerfile()
	if err != nil {
		return c, err
	}

	return c, nil
}

func (s *Stack) AddNodePool(clusterName, poolName string, configurations []Configuration) (np NodePool, err error) {
	var provider, region, version string
	for _, c := range s.Clusters() {
		if c.Name() == clusterName {
			provider = c.Provider
			region = c.Region
			version = c.Version
		}
	}

	if provider == "" || region == "" || version == "" {
		return np, fmt.Errorf("no cluster named %q found", clusterName)
	}

	np.ClusterName = clusterName
	np.PoolName = poolName
	np.Provider = provider
	np.Region = region
	np.Version = version
	np.Configurations = configurations

	for _, enp := range s.NodePools() {
		if enp.ClusterName == np.ClusterName &&
			enp.PoolName == np.PoolName {
			return np, fmt.Errorf("error: node pool %q already exists", enp.Name())
		}
	}

	err = np.Validate(s.cliJSON)
	if err != nil {
		return np, err
	}

	err = s.root.WriteFiles(np.ToHCL())
	if err != nil {
		return np, err
	}

	err = s.root.Write()
	if err != nil {
		return np, err
	}

	return np, nil
}

func (s *Stack) AddService(clusterName, entryName, version string) (svc Service, err error) {
	var foundCluster bool
	for _, c := range s.Clusters() {
		if c.Name() == clusterName {
			foundCluster = true
		}
	}

	if !foundCluster {
		return svc, fmt.Errorf("no cluster named %q found", clusterName)
	}

	catalogEntry, found := s.cliJSON.Catalog[entryName]
	if !found {
		return svc, fmt.Errorf("no entry named %q found in catalog", entryName)
	}

	if version == "" {
		version = "latest"
	}

	catalogVersion, err := catalogEntry.GetReleaseOrLatest(version)
	if err != nil {
		return svc, err
	}

	svc.ClusterName = clusterName
	svc.EntryName = entryName
	svc.Provider = "kustomization"
	svc.Version = catalogVersion.Name
	svc.Configurations = GenerateConfigurations(s.Environments, map[string]cty.Value{})

	for _, esvc := range s.Services() {
		if esvc.ClusterName == svc.ClusterName &&
			esvc.EntryName == svc.EntryName {
			return svc, fmt.Errorf("error: service %q already exists", esvc.Name())
		}
	}

	err = s.root.WriteFiles(svc.ToHCL())
	if err != nil {
		return svc, err
	}

	err = s.root.Write()
	if err != nil {
		return svc, err
	}

	return svc, nil
}

func (s *Stack) Remove(rm string) error {
	toDelete := []string{}
	for mfn, ms := range s.root.Modules {
		for _, m := range ms {
			kind, _, _, err := m.TypeProviderVersion()
			if err != nil {
				return err
			}

			if kind == "cluster" && len(s.Clusters()) == 1 {
				return fmt.Errorf("stacks require one cluster, not removing %q", m.Name)
			}

			pC, err := m.ParentCluster()
			if kind != "cluster" && err != nil {
				return fmt.Errorf("refusing to delete: %s", err)
			}
			if m.Name == rm || pC == rm {
				toDelete = append(toDelete, mfn)
			}
		}
	}

	for pfn, ps := range s.root.Providers {
		for _, p := range ps {
			if p.Name == "kustomization" && p.Alias == rm {
				toDelete = append(toDelete, pfn)
			}
		}
	}

	if len(toDelete) > 0 {
		err := s.root.DeleteFiles(toDelete)
		if err != nil {
			return err
		}

		err = s.root.Write()
		if err != nil {
			return err
		}

		_, err = s.dockerfile()
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("error %q did not match any clusters, node pools or services", rm)
}
