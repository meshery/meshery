---
layout: integration
title: Prometheus
subtitle: Discover and connect to your Prometheus servers and Operators
image: /assets/img/integrations/prometheus/icons/color/prometheus-color.svg
permalink: extensibility/integrations/prometheus
docURL: https://docs.meshery.io/extensibility/integrations/prometheus
description: Guide to configure Prometheus integration with Meshery
integrations-category: Observability and Analysis
integrations-subcategory: Monitoring
registrant: Artifact Hub
components: 
components-count: 0
relationships: 
relationship-count: 0
featureList: [
  "Native support for PromQL",
  "Create custom charts with your own Prometheus queries",
  "Keep charts in-sync with Meshery's panel viewer"
]
howItWorks: "Meshery provides performance reports, including performance test results, node resource metrics etc. so that operators may easily understand the overhead of their service mesh’s control plane and data plane in context of the overhead incurred on nodes running within the cluster. In order to generate performance test reports of service meshes and their workloads, Meshery uses Grafana and/or Prometheus as visualization and metrics systems, respectively. This guide outlines the requirements necessary for Meshery to connect to these systems. The steps may vary depending upon the service mesh and its configuration."
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---

# Configuring Prometheus with Meshery

Prometheus is a widely used open-source monitoring and alerting toolkit designed to record real-time metrics and generate alerts.

## Setup Instructions

1. **Install Prometheus**

You can install Prometheus using Helm in your Kubernetes cluster:

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus
