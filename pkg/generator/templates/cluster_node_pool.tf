{% macro aws() %}
module "{{ name }}" {
  providers = {
    aws = aws.{{ clusterName }}
  }

  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster/node-pool?ref={{ version }}"

  cluster_name = module.{{ clusterName }}.current_metadata["name"]

  {% if configuration_base_key != "apps" %}configuration_base_key = "{{ configuration_base_key }}"{% endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
{% endmacro %}

{% macro azurerm() %}
module "{{ name }}" {
  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster/node-pool?ref={{ version }}"

  cluster_name   = module.{{ clusterName }}.current_metadata["name"]
  resource_group = module.{{ clusterName }}.current_config["resource_group"]

  {% if configuration_base_key != "apps" %}configuration_base_key = "{{ configuration_base_key }}"{% endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
{% endmacro %}

{% macro google() %}
module "{{ name }}" {
  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster/node-pool?ref={{ version }}"

  cluster_metadata = module.{{ clusterName }}.current_metadata

  {% if configuration_base_key != "apps" %}configuration_base_key = "{{ configuration_base_key }}"{% endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
{% endmacro %}

{% if provider == "aws" %}
{{ aws() }}
{% endif %}

{% if provider == "azurerm" %}
{{ azurerm() }}
{% endif %}

{% if provider == "google" %}
{{ google() }}
{% endif %}
