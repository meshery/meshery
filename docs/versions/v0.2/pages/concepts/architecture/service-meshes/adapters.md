---
layout: default
title: Adapters
permalink: concepts/architecture/adapters
type: concepts
redirect_from: architecture/adapters
abstract: "Meshery offers support for more adapters than any other project or product in the world. Meshery uses adapters for managing the various service meshes."
language: en
list: include
---

As the multi-mesh manager, Meshery offers support for more adapters than any other project or product in the world. Meshery uses adapters for managing the various service meshes.

## What are Meshery Adapters?

Adapters allow Meshery to interface with the different service meshes, exposing their differentiated value to users.

Meshery has adapters for managing the following service meshes.
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Adapter Status |  Service Mesh  | Service Mesh Version | Port          |
| :------------: | :------------ | :------------: | :------------ |
{% for adapter in sorted -%}
{% if adapter.project_status -%}
| {{ adapter.project_status }} | <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.version }} | {{ adapter.port }} |
{% endif -%}
{% endfor %}

### Adapter FAQs

#### Is each service mesh adapter made equal?
No, different service mesh adapters are written to expose the unique value of each service mesh. Consequently, they are not equally capable just as each service mesh is not equally capable as the other.

Adapters have a set of operations which are grouped based on predefined operation types. See the [extensibility]({{site.baseurl}}/extensibility) page for more details on adapter operations.

#### How can I create a new adapter?

See the [extensibility]({{site.baseurl}}/extensibility) documentation for details on how new Meshery adapters are made.


#### Can I run more than one instance of the same Meshery adapter?
The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the `meshery-istio` adapter. To do so, modify ~/.meshery/meshery.yaml to include multiple copies of the given adapter.

See the "[Multiple Adapters]({{site.baseurl}}/guides/multiple-adapters)" guide for more information.
