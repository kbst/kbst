package generator

import (
	_ "embed"
	"text/template"
)

//go:embed templates/versions.tf
var templateVersionsString string
var templateVersions *template.Template = template.Must(template.New("templateVersions").Parse(templateVersionsString))

//go:embed templates/variables.tf
var templateVariablesString string
var templateVariables *template.Template = template.Must(template.New("templateVariables").Parse(templateVariablesString))

//go:embed templates/config.auto.tfvars
var templateConfigAutoString string
var templateConfigAuto *template.Template = template.Must(template.New("templateConfigAuto").Parse(templateConfigAutoString))

//go:embed templates/cluster.tf
var templateClusterString string
var templateCluster *template.Template = template.Must(template.New("templateCluster").Parse(templateClusterString))

//go:embed templates/cluster_providers.tf
var templateClusterProvidersString string
var templateClusterProviders *template.Template = template.Must(template.New("templateClusterProviders").Parse(templateClusterProvidersString))

//go:embed templates/cluster_node_pool.tf
var templateClusterNodePoolString string
var templateClusterNodePool *template.Template = template.Must(template.New("templateClusterNodePool").Parse(templateClusterNodePoolString))

//go:embed templates/cluster_service.tf
var templateClusterServiceString string
var templateClusterService *template.Template = template.Must(template.New("templateClusterService").Parse(templateClusterServiceString))
