module "eks_kbst_eu-west-1" {
  providers = {
    aws = "kbst-eu-west-1"
  }

  source = "github.com/kbst/terraform-kubestack//aws/cluster?ref=0.0.0-test.0"

  configuration = {
    apps = {
      number = 5
      string = "testvalue"
    }
    ops = {}
  }
}
