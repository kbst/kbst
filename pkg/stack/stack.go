package stack

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
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
	path         string
	toDelete     []string
	toWrite      []string
	root         *tfhcl.Root
	cliJSON      util.CliJSON
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

func NewStack(r *tfhcl.Root, cj util.CliJSON) *Stack {
	s := &Stack{
		root:    r,
		cliJSON: cj,
	}

	return s
}

func (s *Stack) FromPath(p string) error {
	s.path = p
	err := s.root.Read(s.path)
	if err != nil {
		return err
	}

	bd, ok := s.root.GetVariableValue("base_domain")
	if !ok {
		return fmt.Errorf("value for required var %q not found", "base_domain")
	}

	s.BaseDomain = bd.AsString()

	// if we call FromPath again to refresh
	// we need these to be empty
	s.Clusters = []Cluster{}
	s.NodePools = []NodePool{}
	s.Services = []Service{}

	for _, mf := range s.root.Modules {
		for _, m := range mf {
			kind, provider, version, err := parseKindProviderVersion(m.Source, m.Version)
			if err != nil {
				log.Printf("ignoring module: %q: %s", m.Name, err)
				continue
			}

			if kind != "cluster" {
				continue
			}

			cbk := m.ConfigurationBaseKey
			if cbk == "" {
				cbk = "apps"
			}

			var region string

			switch provider {
			case "aws":
				if _, ok := m.Providers["aws"]; ok {
					for _, providers := range s.root.Providers {
						for _, p := range providers {
							if p.Name == "aws" {
								region = p.Region
							}
						}
					}
				}
			case "azurerm":
				_, r, err := parsePrefixRegion(m.Name)
				if err != nil {
					log.Printf("ignoring module: %q: %s", m.Name, err)
				}
				region = r
			default:
				if v, ok := m.Configuration[cbk]["region"]; ok {
					region = v.AsString()
				}
			}

			c := Cluster{
				tfMod:      &m,
				NamePrefix: m.Configuration[cbk]["name_prefix"].AsString(),
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
		}
	}

	for _, mf := range s.root.Modules {
		for _, m := range mf {
			kind, provider, version, err := parseKindProviderVersion(m.Source, m.Version)
			if err != nil {
				log.Printf("ignoring module: %q: %s", m.Name, err)
				continue
			}

			cbk := m.ConfigurationBaseKey
			if cbk == "" {
				cbk = "apps"
			}

			switch kind {
			case "cluster":
				continue
			case "node_pool":
				var nameSuffix string
				if v, ok := m.Configuration[cbk]["node_pool_name "]; ok {
					nameSuffix = v.AsString()
				} else if v, ok := m.Configuration[cbk]["name "]; ok {
					nameSuffix = v.AsString()
				}

				var clusterName, region string
				for _, c := range s.Clusters {
					if c.Name() == m.ParentCluster {
						clusterName = c.Name()
						region = c.Region
					}
				}

				np := NodePool{
					tfMod:       &m,
					PoolName:    nameSuffix,
					ClusterName: clusterName,
					Provider:    provider,
					Region:      region,
					Version:     version,
				}

				np.Configurations = parseConfiguration(m.ConfigurationBaseKey, m.Configuration)

				s.NodePools = append(s.NodePools, np)
			case "service":
				var entryName string
				if strings.HasSuffix(m.Source, "/kustomization") {
					spl := strings.Split(m.Source, "/")
					entryName = spl[len(spl)-2]
				} else if strings.HasPrefix(m.Name, m.ParentCluster) {
					spl := strings.Split(m.Name, "_")
					entryName = spl[len(spl)-1]
				} else {
					log.Printf("ignoring module: %q: could not detect entry name", m.Name)
					continue
				}

				svc := Service{
					tfMod:       &m,
					EntryName:   entryName,
					ClusterName: m.ParentCluster,
					Provider:    "kustomization",
					Version:     m.Version,
				}

				svc.Configurations = parseConfiguration(m.ConfigurationBaseKey, m.Configuration)

				s.Services = append(s.Services, svc)
			case "elb-dns":
				continue
			default:
				log.Printf("unexpected module: %q: not a kubestack module", m.Name)
				continue
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
	files[path.Join(s.path, "versions.tf")] = fver

	fvar := hclwrite.NewEmptyFile()
	tfhcl.BlockVariable(fvar, "base_domain", "string", "Used to generate fully qualified domain names for all clusters.")
	files[path.Join(s.path, "variables.tf")] = fvar

	ftfvars := hclwrite.NewEmptyFile()
	ftfvars.Body().SetAttributeValue("base_domain", cty.StringVal(s.BaseDomain))
	files[path.Join(s.path, "config.auto.tfvars")] = ftfvars

	for _, c := range s.Clusters {
		for k, v := range c.ToHCL() {
			files[path.Join(s.path, k)] = v
		}
	}

	for _, np := range s.NodePools {
		for k, v := range np.ToHCL() {
			files[path.Join(s.path, k)] = v
		}
	}

	for _, svc := range s.Services {
		for k, v := range svc.ToHCL() {
			files[path.Join(s.path, k)] = v
		}
	}

	return files, nil
}

func (s *Stack) WriteChanges() error {
	current, err := s.Files()
	if err != nil {
		return err
	}

	//
	//
	// determine files to write
	toWrite := make(map[string][]byte)
	for fn, cd := range current {
		if !slices.Contains(s.toWrite, fn) {
			continue
		}

		toWrite[fn] = cd.Bytes()
	}

	// write Dockerfile if changed
	for fn, ed := range s.root.Dockerfiles {
		cd := dockerfile(ed, s.Clusters)
		if !bytes.Equal(ed, cd) {
			toWrite[fn] = cd
		}
	}

	for _, fn := range s.toDelete {
		err := os.Remove(fn)
		if err != nil {
			return err
		}
	}

	for fn, fd := range toWrite {
		mode := os.FileMode(0644)
		fi, err := os.Stat(fn)
		if err == nil {
			mode = fi.Mode()
		}

		err = os.WriteFile(fn, fd, mode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Stack) AddCluster(namePrefix, provider, region, version string, configurations []Configuration) (err error) {
	if version == "" {
		version = "latest"
	}

	frameworkVersion, err := s.cliJSON.Framework.GetReleaseOrLatest(version)
	if err != nil {
		return err
	}

	nc := Cluster{
		NamePrefix:     namePrefix,
		Provider:       provider,
		Region:         region,
		Version:        frameworkVersion.Name,
		Configurations: configurations,
	}

	for _, ec := range s.Clusters {
		if ec.NamePrefix == nc.NamePrefix &&
			ec.Provider == nc.Provider &&
			ec.Region == nc.Region {
			return fmt.Errorf("error: cluster %q already exists", ec.Name())
		}
	}

	nc.Validate(s.cliJSON)

	s.Clusters = append(s.Clusters, nc)

	for k := range nc.ToHCL() {
		s.toWrite = append(s.toWrite, path.Join(s.path, k))
	}

	return nil
}

func (s *Stack) AddNodePool(clusterName, poolName string, configurations []Configuration) (err error) {
	var provider, region, version string
	for _, c := range s.Clusters {
		if c.Name() == clusterName {
			provider = c.Provider
			region = c.Region
			version = c.Version
		}
	}

	if provider == "" || region == "" || version == "" {
		return fmt.Errorf("no cluster named %q found", clusterName)
	}

	nnp := NodePool{
		ClusterName:    clusterName,
		PoolName:       poolName,
		Provider:       provider,
		Region:         region,
		Version:        version,
		Configurations: configurations,
	}

	for _, enp := range s.NodePools {
		if enp.ClusterName == nnp.ClusterName &&
			enp.PoolName == nnp.PoolName {
			return fmt.Errorf("error: node pool %q already exists", enp.Name())
		}
	}

	nnp.Validate(s.cliJSON)

	s.NodePools = append(s.NodePools, nnp)

	for k := range nnp.ToHCL() {
		s.toWrite = append(s.toWrite, path.Join(s.path, k))
	}

	return nil
}

func (s *Stack) AddService(clusterName, entryName, version string) (err error) {
	var foundCluster bool
	for _, c := range s.Clusters {
		if c.Name() == clusterName {
			foundCluster = true
		}
	}

	if !foundCluster {
		return fmt.Errorf("no cluster named %q found", clusterName)
	}

	catalogEntry, found := s.cliJSON.Catalog[entryName]
	if !found {
		return fmt.Errorf("no entry named %q found in catalog", entryName)
	}

	if version == "" {
		version = "latest"
	}

	catalogVersion, err := catalogEntry.GetReleaseOrLatest(version)
	if err != nil {
		return err
	}

	nsvc := Service{
		ClusterName:    clusterName,
		EntryName:      entryName,
		Provider:       "kustomization",
		Version:        catalogVersion.Name,
		Configurations: GenerateConfigurations(s.Environments, map[string]cty.Value{}),
	}

	for _, esvc := range s.Services {
		if esvc.ClusterName == nsvc.ClusterName &&
			esvc.EntryName == nsvc.EntryName {
			return fmt.Errorf("error: service %q already exists", esvc.Name())
		}
	}

	s.Services = append(s.Services, nsvc)

	for k := range nsvc.ToHCL() {
		s.toWrite = append(s.toWrite, path.Join(s.path, k))
	}

	return nil
}

func (s *Stack) Remove(rm string) error {
	madeChange := false

	// clusters
	var newClusters []Cluster
	for _, c := range s.Clusters {
		if c.Name() == rm {
			// we refuse to delete this cluster if it is the only one
			if len(s.Clusters) == 1 {
				return fmt.Errorf("stacks require one cluster, not removing %q", rm)
			}

			madeChange = true

			for k := range c.ToHCL() {
				s.toDelete = append(s.toDelete, path.Join(s.path, k))
			}

			continue
		}
		newClusters = append(newClusters, c)
	}
	s.Clusters = newClusters

	// node pools
	var newNodePools []NodePool
	for _, np := range s.NodePools {
		// we're removing node pools if the
		// node pools name matches or
		// the cluster name matches
		if np.Name() == rm || np.ClusterName == rm {
			madeChange = true

			for k := range np.ToHCL() {
				s.toDelete = append(s.toDelete, path.Join(s.path, k))
			}

			continue
		}
		newNodePools = append(newNodePools, np)
	}
	s.NodePools = newNodePools

	// services
	var newServices []Service
	for _, svc := range s.Services {
		// we're removing services if the
		// services name matches or
		// the cluster name matches
		if svc.Name() == rm || svc.ClusterName == rm {
			madeChange = true

			for k := range svc.ToHCL() {
				s.toDelete = append(s.toDelete, path.Join(s.path, k))
			}

			continue
		}
		newServices = append(newServices, svc)
	}
	s.Services = newServices

	if madeChange {
		return nil
	}

	return fmt.Errorf("error %q did not match any clusters, node pools or services", rm)
}
