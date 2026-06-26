---
title: Overview
display_title: false
categories: [project]
description: Meshery is the self-service engineering platform, enabling collaborative design and operation of cloud and cloud native infrastructure.
aliases:
- /project/overview/
- /project/meshery-overview/
- /getting-started/overview
weight: -10
---

## The open source manager for the infrastructure you run

If you build and run the internal platform your developers ship on, you already know the shape of the problem: many clusters, more than one cloud, and a long tail of cloud native tools that each speak their own dialect. [Meshery](https://meshery.io) is the open source, cloud native manager that gives you one place to design, operate, and characterize all of it.

It is director-level software. In the stack, your orchestrator and its controllers are the individual contributors that do the hands-on work of running workloads and reconciling state. Meshery sits a level above, where a director or executive operates: it sets desired state, coordinates work across teams and environments, enforces standards and policy, and holds a single view of the whole operation, while delegating execution to the layers below.

As a [Cloud Native Computing Foundation](https://www.cncf.io) (CNCF) project, Meshery is a self-service engineering platform. That is the job you are already trying to do: provide paved, governed, self-service infrastructure to the teams you support, without becoming the ticket queue in the middle. Just as important, Meshery is built to bring those engineers together, treating infrastructure as a shared design that teams review and refine the way peers collaborate in a Google Doc.

## What Meshery is to you

Meshery is the management plane you put in front of heterogeneous, Kubernetes-based infrastructure. Underneath, it delegates the actual deployment to your orchestrator, such as Kubernetes. On top, it gives you and your users a consistent way to design, deploy, discover, and operate infrastructure across every cluster and cloud you manage.

What separates Meshery from a pile of scripts and dashboards is its model-driven core. Meshery maintains a registry of [Models]({{< ref "concepts/logical/models/index.md" >}}), Components, Relationships, and Policies, which gives it a semantic understanding of your infrastructure: how components relate, which configurations are valid, and which policies apply in a given context. That understanding is what powers relationship inference, context-aware design, and validation, instead of treating everything as opaque YAML.

Governance and coordination live here too. Meshery's [Policy Engine](https://docs.meshery.io/concepts/logical/policies/) evaluates relationships and enforces context-aware rules as designs take shape, and a Workflow Engine on the roadmap will let Meshery orchestrate and sequence operations across your infrastructure. Both are director-level work: deciding and coordinating what should happen, while the runtime below carries it out. Adding orchestration does not make Meshery a runtime.

## Why you should care

Meshery is worth your attention if any of this sounds like your week:

- **You are paying an integration tax.** Every new tool means more glue code, more bespoke dashboards, and more one-off automation to keep in your head. Meshery's 300+ integrations and single management model absorb much of that.
- **You are drowning in YAML.** Meshery offers visual and collaborative GitOps, so infrastructure can be designed visually and still version-controlled as code.
- **Changes land as diffs nobody can really review.** Meshery turns a change into a shared, visual design that your team reviews together, the way peers comment on and approve a Google Doc, rather than squinting at a raw YAML pull request.
- **Day 2 is where you actually live.** Provisioning is the easy part. Meshery is built for the long tail: discovery and onboarding of what already exists, ongoing configuration management, drift awareness, and performance characterization across releases.
- **Your infrastructure is heterogeneous and it is not getting simpler.** One cluster or many, one cloud or several, Meshery manages them from a single control point.
- **You need self-service with guardrails.** Workspaces, multi-tenancy, and role-based access control let you hand teams isolated, governed access to shared infrastructure rather than fielding their requests one at a time.

## What makes Meshery different

Plenty of tools cover one slice of this. Meshery's value is in combining them on a model-driven foundation:

- **Built-in collaboration.** Meshery is built to bring engineers together, not just to manage machines. Designs are shared artifacts your team works on together, with collaborative design reviews that run the way peers comment on and approve a Google Doc. Infrastructure becomes a conversation among people rather than an opaque diff.
- **A semantic, model-driven core.** Meshery does not just template manifests. Its registry of Models, Relationships, and Policies lets it reason about how components fit together, infer relationships, and apply context-aware guardrails as designs are built.
- **Design, operate, and characterize in one platform.** Meshery spans Day 0 through Day 2. You can adopt, design, provision, discover, operate, and benchmark from one place, with load generation and metrics through Prometheus and Grafana, rather than stitching a separate tool to each step.
- **Breadth without the glue.** 300+ integrations across cloud and cloud native technologies, managed through one consistent model.
- **Open source, customizable, and vendor-neutral.** As a CNCF project, Meshery is a composable building block, not a point of lock-in. It is deeply customizable across branding, identity providers, multi-tenancy, and deployment model, and adapts to internal platforms, branded service-provider offerings, and multi-tenant enterprise deployments alike.

## What is in scope, and what is not

A clear boundary is part of the value, and it follows the same director-level logic: Meshery directs and coordinates, but it does not do the hands-on execution itself. Here is what the project does and does not set out to do.

### In scope

- Lifecycle management of cloud and cloud native infrastructure: provisioning, discovery and onboarding of existing clusters and workloads, and ongoing Day 2 operation
- A model-driven core: a registry of Models, Components, Relationships, and Policies that gives Meshery a semantic understanding of infrastructure and enables context-aware design and validation
- A Policy Engine for relationship inference, validation, and context-aware governance, with a Workflow Engine for orchestration on the roadmap to coordinate and sequence operations
- Built-in collaboration: shared designs, collaborative design reviews, workspaces, multi-tenancy, and role-based access control for teams sharing infrastructure
- Multi-cluster and multi-cloud management from a single control point
- Visual and code-based (GitOps) design of infrastructure and applications
- Performance characterization and load generation, with metrics through Prometheus and Grafana
- 300+ integrations across cloud and cloud native technologies, plus open source extension points for building on top of the platform

### Out of scope

- **A runtime.** Meshery directs and orchestrates work through its Policy Engine and, in the future, a Workflow Engine, then delegates execution to the underlying system, such as Kubernetes. It does not run your workloads or serve traffic itself, it is not a Kubernetes distribution, and it does not install your clusters for you.
- **A data plane.** Meshery is a manager, not a proxy. It is not itself a service mesh, ingress, or load balancer; it designs, deploys, and operates the ones you run.
- **A CI/CD or source-control system.** Meshery complements your pipelines and GitOps workflow. It is not a build, test, or version-control system.
- **An observability backend.** Meshery integrates with Prometheus and Grafana rather than replacing your metrics store, logging, or tracing stack.

## Extensions and customization

Meshery is built to be extended, and its extension points are open source, spanning the UI, APIs, authorization, providers, and schemas. Meshery ships with its own UI and design configurator, and a broad ecosystem of extensions, adapters, and models builds on top of it. The extensions available on Meshery take the platform well beyond its core:

- **Visual, collaborative design and operations.** A visual designer lets teams build infrastructure by dragging from a palette of thousands of versioned Kubernetes and cloud components, with context-aware relationships, freestyle whiteboarding, real-time multi-user collaboration, inline comments and design reviews, and sharing, embedding, and publishing of designs. An operations mode brings every cluster under one roof: validate, dry-run, deploy, clone, and undeploy designs, apply patterns, stream logs, inspect live resources, and open an interactive terminal into pods and containers to debug in real time. Work stays Git-connected, with GitHub Actions and snapshots that turn changes into human-verifiable visual diffs.
- **Identity, access, and multi-tenancy.** A centralized console and identity provider adds a flexible authorization framework with granular role-based access control across organizations, teams, and users, tenant entitlement, and a service-provider-grade organizational hierarchy that scales from a single team to many tenants.
- **A content catalog.** Publish, discover, and reuse cloud native architectures and design patterns, either publicly or privately within your organization.
- **Complete customization.** Meshery is far more than a foundation for an internal developer platform. It can be white-labeled and rebranded with your own identity and custom domain, connected to your own identity providers, tailored through customizable webhooks and a flexible authorization model, and run either as a managed service or fully self-hosted with Helm. Branding, identity, tenancy, and deployment are yours to control, which is what lets service providers and enterprises run Meshery as a branded, multi-tenant platform of their own.

Together, these let you run Meshery not only as an internal platform, but as a branded, multi-tenant offering that you operate as your own.

## Learn Meshery

Meshery includes academies, a hands-on learning platform built into the project. Academies organize material into learning paths such as Mastering Meshery, scenario-based challenges where you deploy and configure real infrastructure, and certifications like the [Certified Meshery Contributor](https://meshery.io/community/certifications) (CMC) program, complete with badges and shareable credentials. Because courses embed live Meshery designs, you practice on real workflows rather than only reading about them. Academies are open source, extensible, and white-labelable, so your organization can publish branded learning paths of its own. See the [academies documentation](https://docs.meshery.io/extensions/academies/) to get started.

## Try Meshery and join the community

Try Meshery in your browser with the [Cloud Native Playground](https://play.meshery.io), with no installation required. When you are ready to run it against your own clusters, the [installation guides](https://docs.meshery.io/installation) cover Docker, Kubernetes, Helm, and more. To get involved, start with the [Community Handbook](https://meshery.io/community) and join the conversation in [Slack](https://slack.meshery.io) or the [discussion forum](https://discuss.meshery.io/).
