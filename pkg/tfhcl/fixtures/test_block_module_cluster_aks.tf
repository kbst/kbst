module "aks_kbst_westeurope" {
  source = "github.com/kbst/terraform-kubestack//azurerm/cluster?ref=0.0.0-test.0"

  configuration = {
    apps = {
      base_domain = var.base_domain
      number      = 5
      string      = "testvalue"
    }
    ops = {}
  }
}
