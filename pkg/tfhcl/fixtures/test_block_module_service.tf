module "test_service" {
  providers = {
    kustomization = kustomization.test_cluster
  }

  source  = "kbst.xyz/catalog/test/kustomization"
  version = "0.0.0-test.0"

  configuration = {
    apps = {}
    ops  = {}
  }
}
