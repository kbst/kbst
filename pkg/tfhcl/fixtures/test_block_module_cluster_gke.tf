module "gke_kbst_europe-west1" {
  source = "github.com/kbst/terraform-kubestack//google/cluster?ref=0.0.0-test.0"

  configuration = {
    apps = {
      number = 5
      string = "testvalue"
    }
    ops = {}
  }
}
