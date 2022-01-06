package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flosch/pongo2/v4"

	_ "embed"

	tfk8s "github.com/jrhouston/tfk8s/contrib/hashicorp/terraform"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type Environment struct {
	IsConfigurationBaseKey bool   `json:"is_configuration_base_key"`
	Name                   string `json:"name"`
	Key                    string `json:"key"`
}

type Stack struct {
	BaseDomain      string        `json:"base_domain"`
	BaseEnvironment string        `json:"base_environment"`
	Environments    []Environment `json:"environments"`
	Modules         []Module      `json:"modules"`
}

func (s *Stack) Unmarshal(d []byte) (err error) {
	err = json.Unmarshal(d, &s)
	if err != nil {
		return err
	}

	ckn := s.configuration_key_names()
	cbk := ckn[s.BaseEnvironment]

	var mod []Module
	for _, m := range s.Modules {
		m.cfgsToCfg(cbk, ckn)

		var cmod []Module
		for _, cm := range m.Children {
			cm.cfgsToCfg(cbk, ckn)

			cmod = append(cmod, cm)
		}
		m.Children = cmod

		mod = append(mod, m)
	}
	s.Modules = mod

	return nil
}

func (s *Stack) configuration_key_names() map[string]string {
	kn := make(map[string]string)
	for _, e := range s.Environments {
		kn[e.Key] = e.Name
	}

	return kn
}

func (s *Stack) Terraform() (files map[string]string, err error) {
	files = make(map[string]string)
	ckn := s.configuration_key_names()
	cbk := ckn[s.BaseEnvironment]

	files["versions.tf"], err = render(templateVersions, nil)
	if err != nil {
		return files, err
	}

	files["variables.tf"], err = render(templateVariables, nil)
	if err != nil {
		return files, err
	}

	files["config.auto.tfvars"], err = render(templateConfigAuto, pongo2.Context{"base_domain": s.BaseDomain})
	if err != nil {
		return files, err
	}

	for _, m := range s.Modules {
		new_files, err := m.toHCL(cbk, s.BaseDomain)
		if err != nil {
			return nil, err
		}

		for k, v := range new_files {
			files[k] = v
		}
	}
	return files, nil
}

func render(t *pongo2.Template, ctx pongo2.Context) (s string, err error) {
	s, err = t.Execute(ctx)
	if err != nil {
		return s, err
	}

	s = strings.TrimSpace(s)
	s += "\n"

	return s, nil
}

type Module struct {
	Name           string                   `json:"name"`
	Provider       string                   `json:"provider"`
	Type           string                   `json:"type"`
	Version        string                   `json:"version"`
	Children       []Module                 `json:"children"`
	Configurations []map[string]interface{} `json:"configurations"`
	Configuration  map[string]map[string]interface{}
}

func (m *Module) cfgsToCfg(cbk string, ckn map[string]string) {
	m.Configuration = make(map[string]map[string]interface{})

	for _, c := range m.Configurations {
		k := ckn[c["env_key"].(string)]
		v := c["data"].(map[string]interface{})

		m.Configuration[k] = v
	}
}

func (m *Module) GetK8sServiceName() string {
	k8sServiceName := map[string]string{
		"aws":     "eks",
		"azurerm": "aks",
		"google":  "gke",
	}

	return k8sServiceName[m.Provider]
}

func (m *Module) toHCL(cbk string, base_domain string) (files map[string]string, err error) {
	files = make(map[string]string)

	region := m.Configuration[cbk]["region"].(string)

	n := fmt.Sprintf(
		"%s_%s_%s",
		m.GetK8sServiceName(),
		m.Configuration[cbk]["name_prefix"].(string),
		region)

	cfg, err := m.cfgToHCL(cbk, "")
	if err != nil {
		return files, fmt.Errorf("cfgToHCL failed: %s", err)
	}

	data := pongo2.Context{
		"name":                   n,
		"provider":               m.Provider,
		"version":                m.Version,
		"configuration_base_key": cbk,
		"configuration":          cfg,
	}

	files[fmt.Sprintf("%s_cluster.tf", n)], err = render(templateCluster, data)
	if err != nil {
		return files, err
	}

	if m.Type == "cluster" {
		pData := pongo2.Context{
			"clusterModule": n,
			"provider":      m.Provider,
			"region":        region,
		}

		files[fmt.Sprintf("%s_%s.tf", n, "providers")], err = render(templateClusterProviders, pData)
		if err != nil {
			return files, err
		}
	}

	for _, cm := range m.Children {

		cn := ""
		var ct *pongo2.Template
		var cData pongo2.Context

		cmcfg, err := cm.cfgToHCL(cbk, n)
		if err != nil {
			return files, fmt.Errorf("cfgToHCL failed: %s", err)
		}

		if cm.Type == "node_pool" {
			var npn string
			if cm.Provider == "azurerm" {
				npn = cm.Configuration[cbk]["node_pool_name"].(string)
			} else {
				npn = cm.Configuration[cbk]["name"].(string)
			}
			cn = fmt.Sprintf("%s_%s_%s", n, cm.Type, npn)
			ct = templateClusterNodePool
			cData = pongo2.Context{
				"name":                   cn,
				"provider":               cm.Provider,
				"version":                cm.Version,
				"clusterName":            n,
				"configuration_base_key": cbk,
				"configuration":          cmcfg,
			}
		}

		if cm.Type == "service" {
			cn = fmt.Sprintf("%s_%s_%s", n, cm.Type, cm.Name)
			ct = templateClusterService
			cData = pongo2.Context{
				"moduleName":             cn,
				"serviceName":            cm.Name,
				"provider":               cm.Provider,
				"version":                strings.TrimPrefix(cm.Version, "v"),
				"providerAlias":          n,
				"configuration_base_key": cbk,
				"configuration":          cmcfg,
			}
		}

		files[fmt.Sprintf("%s.tf", cn)], err = render(ct, cData)
		if err != nil {
			return files, err
		}
	}

	return files, nil
}

func (m *Module) cfgToHCL(cbk string, n string) (hcl string, err error) {
	cfg := m.Configuration

	if m.Type == "cluster" {
		cfg[cbk]["base_domain"] = "var.base_domain"

		if m.Provider == "aws" || m.Provider == "azurerm" {
			delete(cfg[cbk], "region")
		}
	}

	mr := fmt.Sprintf("module.%s.current_config[\"project_id\"]", n)
	if m.Type == "node_pool" {
		if m.Provider == "google" {
			cfg[cbk]["project_id"] = mr

			nl := cfg[cbk]["node_locations"].(string)
			cfg[cbk]["node_locations"] = strings.Split(nl, ",")
		}
	}

	// remove null values
	for _, ov := range m.Configuration {
		for ik, iv := range ov {
			if iv == nil {
				delete(ov, ik)
			}
		}
	}

	d, err := json.Marshal(cfg)
	if err != nil {
		return hcl, fmt.Errorf("json marshall: %s", err)
	}

	hcl, err = formatHCL(d)
	if err != nil {
		return hcl, fmt.Errorf("format hcl: %s", err)
	}

	// remove "" around base_domain var
	hcl = strings.Replace(hcl, "\"var.base_domain\"", "var.base_domain", 1)

	// remove "" around cluster module ref
	hcl = strings.Replace(hcl, fmt.Sprintf("%q", mr), mr, 1)

	return hcl, nil
}

func formatHCL(d []byte) (hcl string, err error) {
	t, err := ctyjson.ImpliedType(d)
	if err != nil {
		return hcl, fmt.Errorf("ctyjson type: %s", err)
	}

	doc, err := ctyjson.Unmarshal(d, t)
	if err != nil {
		return hcl, fmt.Errorf("ctyjson unmarshal: %s", err)
	}

	hcl = tfk8s.FormatValue(doc, 2, true)

	return hcl, nil
}
