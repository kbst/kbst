module "{{.name}}" {
  {{if eq .provider "aws"}}providers = {
    aws = aws.{{.clusterName}}
  } 
  {{end}}
  source = "github.com/kbst/terraform-kubestack//{{.provider}}/cluster/node-pool?ref={{.version}}"

  {{if eq .provider "azurerm" "aws"}}cluster_name = module.{{.clusterName}}.current_metadata["name"]{{end}}
  {{if eq .provider "azurerm"}}resource_group = module.{{.clusterName}}.current_config["resource_group"]{{end}}
  {{if eq .provider "google"}}cluster_metadata = module.{{.clusterName}}.current_metadata{{end}}

  {{if ne .configuration_base_key "apps"}}configuration_base_key = "{{.configuration_base_key}}"{{end}}
  configuration = {{.configuration}}
}
