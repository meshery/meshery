---
title: Karpenter
subtitle: Collaborative and visual infrastructure as design for Karpenter
image: /extensions/models/karpenter/icons/color/karpenter-color.svg
docURL: https://docs.meshery.io/extensibility/integrations/karpenter
description: 
integrations-category: Provisioning
integrations-subcategory: Automation & Configuration
registrant: Artifact Hub
components: 
- name: aws-node-template
  colorIcon: /extensions/models/karpenter/components/aws-node-template/icons/color/aws-node-template-color.svg
  whiteIcon: /extensions/models/karpenter/components/aws-node-template/icons/white/aws-node-template-white.svg
  description: 
- name: provisioner
  colorIcon: /extensions/models/karpenter/components/provisioner/icons/color/provisioner-color.svg
  whiteIcon: /extensions/models/karpenter/components/provisioner/icons/white/provisioner-white.svg
  description: 
components-count: 2
relationships: 
relationship-count: 0
featureList: [
  "Watching for pods that the Kubernetes scheduler has marked as unschedulable",
  "Evaluating scheduling constraints (resource requests, nodeselectors, affinities, tolerations, and topology spread constraints) requested by the pods",
  "Provisioning nodes that meet the requirements of the pods"
]
howItWorks: "Using Meshery and Karpenter, once your Kubernetes cluster and the Karpenter controller are up and running"
howItWorksDetails: "Set up provisioners: By applying a provisioner to Karpenter, you can configure constraints on node provisioning and set timeout values for node expiry or Kubelet configuration values. 

Deploy workloads: When deploying workloads, you can request that scheduling constraints be met to direct which nodes Karpenter provisions for those workloads. "
categories: [integrations]
aliases:
- /extensibility/integrations/karpenter
---
