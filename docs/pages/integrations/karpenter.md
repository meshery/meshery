---
layout: default
title: Karpenter
permalink: integrations/karpenter
type: installation
category: integrations
display-title: "false"
language: en
list: include
image: /assets/img/integrations/karpenter.svg
---

<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>


<!-- This needs replaced with the Category property, not the sub-category.
 #### Category: karpenter -->

### Overview & Features:
1. Karpenter is an open-source node provisioning project built for Kubernetes. Adding Karpenter to a Kubernetes cluster can dramatically improve the efficiency and cost of running workloads on that cluster.

2. Collaborative and visual infrastructure as code for Karpenter

4. 
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.



    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.



5. Watching for pods that the Kubernetes scheduler has marked as unschedulable

6. Evaluating scheduling constraints (resource requests, nodeselectors, affinities, tolerations, and topology spread constraints) requested by the pods

7. Provisioning nodes that meet the requirements of the pods

8. Using Meshery and Karpenter, once your Kubernetes cluster and the Karpenter controller are up and running

9. Set up provisioners: By applying a provisioner to Karpenter, you can configure constraints on node provisioning and set timeout values for node expiry or Kubelet configuration values. 

Deploy workloads: When deploying workloads, you can request that scheduling constraints be met to direct which nodes Karpenter provisions for those workloads. 

