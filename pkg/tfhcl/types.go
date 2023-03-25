package tfhcl

import (
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
	Name                 string         `hcl:"name,label"`
	ProvidersRaw         hcl.Expression `hcl:"providers,optional"`
	Providers            map[string]string
	Source               string         `hcl:"source"`
	Version              string         `hcl:"version,optional"`
	ConfigurationBaseKey string         `hcl:"configuration_base_key,optional"`
	ConfigurationRaw     hcl.Expression `hcl:"configuration"`
	Configuration        map[string]map[string]cty.Value
	Body                 hcl.Body `hcl:",remain"`
}

type Provider struct {
	Name   string   `hcl:"name,label"`
	Region string   `hcl:"region,optional"`
	Body   hcl.Body `hcl:",remain"`
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
