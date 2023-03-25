module "aks_gc0_westeurope_service_nginx" {
  providers = {
    kustomization = kustomization.aks_gc0_westeurope
  }

  source  = "kbst.xyz/catalog/nginx/kustomization"
  version = "1.3.1-kbst.1"

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {}
    apps-dev = {}
    apps-stg = {}
    ops      = {}
  }
}
