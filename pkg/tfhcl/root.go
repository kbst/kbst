package tfhcl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

type Root struct {
	Parser      *hclparse.Parser
	evalContext *hcl.EvalContext
	Variables   map[string][]Variable
	Modules     map[string][]Module
	Providers   map[string][]Provider
	Dockerfiles map[string][]byte
}

func NewRoot() *Root {
	r := Root{
		evalContext: &hcl.EvalContext{
			Variables: map[string]cty.Value{},
			Functions: nil,
		},
		Variables:   make(map[string][]Variable),
		Modules:     make(map[string][]Module),
		Providers:   make(map[string][]Provider),
		Dockerfiles: make(map[string][]byte),
	}
	p := hclparse.NewParser()
	r.Parser = p

	return &r
}

func (r *Root) Read(path string) (err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	diags := hcl.Diagnostics{}

	// variable definitions and defaults
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".tf") {
			continue
		}

		fp := filepath.Join(path, f.Name())
		hclf, diag := r.Parser.ParseHCLFile(fp)
		diags.Extend(diag)

		kb := Blocks{}
		moreDiags := gohcl.DecodeBody(hclf.Body, r.evalContext, &kb)
		diags = append(diags, moreDiags...)

		if !slices.Contains(maps.Keys(r.Variables), fp) {
			r.Variables[fp] = []Variable{}
		}

		for _, vd := range kb.Variables {
			r.Variables[fp] = append(r.Variables[fp], vd)

			r.evalContext.Variables[vd.Name] = vd.Default
		}
	}

	// variable values from tfvar files
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if f.Name() != "terraform.tfvars" && !strings.HasSuffix(f.Name(), ".auto.tfvars") {
			continue
		}

		fp := filepath.Join(path, f.Name())
		hclf, diag := r.Parser.ParseHCLFile(fp)
		diags.Extend(diag)

		vv := make(map[string]cty.Value)
		moreDiags := gohcl.DecodeBody(hclf.Body, r.evalContext, &vv)
		diags = append(diags, moreDiags...)

		for k, v := range vv {
			r.evalContext.Variables[k] = v
		}
	}

	// parse, now that we know the evalContext
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".tf") {
			continue
		}

		fp := filepath.Join(path, f.Name())
		hclf, diag := r.Parser.ParseHCLFile(fp)
		diags.Extend(diag)

		kb := Blocks{}
		moreDiags := gohcl.DecodeBody(hclf.Body, r.evalContext, &kb)
		diags = append(diags, moreDiags...)

		if !slices.Contains(maps.Keys(r.Providers), fp) {
			r.Providers[fp] = []Provider{}
		}
		r.Providers[fp] = append(r.Providers[fp], kb.Providers...)

		for i, mod := range kb.Modules {
			// parse raw module providers
			mod.Providers = make(map[string]string)
			for _, t := range mod.ProvidersRaw.Variables() {
				spl := t.SimpleSplit()
				k := spl.RootName()
				v := spl.Rel[0].(hcl.TraverseAttr).Name
				mod.Providers[k] = v
			}

			// determine ParentCluster
			// based on provider alias for service modules
			if pn, ok := mod.Providers["kustomization"]; ok {
				for _, providers := range r.Providers {
					for _, p := range providers {
						for _, t := range p.KubeconfigRaw.Variables() {
							spl := t.SimpleSplit()
							n := spl.Rel[0].(hcl.TraverseAttr).Name
							if spl.RootName() == "module" && n == pn {
								mod.ParentCluster = n
								break
							}
						}
					}
				}
			}

			// determine ParentCluster
			// based on cluster_name or cluster_metadata for node-pool modules
			for _, t := range mod.ClusterNameRaw.Variables() {
				spl := t.SimpleSplit()
				if spl.RootName() == "module" {
					mod.ParentCluster = spl.Rel[0].(hcl.TraverseAttr).Name
				}
			}

			for _, t := range mod.ClusterMetadataRaw.Variables() {
				spl := t.SimpleSplit()
				if spl.RootName() == "module" {
					mod.ParentCluster = spl.Rel[0].(hcl.TraverseAttr).Name
				}
			}

			// parse raw module configuration
			val, _ := mod.ConfigurationRaw.Value(r.evalContext)

			if !val.IsNull() {
				mod.Configuration = make(map[string]map[string]cty.Value)

				for k, v := range val.AsValueMap() {
					mod.Configuration[k] = make(map[string]cty.Value)

					for ik, iv := range v.AsValueMap() {
						// get the variable's value
						if iv.Type().HasDynamicTypes() {
							if v, ok := r.GetVariableValue(ik); ok {
								iv = v
							}
						}

						if iv.IsNull() || !iv.IsWhollyKnown() {
							continue
						}

						mod.Configuration[k][ik] = iv
					}
				}
			}

			kb.Modules[i] = mod
		}

		if !slices.Contains(maps.Keys(r.Modules), fp) {
			r.Modules[fp] = []Module{}
		}
		r.Modules[fp] = append(r.Modules[fp], kb.Modules...)
	}

	if diags.HasErrors() {
		for _, diag := range diags {
			fmt.Printf("issue parsing hcl: %s\n", diag.Error())
		}
		return diags.Errs()[0]
	}

	dfPath := filepath.Join(path, "Dockerfile")
	dfData, err := os.ReadFile(dfPath)
	if err == nil {
		r.Dockerfiles[dfPath] = dfData
	}

	return nil
}

func (r *Root) GetVariableValue(name string) (value cty.Value, ok bool) {
	v, ok := r.evalContext.Variables[name]
	return v, ok
}

func (r *Root) Write() (err error) {
	for n, f := range r.Parser.Files() {
		err := os.WriteFile(n, f.Bytes, os.FileMode(0644))
		if err != nil {
			return err
		}
	}

	return nil
}
