provider "kustomization" {
  alias = "gke_gc0_europe-west1"

  kubeconfig_raw = module.gke_gc0_europe-west1.kubeconfig
}

locals {
  gke_gc0_europe-west1_kubeconfig = yamldecode(module.gke_gc0_europe-west1.kubeconfig)
}

provider "kubernetes" {
  alias = "gke_gc0_europe-west1"

  host                   = local.gke_gc0_europe-west1_kubeconfig["clusters"][0]["cluster"]["server"]
  cluster_ca_certificate = base64decode(local.gke_gc0_europe-west1_kubeconfig["clusters"][0]["cluster"]["certificate-authority-data"])
  token                  = local.gke_gc0_europe-west1_kubeconfig["users"][0]["user"]["token"]
}
