{% macro aws() %}
module "{{ name }}" {
  providers = {
    aws = aws.{{ name }}
  }

  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster/node-pool?ref={{ version }}"

  {% if configuration_base_key != "apps" %}configuration_base_key = "{{ configuration_base_key }}"{% endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
{% endmacro %}

{% macro default() %}
module "{{ name }}" {
  source = "github.com/kbst/terraform-kubestack//{{ provider }}/cluster/node-pool?ref={{ version }}"

  {% if configuration_base_key != "apps" %}configuration_base_key = "{{ configuration_base_key }}"{% endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
{% endmacro %}

{% if provider == "aws" %}
{{ aws() }}
{% else %}
{{ default() }}
{% endif %}
