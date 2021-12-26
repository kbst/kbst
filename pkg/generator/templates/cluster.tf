module "{{.name}}" {
  {{if eq .provider "aws"}}providers = {
    aws = aws.{{.name}}
  } 
  {{end}}
  source = "github.com/kbst/terraform-kubestack//{{.provider}}/cluster?ref=v0.16.0-beta.0"
  {{if ne .configuration_base_key "apps"}}configuration_base_key = "{{.configuration_base_key}}"{{end}}
  configuration = {{.configuration}}
}
