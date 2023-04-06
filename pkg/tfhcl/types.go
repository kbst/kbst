package tfhcl

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type Blocks struct {
	Modules   []Module    `hcl:"module,block"`
	Providers []Provider  `hcl:"provider,block"`
	Variables []Variable  `hcl:"variable,block"`
	Terraform []Terraform `hcl:"terraform,block"`
	Locals    []Locals    `hcl:"locals,block"`

	// kubestack unused
	Resources   []Resource   `hcl:"resource,block"`
	DataSources []DataSource `hcl:"data,block"`
	Output      []Output     `hcl:"output,block"`
}

type Module struct {
	Name               string         `hcl:"name,label"`
	ProvidersRaw       hcl.Expression `hcl:"providers,optional"`
	Providers          map[string]Provider
	Source             string         `hcl:"source"`
	Version            string         `hcl:"version,optional"`
	ClusterNameRaw     hcl.Expression `hcl:"cluster_name,optional"`
	ClusterMetadataRaw hcl.Expression `hcl:"cluster_metadata,optional"`
	MetadataFQDNRaw    hcl.Expression `hcl:"metadata_fqdn,optional"`

	ConfigurationBaseKey string         `hcl:"configuration_base_key,optional"`
	ConfigurationRaw     hcl.Expression `hcl:"configuration"`
	Configuration        map[string]map[string]cty.Value
	Body                 hcl.Body `hcl:",remain"`
}

type Provider struct {
	Name          string         `hcl:"name,label"`
	Alias         string         `hcl:"alias,optional"`
	Region        string         `hcl:"region,optional"`
	KubeconfigRaw hcl.Expression `hcl:"kubeconfig_raw,optional"`
	Body          hcl.Body       `hcl:",remain"`
}

type Variable struct {
	Name        string         `hcl:"name,label"`
	Type        hcl.Expression `hcl:"type,optional"`
	Description string         `hcl:"description,optional"`
	Default     cty.Value      `hcl:"default,optional"`
	Body        hcl.Body       `hcl:",remain"`
}

type Terraform struct {
	Backend           []Backend          `hcl:"backend,block"`
	RequiredProviders []RequiredProvider `hcl:"required_providers,block"`
	RequiredVersion   string             `hcl:"required_version,optional"`
}

type RequiredProvider struct {
	Body hcl.Body `hcl:",remain"`
}

type Backend struct {
	Name string   `hcl:"name,label"`
	Body hcl.Body `hcl:",remain"`
}

type Locals struct {
	Body hcl.Body `hcl:",remain"`
}

type Resource struct {
	Type string   `hcl:"type,label"`
	Name string   `hcl:"name,label"`
	Body hcl.Body `hcl:",remain"`
}

type DataSource struct {
	Type string   `hcl:"type,label"`
	Name string   `hcl:"name,label"`
	Body hcl.Body `hcl:",remain"`
}

type Output struct {
	Name string   `hcl:"name,label"`
	Body hcl.Body `hcl:",remain"`
}

func (m *Module) ConfigurationBaseKeyOrDefault() string {
	cbk := m.ConfigurationBaseKey
	if cbk == "" {
		cbk = "apps"
	}

	return cbk
}

func (m *Module) TypeProviderVersion() (t, p, v string, err error) {
	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//") {
		v = strings.Split(m.Source, "?ref=")[1]
	}

	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//aws/cluster") {
		t = "cluster"
		p = "aws"
	}

	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//aws/cluster/elb-dns") {
		t = "elb-dns"
	}

	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//google/cluster") {
		t = "cluster"
		p = "google"
	}

	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//azurerm/cluster") {
		t = "cluster"
		p = "azurerm"
	}

	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//aws/cluster/node-pool") {
		t = "node_pool"
	}

	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//google/cluster/node-pool") {
		t = "node_pool"
	}

	if strings.HasPrefix(m.Source, "github.com/kbst/terraform-kubestack//azurerm/cluster/node-pool") {
		t = "node_pool"
	}

	if strings.HasPrefix(m.Source, "kbst.xyz/catalog") {
		t = "service"
		p = "kustomization"
		v = m.Version
	}

	if t != "" && p != "" && v != "" {
		return t, p, v, nil
	}

	return t, p, v, fmt.Errorf("could not detect type, provider and version for: source: %q, version: %q", m.Source, m.Version)
}

func (m *Module) Region() (region string, err error) {
	_, provider, _, err := m.TypeProviderVersion()
	if err != nil {
		return region, err
	}

	switch provider {
	case "aws":
		if p, ok := m.Providers["aws"]; ok {
			region = p.Region
		}
	case "azurerm":
		nspl := strings.Split(m.Name, "_")

		if len(nspl) < 3 {
			return region, fmt.Errorf("can not parse region from %q", m.Name)
		}

		region = nspl[2]
	default:
		if v, ok := m.Configuration[m.ConfigurationBaseKeyOrDefault()]["region"]; ok {
			region = v.AsString()
		} else {
			return region, fmt.Errorf("module %q has no region configuration attribute", m.Name)
		}
	}

	return region, nil
}

func (m *Module) NamePrefix() (region string, err error) {
	if v, ok := m.Configuration[m.ConfigurationBaseKeyOrDefault()]["name_prefix"]; ok {
		return v.AsString(), nil
	}
	return "", fmt.Errorf("module %q has no name_prefix configuration attribute", m.Name)
}

func (m *Module) NodePoolName() (name string, err error) {
	cbk := m.ConfigurationBaseKeyOrDefault()

	if v, ok := m.Configuration[cbk]["node_pool_name"]; ok {
		name = v.AsString()
	} else if v, ok := m.Configuration[cbk]["name"]; ok {
		name = v.AsString()
	} else {
		return name, fmt.Errorf("module %q has no node_pool_name or name configuration attribute", m.Name)
	}

	return name, nil
}

func (m *Module) ParentCluster() (string, error) {
	// provider alias
	if p, ok := m.Providers["kustomization"]; ok {
		for _, t := range p.KubeconfigRaw.Variables() {
			spl := t.SimpleSplit()
			n := spl.Rel[0].(hcl.TraverseAttr).Name
			if spl.RootName() == "module" && n == p.Alias {
				return n, nil
			}
		}
	}

	// cluster_name
	for _, t := range m.ClusterNameRaw.Variables() {
		spl := t.SimpleSplit()
		if spl.RootName() == "module" {
			return spl.Rel[0].(hcl.TraverseAttr).Name, nil
		}
	}

	// cluster_metadata
	for _, t := range m.ClusterMetadataRaw.Variables() {
		spl := t.SimpleSplit()
		if spl.RootName() == "module" {
			return spl.Rel[0].(hcl.TraverseAttr).Name, nil
		}
	}

	// metadata_fqdn
	for _, t := range m.MetadataFQDNRaw.Variables() {
		spl := t.SimpleSplit()
		if spl.RootName() == "module" {
			return spl.Rel[0].(hcl.TraverseAttr).Name, nil
		}
	}

	return "", fmt.Errorf("no parent cluster found for %q", m.Name)
}
