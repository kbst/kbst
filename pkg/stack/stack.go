package stack

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type Stack struct {
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
	err := s.root.Read(p)
	if err != nil {
		return err
	}

	bd, exists := s.root.VariableValues["base_domain"]
	if !exists {
		return fmt.Errorf("value for required var %q not found", "base_domain")
	}

	s.BaseDomain = bd.AsString()

	for _, mf := range s.root.Modules {
		for _, m := range mf {
			prefix, region, err := parsePrefixRegion(m.Name)
			if err != nil {
				log.Printf("ignoring module: %q: %s", m.Name, err)
				continue
			}

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
				clusterName, nameSuffix := parseNodePoolClusteNameNameSuffix(m.Name)
				np := NodePool{
					PoolName:    nameSuffix,
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
				log.Printf("ignoring module: %q: not a kubestack module", m.Name)
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
		}
	}

	for _, np := range s.NodePools {
		for k, v := range np.ToHCL() {
			files[k] = v
		}
	}

	for _, svc := range s.Services {
		for k, v := range svc.ToHCL() {
			files[k] = v
		}
	}

	if s.root != nil {
		for k := range s.root.Parser.Sources() {
			_, found := files[k]
			if found {
				continue
			}

			// add Terraform files we do not modify
			// with empty body, so we can exclude them
			// in WriteChanges below
			files[k] = hclwrite.NewFile()
		}
	}

	return files, nil
}

func (s *Stack) WriteChanges() error {
	existing := s.root.Parser.Sources()

	current, err := s.Files()
	if err != nil {
		return err
	}

	//
	//
	// determine files to delete
	toDelete := []string{}
	for fn := range existing {
		_, found := current[fn]
		if !found {
			toDelete = append(toDelete, fn)
		}
	}

	//
	//
	// determine files to write
	toWrite := make(map[string][]byte)
	for fn, cd := range current {
		// if cd is empty, we don't touch this file
		if len(cd.Bytes()) == 0 {
			continue
		}

		// if the file does not exist yet
		// or the data has changed
		// add the name to the list of files to write
		ed, found := existing[fn]
		if !found || !bytes.Equal(ed, cd.Bytes()) {
			toWrite[fn] = cd.Bytes()
		}
	}

	// write Dockerfile if changed
	for fn, ed := range s.root.Dockerfiles {
		cd := dockerfile(ed, s.Clusters)
		if !bytes.Equal(ed, cd) {
			toWrite[fn] = cd
		}
	}

	for _, fn := range toDelete {
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

	nc.Validate(s.cliJSON)

	s.Clusters = append(s.Clusters, nc)

	return nil
}

func (s *Stack) AddNodePool(clusterName, nameSuffix string, configurations []Configuration) (err error) {
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
		PoolName:       nameSuffix,
		Provider:       provider,
		Region:         region,
		Version:        version,
		Configurations: configurations,
	}

	nnp.Validate(s.cliJSON)

	s.NodePools = append(s.NodePools, nnp)

	return nil
}

func (s *Stack) AddService(clusterName, entryName, version string) (err error) {
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

	s.Services = append(s.Services, Service{
		ClusterName:    clusterName,
		EntryName:      entryName,
		Provider:       "kustomization",
		Version:        catalogVersion.Name,
		Configurations: GenerateConfigurations(s.Environments, map[string]cty.Value{}),
	})

	return nil
}

func (s *Stack) Remove(rm string) error {
	// if we're removing a cluster
	for ic, c := range s.Clusters {
		if c.Name() == rm {
			if len(s.Clusters) == 1 {
				return fmt.Errorf("stacks require one cluster, not removing %q", rm)
			}

			s.Clusters = slices.Delete(s.Clusters, ic, ic+1)

			// if we are removing a cluster
			// also remove all its node pools
			for inp, np := range s.NodePools {
				if np.ClusterName == rm {
					s.NodePools = slices.Delete(s.NodePools, inp, inp+1)
				}
			}

			// and services
			for isvc, svc := range s.Services {
				if svc.ClusterName == rm {
					s.Services = slices.Delete(s.Services, isvc, isvc+1)
				}
			}

			return nil
		}
	}

	// if we're removing a node pool
	for inp, np := range s.NodePools {
		if np.Name() == rm {
			s.NodePools = slices.Delete(s.NodePools, inp, inp+1)

			return nil
		}
	}

	// if we're removing a service
	for isvc, svc := range s.Services {
		if svc.Name() == rm {
			s.Services = slices.Delete(s.Services, isvc, isvc+1)

			return nil
		}
	}

	return fmt.Errorf("error %q did not match any clusters, node pools or services", rm)
}
