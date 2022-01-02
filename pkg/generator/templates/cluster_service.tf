module "{{ moduleName }}" {
  providers = {
    kustomization = kustomization.{{ providerAlias }}
  }

  source  = "kbst.xyz/catalog/{{ serviceName }}/{{ provider }}"
  version = "{{ version }}"
{% if configuration_base_key != "apps" %}
  configuration_base_key = "{{ configuration_base_key }}"
{%- endif %}
  configuration = {% autoescape off %}{{ configuration }}{% endautoescape %}
}
