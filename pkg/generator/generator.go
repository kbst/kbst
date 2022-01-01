package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/jinzhu/copier"

	_ "embed"

	tfk8s "github.com/jrhouston/tfk8s/contrib/hashicorp/terraform"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type Environment struct {
	IsConfigurationBaseKey bool   `json:"is_configuration_base_key"`
	Name                   string `json:"name"`
}

type Stack struct {
	BaseDomain   string        `json:"base_domain"`
	Environments []Environment `json:"environments"`
	Modules      []Module      `json:"modules"`
}

func (s *Stack) Unmarshal(d []byte) (err error) {
	err = json.Unmarshal(d, &s)
	if err != nil {
		return err
	}

	return nil
}

func (s *Stack) configuration_base_key() (key string) {
	for _, e := range s.Environments {
		if e.IsConfigurationBaseKey {
			key = e.Name
		}
	}

	return key
}

func (s *Stack) Terraform() (files map[string]string, err error) {
	files = make(map[string]string)
	cbk := s.configuration_base_key()

	files["versions.tf"], err = render(templateVersions, nil)
	if err != nil {
		return files, err
	}

	files["variables.tf"], err = render(templateVariables, nil)
	if err != nil {
		return files, err
	}

	files["config.auto.tfvars"], err = render(templateConfigAuto, s)
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

func render(t *template.Template, d interface{}) (s string, err error) {
	var tpl bytes.Buffer
	err = t.Execute(&tpl, d)
	if err != nil {
		return s, err
	}

	s = tpl.String()

	return s, nil
}

type Module struct {
	Name          string                 `json:"name"`
	Provider      string                 `json:"provider"`
	Type          string                 `json:"type"`
	Version       string                 `json:"version"`
	Children      []Module               `json:"children"`
	Configuration map[string]interface{} `json:"configuration"`
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

	n := fmt.Sprintf(
		"%s_%s_%s",
		m.GetK8sServiceName(),
		m.Configuration["name_prefix"].(string),
		m.Configuration["region"].(string))

	cfg, err := m.cfgToHCL()
	if err != nil {
		return files, fmt.Errorf("cfgToHCL failed: %s", err)
	}

	data := map[string]string{
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
		pData := map[string]string{
			"clusterModule": n,
			"provider":      m.Provider,
			"region":        m.Configuration["region"].(string),
		}

		files[fmt.Sprintf("%s_%s.tf", n, "providers")], err = render(templateClusterProviders, pData)
		if err != nil {
			return files, err
		}
	}

	for _, cm := range m.Children {

		cn := ""
		var ct *template.Template
		var cData map[string]string

		cmcfg, err := cm.cfgToHCL()
		if err != nil {
			return files, fmt.Errorf("cfgToHCL failed: %s", err)
		}

		if cm.Type == "node_pool" {
			cn = fmt.Sprintf("%s_%s_%s", n, cm.Type, cm.Configuration["name"].(string))
			ct = templateClusterNodePool
			cData = map[string]string{
				"name":                   cn,
				"provider":               cm.Provider,
				"version":                cm.Version,
				"clusterName":            n,
				"configuration_base_key": cbk,
				"configuration":          cmcfg,
			}
		}

		if cm.Type == "service" {
			cn = fmt.Sprintf("%s_cluster_service_%s", n, cm.Name)
			ct = templateClusterService
			cData = map[string]string{
				"moduleName":             cn,
				"serviceName":            cm.Name,
				"provider":               cm.Provider,
				"version":                cm.Version,
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

func (m *Module) cfgToHCL() (hcl string, err error) {
	var cfg map[string]interface{}
	copier.Copy(&cfg, &m.Configuration)

	if m.Type == "cluster" {
		cfg["base_domain"] = "var.base_domain"

		if m.Provider == "aws" {
			delete(cfg, "region")
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

	hcl = tfk8s.FormatValue(doc, 2)

	return hcl, nil
}
