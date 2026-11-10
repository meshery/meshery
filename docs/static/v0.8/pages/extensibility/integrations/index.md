<!-- ---
layout: page
title: Integrations
abstract: Integrations with other services.
language: en
permalink: extensibility/integrations
type: integrations
display-title: "true"
language: en
list: exclude
abstract: Integrations with other platforms and services.
---
Meshery provides 220+ built-in integrations which refer to the supported connections and interactions between Meshery and various cloud native platforms, tools, and technologies. Meshery's approach is Kubernetes-native which means you can easily incorporate Meshery into your existing workflow without additional setup or integration effort. 

{% assign sorted_index = site.pages | sort: "name" | alphabetical %}
{% assign total = sorted_index | size %}
{% capture totalled %}

### All Integrations by Name ({{ total }})

{% endcapture %}
{{totalled}}

<!--
UNCOMMENT WHEN INTEGRATIONS COLLECTION IS READY
### All Integrations by Name ({{ site.integrations.size }}) 

<ul>
  {% assign sorted_index = site.pages | where: "type", "extensibility" | sort: "name" | alphabetical %}

    {% for item in sorted_index %}
    {% if item.type=="extensibility" and item.category=="integrations" and item.list=="include" and item.language == "en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul> -->
