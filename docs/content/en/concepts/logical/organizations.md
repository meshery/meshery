---
title: Organizations
description: Organizations are the unit of tenancy in Meshery — they group users together and own all of the resources those users create.
aliases:
- /concepts/organizations/
---

An **Organization** is the unit of tenancy in Meshery. Organizations group users together and **own all of the resources those users create** — Workspaces, Environments, Designs, Connections, and more. Where a [Workspace]({{< ref "concepts/logical/workspaces.md" >}}) is your team's collaboration space, the Organization is the top-level container that everything else lives inside of, and the boundary that keeps one tenant's resources isolated from another's.

{{% alert color="dark" title="The unit of tenancy" %}}
Organizations are the outermost boundary of ownership and access in Meshery. When using a Remote Provider, every Workspace, Environment, Design, and Connection belongs to exactly one Organization.
{{% /alert %}}

## Organizations and Providers

Whether you have Organizations at all depends on which [Provider]({{< ref "reference/extensibility/providers/index.md" >}}) Meshery is connected to:

- With the **Local** provider, Meshery runs in single-user mode. There is no authentication and no multi-tenancy, so Organizations, teams, and shared ownership do not come into play.
- With a **Remote Provider** (such as [Meshery Cloud](https://cloud.meshery.io)), Meshery runs in multi-user mode. The Remote Provider supplies identity, and **Organizations become the structure through which users, ownership, and access control are managed.**

In other words, Organizations are a capability that a Remote Provider extends Meshery with. Different Remote Providers can offer richer or simpler organization models — hierarchical organizations, teams as user groups, fine-grained roles — through Meshery's [provider extensibility]({{< ref "reference/extensibility/providers/index.md" >}}) framework.

## Key Features

- **Tenancy boundary** — an Organization isolates its resources from every other Organization. Members of one Organization cannot see or manage another's resources unless access is explicitly shared.
- **Resource ownership** — Workspaces, Environments, Designs, and Connections are all owned by an Organization, not by an individual user, so work survives changes in team membership.
- **Membership** — users belong to one or more Organizations. A user's permissions in one Organization are independent of their permissions in another.
- **Access control** — with a Remote Provider, Organizations carry roles and permissions that govern what each member may do with the Organization's resources.
- **Identity and branding** — a Remote Provider can extend an Organization with its own identity provider, custom domain, and branding (see [With a Remote Provider](#organizations-with-a-remote-provider) below).

## Where Organizations fit

Organizations sit at the top of Meshery's logical hierarchy:

- An **Organization** groups **users** (and, with a Remote Provider, **teams**).
- The Organization **owns** its **Workspaces**, **Environments**, **Designs**, and **Connections**.
- **Teams** are granted access to **Workspaces**.
- **Workspaces** bring together **Environments** (groupings of **Connections**) and the **Designs** deployed against them.

So a typical path from the top down is: *Organization → Team → Workspace → Environment → Connection*, with Designs deployed into the infrastructure a Workspace can reach.

## Key Relationships

### Teams

- Teams are groups of users **within** an Organization, used to grant access to Workspaces and their resources.
- Teams are a Remote Provider extension; the depth of team and role support varies by provider.

### Workspaces

- Workspaces are owned by an Organization and serve as the collaboration point where teams do their work.

See "[Workspaces]({{< ref "concepts/logical/workspaces.md" >}})" for more information.

### Environments

- Environments group Connections (Kubernetes clusters, Prometheus instances, and so on) and are owned by an Organization.

See "[Environments]({{< ref "concepts/logical/environments.md" >}})" for more information.

### Designs

- Designs — reusable, declarative characterizations of your infrastructure — are owned by an Organization and deployed within the context of a Workspace.

See "[Designs]({{< ref "concepts/logical/designs.md" >}})" for more information.

### Connections

- Connections are the managed or discovered infrastructure resources (Kubernetes clusters, Prometheus instances, and so on) that Meshery works with. They are grouped by Environments and owned by an Organization.

See "[Connections]({{< ref "concepts/logical/connections/index.md" >}})" for more information.

## Organizations with a Remote Provider

A Remote Provider can extend an Organization well beyond simple grouping. Using [Meshery Cloud](https://cloud.meshery.io) as the reference implementation, an Organization can be configured with:

- **Its own identity provider** — email-and-password and social sign-in (Google, GitHub), or the Organization's own OAuth/OIDC single sign-on.
- **A custom domain and branding** — its own URL, logo, colors, and login screen ("white-labeling").
- **Roles, teams, and fine-grained permissions** — extensible RBAC returned by the Remote Provider as part of the user's token.
- **Controlled membership** — open self-service sign-up or invitation-only joining, optionally restricted to an email domain.

How those choices combine differs based on your connected Remote Provider. Experiences vary from a shared, hosted experience to a fully white-labeled deployment on the Organization's own domain. See your respective Remote Provider's documentation.

## Best Practices

- Use separate Organizations to isolate tenants — distinct customers, partners, or business units — so their resources and access never overlap.
- Keep ownership at the Organization level (rather than relying on individual users) so work persists as team membership changes.
- Depending upon your Remote Provider, you may assign access through teams and roles rather than per-user, and potentially customize your organization's branding and identity.
- When graduating from the Local provider to a Remote Provider, create an Organization first so resource ownership is established from the start.

Organizations give Meshery a clear, isolated home for every team's work. Paired with a Remote Provider, they become the foundation for identity, access control, and a branded, multi-tenant experience.
