module "gke_gc0_europe-west1_service_tektoncd" {
  providers = {
    kustomization = kustomization.gke_gc0_europe-west1
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
