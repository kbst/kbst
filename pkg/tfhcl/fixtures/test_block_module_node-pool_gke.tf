module "gke_kbst_europe-west1_node_pool_extra" {
  source = "github.com/kbst/terraform-kubestack//google/cluster/node-pool?ref=0.0.0-test.0"

  configuration = {
    apps = {
      number = 5
      string = "testvalue"
    }
    ops = {}
  }
}
