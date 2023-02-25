module "aks_gc0_westeurope" {
  source = "github.com/kbst/terraform-kubestack//azurerm/cluster?ref=v0.18.1-beta.0"

  configuration = {
    apps = {
      base_domain                  = var.base_domain
      availability_zones           = "1,2,3"
      default_node_pool_max_count  = 9
      default_node_pool_min_count  = 3
      default_node_pool_node_count = 3
      default_node_pool_vm_size    = "Standard_D2_v4"
      name_prefix                  = "gc0"
      resource_group               = "terraform-kubestack-testing"
    }
    ops = {}
  }
}
