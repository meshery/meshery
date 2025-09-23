---
layout: default
title: Troubleshooting Out-of-Cluster Meshery Installation Connectivity
abstract: Troubleshoot connectivity issues for out-of-cluster Meshery installations.
permalink: guides/troubleshooting/out-of-cluster-connectivity
type: guides
category: troubleshooting
language: en
---

When you add a cluster to the Meshery Server, you might see this warning:

<img class="center" style="width:100%; height:auto;" src="{{site.baseurl}}/assets/img/troubleshoot/out-of-cluster-connectivity-error.png" alt="Out-of-cluster connectivity warning" />

> **⚠️ No reachable contexts found in the uploaded kubeconfig `file_name`**

---

## 🌐 What Does "Out-of-Cluster" Mean?

An _out-of-cluster installation_ means the **Meshery Server** runs **outside your target environment** (for example, as a standalone binary, in a Docker container, or in another Kubernetes cluster).

The **Meshery Operator** and workloads still run **inside your target cluster**.

In this setup, the Meshery Server must **directly reach your cluster’s API server URL**.

---

## 🔍 Debugging Steps

Follow these steps to fix the connectivity issue:

### 1. Open your kubeconfig file

By default, the file is here:

```bash
~/.kube/config
```

### 2. Find the Cluster API Server URL

Look for the `server` field in your kubeconfig:

```yaml
clusters:
  - cluster:
      server: https://<ip-address>:<port>
```

This field contains the Kubernetes API server URL.

### 3. Test Connectivity from Meshery Server

- If Meshery runs **on your host machine**, test access with:

  ```bash
  curl -k https://<ip-address>:<port>
  ```

- If Meshery runs **inside a container or pod**, connect to it and test from inside:

  ```bash
  kubectl exec -it <meshery-server-pod> -- curl -k https://<ip-address>:<port>
  ```

If the tests fail, Meshery cannot communicate with the cluster.

### 4. Configure Firewall and Network Rules

Make sure you:

- Allow **outbound** traffic from the Meshery Server host/container to the Kubernetes API server IP and port.
- Allow **inbound** responses from the API server back to Meshery.
- Verify that your load balancer or cloud provider endpoint is reachable from where Meshery runs.

---

Once the Meshery Server can reach the Kubernetes API server, the warning will disappear when you add the cluster.
