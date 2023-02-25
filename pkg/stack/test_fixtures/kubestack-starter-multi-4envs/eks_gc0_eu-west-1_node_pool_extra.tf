module "eks_gc0_eu-west-1_node_pool_extra" {
  providers = {
    aws = aws.eks_gc0_eu-west-1
  }

  source = "github.com/kbst/terraform-kubestack//aws/cluster/node-pool?ref=v0.18.1-beta.0"

  cluster_name = module.eks_gc0_eu-west-1.current_metadata["name"]

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {
      desired_capacity = 3
      instance_types   = "t3a.xlarge"
      max_size         = 9
      min_size         = 3
      name             = "extra"
    }
    apps-stg = {}
    ops      = {}
    apps-dev = {}
  }
}
