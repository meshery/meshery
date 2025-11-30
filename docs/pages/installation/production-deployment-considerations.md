---
layout: default
title: Production Deployment Considerations
permalink: installation/prod-deployment-considerations
type: guides
category: installation
display-title: "false"
language: en
abstract: Considerations when deploying Meshery in a production environment
---

# Production Deployment Considerations

This document outlines key factors to consider when deploying Meshery in a production environment. 

Meshery, a powerful cloud native management plane, requires careful planning and configuration to ensure reliability, scalability, security, and observability through effective monitoring.

## Infrastructure & Resource Requirements

It is important to consider the following infrastructure and networking aspects when deploying Meshery to a production environment.

### High Availability

For production environments, deploy Meshery in a highly available (HA) configuration.
This requires both:

1. **A Highly Available Infrastructure**: deployed across multiple nodes or availability zones, with load balancing to eliminate single points of failure. For example with self-managed Kubernetes cluster, this means running a minimum of three stacked control plane nodes to satisfy the [High Availability topology](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/ha-topology/).

Learn more about creating an [HA Cluster](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/high-availability/).

2. **A Highly Available Meshery deployment**: run at least three replicas of each Meshery workload, such as `deployment/meshery` and `deployment/meshery-operator`. The Kubernetes Scheduler will distribute these replicas across available nodes. If you want to customize scheduling or assign pods to specific nodes, refer to the [Pod assignment guide](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/)

For example, when installing Meshery with its [Helm Chart](https://github.com/meshery/meshery/tree/master/install/kubernetes/helm/meshery) set `replicaCount` to `3`:

<pre class="codeblock-pre">
<div class="codeblock">
<div class="clipboardjs">
helm repo add meshery https://meshery.io/charts/
helm repo update

helm install meshery meshery/meshery --namespace meshery --create-namespace --set replicaCount=3
</div></div>
</pre>


### Scalability

For mission-critical production systems, plan for Meshery to handle increased traffic on its workloads as your system grows.

You may use [KEDA](https://keda.sh) (a Kubernetes Event-driven Autoscaling) or create `HorizontalPodAutoscaler` resources targeting Meshery workloads, such as `deployment/meshery` and `deployment/meshery-operator`, to autoscale up or down based on resource usage or traffic-related metrics such as `req/s`.

Learn more about the [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/).

Also, consider using a scalable infrastructure that can automatically adjust resources based on demand, such as a _Cluster Autoscaler_ provided by your cloud vendor.


### Networking & Security

It is recommended to configure proper network connectivity for Meshery to communicate with Kubernetes clusters, and other relevant components. 
Ensure secure access and consider using [network policies](https://kubernetes.io/docs/concepts/services-networking/network-policies/) or firewalls to restrict traffic.

Make sure to keep the following ports open for inbound and outbound connectivity.  

```bash
Protocol: TCP
Ports: 4222, 31214, 6222, 30315, 8222, 30791, 7777, 30400, 7422, 32493, 7522, 30337, 10000, 9081, 30806, 80, 443

Protocol: gRPC
Ports: 10005, 10002, 10012, 10001, 10007, 10010, 10013, 10004, 10011, 10006
```
Check what ports are required for your use-case and more details in [Meshery's Network Ports]({{site.baseurl}}/concepts/architecture#network-ports)



If you are using **[Emissary Ingress](https://github.com/emissary-ingress/emissary)** make sure to configure it for secure WebSocket support by defining a `Mapping` resource and ensuring TLS is handled correctly as in this [guide](https://emissary-ingress.dev/docs/3.10/howtos/websockets/)

Another security consideration is to **avoid using the `Local Provider`**, since it allows unauthenticated sessions. Instead, preselect a `Remote Provider`, which enforces authentication and authorization, securing the public-facing Meshery UI. Using a Remote Provider also enables Meshery’s multi-user mode, allowing collaboration.

Below is an image of selecting provider in Meshery's UI, Local Provider is named **None** while Remote Provider is default to **Meshery** or **Layer5**

<a href="{{ site.baseurl }}/assets/img/providers/provider_screenshot_new.png">
<img src="{{ site.baseurl }}/assets/img/providers/provider_screenshot_new.png" width="80%" height="auto" alt="Meshery's providers" /></a>


Remote providers also can be implemented by anyone or an organization that wishes to integrate with Meshery.

if you're [building a provider]({{site.baseurl}}/extensibility/providers#building-a-provider) please **configure OAuth Callback URL** for your connected remote provider, this is helpful when you are deploying Meshery behind an ingress gateway or reverse proxy. You can specify a custom, redirect endpoint for the connected Remote Provider, check this [example]({{site.baseurl}}/extensibility/providers#example-deployment-scenario-meshery-cloud-istio-and-azure-kubernetes-service).

Learn more about [Providers]({{site.baseurl}}/extensibility/providers)



## Monitoring & Observability

Ensuring Meshery’s reliability in production requires a strong monitoring and observability strategy. This includes comprehensive monitoring, centralized logging, and proactive alerting.

### Comprehensive Monitoring

Set up robust monitoring to track Meshery’s health, performance, and resource utilization. Key aspects include:

  - Pod health: Monitor the status of Meshery pods (e.g., `Running` or `Error`).
  - Performance metrics: Track `Total CPU %`, `System CPU %`, `RSS Memory`, and `Memory Cache` on pods/containers.
  - Resource utilization: Observe CPU and memory usage against configured resource requests and limits.
  - Service-level metrics: Monitor Meshery API response times, request per second `req/s`, latency, error rates, and overall system load.
  - Distributed traces: Capture meshery-server traces and spans; correlate traces with logs for deeper insights.

### Centralized Logging

Aggregate and store Meshery logs in a centralized system for efficient analysis and troubleshooting.

Filter logs by service (`meshery`, `meshery-operator`, `meshery-broker`, …), `status_code`, `http.method` (GET, POST, …), or environment (`production`, `staging`, …).

Use log analysis tools to identify issues, detect anomalies, and recognize logs' patterns.

### Alerting

Configure alerting systems to notify administrators of critical events, such as:

  - Service downtime.
  - Infrastructure outage or nodes not in `Ready` state
  - Performance issues and bottlenecks.
  - High error rates and HTTP `500` `501` `502` `503` responses.
  - Resource exhaustion or abnormal usage patterns.

### Tooling

To achieve Monitoring and Observability, consider leveraging tools such as OpenTelemetry, Prometheus, Grafana, Fluent Bit, Jaeger, Datadog or Dynatrace for collecting, correlating, and visualizing metrics, traces and logs.

