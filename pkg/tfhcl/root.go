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
	Parser         *hclparse.Parser
	evalContext    *hcl.EvalContext
	Variables      map[string][]Variable
	VariableValues map[string]cty.Value
	Modules        map[string][]Module
	Providers      map[string][]Provider
	Dockerfiles    map[string][]byte
}

func NewRoot() *Root {
	r := Root{
		evalContext: nil,
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
		r.Variables[fp] = append(r.Variables[fp], kb.Variables...)

		for i, mod := range kb.Modules {
			// parse raw module providers
			mod.Providers = make(map[string]string)
			for _, t := range mod.ProvidersRaw.Variables() {
				spl := t.SimpleSplit()
				k := spl.RootName()
				v := spl.Rel[0].(hcl.TraverseAttr).Name
				mod.Providers[k] = v
			}

			// parse raw module configuration
			val, _ := mod.ConfigurationRaw.Value(r.evalContext)

			if !val.IsNull() {
				mod.Configuration = make(map[string]map[string]cty.Value)
				mod.Configuration = make(map[string]map[string]cty.Value)

				mod.Configuration = make(map[string]map[string]cty.Value)

				for k, v := range val.AsValueMap() {
					mod.Configuration[k] = make(map[string]cty.Value)

					for ik, iv := range v.AsValueMap() {
						if iv.IsNull() && !iv.IsWhollyKnown() {
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

		if !slices.Contains(maps.Keys(r.Providers), fp) {
			r.Providers[fp] = []Provider{}
		}
		r.Providers[fp] = append(r.Providers[fp], kb.Providers...)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if f.Name() != "config.auto.tfvars" {
			continue
		}

		fp := filepath.Join(path, f.Name())
		hclf, diag := r.Parser.ParseHCLFile(fp)
		diags.Extend(diag)

		vv := make(map[string]cty.Value)
		moreDiags := gohcl.DecodeBody(hclf.Body, nil, &vv)
		diags = append(diags, moreDiags...)

		r.VariableValues = vv
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

func (r *Root) Write() (err error) {
	for n, f := range r.Parser.Files() {
		err := os.WriteFile(n, f.Bytes, os.FileMode(0644))
		if err != nil {
			return err
		}
	}

	return nil
}
