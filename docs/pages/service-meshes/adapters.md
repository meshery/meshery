---
layout: page
title: Adapters
permalink: service-meshes/adapters
---

# Service Mesh Adapters

### What are Meshery Adapters?

Adapters allow Meshery to interface with the different service meshes.

Meshery has adapters for managing the following service meshes.

| Service Mesh  |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; Status        |Port          |
| :------------ | :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status -%}
| [Meshery Adapter for {{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;{{ adapter.project_status }} | {{ adapter.port }}
{% endif -%}
{% endfor %}

<br>
### Adapter FAQs

#### Is each service mesh adapter made equal?
No, different service mesh adapters are written to expose the unique value of each service mesh. Consequently, they are not equally capable just as each service mesh is not equally capable as the other.

Adapters have a set of operations which are grouped based on predefined operation types. See the [extensibility](/docs/extensibility) page for more details on adapter operations.

#### How can I create a new adapter?
See the [extensibility](/docs/extensibility) documentation for details on how new Meshery adapters are made.

#### Can I run more than one instance of the same Meshery adapter?
The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the `meshery-istio` adapter. To do so, modify ~/.meshery/meshery.yaml to include multiple copies of the given adapter.

See the "[Multiple Adapters](/docs/guides/multiple-adapters)" guide for more information.
