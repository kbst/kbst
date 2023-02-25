module "gke_gc0_europe-west1" {
  providers = {
    kubernetes = kubernetes.gke_gc0_europe-west1
  }

  source = "github.com/kbst/terraform-kubestack//google/cluster?ref=v0.18.1-beta.0"

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {
      base_domain                = var.base_domain
      cluster_initial_node_count = 1
      cluster_machine_type       = "e2-standard-8"
      cluster_max_node_count     = 3
      cluster_min_master_version = "1.20"
      cluster_min_node_count     = 1
      cluster_node_locations     = "europe-west1-b,europe-west1-c,europe-west1-d"
      name_prefix                = "gc0"
      project_id                 = "terraform-kubestack-testing"
      region                     = "europe-west1"
    }
    ops      = {}
    apps-dev = {}
    apps-stg = {}
  }
}
