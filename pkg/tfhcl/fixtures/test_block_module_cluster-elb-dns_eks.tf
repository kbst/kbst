module "eks_kbst_eu-west-1_ingress_elb_dns" {
  providers = {
    aws = "kbst-eu-west-1"
  }

  source = "github.com/kbst/terraform-kubestack//aws/cluster/elb-dns?ref=0.0.0-test.0"

  configuration = {
    apps = {
      number = 5
      string = "testvalue"
    }
    ops = {}
  }
}
