---
layout: default
title: Genişletilebilirlik
permalink: tr/extensibility
type: Reference
abstract: 'Meshery mimarisi genişletilebilir. Meshery, <a href="extensibility#adapters"> adaptörler </a>, <a href="extensibility#load-generators"> yük oluşturucular </a> ve <a href aracılığıyla farklı hizmet ağlarıyla çalışmak için birkaç uzantı noktası sağlar = "genişletilebilirlik # sağlayıcıları"> sağlayıcılar </a>.'
redirect_from: reference/extensibility
language: tr
lang: tr
categories: tr
---

Meshery has an extensible architecture with several extension points. Meshery provides several extension points for working with different service meshes via [adapters](#adapters), different [load generators](#load-generators) and different [providers](#providers). Meshery also offers a REST API.

**Guiding Principles for Extensibility**

The following principles are upheld in the design of Meshery's extensibility.

1. Recognize that different deployment environments have different systems to integrate with.
1. Offer a default experience that provides the optimal user experience.

## Extension Points

Meshery is not just an application. It is a set of microservices where the central component is itself called Meshery. Integrators may extend Meshery by taking advantage of designated Extension Points. Extension points come in various forms and are available through Meshery’s architecture.

![Meshery Extension Points]({{site.baseurl}}/assets/img/architecture/meshery_extension_points.svg)

_Figure: Extension points available throughout Meshery_

The following points of extension are currently incorporated into Meshery.

## Types of Extension Points

1. [Providers]({{site.baseurl}}/extensibility/providers)
1. [Load Generators]({{site.baseurl}}/extensibility/load-generators)
1. [Adapters]({{site.baseurl}}/extensibility/adapters)
1. [REST API](/extensibility/api#rest)
1. [GraphQL API](/extensibility/api#graphql)
