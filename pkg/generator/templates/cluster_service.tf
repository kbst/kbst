module "{{.moduleName}}" {
  providers = {
    kustomization = kustomization.{{.providerAlias}}
  }

  source  = "kbst.xyz/catalog/{{.serviceName}}/{{.provider}}"
  version = "{{.version}}"
  {{if ne .configuration_base_key "apps"}}configuration_base_key = "{{.configuration_base_key}}"{{end}}
  configuration = {{.configuration}}
}
