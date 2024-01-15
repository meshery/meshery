---
layout: default
title: Using Multiple Adapters
permalink: guides/multiple-adapters
type: guides
language: en
abstract: Meshery is capable of running zero or more adapters. Meshery offers many features without the need for adapters. Adapters are optional components that enhance and extend Meshery's core functionality.
---

## Advanced Configuration

Meshery is capable of running zero or more adapters. Meshery offers many features without the need for adapters. Adapters are optional components that enhance and extend Meshery's core functionality.

### Modifying the default adapter deployment configuration

The number of adapters, type of adapters, where they are deployed, how they are named and what port they are exposed on are all configurable deployment options. To modify the default configuration, find `~/.meshery/meshery.yaml` on your system. `~/.meshery/meshery.yaml` is a Docker Compose file.

#### Configuration: Running fewer Meshery adapters

To customize which Meshery Adapters are used in which deployments, customize your contexts in your meshconfig.

*Recommended:*
Configure your meshconfig, using `mesheryctl system context` to customize which Meshery Adapters are used in which deployments.

*Alternative:*
Alternatively, directly modify the `~/.meshery/meshery.yaml` configuration file, remove the entry(ies) of the adapter(s) you are removing from your deployment.

#### Configuration: Running more than one instance of the same Meshery adapter

The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the Meshery Adapter for NGINX Service Mesh. To do so, you can use either of Meshery's clients or to modify your Meshery deployment:
 - Using `mesheryctl`, modify `~/.meshery/meshery.yaml` to include multiple copies of the given adapter.
 - Using Meshery UI, navigate to the Settings page and enter the host and port of your additional adapter.

#### Configuration: Choosing an adapter while installing Mesheryctl

While installing mesheryctl using bash installation script, we can choose which adapter to be loaded.
This is done by passing ADAPTERS environment variable to meshery bash script.

*For e.g.* 
`curl -L https://meshery.io/install | ADAPTERS=consul PLATFORM=kubernetes bash -` installs mesheryctl and starts Meshery Server in your connected Kubernetes cluster deploying only the Meshery Adapter for Consul and not the rest of Meshery's adapters.

<h5>Demo of Meshery managing deployments across multiple Kubernetes clusters:</h5>

<iframe class="container" width="560" height="315" src="https://www.youtube.com/embed/yWPu3vq4vEs?start=5041" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

See on YouTube: [Cloud Native Austin Virtual Meetup: April 2020](https://youtu.be/yWPu3vq4vEs?t=5041&list=PL3A-A6hPO2IOpTbdH89qR-4AE0ON13Zie)
