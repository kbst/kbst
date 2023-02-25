module "eks_gc0_eu-west-1_service_prometheus" {
  providers = {
    kustomization = kustomization.eks_gc0_eu-west-1
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
