package generator

import (
	_ "embed"

	"github.com/flosch/pongo2/v4"
)

//go:embed templates/versions.tf
var templateVersionsString string
var templateVersions *pongo2.Template = pongo2.Must(pongo2.FromString(templateVersionsString))

//go:embed templates/variables.tf
var templateVariablesString string
var templateVariables *pongo2.Template = pongo2.Must(pongo2.FromString(templateVariablesString))

//go:embed templates/config.auto.tfvars
var templateConfigAutoString string
var templateConfigAuto *pongo2.Template = pongo2.Must(pongo2.FromString(templateConfigAutoString))

//go:embed templates/cluster.tf
var templateClusterString string
var templateCluster *pongo2.Template = pongo2.Must(pongo2.FromString(templateClusterString))

//go:embed templates/cluster_providers.tf
var templateClusterProvidersString string
var templateClusterProviders *pongo2.Template = pongo2.Must(pongo2.FromString(templateClusterProvidersString))

//go:embed templates/cluster_node_pool.tf
var templateClusterNodePoolString string
var templateClusterNodePool *pongo2.Template = pongo2.Must(pongo2.FromString(templateClusterNodePoolString))

//go:embed templates/cluster_service.tf
var templateClusterServiceString string
var templateClusterService *pongo2.Template = pongo2.Must(pongo2.FromString(templateClusterServiceString))
