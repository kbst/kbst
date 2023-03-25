provider "azurerm" {
  features {
  }
}

provider "kustomization" {
  alias = "aks_gc0_westeurope"

  kubeconfig_raw = module.aks_gc0_westeurope.kubeconfig
}

