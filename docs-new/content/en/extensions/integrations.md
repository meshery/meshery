---
title: Meshery Integrations
description: Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via adapters, load generators and providers
aliases: 
- extensibility/integrations
params:
  suggested_reading: false
---

Meshery provides 220+ built-in integrations which refer to the supported connections and interactions between Meshery and various cloud native platforms, tools, and technologies. Meshery's approach is Kubernetes-native which means you can easily incorporate Meshery into your existing workflow without additional setup or integration effort.

{% assign sorted_index = site.models | sort: "name" | alphabetical %}
{% assign total = sorted_index | size %}
{% capture totalled %}

### All Integrations by Name ({{ total }})

{% endcapture %}
{{totalled}}

Optionally, you can [navigate all integrations visually](https://meshery.io/integrations).

{{< integrations-list >}}
