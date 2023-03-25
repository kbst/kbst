module "aks_kbst_westeurope" {
  source = "github.com/kbst/terraform-kubestack//aks/cluster?ref=0.0.0-test.0"

  configuration = {
    apps = {
      number = 5
      string = "testvalue"
    }
    ops = {}
  }
}
