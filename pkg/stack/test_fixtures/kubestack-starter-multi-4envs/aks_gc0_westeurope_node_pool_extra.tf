module "aks_gc0_westeurope_node_pool_extra" {
  source = "github.com/kbst/terraform-kubestack//azurerm/cluster/node-pool?ref=v0.18.1-beta.0"

  cluster_name   = module.aks_gc0_westeurope.current_metadata["name"]
  resource_group = module.aks_gc0_westeurope.current_config["resource_group"]

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {
      max_count      = 9
      min_count      = 3
      node_count     = 3
      node_pool_name = "extra"
      vm_size        = "Standard_D2_v4"
    }
    apps-dev = {}
    apps-stg = {}
    ops      = {}
  }
}
