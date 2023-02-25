module "eks_gc0_eu-west-1_service_nginx" {
  providers = {
    kustomization = kustomization.eks_gc0_eu-west-1
  }

  source  = "kbst.xyz/catalog/nginx/kustomization"
  version = "1.3.1-kbst.1"

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {}
    apps-stg = {}
    ops      = {}
    apps-dev = {}
  }
}
