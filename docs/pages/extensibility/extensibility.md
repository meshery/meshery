---
layout: default
title: Extensibility
permalink: extensibility
type: Extensibility
abstract: 'Meshery architecture is extensible. Meshery provides several extension points for working with different service meshes via <a href="extensibility#adapters">adapters</a>, <a href="extensibility#load-generators">load generators</a> and <a href="extensibility#providers">providers</a>.'
redirect_from:
  - reference/extensibility
  - extensibility/
---

Meshery has an extensible architecture with several extension points. Meshery provides several extension points for working with different service meshes via [adapters]({{site.baseurl}}/extensibility/adapters), different [load generators]({{site.baseurl}}/extensibility/load-generators) and different [providers]({{site.baseurl}}/extensibility/providers). Meshery also offers a REST API.

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

1. [Adapters]({{site.baseurl}}/extensibility/adapters)
   - Messaging Framework (CloudEvents and NATS)
1. [GraphQL API](/extensibility/api#graphql)
1. [Load Generators]({{site.baseurl}}/extensibility/load-generators)
1. Patterns (Open Application Model (OAM))
1. [Providers]({{site.baseurl}}/extensibility/providers)
1. [REST API](/extensibility/api#rest)
1. [UI Plugins](extensibility/ui)