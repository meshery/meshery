---
layout: default
title: Adapters
permalink: concepts/architecture/adapters
type: components
redirect_from: architecture/adapters
abstract: "Adapters extend Meshery's management capabilities in any number of ways, including lifecycle, configuration, performance, governance, identity..."
language: en
list: include
---

## What are Meshery Adapters?

Part of Meshery's extensibility as a platform, Meshery Adapters are purpopse-built to address an area in need of management that is either considered optional to the platform and/or is considered an area in which additional depth of control is needed. Adapters extend Meshery's management capabilities in any number of ways, including lifecycle, configuration, performance, governance, identity and so on. Meshery Adapters come in different form factors, and depending on their purpose, deliver different sets or capabilities. Each Adapter registers its capabilities with Meshery Server. Meshery Server, in-turn, exposes those capabilities for you to control.

## Meshery Adapters for Lifecycle Management

Adapters that extend Meshery's lifecycle management capabilities for infrastructure do so, by offering an infrastructure-specific interface to increase the depth of control that Meshery has over a particular technology. Meshery uses adapters to offer choice of load generator (for performance management) and for managing different layers of your infrastructure. Lifecycle adapters allow Meshery to interface with the different cloud native infrastructure, exposing their differentiated value to users.

Meshery has lifecycle adapters for managing the following cloud native infrastructure.
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Adapter Status | Adapter | Port | Earliest Version supported |
| :------------: | :----------: | :--: | :------------------------: |
{% for adapter in sorted -%}
{% if adapter.project_status -%}
| {{ adapter.project_status }} | <img src="{{ adapter.image }}" style="width:20px" data-logo-for-dark="{{ adapter.white_image }}" data-logo-for-light="{{ adapter.image }}" id="logo-dark-light" loading="lazy"/> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.port }} | {{adapter.earliest_version}} |
{% endif -%}
{% endfor %}

## Meshery Adapters for Performance Management

_v0.8.0 Roadmap:_ The `meshery-perf` adapter externalizes Nighthawk as an Meshery component.

Meshery Server allows users to generate traffic load tests using Nighthawk, fortio, and wrk2. Using the `meshery-perf` adapter, you can schedule, control, and execute performance tests.

Run the `meshery-perf` adapter as an externalized load generator when you: 

1. Need a smaller sized container image for Meshery. Nighthawk binaries are dynamically linked (C++) and they need other dependencies to work. This causes bloat in Meshery Server’s image which doesn’t need them.
1. Need *adaptive load control* of your performance tests, controlling the variability by which the system under test receives load. Use Meshery Server to run adaptive load tests.
1. Need *distributed load testing* and the ability to horizontally scale Nighthawk, using Nighthawk’s execution forwarding service and results sink.

## Adapter Deployment and Registration

Like every Meshery component, Meshery Adapters use MeshKit.

### Adapter FAQs

#### Is each Meshery adapter made equal?

No, different Meshery adapters are written to expose the unique value of each cloud native infrastructure. Consequently, they are not equally capable just as each cloud native infrastructure is not equally capable as the other. Each Adapter has a set of operations which are grouped based on predefined operation types. See the [extensibility]({{site.baseurl}}/extensibility) page for more details on adapter operations.

#### How can I create a new adapter?

Yes, see the [extensibility]({{site.baseurl}}/extensibility) documentation for details how to create a new Meshery Adapter. See the Meshery Adapter Template repository as boilerplate for your new adapter.

#### Do adapters have to be written in Golang?

No. Adapters much interface with Meshery Server via gRPC. What language is used in that adapter is the perogative of a given adapter's maintainers.

#### Can I run more than one instance of the same Meshery adapter?

Yes. The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the `meshery-istio` adapter. To do so, modify ~/.meshery/meshery.yaml to include multiple copies of the given adapter.

See the "[Multiple Adapters]({{site.baseurl}}/guides/multiple-adapters)" guide for more information.
