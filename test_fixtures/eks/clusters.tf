module "eks_zero" {
  providers = {
    aws = aws.eks_zero
  }

  source = "github.com/kbst/terraform-kubestack//aws/cluster?ref=v0.10.0-beta.0"

  configuration = var.clusters["eks_zero"]
}
