---
title: Adapters
description: "Adapters extend Meshery's management capabilities in any number of ways, including lifecycle, configuration, performance, governance, identity..."
aliases:
- /architecture/adapters/
---

## What are Meshery Adapters?

Part of Meshery's extensibility as a platform, Meshery Adapters are purpose-built to address an area in need of management that is either considered optional to the platform and/or is considered an area in which additional depth of control is needed. Adapters extend Meshery's management capabilities in any number of ways, including lifecycle, configuration, performance, governance, identity and so on. Meshery Adapters come in different form factors, and depending on their purpose, deliver different sets or capabilities. Each Adapter registers its capabilities with Meshery Server. Meshery Server, in-turn, exposes those capabilities for you to control.

## Meshery Adapters for Lifecycle Management

Adapters that extend Meshery's lifecycle management capabilities for infrastructure do so, by offering an infrastructure-specific interface to increase the depth of control that Meshery has over a particular technology. Meshery uses adapters to offer choice of load generator (for performance management) and for managing different layers of your infrastructure. Adapters allow Meshery to interface with the different cloud native infrastructure, exposing their differentiated value to users.

Meshery has adapters for managing the following cloud and cloud native infrastructure.

{{< adapters-table >}}

Each adapter is a separately deployable container. Adapters behave largely similarly, but they intentionally differ so they can mold to the infrastructure being managed. That differentiation is why each infrastructure project has its own adapter rather than being driven only through generic Kubernetes APIs.

Individual adapter pages document only what is specific to that infrastructure. How adapters install infrastructure (Day 1) and how Meshery applies ongoing configuration (Day 2) is described here.

### Day 1: Installing infrastructure

When an adapter installs the given infrastructure, it uses that project's own packaging and tooling. In their current form, adapters have a few methods of installation:

1. Helm Chart
2. Kubernetes manifests
3. CLI-based installation
4. Roadmap: OpenTofu / Terraform

When installing infrastructure, each **supported** installation method is tried in the order defined by that adapter. The adapter falls back to subsequent methods if a method fails. Which methods an adapter implements, and in which order it tries them, is an adapter-specific choice. See the adapter's page for what that adapter supports today.

Meshery does not invent a parallel installer when the project already ships Helm charts, manifests, or a CLI. The adapter wraps those mechanisms, streams status back to Meshery Server, and exposes version selection in the UI and CLI.

### Day 2: Ongoing configuration through Designs

Ongoing configuration management of the given infrastructure is performed through the deployment of [Meshery Designs]({{< ref "guides/configuration-management/working-with-designs/_index.md" >}}), not by re-running the Day 1 installer for every change.

1. Meshery Server receives a request to deploy a design.
2. Server itemizes each component in the design.
3. Server consults its internal registry to identify the entity (the **registrant**) that registered the component's model in the first place (for example, a Meshery Adapter).
4. Server parcels out requests for change to each registrant whose components are implicated in the design. That may be one or more adapters, and/or direct communication with one or more Kubernetes clusters.

Adapters are registrants: they register models with Meshery Server so the server knows where to send provisioning requests. Components of those registered models are provisioned by the adapter over gRPC. See [Deployment Engine]({{< ref "concepts/architecture/deployment-engine/index.md" >}}) for how Meshery chooses between an adapter and direct cluster communication.

## Meshery Adapters for Performance Management

Meshery Server allows users to generate traffic load tests using fortio.

## Adapter Deployment and Registration

Like every Meshery component, Meshery Adapters use MeshKit.

When an Adapter registers its models with Meshery Server, it also tells the server where to reach it. That is what makes an Adapter more than a source of definitions: at deployment time, any component belonging to a model an Adapter registered is provisioned *by that Adapter*, over gRPC, rather than by Meshery Server. See [Deployment Engine]({{< ref "concepts/architecture/deployment-engine/index.md" >}}) for how Meshery chooses between the two.

### Adapter FAQs

#### Is each Meshery adapter made equal?

No, different Meshery adapters are written to expose the unique value of each cloud native infrastructure. Consequently, they are not equally capable just as each cloud native infrastructure is not equally capable as the other. Each Adapter has a set of operations which are grouped based on predefined operation types. See the [extensibility]({{< ref "reference/extensibility/_index.md" >}}) page for more details on adapter operations.

#### How do adapters install and configure infrastructure?

Day 1 installation uses the infrastructure project's own Helm charts, Kubernetes manifests, and/or CLI, tried in an adapter-defined order with fallback. OpenTofu / Terraform support is on the roadmap. Day 2 configuration is applied by deploying Meshery Designs: Meshery Server itemizes the design, looks up each component's registrant in its registry, and parcels change requests to the responsible adapter(s) and/or Kubernetes clusters. See [Day 1: Installing infrastructure](#day-1-installing-infrastructure) and [Day 2: Ongoing configuration through Designs](#day-2-ongoing-configuration-through-designs).

#### How can I create a new adapter?

Yes, see the [extensibility]({{< ref "reference/extensibility/_index.md" >}}) documentation for details how to create a new Meshery Adapter. See the Meshery Adapter Template repository as boilerplate for your new adapter.

#### Do adapters have to be written in Golang?

No. Adapters much interface with Meshery Server via gRPC. What language is used in that adapter is the prerogative of a given adapter's maintainers.

#### Can I run more than one instance of the same Meshery adapter?

Yes. The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the `meshery-istio` adapter. To do so, modify ~/.meshery/meshery.yaml to include multiple copies of the given adapter.

See the "[Multiple Adapters]({{< ref "installation/advanced/multiple-adapters.md" >}})" guide for more information.
