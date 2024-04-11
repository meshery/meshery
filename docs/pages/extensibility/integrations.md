---
layout: default
title: "Extensibility: Meshery Integrations"
permalink: extensibility/integrations
type: Extensibility
abstract: "Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via adapters, load generators and providers"
language: en
#redirect_from: extensibility
suggested-reading: false
---

Meshery provides 220+ built-in integrations which refer to the supported connections and interactions between Meshery and various cloud native platforms, tools, and technologies. Meshery's approach is Kubernetes-native which means you can easily incorporate Meshery into your existing workflow without additional setup or integration effort.

{% assign sorted_index = site.models | sort: "name" | alphabetical %}
{% assign total = sorted_index | size %}
{% capture totalled %}

### All Integrations by Name ({{ total }})

{% endcapture %}
{{totalled}}

Optionally, you can [navigate all integrations visually](https://meshery.io/integrations).

<!--
UNCOMMENT WHEN INTEGRATIONS COLLECTION IS READY
### All Integrations by Name ({{ site.integrations.size }}) 
update: @vishalvivekm: integration collection is now models accessible via site.models and directory is docs/_models/
total: {{ site.models.size }}
-->

<ul>
    {% for item in sorted_index %}
    {% if item.type=="extensibility" and item.category=="integrations" and item.list=="include" and item.language == "en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>
