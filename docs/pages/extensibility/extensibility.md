---
layout: default
title: Extensibility
permalink: extensibility
type: Extensibility
abstract: 'Meshery has an extensible architecture with several different types of extension points.'
# redirect_from:
#   - reference/extensibility
#   - extensibility/
language: en
list: exclude
---

Meshery has an extensible architecture with several different types of extension points.

 <!-- via [adapters]({{site.baseurl}}/extensibility/adapters), different [load generators]({{site.baseurl}}/extensibility/load-generators) and different [providers]({{site.baseurl}}/extensibility/providers). Meshery also offers a REST API. -->

## Extension Points

Meshery is not just an application. It is a set of microservices where the central component is itself called Meshery. Integrators may extend Meshery by taking advantage of designated Extension Points. Extension points come in various forms and are available through Meshery’s architecture.

![Meshery Extension Points]({{site.baseurl}}/assets/img/architecture/meshery_extension_points.svg)

_Figure: Extension points available throughout Meshery_

<!-- 
1. [Adapters]({{site.baseurl}}/extensibility/adapters)
   -  Messaging Framework (CloudEvents and NATS) 
1. [GraphQL API](/extensibility/api#graphql)
1. [Load Generators]({{site.baseurl}}/extensibility/load-generators)
1. [Providers]({{site.baseurl}}/extensibility/providers)
1. [REST API](/extensibility/api#rest)
1. [UI Plugins](extensibility/ui)
1. [Integrations](/extensibility/integrations)
1. [Extensions](/extensibility/extensions) 
-->

## Types of Extension Points

The following points of extension are currently incorporated into Meshery.

{% assign sorted = site.pages | sort: "extensibility" %}

<ul>
    {% for item in sorted %}
    {% if item.type=="Extensibility" and item.list!="exclude" and item.language !="es"  -%}
      <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
      {% if item.abstract != " " %}
              - {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
    {% endfor %}
</ul>

## As an Extension Provider, How to Verify Integration’s Compatibility?

When Meshery updates its Golang version, it’s essential for extension providers to verify that their Golang-based integrations remain compatible. Follow these steps to ensure your plugin works with the updated Golang version.

1. Checkout the [Meshery](https://github.com/meshery/meshery) repository.
2. Update `/go.mod` to new Golang version.
3. Update `/install/Makefile.core.mk`. Change $GOVERSION to new Golang version.
4. Run `make server`.
5. Update your Golang-based plugin to use the new Golang version and Build.
6. Run your upgraded Meshery Server and your upgraded extension together.
7. Check for errors in Meshery Server logs.
8. Validate your plugin’s functionality.
