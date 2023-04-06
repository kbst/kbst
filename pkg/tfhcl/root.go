package tfhcl

import (
	"bytes"
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
	Path        string
	Parser      *hclparse.Parser
	evalContext *hcl.EvalContext
	Variables   map[string][]Variable
	Modules     map[string][]Module
	Providers   map[string][]Provider
	toWrite     map[string][]byte
	toDelete    []string
}

func NewRoot(path string) *Root {
	r := Root{
		Path: path,
	}

	r.clear()

	return &r
}

func (r *Root) clear() {
	r.Parser = hclparse.NewParser()

	r.evalContext = &hcl.EvalContext{
		Variables: map[string]cty.Value{},
		Functions: nil,
	}
	r.Variables = make(map[string][]Variable)
	r.Modules = make(map[string][]Module)
	r.Providers = make(map[string][]Provider)
	r.toWrite = make(map[string][]byte)
	r.toDelete = make([]string, 0)
}

func (r *Root) Read() (err error) {
	// if we re-read, clear the data
	r.clear()

	files, err := os.ReadDir(r.Path)
	if err != nil {
		return err
	}

	diags := hcl.Diagnostics{}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		fp := filepath.Join(r.Path, f.Name())
		_, diag := r.Parser.ParseHCLFile(fp)
		diags.Extend(diag)
	}

	return r.decode()
}

func (r *Root) decode() (err error) {
	diags := hcl.Diagnostics{}

	// variable definitions, defaults and providers
	for k, v := range r.Parser.Files() {
		if !strings.HasSuffix(k, ".tf") {
			continue
		}

		kb := Blocks{}
		moreDiags := gohcl.DecodeBody(v.Body, r.evalContext, &kb)
		diags = append(diags, moreDiags...)

		if !slices.Contains(maps.Keys(r.Variables), k) {
			r.Variables[k] = []Variable{}
		}

		for _, vd := range kb.Variables {
			r.Variables[k] = append(r.Variables[k], vd)
			r.evalContext.Variables[vd.Name] = vd.Default
		}

		if !slices.Contains(maps.Keys(r.Providers), k) {
			r.Providers[k] = []Provider{}
		}
		r.Providers[k] = append(r.Providers[k], kb.Providers...)
	}

	// variable values from tfvar files
	for k, v := range r.Parser.Files() {
		if k != "terraform.tfvars" && !strings.HasSuffix(k, ".auto.tfvars") {
			continue
		}

		vv := make(map[string]cty.Value)
		moreDiags := gohcl.DecodeBody(v.Body, r.evalContext, &vv)
		diags = append(diags, moreDiags...)

		for ik, iv := range vv {
			r.evalContext.Variables[ik] = iv
		}
	}

	// parse, now that we know the evalContext
	for k, v := range r.Parser.Files() {
		if !strings.HasSuffix(k, ".tf") {
			continue
		}

		kb := Blocks{}
		moreDiags := gohcl.DecodeBody(v.Body, r.evalContext, &kb)
		diags = append(diags, moreDiags...)

		for i, mod := range kb.Modules {
			// parse raw module providers
			mod.Providers = make(map[string]Provider)
			for _, t := range mod.ProvidersRaw.Variables() {
				spl := t.SimpleSplit()
				k := spl.RootName()
				v := spl.Rel[0].(hcl.TraverseAttr).Name

				for _, pf := range r.Providers {
					for _, p := range pf {
						if p.Name == k && p.Alias == v {
							mod.Providers[p.Name] = p
						}
					}
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

		if !slices.Contains(maps.Keys(r.Modules), k) {
			r.Modules[k] = []Module{}
		}
		r.Modules[k] = append(r.Modules[k], kb.Modules...)
	}

	if diags.HasErrors() {
		for _, diag := range diags {
			fmt.Printf("issue parsing hcl: %s\n", diag.Error())
		}
		return diags.Errs()[0]
	}

	return nil
}

func (r *Root) WriteFiles(data map[string][]byte) error {
	for k, v := range data {
		fp := filepath.Join(r.Path, k)
		r.toWrite[fp] = v
	}

	return nil
}

func (r *Root) DeleteFiles(paths []string) error {
	r.toDelete = append(r.toDelete, paths...)

	return nil
}

func (r *Root) GetVariableValue(name string) (value cty.Value, ok bool) {
	v, ok := r.evalContext.Variables[name]
	return v, ok
}

func (r *Root) SetVariableValue(name string, value cty.Value) {
	r.evalContext.Variables[name] = value
}

func (r *Root) Write() (err error) {
	for n, f := range r.toWrite {

		var ef []byte
		exists := false
		mode := os.FileMode(0644)
		if fi, err := os.Stat(n); err == nil {
			exists = true
			mode = fi.Mode()
			ef, err = os.ReadFile(n)
			if err != nil {
				return err
			}
		}

		if !exists || !bytes.Equal(ef, f) {
			if slices.Contains(r.toDelete, n) {
				// no point writing a file
				// if we delete it later anyway
				continue
			}
			err = os.WriteFile(n, f, mode)
			if err != nil {
				return err
			}
		}
	}

	for _, n := range r.toDelete {
		err := os.Remove(n)
		if err != nil {
			return err
		}
	}

	return r.Read()
}
