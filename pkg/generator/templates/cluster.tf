{% macro aws() %}
module "{{ name }}" {
  providers = {
    aws        = aws.{{ name }}
    kubernetes = kubernetes.{{ name }}
  }

  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster?ref={{ version }}"
{% if configuration_base_key != "apps" %}
  configuration_base_key = "{{ configuration_base_key }}"
{%- endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
{% endmacro %}

{% macro google() %}
module "{{ name }}" {
  providers = {
    kubernetes = kubernetes.{{ name }}
  }

  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster?ref={{ version }}"
{% if configuration_base_key != "apps" %}
  configuration_base_key = "{{ configuration_base_key }}"
{%- endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
{% endmacro %}

{% macro azurerm() %}
module "{{ name }}" {
  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster?ref={{ version }}"
{% if configuration_base_key != "apps" %}
  configuration_base_key = "{{ configuration_base_key }}"
{%- endif %}
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