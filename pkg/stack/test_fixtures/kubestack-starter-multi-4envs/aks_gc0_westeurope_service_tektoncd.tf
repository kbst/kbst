module "aks_gc0_westeurope_service_tektoncd" {
  providers = {
    kustomization = kustomization.aks_gc0_westeurope
  }

  source  = "kbst.xyz/catalog/tektoncd/kustomization"
  version = "0.42.0-kbst.0"

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {}
    apps-stg = {}
    ops      = {}
    apps-dev = {}
  }
}
