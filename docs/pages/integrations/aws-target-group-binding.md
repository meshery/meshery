---
layout: enhanced
title: AWS Target Group Binding
subtitle: Collaborative and visual infrastructure as code for AWS Target Group Binding
image: /assets/img/integrations/aws-target-group-binding/icons/color/aws-target-group-binding-color.svg
permalink: extensibility/integrations/aws-target-group-binding
docURL: https://docs.meshery.io/extensibility/integrations/aws-target-group-binding
description: 
integrations-category: Cloud Native Network
integrations-subcategory: Service Proxy
registrant: artifacthub
components: 
- name: target-group-binding
  colorIcon: assets/img/integrations/aws-target-group-binding/components/target-group-binding/icons/color/target-group-binding-color.svg
  whiteIcon: assets/img/integrations/aws-target-group-binding/components/target-group-binding/icons/white/target-group-binding-white.svg
  description: 
featureList: [
  "Expose your pods using an existing ALB TargetGroup or NLB TargetGroup.",
  "Allows you to provision the load balancer infrastructure completely outside of Kubernetes.",
  "Used by AWS LoadBalancer controller internally to  support the functionality for Ingress and Service resource as well."
]
howItWorks: "Collaborative Infrastructure as Code"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
display-title: "false"
---
<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>

<p>
Povision a TargetGroupBinding resource that can be combined with an externally provisioned AWS Loadbalancer to provide an ingress route into the defined kubernetes service.
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
