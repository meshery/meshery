---
layout: default
title: Concepts
permalink: concepts/logical
language: en
list: exclude
abstract: Concepts for understanding Meshery's various features and components.
---

Concepts for understanding Meshery's various features and components.

As a cloud-native management plane, Meshery is able to empower its users with a wide range of tools that provide support for the majority of the systems in the cloud and cloud native ecosystems. AMeshery abstracts away the system specific requirements and help the users focus on getting things done.

The logical concept include in Meshery establish a set of constructs with clearly-defined boundaries and vigil eye to extensibility. These logical constucts set a foundation for the project to build upon and provide a consistent way of relating between multiple components. The logical concepts are:

1. Versioned
1. Extensible
1. Composable
1. Portable
1. Interoperable
1. Configurable
1. Documented
1. Testable
1. Maintainable
1. Secure (v0.9.0)
1. Observable (v0.1.0)

Every construct will be represented in multiple forms.

**Schema** (static) - the skeletal structure representing a logical view of the size, shape, characteristics of a construct.
*Example: Component schema found in github.com/meshery/schemas*
**Definition** (static) - An implementation of the Schema containing specific configuration for the construct at-hand.
*Example: Component definition generically describing a Kubernetes Pod*
**Declaration** (static) - A defined construct; A specific deof the Definition.
*Example: Component configuration of an NGINX container as a Kubernetes Pod*
**Instance** (dynamic) - A realized construct (deployed/discovered); An instantiation of the Declaration.
*Example: NGINX-as234z2 pod running in cluster*

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