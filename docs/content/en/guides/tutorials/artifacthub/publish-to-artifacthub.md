---
title: Publishing Meshery Designs to Artifact Hub
model: artifacthub
params:
    kind: design
categories: [tutorials]
description: Step by step example for how to export Meshery Designs and publish them to an Artifact Hub repository.
aliases:
- /guides/tutorials/publish-to-artifacthub
---

In this tutorial, we'll see how to export a design from Meshery, which we will use to populate an Artifact Hub repository.

## Prerequisites

1. Access to Meshery. [Install Meshery](/installation) or access the [Meshery Playground](https://play.meshery.io).
2. An Artifact Hub repository.

## Steps

1. Navigate to Designs
- Locally: http://localhost:9081/configuration/designs?view=grid
- Playground: https://playground.meshery.io/configuration/designs?view=grid

![](/guides/tutorials/images/publish-to-artifacthub/designs.png)

2. Click to download your design
![](/guides/tutorials/images/publish-to-artifacthub/design-download.png)

3. Prepare your Artifact Hub repository

You will need to have an Artifact hub repository already created with `Kind` as `Meshery Designs`. See [Artifact Hub documentation](https://artifacthub.io/docs/topics/repositories/meshery-designs/) for more information on managing repositories.

### Push Design to Artifact Hub repository

At this point, you should have downloaded your design as a `tar` archive. This archive contains a second archive which holds the metadata files to publish to your Artifact Hub repository. Expand the first `tar` archive and locate the resulting files. Expand the second archive you find there and you will have a folder containing two files: `artifacthub-pkg.yml` and `design.yml`. Move these files to your prepared Artifact Hub repository and push it upstream.

### Verify repository in Artifact Hub

Once the files are pushed to the Artifact Hub repo you will need to wait until Artifact Hub indexes it. You can verify the status of the repository in the Artifact Hub control panel.
