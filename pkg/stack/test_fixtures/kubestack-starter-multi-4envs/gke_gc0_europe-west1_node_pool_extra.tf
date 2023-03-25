module "gke_gc0_europe-west1_node_pool_extra" {
  source = "github.com/kbst/terraform-kubestack//google/cluster/node-pool?ref=v0.18.1-beta.0"

  cluster_metadata = module.gke_gc0_europe-west1.current_metadata

  configuration_base_key = "apps-prd"
  configuration = {
    apps-prd = {
      location           = module.gke_gc0_europe-west1.current_config["region"]
      project_id         = module.gke_gc0_europe-west1.current_config["project_id"]
      initial_node_count = 1
      machine_type       = "e2-standard-8"
      max_node_count     = 3
      min_node_count     = 1
      name               = "extra"
    }
    apps-stg = {}
    ops      = {}
    apps-dev = {}
  }
}
