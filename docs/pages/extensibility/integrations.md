---
layout: default
title: "Extensibility: Meshery Integrations"
permalink: extensibility/integrations
type: Extensibility
abstract: 'Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via <a href="extensibility#adapters">adapters</a>, <a href="extensibility#load-generators">load generators</a> and <a href="extensibility#providers">providers</a>.'
language: en
#redirect_from: extensibility
---

Meshery provides 220+ built-in integrations which refer to the supported connections and interactions between Meshery and various cloud native platforms, tools, and technologies.
Meshery's approach is Kubernetes-native which means you can easily incorporate Meshery into your existing workflow without additional setup or integration effort.

**Integrations List**

<ul>
{% for integration in site.data.integrations.integrations %}
    <li> <strong>{{ integration.name }}</strong> <br> {{ integration.desc }}</li>
{% endfor %}
</ul>
