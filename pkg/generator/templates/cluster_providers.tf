{% macro azurerm() %}
provider "azurerm" {
  features {}
}
{% endmacro %}

{% macro aws() %}
provider "aws" {
  alias = "{{ clusterModule }}"

  region = "{{ region }}"
}
{% endmacro %}

{% if provider == "aws" %}{{ aws() }}{% endif %}{% if provider == "azurerm" %}{{ azurerm() }}{% endif %}
provider "kustomization" {
  alias = "{{ clusterModule }}"

  kubeconfig_raw = module.{{ clusterModule }}.kubeconfig
}
{% if provider != "azurerm" %}
locals {
  {{ clusterModule }}_kubeconfig = yamldecode(module.{{ clusterModule }}.kubeconfig)
}

provider "kubernetes" {
  alias = "{{ clusterModule }}"

  host                   = local.{{ clusterModule }}_kubeconfig["clusters"][0]["cluster"]["server"]
  cluster_ca_certificate = base64decode(local.{{ clusterModule }}_kubeconfig["clusters"][0]["cluster"]["certificate-authority-data"])
  token                  = local.{{ clusterModule }}_kubeconfig["users"][0]["user"]["token"]
}
{% endif %}
