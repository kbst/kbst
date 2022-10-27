module "test_service" {
  source  = "test_source"
  version = "0.0.0-test.0"

  configuration_base_key = "apps-prod"
  configuration = {
    apps-prod = {}
    apps      = {}
    ops       = {}
  }
}
