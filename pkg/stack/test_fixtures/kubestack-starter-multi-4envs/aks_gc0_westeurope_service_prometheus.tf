module "aks_gc0_westeurope_service_prometheus" {
  providers = {
    kustomization = kustomization.aks_gc0_westeurope
  }

  source  = "kbst.xyz/catalog/prometheus/kustomization"
  version = "0.61.0-kbst.0"

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {}
    apps-dev = {}
    apps-stg = {}
    ops      = {}
  }
}
