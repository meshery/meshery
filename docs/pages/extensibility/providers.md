---
layout: default
title: "Extensibility: Providers"
permalink: extensibility/providers
type: Extensibility
#redirect_from: architecture/adapters
abstract: "Meshery uses adapters to enrich the level of depth by which it manages cloud native infrastructure."
language: en
list: include
---

Meshery offers Providers as a point of extensibility. With a built-in Local Provider (named "None"), Meshery Remote Providers are designed to be pluggable. Remote Providers offer points of extension to users / integrators to deliver enhanced functionality, using Meshery as a platform. A specific provider can be enforced in a Meshery instance by passing the name of the provider with the env variable PROVIDER.

1. **Extensibility points offer clean separation of Meshery's core functionality versus plugin functionality.**
   - Meshmap is an example of a feature to be delivered via Remote Provider.
1. **Remote Providers should be able to offer custom RBAC, custom UI components, and custom backend components**
   - Dynamically loadable frameworks need to be identified or created to serve each of these purposes.

### Design Principles: Meshery Remote Provider Framework

Meshery's Remote Provider extensibility framework is designed to enable:

1. **Pluggable UI Functionality:**

   - Out-of-tree custom UI components with seamless user experience.
   - A system of remote retrieval of extension packages (ReactJS components and Golang binaries).

1. **Pluggable Backend Functionality:**

   - Remote Providers have any number of capabilities unbeknownst to Meshery.

1. **Pluggable AuthZ**
   - Design an extensible role based access control system such that Remote Providers can determine their own set of controls. Remote Providers to return JWTs with custom roles, permission keys and permission keychains.

![Providers](/assets/img/providers/provider_screenshot.png)

### What functionality do Providers perform?

What a given Remote Provider offers might vary broadly between providers. Meshery offers extension points that Remote Providers are able to use to inject different functionality - functionality specific to that provider.

- **Authentication and Authorization**
  - Examples: session management, two factor authentication, LDAP integration.
- **Long-Term Persistence**
  - Examples: Storage and retrieval of performance test results.
  - Examples: Storage and retrieval of user preferences.
- **Enhanced Visualization**
  - Examples: Creation of a visual service mesh topology.
  - Examples: Different charts (metrics), debug (log viewer), distributed trace explorers.
- **Reporting**
  - Examples: Using Meshery's GraphQL server to compose new dashboards.

## Types of providers

Two types of providers are defined in Meshery: `local` and `remote`. The Local provider is built-into Meshery. Remote providers are may be implemented by anyone or organization that wishes to integrate with Meshery. Any number of Remote providers may be available in your Meshery deployment.

### Remote Providers

Use of a Remote Provider, puts Meshery into multi-user mode and requires user authentication. Use a Remote provider when your use of Meshery is ongoing or used in a team environment (used by multiple people).

Name: **“Meshery”** (default)

- Enforces user authentication.
- Long-term term persistence of test results.
- Save environment setup.
- Retrieve performance test results.
- Retrieve conformance test results.
- Free to use.

### Local Provider

Use of the Local Provider, "None", puts Meshery into single-user mode and does not require authentication. Use the Local provider when your use of Meshery is intended to be shortlived.

Name: **“None”**

- No user authentication.
- Container-local storage of test results. Ephemeral.
- Environment setup not saved.
- No performance test result history.
- No conformance test result history.
- Free to use.

## How to Build a Remote Provider

See [Building a Remote Provider]({{site.baseurl}}/extensibility/providers/building-a-remote-provider}}) for more information.