module "aks_kbst_westeurope_node_pool_extra" {
  source = "github.com/kbst/terraform-kubestack//azurerm/cluster/node-pool?ref=0.0.0-test.0"

  configuration = {
    apps = {
      base_domain = var.base_domain
      number      = 5
      string      = "testvalue"
    }
    ops = {}
  }
}
