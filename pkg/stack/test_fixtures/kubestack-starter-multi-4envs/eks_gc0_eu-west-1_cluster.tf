module "eks_gc0_eu-west-1" {
  providers = {
    aws        = aws.eks_gc0_eu-west-1
    kubernetes = kubernetes.eks_gc0_eu-west-1
  }

  source = "github.com/kbst/terraform-kubestack//aws/cluster?ref=v0.18.1-beta.0"

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {
      base_domain                = var.base_domain
      cluster_availability_zones = "eu-west-1a,eu-west-1b,eu-west-1c"
      cluster_desired_capacity   = 3
      cluster_instance_type      = "t3a.xlarge"
      cluster_max_size           = 9
      cluster_min_size           = 3
      name_prefix                = "gc0"
    }
    ops      = {}
    apps-dev = {}
    apps-stg = {}
  }
}
