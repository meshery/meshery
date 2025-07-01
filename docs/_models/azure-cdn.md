---
layout: integration
title: Azure CDN
subtitle: Collaborative and visual infrastructure as design for Azure CDN
image: /assets/img/integrations/azure-cdn/icons/color/azure-cdn-color.svg
permalink: extensibility/integrations/azure-cdn
docURL: https://docs.meshery.io/extensibility/integrations/azure-cdn
description: 
integrations-category: Cloud Native Storage
integrations-subcategory: Content Delivery Network
registrant: GitHub
components: 
- name: afd-custom-domain
  colorIcon: assets/img/integrations/azure-cdn/components/afd-custom-domain/icons/color/afd-custom-domain-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/afd-custom-domain/icons/white/afd-custom-domain-white.svg
  description: 
- name: afd-endpoint
  colorIcon: assets/img/integrations/azure-cdn/components/afd-endpoint/icons/color/afd-endpoint-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/afd-endpoint/icons/white/afd-endpoint-white.svg
  description: 
- name: afd-origin-group
  colorIcon: assets/img/integrations/azure-cdn/components/afd-origin-group/icons/color/afd-origin-group-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/afd-origin-group/icons/white/afd-origin-group-white.svg
  description: 
- name: afd-origin
  colorIcon: assets/img/integrations/azure-cdn/components/afd-origin/icons/color/afd-origin-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/afd-origin/icons/white/afd-origin-white.svg
  description: 
- name: profile
  colorIcon: assets/img/integrations/azure-cdn/components/profile/icons/color/profile-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/profile/icons/white/profile-white.svg
  description: 
- name: profiles-endpoint
  colorIcon: assets/img/integrations/azure-cdn/components/profiles-endpoint/icons/color/profiles-endpoint-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/profiles-endpoint/icons/white/profiles-endpoint-white.svg
  description: 
- name: route
  colorIcon: assets/img/integrations/azure-cdn/components/route/icons/color/route-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/route/icons/white/route-white.svg
  description: 
- name: rule
  colorIcon: assets/img/integrations/azure-cdn/components/rule/icons/color/rule-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/rule/icons/white/rule-white.svg
  description: 
- name: rule-set
  colorIcon: assets/img/integrations/azure-cdn/components/rule-set/icons/color/rule-set-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/rule-set/icons/white/rule-set-white.svg
  description: 
- name: secret
  colorIcon: assets/img/integrations/azure-cdn/components/secret/icons/color/secret-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/secret/icons/white/secret-white.svg
  description: 
- name: security-policy
  colorIcon: assets/img/integrations/azure-cdn/components/security-policy/icons/color/security-policy-color.svg
  whiteIcon: assets/img/integrations/azure-cdn/components/security-policy/icons/white/security-policy-white.svg
  description: 
components-count: 11
relationships: 
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of Reule Set CDN(parent component) is patched with the configuration of Rule(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of AfdEndpoint CDN(parent component) is patched with the configuration of Route(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of AfdOriginGroup CDN(parent component) is patched with the configuration of AfdOrigin(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of CDN Profile(parent component) is patched with the configuration of AfdCustomDomain(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of CDN Profile(parent component) is patched with the configuration of AfdEndpoint(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of CDN Profile(parent component) is patched with the configuration of AfdOriginGroup(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of CDN Profile(parent component) is patched with the configuration of ProfilesEndpoint(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of CDN Profile(parent component) is patched with the configuration of RuleSet(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of CDN Profile(parent component) is patched with the configuration of Secret(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of CDN Profile(parent component) is patched with the configuration of SecurityPolicy(child component). "
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between Route and AfdOriginGroup(azure-cdn)"
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between Route and RuleSet(azure-cdn)"
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between SecurityPolicy and AfdEndpoint(azure-cdn)"
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between SecurityPolicy and WebApplicationFirewallPolicy(azure-network)"
relationship-count: 14
featureList: [
  "Drag-n-drop cloud native infrastructure designer to configure, model, and deploy your workloads.",
  "Invite anyone to review and make changes to your private designs.",
  "Ongoing synchronization of Kubernetes configuration and changes across any number of clusters."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
