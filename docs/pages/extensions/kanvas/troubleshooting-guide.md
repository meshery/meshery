---
layout: default
title: Kanvas Troubleshooting Guide
permalink: extensions/kanvas/troubleshooting-guide
language: en
abstract: Troubleshooting Kanvas errors
type: extensions
category: kanvas
list: include
---

## Navigating to the Kanvas on local server(using `make server`)

### Error:

![Kanvas Error](/assets/img/kanvas/kanvas-local-error.png)

### Temporary Fix:
```
git remote update
make ui-setup
make ui-build
make server
```

### Cause:

Extensions like Kanvas aren't available in local builds of Meshery Server unless you have separately built that extension.
Kanvas is currently in a closed repository available to community members who have been invited to join our core team.

If you continue contribute to the Meshery project, you will have the chance to be invited to join the core team and gain access to the Kanvas repository.


### How Temporary Fix Works:

Pre-built packages of extensions, like Kanvas are readily available and compatible with each release tag of Meshery. If you checkout and build Meshery locally on a stable release tag or if you checkout and built Meshery locally just after a release (like today's release), you will find that the pre-built package of different extensions download and run.

The further in time (further in code changes) away that your local Meshery build is away from a stable release tag, the lower your chances of a pre-built extension package working. 