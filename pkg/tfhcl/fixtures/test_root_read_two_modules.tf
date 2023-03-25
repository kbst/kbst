module "test_mod1" {
  source  = "test_source"
  version = "0.0.0-test.0"

  test = var.testvar
}

module "test_mod2" {
  source  = "test_source"
  version = "0.0.0-test.0"

  test = var.testvar
}
