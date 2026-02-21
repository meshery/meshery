---
layout: default
title: "Meshery Integrations"
permalink: /extensions/integrations
type: extensions
abstract: "Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via adapters, load generators, and providers."
language: en
redirect_from:
  - /extensibility/integrations
suggested-reading: false
list: include
---

Meshery provides 220+ built-in integrations which refer to the supported connections and interactions between Meshery and various cloud native platforms, tools, and technologies. Meshery's approach is Kubernetes-native which means you can easily incorporate Meshery into your existing workflow without additional setup or integration effort.

Optionally, you can [navigate all integrations visually](https://meshery.io/integrations).

{% comment %}
Integration selection criteria:

- Models with layout: "integration"
- Models with an `integrations-category` key
- Models with type: "extensions" and category: "integrations"
  From these, entries with `isAnnotation: true` are displayed separately.
  {% endcomment %}

{%- assign integrations_all = "" | split: "" -%}

{%- for m in site.models -%}
{% if m.layout == "integration" %}
{%- assign integrations_all = integrations_all | push: m -%}
{% elsif m.integrations-category %}
{%- assign integrations_all = integrations_all | push: m -%}
{% elsif m.type == "extensions" and m.category == "integrations" %}
{%- assign integrations_all = integrations_all | push: m -%}
{% endif %}
{%- endfor -%}

{%- assign annotations = integrations_all | where: "isAnnotation", true -%}
{%- assign main_integrations = integrations_all | where_exp: "i", "i.isAnnotation != true" -%}

{%- assign sorted_integrations = main_integrations | sort: "title" -%}
{%- assign sorted_annotations = annotations | sort: "title" -%}
{%- assign total = integrations_all | size -%}

### All Integrations by Name ({{ total }})

#### Full Service Integrations ({{ sorted_integrations | size }})

These are the primary integrations and extensions that Meshery models directly.

<ul>
  {% for item in sorted_integrations %}
    <li>
      <a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract %} - {{ item.abstract }}{% endif %}
    </li>
  {% endfor %}
</ul>

{% if sorted_annotations.size > 0 %}

#### Supporting Annotations ({{ sorted_annotations | size }})

The following entries are **annotation-based models**. They represent metadata, labels, or lightweight objects that support Meshery’s full integrations, but are not standalone services themselves.

<ul>
  {% for ann in sorted_annotations %}
    <li>
      <a href="{{ site.baseurl }}{{ ann.url }}">{{ ann.title }}</a>
      {% if ann.abstract %} - {{ ann.abstract }}{% endif %}
    </li>
  {% endfor %}
</ul>
{% endif %}
