---
title: Extensibility
description: Meshery has an extensible architecture with several different types of extension points.
categories: [extensibility]
aliases:
- /extensibility
---

Meshery has an extensible architecture with several different types of extension points.


## Extension Points

Meshery is not just an application. It is a set of microservices where the central component is itself called Meshery. Integrators may extend Meshery by taking advantage of designated Extension Points. Extension points come in various forms and are available through Meshery’s architecture.

![Meshery Extension Points](./images/meshery_extension_points.svg)

_Figure: Extension points available throughout Meshery_

Extension points in Meshery come in different shapes with the contracts of their boundaries ranging in their expression and depth of behavioral augmentation allowed; ranging from loosely-defined to strictly validated and enforced. 

## List of Extensions

A browsable collection of various Meshery extensions is available at [https://meshery.io/extensions](https://meshery.io/extensions).

## Types of Extension Points

The following points of extension are currently incorporated into Meshery.

| Extension point | What it extends | When it is applied |
| :--- | :--- | :--- |
| **[Adapters]({{< ref "extensions/adapters/_index.md" >}})** | Per-technology lifecycle, configuration, and performance operations, reached from Meshery Server over gRPC. | Runtime, opt-in |
| **[Authorization keys]({{< ref "reference/extensibility/authorization/index.md" >}})** | The keys, keychains, and roles that gate features in Meshery UI. | Runtime |
| **[Build-time extensibility]({{< ref "reference/extensibility/build-time.md" >}})** | Configuration, data, and packages baked into a custom Meshery container image. | Image build |
| **[Load generators]({{< ref "reference/extensibility/load-generators.md" >}})** | The engines behind Meshery's performance management. | Runtime |
| **[Models and Integrations]({{< ref "extensions/models/_index.md" >}})** | The registry of components and relationships Meshery designs and operates. | Runtime |
forms. | Runtime |
| **[Providers]({{< ref "reference/extensibility/providers/index.md" >}})** | Identity, authorization, durable persistence, and the extension package Meshery loads. | Runtime |
| **[REST and GraphQL APIs]({{< ref "reference/extensibility/api.md" >}})** | Programmatic access for external systems and automation. | Runtime |
| **[Schema annotations]({{< ref "reference/extensibility/schemas.md" >}})** | Model and component behavior expressed through `x-annotations`. | Design time |
| **[UI extension points]({{< ref "reference/extensibility/ui.md" >}})** | Meshery UI: navigator, account, user preferences, collaborator, and RJSF 

When you extend Meshery, also see
[Ensuring Extension Compatibility]({{< ref "reference/extensibility/verify-compatibility.md" >}})
for keeping an extension aligned with the platform version it runs against.

## Security and trust model

Extension points are a capability, not a sandbox. **An extension you enable runs inside the
trust boundary of the deployment that enables it**, and Meshery provides no privilege
separation between an extension and itself:

- UI extension components load into Meshery UI's own browser origin and JavaScript context,
  with access to the DOM, the session cookie, and any API the signed-in user can call.
- A Provider extension package may carry a server-side plugin that is loaded in-process by
  Meshery Server and handed the datastore, the Broker connection, the MeshSync channel, and
  the Kubernetes connection tracker.
- Adapters hold their own cluster credentials and, by default, share the `meshery-server`
  ServiceAccount.

For operators, this means selecting a Remote Provider and enabling an adapter are security
decisions, not just functional ones. The production guidance - what to evaluate before
enabling an extension, how to pin an extension package, how to scope an adapter, and how to
remove one - lives in
[Trusting an extension]({{< ref "installation/production/security-hardening.md#trusting-an-extension" >}})
within the
[Production Deployment]({{< ref "installation/production/_index.md" >}}) set.

For extension authors, the security of an extension published outside the
`meshery-extensions` GitHub organization is the responsibility of its author and of the
operator who enables it. Report a suspected vulnerability through the process in
[Security Vulnerabilities]({{< ref "project/security-vulnerabilities.md" >}}).
