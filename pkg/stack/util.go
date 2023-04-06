package stack

import (
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/zclconf/go-cty/cty"
)

func parseConfiguration(cbk string, cm map[string]map[string]cty.Value) (cfgs []Configuration) {
	for name, attrs := range cm {
		cfg := Configuration{
			EnvironmentKey: name,
			Attributes:     attrs,
		}

		if cbk == "" {
			cbk = "apps"
		}

		if name == cbk {
			cfgs = append([]Configuration{cfg}, cfgs...)
			continue
		}

		cfgs = append(cfgs, cfg)
	}

	return cfgs
}

func GenerateConfigurations(envs []Environment, baseCfg map[string]cty.Value) []Configuration {
	cfgs := []Configuration{}
	for _, env := range envs {
		attrs := make(map[string]cty.Value)
		if env.IsBaseKey {
			attrs = baseCfg
		}

		cfg := Configuration{
			EnvironmentKey: env.Key,
			Attributes:     attrs,
		}

		cfgs = append(cfgs, cfg)
	}

	return cfgs
}

func convertToTfhclConfiguration(in []Configuration) (out []tfhcl.Configuration) {
	for _, cfg := range in {
		out = append(out, tfhcl.Configuration{
			EnvironmentKey: cfg.EnvironmentKey,
			Attributes:     cfg.Attributes,
		})
	}

	return out
}
