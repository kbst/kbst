module "aks_zero" {
  source = "github.com/kbst/terraform-kubestack//azurerm/cluster?ref=v0.10.0-beta.0"

  configuration = var.clusters["aks_zero"]
}

module "eks_zero" {
  providers = {
    aws = aws.eks_zero
  }

  source = "github.com/kbst/terraform-kubestack//aws/cluster?ref=v0.10.0-beta.0"

  configuration = var.clusters["eks_zero"]
}

module "gke_zero" {
  source = "github.com/kbst/terraform-kubestack//google/cluster?ref=v0.10.0-beta.0"

  configuration = var.clusters["gke_zero"]
}
