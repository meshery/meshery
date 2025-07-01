---
layout: default
title: Enabling Extensions for Local Development
permalink: guides/troubleshooting/enabling-extensions-locally
type: guides
category: troubleshooting
language: en
abstract: A guide on how to build and enable Meshery extensions for use in a local development environment.
list: include
---

If you've set up your local Meshery development environment and an extension isn't showing up, this guide is here to help. Before you start debugging, the most important first step is to understand the **type** of the component you are troubleshooting.

Meshery has different types of extensions that integrate in different ways. This means their failure modes and troubleshooting methods are also very different.

## Identifying the Extension's Integration Type

Before diving into troubleshooting, you need to identify how your extension connects to Meshery. Here's how to determine its type:

### Integration Method: Tightly-Coupled vs. Loosely-Coupled

- **Loosely-Coupled Extensions:**
  - Communicate with Meshery Server through standard APIs like gRPC.
  - Example: [Meshery Adapters](https://docs.meshery.io/concepts/architecture/adapters) like the [meshery-istio](https://github.com/meshery-extensions/meshery-istio) adapter.

- **Tightly-Coupled Plugins:**
  - Rely on shared software libraries and exact package versions with the Meshery Server.
  - Are highly sensitive to changes in their environment, especially to dependency versions.
  - Example: [Meshery Kanvas](https://kanvas.new/), a visual designer for Kubernetes and cloud native applications. 


### Source Code Availability: Open vs. Closed

- **Open-source components** have publicly available source code.
  - You'll find these in public repositories under the [meshery-extensions](https://github.com/orgs/meshery-extensions/repositories?type=all) GitHub organization.
- **Closed-source components** have restricted access to their source code.
  - These repositories are private and require an invitation to access.

> Learn about which repositories [require program participation](https://layer5.io/community/handbook/repository-overview).

## Troubleshooting by Extension Type

### Open-Source and Loosely-Coupled Extensions

For open-source extensions like Meshery Adapters, you can use standard technical debugging:

- **Common Issues**: Networking problems, port conflicts, or component-specific errors.
- **Debugging Steps**: Use standard tools like `docker ps` and `docker logs` to identify the cause of the issue.

### Closed-Source and Tightly-Coupled Extensions
{% include alert.html type="info" title="Closed-Source Issues: No Code Fix" content="Problems with these extensions aren't typical code bugs. The core issue is lack of source code access, which prevents standard debugging." %}

#### Why you can't fix this locally

You don't have permission to access the private code repository. This means:
- You **cannot build the extension from source**
- Your only option is using a [pre-built package](https://github.com/layer5labs/meshery-extensions-packages)

#### The dependency mismatch problem

Tightly-coupled plugins require **exact package matches** between the extension and your Meshery server. When using pre-built packages:
- They might work if dependencies happen to match your local environment
- Any update to your Meshery server will likely break compatibility
- Failures often occur without clear error messages

> Without source access, these dependency issues are practically unresolvable.

#### How to move forward

If you want to work with extensions like Kanvas:
1. Engage with the community 
2. Contribute consistently to Meshery's open-source areas
3. Gain access through demonstrated commitment

> Already have Kanvas permissions? See this [forum guide](https://discuss.layer5.io/t/unable-to-setup-kanvas-locally/6431) for specific fixes.
