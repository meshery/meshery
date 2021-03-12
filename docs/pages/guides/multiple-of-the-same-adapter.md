---
layout: default
title: Using Multiple Adapters
permalink: guides/multiple-adapters
type: Guides
---

## Advanced Configuration

Meshery is capable of running zero or more service mesh adapters. Without any service mesh adapters, some but not all of Meshery's features will continue to function (e.g. performance testing of workloads not running on a service mesh).

### Modifying the default adapter deployment configuration
The number of adapters, type of adapters, where they are deployed, how they are named and what port they are exposed on are all configurable deployment options. To modify the default configuration, find *~/.meshery/meshery.yaml* on your system. *~/.meshery/meshery.yaml* is a Docker Compose file.

#### Configuration: Running fewer Meshery adapters
In the *~/.meshery/meshery.yaml* configuration file, remove the entry(ies) of the adapter(s) you are removing from your deployment.

#### Configuration: Running more than one instance of the same Meshery adapter

The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the *meshery-istio* adapter. To do so, modify *~/.meshery/meshery.yaml* to include multiple copies of the given adapter.

Demo of Meshery managing service mesh deployments across multiple clusters:

<iframe class="container" width="560" height="315" src="https://www.youtube.com/embed/yWPu3vq4vEs?start=5041" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

See on YouTube: [Cloud Native Austin Virtual Meetup: April 2020](https://youtu.be/yWPu3vq4vEs?t=5041&list=PL3A-A6hPO2IOpTbdH89qR-4AE0ON13Zie)