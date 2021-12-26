{{if eq .provider "azurerm"}}provider "azurerm" {
  features {}
}{{else if eq .provider "aws"}}provider "aws" {
  alias = "{{.clusterModule}}"

  region = "{{.region}}"
}
{{end}}
provider "kustomization" {
  alias          = "{{.clusterModule}}"
  kubeconfig_raw = module.{{.clusterModule}}.kubeconfig
}

locals {
  {{.clusterModule}}_kubeconfig = yamldecode(module.{{.clusterModule}}.kubeconfig)
}

provider "kubernetes" {
  alias = "{{.clusterModule}}"

  host                   = local.{{.clusterModule}}_kubeconfig["clusters"][0]["cluster"]["server"]
  cluster_ca_certificate = base64decode(local.{{.clusterModule}}_kubeconfig["clusters"][0]["cluster"]["certificate-authority-data"])

  exec {
    api_version = local.{{.clusterModule}}_kubeconfig["users"][0]["user"]["exec"]["apiVersion"]
    args        = local.{{.clusterModule}}_kubeconfig["users"][0]["user"]["exec"]["args"]
    command     = local.{{.clusterModule}}_kubeconfig["users"][0]["user"]["exec"]["command"]
  }
}
