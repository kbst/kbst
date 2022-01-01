module "{{.name}}" {
  {{if eq .provider "aws"}}providers = {
    aws = aws.{{.name}}
  } 
  {{end}}
  source = "github.com/kbst/terraform-kubestack//{{.provider}}/cluster?ref={{.version}}"
  {{if ne .configuration_base_key "apps"}}configuration_base_key = "{{.configuration_base_key}}"{{end}}
  configuration = {{.configuration}}
}
