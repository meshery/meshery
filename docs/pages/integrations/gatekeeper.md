---
layout: enhanced
title: OPA Gatekeeper
subtitle: Collaborative and visual infrastructure as code for Frinx Machine
image: /assets/img/integrations/gatekeeper/icons/color/gatekeeper-color.svg
permalink: extensibility/integrations/gatekeeper
docURL: https://docs.meshery.io/extensibility/integrations/gatekeeper
description: 
integrations-category: Security & Compliance
integrations-subcategory: Security & Compliance
registrant: Artifact Hub
components: 
- name: assign
  colorIcon: assets/img/integrations/gatekeeper/components/assign/icons/color/assign-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/assign/icons/white/assign-white.svg
  description: 
- name: assign-metadata
  colorIcon: assets/img/integrations/gatekeeper/components/assign-metadata/icons/color/assign-metadata-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/assign-metadata/icons/white/assign-metadata-white.svg
  description: 
- name: constraint-template
  colorIcon: assets/img/integrations/gatekeeper/components/constraint-template/icons/color/constraint-template-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/constraint-template/icons/white/constraint-template-white.svg
  description: 
- name: expansion-template
  colorIcon: assets/img/integrations/gatekeeper/components/expansion-template/icons/color/expansion-template-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/expansion-template/icons/white/expansion-template-white.svg
  description: 
- name: modify-set
  colorIcon: assets/img/integrations/gatekeeper/components/modify-set/icons/color/modify-set-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/modify-set/icons/white/modify-set-white.svg
  description: 
- name: mutator-pod-status
  colorIcon: assets/img/integrations/gatekeeper/components/mutator-pod-status/icons/color/mutator-pod-status-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/mutator-pod-status/icons/white/mutator-pod-status-white.svg
  description: 
- name: provider
  colorIcon: assets/img/integrations/gatekeeper/components/provider/icons/color/provider-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/provider/icons/white/provider-white.svg
  description: 
- name: assign-image
  colorIcon: assets/img/integrations/gatekeeper/components/assign-image/icons/color/assign-image-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/assign-image/icons/white/assign-image-white.svg
  description: 
- name: expansion-template-pod-status
  colorIcon: assets/img/integrations/gatekeeper/components/expansion-template-pod-status/icons/color/expansion-template-pod-status-color.svg
  whiteIcon: assets/img/integrations/gatekeeper/components/expansion-template-pod-status/icons/white/expansion-template-pod-status-white.svg
  description: 
featureList: [
  "Native Kubernetes CRDs for instantiating the policy library (aka constraints)",
  "An extensible, parameterized policy library",
  "Ongoing synchronization of Kubernetes configuration and changes across any number of clusters."
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
Compared to using OPA with its sidecar kube-mgmt (aka Gatekeeper v1.0), Gatekeeper introduces the following functionality:
</p>
<p>
    Connect GitHub with Meshery and import selectively import your existing Helm Charts, Docker Compose applications, and Kubernetes manifests.Visually configure and customize your cloud native infrastructure.
    Save and share your design patterns to GitHub using either public or private repositories.
</p>
<p>
    Learn more about <a href="/blog/service-mesh-specifications/pipelining-service-mesh-specifications">pipelining service mesh specifications</a> and using Service Mesh Interface and Service Mesh Performance specs on your CI/CD pipelines with Meshery's GitHub Actions.
</p>
