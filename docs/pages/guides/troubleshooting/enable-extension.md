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

### Identifying the Extension's Integration Type

Before diving into troubleshooting, you need to identify how your extension connects to Meshery. Here's how to determine its type:

#### Integration Method: Tightly-Coupled vs. Loosely-Coupled

- **Tightly-Coupled Plugins** rely on shared software libraries and package versions with the Meshery Server. These tightly-coupled plugins can be more brittle than other types of extensions, meaning they are highly sensitive to changes in their environment, especially to dependency versions.
  - Example: [Meshery Kanvas](https://docs.meshery.io/extensions/kanvas)
  
- **Loosely-Coupled Extensions** communicate with Meshery Server through standard APIs like gRPC.
  - Example: [Meshery Adapters](https://docs.meshery.io/concepts/architecture/adapters) like the [`meshery-istio`](https://github.com/meshery-extensions/meshery-istio) adapter.

#### Source Code Availability: Open vs. Closed

- **Open-source components** have publicly available source code.
  - You'll find these in public repositories under the [meshery-extensions](https://github.com/orgs/meshery-extensions/repositories?type=all) GitHub organization.
- **Closed-source components** have restricted access to their source code.
  - These repositories are private and require an invitation to access.

> Learn about which repositories require program participation: https://layer5.io/community/handbook/repository-overview

### Troubleshooting by Extension Type

#### Open-Source and Loosely-Coupled Extensions

For open-source extensions like Meshery Adapters, you can use standard technical debugging:

- **Common Issues**: Networking problems, port conflicts, or component-specific errors.
- **Debugging Steps**: Use standard tools like `docker ps` and `docker logs` to identify the root cause of the issue.

#### Closed-Source and Tightly-Coupled Extensions

{% include alert.html type="info" title="Closed-Source Code Issues Cannot Be Resolved" content="Problems with closed-source extensions aren't typical technical bugs you can fix with code. The issue stems from a single root cause: lack of access to the source code." %}

##### The Root Cause: Lack of Source Code Access

The fundamental reason a closed-source extension fails in a local development environment is that you don't not have the necessary permissions to access its private code repository.

This lack of access prevents you from building the extension from source yourself. This leads directly to a critical technical problem:

##### The Result: Unresolvable Dependency Issues

Because you cannot build the extension from source, your only alternative for local use is to rely on a [pre-built package](https://github.com/layer5labs/meshery-extensions-packages). However, this introduces the critical technical challenge of a **Dependency Package Mismatch**.

The requirement for these tightly-coupled plugins is extremely strict: they must use the **exact same packages** as your Meshery server. Think of the plugin and the server as two precise gears that must be identical to work together.

You might find that a pre-built package works on your local machine, but only if it so happens that its dependencies perfectly match the versions in your local Meshery Server environment.

However, this is not a reliable method for development. Any update to your local Meshery instance can break this fragile match, causing the extension to fail to load, often without a clear error message. This is why, without the ability to build a fresh copy from source, these dependency issues are considered practically unresolvable for a consistent development workflow.

##### What Are the Next Steps?

If you want to work with closed-source extensions like Kanvas, the solution is community engagement:

- **Community Contribution Path**: To access features like Kanvas, you'll need to become more involved in the community.
  - The path to joining involves demonstrating your commitment and consistency. A good first step is to contribute to other open-source areas of Meshery for a period of time, showing your interest and growing your expertise within the community. This 'soft' contribution is the key to unlocking the 'hard' access to more advanced projects.

> If you have Kanvas permissions but still encounter problems, check this forum reply for technical troubleshooting steps: https://discuss.layer5.io/t/unable-to-setup-kanvas-locally/6431

