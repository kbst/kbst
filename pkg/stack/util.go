package stack

import (
	"strings"

	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/zclconf/go-cty/cty"
)

func parsePrefixRegion(n string) (prefix, region string) {
	nspl := strings.Split(n, "_")

	return nspl[1], nspl[2]
}

func parseKindProviderVersion(s, v string) (kind, provider, version string) {
	if strings.HasPrefix(s, "github.com/kbst/terraform-kubestack//") {
		version = strings.Split(s, "?ref=")[1]
	}

	if strings.HasPrefix(s, "github.com/kbst/terraform-kubestack//aws/cluster") {
		kind = "cluster"
		provider = "aws"
	}

	if strings.HasPrefix(s, "github.com/kbst/terraform-kubestack//google/cluster") {
		kind = "cluster"
		provider = "google"
	}

	if strings.HasPrefix(s, "github.com/kbst/terraform-kubestack//azurerm/cluster") {
		kind = "cluster"
		provider = "azurerm"
	}

	if strings.HasPrefix(s, "github.com/kbst/terraform-kubestack//aws/cluster/node-pool") {
		kind = "node_pool"
	}

	if strings.HasPrefix(s, "github.com/kbst/terraform-kubestack//google/cluster/node-pool") {
		kind = "node_pool"
	}

	if strings.HasPrefix(s, "github.com/kbst/terraform-kubestack//azurerm/cluster/node-pool") {
		kind = "node_pool"
	}

	if strings.HasPrefix(s, "kbst.xyz/catalog") {
		kind = "service"
		provider = "kustomization"
		version = v
	}

	return kind, provider, version
}

func parseNodePoolClusteNameNameSuffix(s string) (clusterName, nameSuffix string) {
	spl := strings.Split(s, "_node_pool_")
	return spl[0], spl[1]
}

func parseServiceClusteNameEntryName(s string) (clusterName, entryName string) {
	spl := strings.Split(s, "_service_")
	return spl[0], spl[1]
}

func parseConfiguration(cbk string, cm map[string]map[string]cty.Value) (cfgs []Configuration) {
	for name, attrs := range cm {
		cfg := Configuration{
			EnvironmentKey: name,
			Attributes:     attrs,
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
