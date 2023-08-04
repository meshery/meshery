---
layout: default
title: "Extensibility: Meshery Integrations"
permalink: extensibility/integrations
type: Extensibility
abstract: 'Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via <a href="extensibility#adapters">adapters</a>, <a href="extensibility#load-generators">load generators</a> and <a href="extensibility#providers">providers</a>.'
language: en
#redirect_from: extensibility
---

Meshery provides 220+ built-in integrations which refer to the supported connections and interactions between Meshery and various cloud native platforms, tools, and technologies. Meshery's approach is Kubernetes-native which means you can easily incorporate Meshery into your existing workflow without additional setup or integration effort. 

## Integrations

Optionally, you can [navigate all integrations visually](https://meshery.io/integrations).

{% assign integrations = site.data.integrations.integrations | sort: "name" %}
{% for integration in integrations %}

- **{{ integration.name }}** - {{ integration.desc }}

{% endfor %}
