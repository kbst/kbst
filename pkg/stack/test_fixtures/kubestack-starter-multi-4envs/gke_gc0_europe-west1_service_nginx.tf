module "gke_gc0_europe-west1_service_nginx" {
  providers = {
    kustomization = kustomization.gke_gc0_europe-west1
  }

  source  = "kbst.xyz/catalog/nginx/kustomization"
  version = "1.3.1-kbst.1"

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {}
    ops      = {}
    apps-dev = {}
    apps-stg = {}
  }
}
