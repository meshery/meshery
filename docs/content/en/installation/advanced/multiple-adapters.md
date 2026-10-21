---
title: Using Multiple Adapters
aliases: 
- /guides/multiple-adapters
categories: [installation]
description: Choose which Meshery Adapters to run, run more than one instance of an adapter, or run none at all.
---

## Advanced Configuration

Meshery can run zero or more adapters. Many Meshery features need no adapter at all. Adapters are optional components that enhance and extend Meshery's core functionality.

### Modifying the default adapter deployment configuration

You can configure how many adapters run, which ones, where they run, what they are named, and which port each one listens on. On Docker, `mesheryctl` writes the deployment to `~/.meshery/meshery.yaml`, a Docker Compose file. It downloads that file again and rewrites its services from your current `mesheryctl` context every time you run `mesheryctl system start`, so manual edits to it do not last. Configure adapters through your context instead.

#### Configuration: Running fewer Meshery adapters

Each `mesheryctl` context in your meshconfig (`~/.meshery/config.yaml`) lists the adapters, which `mesheryctl` calls components, to deploy with Meshery. To choose a different set, create a context that lists only the adapters you want, and make it current:

{{< code code="mesheryctl system context create my-context --components meshery-istio --components meshery-linkerd --set" >}}

`mesheryctl` can deploy `meshery-istio`, `meshery-linkerd`, `meshery-consul`, `meshery-nsm`, `meshery-kuma`, `meshery-traefik-mesh`, and `meshery-cilium`. To run no adapters, create a context without `--components`.

#### Configuration: Running more than one instance of the same Meshery adapter

You can run more than one instance of the same Meshery Adapter, for example two instances of the Meshery Adapter for Istio. A `mesheryctl` context deploys at most one instance of each adapter, so run the additional instance yourself, then register its host and port with Meshery Server in one of these ways:
- In Meshery UI, go to **Lifecycle** > **Adapters** and enter the host and port of the additional adapter.
- Add its `host:port` to Meshery Server's `ADAPTER_URLS` environment variable. See [Server Environment Variables]({{< ref "installation/advanced/environment-variables.md" >}}).

#### Configuration: Choosing an adapter while installing Mesheryctl

When you install `mesheryctl` with the Bash installation script, you can choose which adapters to load by passing the `ADAPTERS` environment variable to the script.

*For example:*
`curl -L https://meshery.io/install | ADAPTERS=consul PLATFORM=kubernetes bash -` installs `mesheryctl`, creates a context named `custom` that lists only the Meshery Adapter for Consul, and starts Meshery Server in your current Kubernetes cluster with that adapter.

<h5>Demo of Meshery managing deployments across multiple Kubernetes clusters:</h5>

{{< youtube id="yWPu3vq4vEs" start="5041" class="yt-embed-container" >}}

See on YouTube: [Cloud Native Austin Virtual Meetup: April 2020](https://youtu.be/yWPu3vq4vEs?t=5041&list=PL3A-A6hPO2IOpTbdH89qR-4AE0ON13Zie)
