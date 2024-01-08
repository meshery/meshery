---
layout: default
title: Concepts
permalink: concepts/logical
language: en
list: exclude
abstract: Concepts for understanding Meshery's various features and components.
---

As a cloud-native management plane, Meshery empowers you with a wide range of tools that provide support for the majority of the systems in the cloud and cloud native ecosystems. Meshery abstracts away the system specific requirements and help you focus on getting things done.

The logical concepts included in Meshery establish a set of constructs with clearly-defined boundaries, each of which is extensible. These contructs set a foundation for the project to build upon and provide a consistent way of relating between multiple components. The logical concepts are:

1. Versioned (see [Schemas](https://github.com/meshery/schemas))
1. Extensible (see [Extension Points](/extensibility)
1. Composable (see Patterns)
1. Portable (see Export/Import)
1. Interoperable (see [Compatibility Matrix](/installation/compatibility-matrix))
1. Configurable (see [Lifecycle Management](/tasks/lifecycle-management))
1. Documented (_you are here_)
1. Testable
1. Maintainable
1. Secure (v0.9.0)
1. Observable (v0.1.0)

Every construct is represented in multiple forms:

- **Schema** (static) - the skeletal structure representing a logical view of the size, shape, characteristics of a construct.
  - *Example: Component schema found in github.com/meshery/schemas*
- **Definition** (static) - An implementation of the Schema containing specific configuration for the construct at-hand.
  - *Example: Component definition generically describing a Kubernetes Pod*
- **Declaration** (static) - A defined construct; A specific deof the Definition.
  - *Example: Component configuration of an NGINX container as a Kubernetes Pod*
- **Instance** (dynamic) - A realized construct (deployed/discovered); An instantiation of the Declaration.
  - *Example: NGINX-as234z2 pod running in cluster*

{% assign sorted_pages = site.pages | sort: "name" %}

## Logical Concepts

<ul>
    {% for item in sorted_pages %}
    {% if item.type=="concepts" and item.language=="en" -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
        -  {{ item.abstract }}
      {% endif %}
      </li>
      {% endif %}
    {% endfor %}
</ul>

[![Meshery Extension Points]({{site.baseurl}}/assets/img/architecture/meshery_extension_points.svg)]({{site.baseurl}}/assets/img/architecture/meshery_extension_points.svg)

_Figure: Extension points available throughout Meshery_