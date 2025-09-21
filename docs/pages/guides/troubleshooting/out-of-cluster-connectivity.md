---
layout: default
title: Troubleshooting Out-of-Cluster Meshery Installation Connectivity
abstract: Troubleshoot connectivity issues for out-of-cluster Meshery installations.
permalink: guides/troubleshooting/out-of-cluster-connectivity
type: guides
category: troubleshooting
language: en
---

You may encounter the following warning when adding a cluster to the Meshery Server:

<img class="center" style="width:100%; height:auto;" src="{{site.baseurl}}/assets/img/troubleshoot/out-of-cluster-connectivity-error.png" alt="Out-of-cluster connectivity warning" />

> **⚠️ No reachable contexts found in the uploaded kubeconfig `file_name`**

This indicates that the **Meshery Server**, installed **out-of-cluster**, cannot access the Kubernetes cluster’s API server URL.

---

## 🔍 Debugging Steps

Follow these steps to resolve the connectivity issue:

### 1. Open your kubeconfig file

The default location is:

```bash
~/.kube/config
```

### 2. Identify the Cluster API Server URL

Inside the kubeconfig, look for a section like this:

```yaml
clusters:
  - cluster:
      server: https://<ip-address>:<port>
```

The `server` field is the URL of your Kubernetes API server.

### 3. Verify Meshery’s Access

Ensure that the Meshery Server can reach this URL. If it cannot:

- **Check network connectivity** from the Meshery Server container to the cluster.
- **Update the API server URL** in the kubeconfig if necessary.
- **Ensure proper firewall rules** allow Meshery to reach the Kubernetes API.

> Once Meshery can access the API server, the warning should disappear when adding the cluster.
