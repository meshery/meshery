---
layout: default
title: Troubleshooting Meshery Installations
description: Troubleshoot Meshery installation and deployment
permalink: guides/troubleshooting/installation
type: Guides
category: troubleshooting
language: en
---

## Meshery's Preflight Checks

Anytime a `mesheryctl system` command is executed, a series of preflight checks are run. An attempt will be made to connect to the Kubernetes cluster configured in the user's kubeconfig as their current-context .

1. Check whether `mesheryctl` can initialize a Kubernetes client.

   Situation: `mesheryctl` fails to query for pods in the default namespace of the user's current Kubernetes context.

2. Remove `~/.meshery` to reinitialize Meshery

   Situation: Unable to start Meshery Server with `make run-local` due to error of `key/value size is invalid`

## Setting up Meshery using Kind or Minikube

The difficulty with Minikube and Kind clusters is that they typically don't support LoadBalancer service networking by default. Meshery UI and Meshery Broker are configured for LoadBalancer service networking by default. There are a number of solutions this overcoming this challenge. Here are a few methods:

1. Use the MetalLB Minikube add-on that provides load balancing. `minikube addons enable metallb`

   MetalLB setup: [link](https://kubebyexample.com/learning-paths/metallb/install)

2. Use Minikube tunnel to expose services. `minikube tunnel`.

   Docs: [link](https://minikube.sigs.k8s.io/docs/handbook/accessing/#using-minikube-tunnel)

   A simpler way to resolve this issue can be `port-forwarding`. Run the following command in terminal:

   `kubectl port-forward service/meshery 9081:9081 -n meshery`

3. For `kind`, you can prefer installing MetalLB with a custom configmap.

   Docs: [link](https://kind.sigs.k8s.io/docs/user/loadbalancer/)

## Meshery Operator

By default, Meshery Operator is installed in all the connected clusters automatically once Meshery server detects those clusters. The operator can manually be turned off on particular cluster from the settings page.

### Disabling the operator

The env variable DISABLE_OPERATOR=true can be used to signal Meshery server to not install operator in any of the clusters at any point in time after starting. While using Meshery server locally, the `make server-without-operator` should be used to start Meshery in disabled operator mode.

### Meshery Broker

Example of a healthy Meshery Broker server with an actively connected (subscribed) Meshery Server:

```
➜  ~ kubectl logs -n meshery meshery-broker-0 nats
[8] 2021/09/08 21:46:03.070952 [INF] Starting nats-server version 2.1.9
[8] 2021/09/08 21:46:03.070982 [INF] Git commit [7c76626]
[8] 2021/09/08 21:46:03.071308 [INF] Starting http monitor on 0.0.0.0:8222
[8] 2021/09/08 21:46:03.071370 [INF] Listening for client connections on 0.0.0.0:4222
[8] 2021/09/08 21:46:03.071512 [INF] Server id is NAAYJNX4LDDNXW5UE7IP7PRQR2W2JP546XSFNUWQQHN7JYY27RG47KSG
[8] 2021/09/08 21:46:03.071516 [INF] Server is ready
```

For details about the state of the Meshery Server subscription see the http monitor port on Meshery Broker.

### Meshery Unable to Connect to Kubernetes

Meshery is unable to detect the Kubernetes connection running on your local system, even after manually uploading the `.kube config` file

When deploying Meshery out-of-cluster, verify your kubeconfig's contexts and the ability for Meshery Server to reach Kubernetes cluster API from whatever host and network that Meshery Server is being deployed on.

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">kubectl config get-contexts</div></div>
</pre>

If you're using Docker Destkop, consider whether you need to change your current Kubernetes context to `docker-desktop`.

<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">kubectl config use-context
docker-desktop</div></div>
</pre>

## Meshery Remote Providers

Once Meshery is installed, the remote provider "Meshery" can be chosen from UI or by using the command `mesheryctl system login`:

![Providers](/assets/img/providers/provider_screenshot.png)

```bash
➜  ~ mesheryctl system login
Use the arrow keys to navigate: ↓ ↑ → ←
? Select a Provider:
  ▸ Meshery
    None
```

If you cannot see "Meshery" Remote Provider and find such error logs in Meshery Server's logs (`mesheryctl system logs`), please make sure that Meshery Server is able to reach "https://meshery.layer5.io" in order to initialize the "Meshery" Remote Provider.

```bash
time="2021-11-10T11:05:30Z" level=error msg="[Initialize Provider]: Failed to get capabilities Get \"https://meshery.layer5.io/v0.5.71/capabilities?os=meshery\": dial tcp 3.140.89.205:443: i/o timeout"
```

For more details about Meshery Providers:

- [Extensibility: Providers](/extensibility/providers)

## See Also

- [Meshery Error Code Reference](/reference/error-codes)

