module "test_service" {
  providers = {
    kustomization = kustomization.test_cluster
  }

  source  = "test_source"
  version = "0.0.0-test.0"

  configuration = {
    apps = {}
    ops  = {}
  }
}
