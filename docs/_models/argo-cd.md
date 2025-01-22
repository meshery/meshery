---
layout: integration
title: Argo CD Applications
subtitle: Collaborative and visual infrastructure as design for Argo CD Applications
image: /assets/img/integrations/argo-cd/icons/color/argo-cd-color.svg
permalink: extensibility/integrations/argo-cd
docURL: https://docs.meshery.io/extensibility/integrations/argo-cd
description: 
integrations-category: App Definition and Development
integrations-subcategory: Continuous Integration & Delivery
registrant: Artifact Hub
components: 
- name: argo-cd-extension
  colorIcon: assets/img/integrations/argo-cd/components/argo-cd-extension/icons/color/argo-cd-extension-color.svg
  whiteIcon: assets/img/integrations/argo-cd/components/argo-cd-extension/icons/white/argo-cd-extension-white.svg
  description: 
- name: application
  colorIcon: assets/img/integrations/argo-cd/components/application/icons/color/application-color.svg
  whiteIcon: assets/img/integrations/argo-cd/components/application/icons/white/application-white.svg
  description: 
- name: application-set
  colorIcon: assets/img/integrations/argo-cd/components/application-set/icons/color/application-set-color.svg
  whiteIcon: assets/img/integrations/argo-cd/components/application-set/icons/white/application-set-white.svg
  description: 
- name: app-project
  colorIcon: assets/img/integrations/argo-cd/components/app-project/icons/color/app-project-color.svg
  whiteIcon: assets/img/integrations/argo-cd/components/app-project/icons/white/app-project-white.svg
  description: 
components-count: 4
relationships: 
- type: "Parent"
  kind: "Hierarchical"
  description: "A parent-child relationship implies the requirement of the parent component before the child component can be created. For example, a Namespace in Kubernetes can be a parent of Pods within that Namespace. Similarly, in Argo CD, an AppProject represents a logical grouping of Applications. Applications and ApplicationSets reference their AppProject by name."
- type: "Parent"
  kind: "Hierarchical"
  description: "A parent-child relationship implies the requirement of the parent component before the child component can be created. For example, a Namespace in Kubernetes can be a parent of Pods within that Namespace. Similarly, in Argo CD, an AppProject represents a logical grouping of Applications. Applications and ApplicationSets reference their AppProject by name."
relationship-count: 2
featureList: [
  "Gain mastery over Argo as a workflow engine for Kubernetes",
  "Design application definitions using the intuitive, context-aware visual designer, MeshMap.",
  "Visualize your declarative, GitOps continuous delivery tool for Kubernetes."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
