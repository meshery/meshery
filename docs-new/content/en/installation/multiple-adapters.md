---
title: "Using Multiple Meshery Adapters"
description: "Running multiple instances of the same Meshery adapter"
weight: 60
aliases:
  - /installation/platforms/multiple-of-the-same-adapter
---

# Using Multiple Meshery Adapters

You can run multiple instances of the same Meshery adapter to manage multiple service mesh instances or for high availability.

## Configuring Multiple Adapters

### Using mesheryctl

When starting Meshery, you can specify multiple adapter instances in your configuration:

{{< code >}}
mesheryctl system config --adapter-urls "localhost:10000,localhost:10001"
{{< /code >}}

### Using Helm

Modify your Helm values to include multiple adapter deployments:

{{< code >}}
helm install meshery meshery/meshery \
  --set adapters.replicas=2 \
  --namespace meshery
{{< /code >}}

## Use Cases

- **Multi-cluster management**: Connect different adapters to different clusters
- **High availability**: Run redundant adapter instances for fault tolerance
- **Testing**: Test different adapter configurations simultaneously

## Verifying Configuration

Check the running adapters:

{{< code >}}
mesheryctl system status
{{< /code >}}

{{< related-discussions tag="meshery" >}}
