---
title: Installation
description: "Installation procedures for deploying Meshery with mesheryctl"
weight: 20
---

<h1>Installation</h1>

<h2>Supported Platforms</h2>

Meshery deploys as a set of Docker containers to either a Docker host or a Kubernetes cluster. A given deployment of Meshery can be described as either an _in-cluster_ or an _out-of-cluster_ deployment. Meshery runs as a standalone management plane on a Docker host (_out-of-cluster_) or within a Kubernetes cluster (_in-cluster_). See the complete list of supported platforms below.

<h3>Install <code>mesheryctl</code></h3>

- [Bash](/installation/linux-mac/bash) - Install Meshery CLI on Linux or MacOS with Bash
- [Brew](/installation/linux-mac/brew) - Install Meshery CLI on Linux or macOS with Homebrew
- [Scoop](/installation/windows/scoop) - Install Meshery CLI on Windows with Scoop

<h3>Install on Kubernetes</h3>

- [AKS](/installation/kubernetes/aks) - Manage your AKS clusters with Meshery. Deploy Meshery in AKS in-cluster or out-of-cluster.
- [Codespaces](/installation/codespaces) - Build and contribute to Meshery using GitHub Codespaces
- [EKS](/installation/kubernetes/eks) - Install Meshery on Elastic Kubernetes Service. Deploy Meshery in EKS in-cluster or outside of EKS out-of-cluster.
- [GKE](/installation/kubernetes/gke) - Install Meshery on Google Kubernetes Engine. Deploy Meshery in GKE in-cluster or outside of GKE out-of-cluster.
- [Helm](/installation/kubernetes/helm) - Install Meshery on Kubernetes using Helm. Deploy Meshery in Kubernetes in-cluster.
- [KinD](/installation/kubernetes/kind) - Install Meshery on KinD. Deploy Meshery in KinD in-cluster or outside of KinD out-of-cluster.
- [Kubernetes](/installation/kubernetes) - Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.
- [KubeSphere](/installation/kubernetes/kubesphere) - Install Meshery on KubeSphere
- [Minikube](/installation/kubernetes/minikube) - Install Meshery on Minikube. Deploy Meshery in Minikube in-cluster or outside of Minikube out-of-cluster.

<h3>Install on Docker</h3>

- [Docker](/installation/docker) - Install Meshery on Docker
- [Docker Extension](/installation/docker/docker-extension) - Install Docker Extension for Meshery
- [GitHub Codespaces](/installation/codespaces) - Install Meshery in GitHub Codespaces

<h3>Additional Resources</h3>

- [Upgrading Meshery](/installation/upgrades) - How to upgrade Meshery
- [Using Multiple Adapters](/installation/multiple-adapters) - Running multiple Meshery adapters
- [Compatibility Matrix](/installation/compatibility-matrix) - Meshery compatibility with platforms and versions
