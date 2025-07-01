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

- **Tightly-Coupled Plugins** rely on shared software libraries and package versions with the Meshery Server.
  - These tightly-coupled plugins can be more brittle than other types of extensions, meaning they are highly sensitive to changes in their environment, especially to dependency versions.
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

{% include alert.html type="info" title="闭源代码问题不能解决" content="Problems with closed-source, tightly-coupled extensions aren't typical technical bugs you can fix with code. Standard debugging methods won't be effective here." %}

If you've identified the component as a tightly-coupled plugin, the reason for failure is almost always one of the following:

- **Permission Issues:** Because extensions like Kanvas can be closed-source, your user account may not have the necessary permissions to access its code repository. Without access, your local Meshery Server cannot build and load it.

- **Dependency Package Mismatch:** This is the most critical technical reason for failure for tightly-coupled extensions.

  For these plugins, you must use the **exact same packages** in your Meshery server as in your plugins. This requirement is extremely strict—think of the plugin and the server as two precise gears that must be identical to work together. Even a minor version difference in a single shared library can cause the entire extension to fail to load, often without a clear error message.

  This package mismatch issue is especially relevant when dealing with pre-built packages. While for some extensions you can download them from the [Layer5 Labs Meshery Extensions Packages repository](https://github.com/layer5labs/meshery-extensions-packages) to save time, you must be extremely careful. **These pre-built packages must still perfectly match the package versions used by your local Meshery server to work correctly.**

  When they don't match, the extension will likely never load successfully.

##### What Are the Next Steps?

If you want to work with closed-source extensions like Kanvas, the solution is community engagement:

- **Community Contribution Path**: To access features like Kanvas, you'll need to become more involved in the community.
  - The path to joining involves demonstrating your commitment and consistency. A good first step is to contribute to other open-source areas of Meshery for a period of time, showing your interest and growing your expertise within the community. This 'soft' contribution is the key to unlocking the 'hard' access to more advanced projects.

> If you have Kanvas permissions but still encounter problems, check this forum reply for technical troubleshooting steps: https://discuss.layer5.io/t/unable-to-setup-kanvas-locally/6431

