---
title: Karpenter
subtitle: Collaborative and visual infrastructure as design for Karpenter
image: /extensions/models/karpenter/icons/color/karpenter-color.svg
aliases: /extensibility/integrations/karpenter
docURL: https://docs.meshery.io/extensibility/integrations/karpenter
description: 
integrations-category: Provisioning
integrations-subcategory: Automation & Configuration
registrant: Artifact Hub
components: 
- name: aws-node-template
  colorIcon: content/en/extensions/models/karpenter/components/aws-node-template/icons/color/aws-node-template-color.svg
  whiteIcon: content/en/extensions/models/karpenter/components/aws-node-template/icons/white/aws-node-template-white.svg
  description: 
- name: provisioner
  colorIcon: content/en/extensions/models/karpenter/components/provisioner/icons/color/provisioner-color.svg
  whiteIcon: content/en/extensions/models/karpenter/components/provisioner/icons/white/provisioner-white.svg
  description: 
- name: node-claim
  colorIcon: content/en/extensions/models/karpenter/components/node-claim/icons/color/node-claim-color.svg
  whiteIcon: content/en/extensions/models/karpenter/components/node-claim/icons/white/node-claim-white.svg
  description: 
- name: node-pool
  colorIcon: content/en/extensions/models/karpenter/components/node-pool/icons/color/node-pool-color.svg
  whiteIcon: content/en/extensions/models/karpenter/components/node-pool/icons/white/node-pool-white.svg
  description: 
- name: oci-node-class
  colorIcon: content/en/extensions/models/karpenter/components/oci-node-class/icons/color/oci-node-class-color.svg
  whiteIcon: content/en/extensions/models/karpenter/components/oci-node-class/icons/white/oci-node-class-white.svg
  description: 
components-count: 5
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
---
