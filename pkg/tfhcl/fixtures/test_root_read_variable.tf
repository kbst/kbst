module "test_service" {
  source  = "test_source"
  version = "0.0.0-test.0"

  test = var.testvar
}

variable "testvar" {
  type = string
}
