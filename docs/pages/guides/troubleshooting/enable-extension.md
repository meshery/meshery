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

If you've set up your local Meshery development environment, but you notice an extension isn't showing up or working as expected, this guide is here to help. Before you start debugging, the most important first step is to **understand the type** of the component you are troubleshooting.

Meshery has different types of extensions that integrate in different ways. This means their failure modes and troubleshooting methods are also very different.

## Extension Types

### 1. Tightly-Coupled Frontend UI Plugins

This scenario describes extensions that are deeply integrated into the Meshery UI and share a close relationship with the Meshery Server (e.g., [Meshery Kanvas](https://docs.meshery.io/extensions/kanvas)).

### 2. Loosely-Coupled Backend Extensions

This scenario describes extensions that run as independent processes and communicate with Meshery Server over a standard API.

**Example: [Meshery Adapters](https://docs.meshery.io/concepts/architecture/adapters)** (like the [`meshery-istio`](https://github.com/meshery-extensions/meshery-istio) adapter)

#### How to Identify extesnion Type

You can identify this type of component by **Codebase Status:** It is **closed-source**. You will not be able to find its source code in a public repository on the GitHub organization.
https://github.com/orgs/meshery-extensions/repositories?type=all 是这个仓库吗？我不太确定耶

## Common Reasons for Failure

If you've identified the component as a tightly-coupled plugin, the reason for failure is almost always one of the following:

* **Permission Issues (The Main Reason):** Because Kanvas is a closed-source project, if your user account does not have the necessary permissions to access its code repository, your local Meshery Server cannot build and load it. 
* **Dependency Package Mismatch:** A tightly-coupled plugin like Kanvas requires that it and the Meshery Server use the **exact same versions of their shared software libraries (packages).** Think of them as two precise gears that must be the same size to work together. If there is even a small version mismatch between your local Meshery Server and the plugin, they will not engage correctly, and the feature will fail to load.

如果你想加入核心团队，查看xxx

如果你有kanvas的权限的话，trouble shooting的步骤可以查看这个forumu回答：https://discuss.layer5.io/t/unable-to-setup-kanvas-locally/6431/2

#### What Are the Next Steps?

The first thing to understand is that this is not a standard technical bug you can fix with code. Trying to debug it with typical methods will not be effective. The correct path forward is to:

* **Understand the Situation:** This is not a standard technical bug you can fix with code.
* **Follow the Contribution Path:** If you are interested in contributing to features like Kanvas, the path involves becoming a more deeply involved community member.
  * Learn about the process for joining the development team. You can find more information in the [Community Handbook](？？？？).

#### Common Reasons for Failure

If an adapter is failing, it is a technical problem that you can debug. This is a technical problem that you can debug and solve yourself. Here are the steps to follow:

1. **Check if the container is running:** 
2. **Check the logs:** 

| **Component Type** | Tightly-Coupled Frontend Plugin | Loosely-Coupled Backend Extension |
| :-------------------- | :------------------------------------------------------------- | :----------------------------------------------------------- |
| **Source Model** | Closed-Source | Open-Source |
| **Integration** | Shared `package` dependencies | Standard `gRPC` API |
| **Log Location** | Browser Developer Console | Docker Container Logs (`docker logs`) |
| **Common Problem** | Permissions, dependency version mismatch | Networking, port conflicts, internal errors |