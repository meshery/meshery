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

When running a local instance of Meshery, you might try to navigate to an extension and find that it won't load. This is expected behavior for some setups.

[![Kanvas Error](/assets/img/trouble-shooting/local-error.png)](/assets/img/trouble-shooting/local-error.png)

### Why This Happens

Your local development environment, started with `make server`, is primarily set up to run the core Meshery components. Many extensions, especially those with complex user interfaces, are developed as separate projects. To keep your local setup fast and lean, these extensions aren't automatically built and included.

You need to explicitly tell your environment to build the extension or integrate a pre-compiled version of it.

### How to Enable an Extension

There are generally two ways to get an extension running locally, depending on whether it's open-source or closed-source:

- **Building from Source:** For most open-source extensions, you can clone their source code repository and run a build command. This compiles the extension from scratch and integrates it into your Meshery UI.
- **Integrating a Pre-built Package:** For some extensions (like Kanvas) or simply for convenience, you can run a command that fetches a ready-to-use, pre-compiled package and plugs it into your local build.

Let's use Kanvas as an example of this second method.

### Example: Enabling the Kanvas Extension

[Kanvas](https://kanvas.new/) is a powerful but closed-source extension. This means you can't build it from source code yourself. Instead, you'll integrate its pre-built package.

> If you continue to contribute to the Meshery project, you will have the chance to be invited to join the core team and gain access to the Kanvas repository.

#### The Fix for Kanvas

If you can't access Kanvas locally, run the following commands from the root of your Meshery repository. The order is important.

1. **`git remote update`**
2. **`make ui-setup`**
 - Prepares your UI development environment and dependencies.
3. **`make ui-build`**
 - This is the key step. It builds the main Meshery UI and, in the process, fetches and integrates the pre-built Kanvas package.
4. **`make server`**
 - Restarts your local Meshery server with the newly included extension.

Once these steps are complete, refresh your browser. The Kanvas extension should now be available.

#### How This Fix Works

The Meshery project provides stable, pre-built packages for extensions like Kanvas that are compatible with each official release. The `make ui-build` command is smart enough to find the correct version of the Kanvas package that matches your local Meshery version, download it, and integrate it for you.

For the best results, try to keep your local `master` branch  up-to-date with the main repository. If your local code has diverged too far from a recent release, you might encounter compatibility issues with the pre-built package.