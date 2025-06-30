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

If you've set up your local Meshery development environment and an extension isn't showing up or working as expected, this guide is here to help. Before you start debugging, the most important first step is to understand the **type** of the component you are troubleshooting.

Meshery has different types of extensions that integrate in different ways. This means their failure modes and troubleshooting methods are also very different.

### Identifying the Extension's Integration Type

Before diving into troubleshooting, you need to identify how your extension connects to Meshery. Here's how to determine its type:

#### Integration Method: Tightly-Coupled vs. Loosely-Coupled

- **Tightly-Coupled Plugins** rely on shared software libraries and package versions with the Meshery Server
  - These tightly-coupled plugins can be more brittle than other types of extensions, meaning they are highly sensitive to changes in their environment, especially to dependency versions
  - Example: [Meshery Kanvas](https://docs.meshery.io/extensions/kanvas)
- **Loosely-Coupled Extensions** communicate with Meshery Server through standard APIs like gRPC
  - Example: [Meshery Adapters](https://docs.meshery.io/concepts/architecture/adapters) like the [`meshery-istio`](https://github.com/meshery-extensions/meshery-istio) adapter

#### Source Code Availability: Open vs. Closed

- **Open-source components** have publicly available source code
  - You'll find these in public repositories under the [meshery-extensions](https://github.com/orgs/meshery-extensions/repositories?type=all) GitHub organization
- **Closed-source components** have restricted access to their source code
  - 这些仓库是private的，需要收到邀请才可以获得

### Troubleshooting by Extension Type

#### Closed-Source Extensions

{% include alert.html type="info" title="Understanding the Issue" content="Closed-source extension problems aren't typical technical bugs you can fix with code. Standard debugging methods won't be effective here." %}

除了没办法本地布置以外，If you want to working with closed-source extensions like Kanvas, the solution is community engagement:

- **Community Contribution Path**: To access features like Kanvas, you'll need to become more involved in the community
  - The path to joining involves demonstrating your commitment and consistency. A good first step is to contribute to other open-source areas of Meshery for a period of time, showing your interest and growing your expertise within the community. This 'soft' contribution is the key to unlocking the 'hard' access to more advanced projects

> Learn about 哪些仓库需要program participation required：https://layer5.io/community/handbook/repository-overview

{% include alert.html type="note" title="Kanvas Access Issues" content="If you have Kanvas permissions but still encounter problems, check this forum discussion: https://discuss.layer5.io/t/unable-to-setup-kanvas-locally/6431" %}

> 如果你有kanvas权限，但是仍然遇到了问题，查看这个论坛的回复：https://discuss.layer5.io/t/unable-to-setup-kanvas-locally/6431

#### Open-Source Extensions

For open-source extensions like Meshery Adapters, you can use standard technical debugging:

- **Common Issues**: Networking problems, port conflicts, or component-specific errors
- **Debugging Steps**: Use standard tools like `docker ps` and `docker logs` to identify the root cause

### Package Mismatch Issues

{% include alert.html type="warning" title="Critical: Package Version Matching" content="For tightly-coupled plugins like Kanvas, you must use the exact same packages in your Meshery server as in your plugins. This requirement is extremely strict - think of the plugin and server as two precise gears that must be identical to work together. Even a minor version difference in a single shared library can cause the entire extension to fail to load, often without a clear error message." %}

When package versions don't match between your Meshery server and tightly-coupled plugins, the extension will never load successfully. This is one of the most common causes of extension failures in local development environments.

### Pre-built Packages

For some extensions, pre-built packages are available for download from the Layer5 Labs Meshery Extensions Packages repository. These packages can save you from having to build extensions from source code, but they must still match your Meshery server's package versions exactly.