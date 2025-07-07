---
layout: default
title: Troubleshooting Extensions for Local Development
permalink: guides/troubleshooting/enabling-extensions-locally
type: guides
category: troubleshooting
language: en
abstract: A guide to troubleshooting Meshery extensions in a local development environment.
list: include
---

If you've set up your local Meshery development environment and an extension isn't showing up, this guide is here to help. Before you start debugging, the most important first step is to understand the **type** of the component you are troubleshooting.

Meshery has different types of extensions that integrate in different ways. This means their failure modes and troubleshooting methods are also very different.

## Step 1: Locate the Extension's Repository

Most Meshery extensions are open source and their repositories can be found in the [meshery-extensions GitHub organization](https://github.com/orgs/meshery-extensions/repositories?type=all). Your first step is to find the specific repository for the extension you are troubleshooting.

If you can find the repository and it is public, proceed with the steps for [Loosely-Coupled Extensions](#loosely-coupled-extensions).

If you cannot find the repository or it is private, the issue is likely related to a dependency mismatch, a common challenge with [Tightly-Coupled Extensions](#tightly-coupled-extensions).

## Step 2: Follow the Troubleshooting Path

### Loosely-Coupled Extensions

These extensions communicate with Meshery Server through standard APIs. Issues are typically easier to diagnose.
- Example: [Meshery Adapters](/concepts/architecture/adapters) like the [meshery-istio](https://github.com/meshery-extensions/meshery-istio) adapter.
- Common Issues: Networking problems, port conflicts, or component-specific errors.
- Debugging Steps: Use standard tools like `docker ps` to check if the container is running and `docker logs <container-id>` to inspect its logs for errors.

### Tightly-Coupled Extensions

These extensions are highly sensitive to their environment because they require exact package versions to match the Meshery Server.
- Example: [Meshery Kanvas](https://kanvas.new/), a visual designer for Kubernetes and cloud native applications.

{% include alert.html type="info" title="Resolving Dependency Mismatches" content="Problems with tightly-coupled extensions are often **not traditional code bugs** but rather version conflicts between the extension and your local Meshery environment." %}

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

#### Next Steps

Given that you cannot build the extension from source, a solution typically requires seeking assistance.

While you can attempt to experiment by manually testing different versions of packages from the [meshery-extensions-packages](https://github.com/layer5labs/meshery-extensions-packages) repository, this trial-and-error method is not guaranteed to work.

The most reliable path forward is to:
- **Seek assistance:** Ask for guidance in the [Layer5 community](/project/community).
- **Check for requirements:** To check access requirements or find maintainers, consult the [repository overview](https://layer5.io/community/handbook/repository-overview).