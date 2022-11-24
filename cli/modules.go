package cli

import (
	"github.com/kbst/kbst/pkg/stack"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/zclconf/go-cty/cty"
)

type Modules struct {
	Catalog   map[string]util.Entry
	Framework util.Entry
	CloudInfo util.CloudInfo
}

func (m Modules) AddService(entry string, release string, gitRef string, path string) (err error) {
	r := tfhcl.NewRoot()
	s := stack.NewStack(r, util.CliJSON{})
	err = s.FromPath(path)
	if err != nil {
		return err
	}

	version, err := m.Catalog[entry].GetReleaseOrLatest("latest")
	if err != nil {
		return err
	}

	cfgs := stack.GenerateConfigurations(s.Environments, map[string]cty.Value{})

	for _, c := range s.Clusters {
		s.Services = append(s.Services, stack.Service{
			EntryName:      entry,
			ClusterName:    c.Name(),
			Provider:       "kustomization",
			Version:        version.Name,
			Configurations: cfgs,
		})
	}

	return nil
}
