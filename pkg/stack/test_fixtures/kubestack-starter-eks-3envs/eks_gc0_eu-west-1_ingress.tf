module "eks_gc0_eu-west-1_nginx" {
  providers = {
    kustomization = kustomization.eks_gc0_eu-west-1
  }
  source  = "kbst.xyz/catalog/nginx/kustomization"
  version = "1.3.1-kbst.1"

  configuration_base_key = "apps-prod"
  configuration = {
    apps-prod = {}

    apps = {}

    ops = {}
  }
}

module "eks_gc0_eu-west-1_dns_zone" {
  providers = {
    aws        = aws.eks_gc0_eu-west-1
    kubernetes = kubernetes.eks_gc0_eu-west-1
  }

  source = "github.com/kbst/terraform-kubestack//aws/cluster/elb-dns?ref=v0.18.1-beta.0"

  ingress_service_name      = "ingress-nginx-controller"
  ingress_service_namespace = "ingress-nginx"

  metadata_fqdn = module.eks_gc0_eu-west-1.current_metadata["fqdn"]

  depends_on = [module.eks_gc0_eu-west-1, module.eks_gc0_eu-west-1_nginx]
}
