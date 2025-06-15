---
layout: integration
title: Azure API Management
subtitle: Collaborative and visual infrastructure as design for Azure API Management
image: /assets/img/integrations/azure-api-management/icons/color/azure-api-management-color.svg
permalink: extensibility/integrations/azure-api-management
docURL: https://docs.meshery.io/extensibility/integrations/azure-api-management
description: 
integrations-category: App Definition and Development
integrations-subcategory: API Gateway
registrant: GitHub
components: 
- name: api
  colorIcon: assets/img/integrations/azure-api-management/components/api/icons/color/api-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/api/icons/white/api-white.svg
  description: 
- name: api-version-set
  colorIcon: assets/img/integrations/azure-api-management/components/api-version-set/icons/color/api-version-set-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/api-version-set/icons/white/api-version-set-white.svg
  description: 
- name: authorization-provider
  colorIcon: assets/img/integrations/azure-api-management/components/authorization-provider/icons/color/authorization-provider-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/authorization-provider/icons/white/authorization-provider-white.svg
  description: 
- name: authorization-providers-authorization
  colorIcon: assets/img/integrations/azure-api-management/components/authorization-providers-authorization/icons/color/authorization-providers-authorization-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/authorization-providers-authorization/icons/white/authorization-providers-authorization-white.svg
  description: 
- name: authorization-providers-authorizations-access-policy
  colorIcon: assets/img/integrations/azure-api-management/components/authorization-providers-authorizations-access-policy/icons/color/authorization-providers-authorizations-access-policy-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/authorization-providers-authorizations-access-policy/icons/white/authorization-providers-authorizations-access-policy-white.svg
  description: 
- name: backend
  colorIcon: assets/img/integrations/azure-api-management/components/backend/icons/color/backend-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/backend/icons/white/backend-white.svg
  description: 
- name: named-value
  colorIcon: assets/img/integrations/azure-api-management/components/named-value/icons/color/named-value-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/named-value/icons/white/named-value-white.svg
  description: 
- name: policy
  colorIcon: assets/img/integrations/azure-api-management/components/policy/icons/color/policy-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/policy/icons/white/policy-white.svg
  description: 
- name: policy-fragment
  colorIcon: assets/img/integrations/azure-api-management/components/policy-fragment/icons/color/policy-fragment-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/policy-fragment/icons/white/policy-fragment-white.svg
  description: 
- name: product-api
  colorIcon: assets/img/integrations/azure-api-management/components/product-api/icons/color/product-api-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/product-api/icons/white/product-api-white.svg
  description: 
- name: product-policy
  colorIcon: assets/img/integrations/azure-api-management/components/product-policy/icons/color/product-policy-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/product-policy/icons/white/product-policy-white.svg
  description: 
- name: product
  colorIcon: assets/img/integrations/azure-api-management/components/product/icons/color/product-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/product/icons/white/product-white.svg
  description: 
- name: service
  colorIcon: assets/img/integrations/azure-api-management/components/service/icons/color/service-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/service/icons/white/service-white.svg
  description: 
- name: subscription
  colorIcon: assets/img/integrations/azure-api-management/components/subscription/icons/color/subscription-color.svg
  whiteIcon: assets/img/integrations/azure-api-management/components/subscription/icons/white/subscription-white.svg
  description: 
components-count: 14
relationships: 
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of API(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of ApiVersionSet(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of AuthorizationProvider(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of Backend(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of NamedValue(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of Policy(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of PolicyFragment(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of Product(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of APIM Service(parent component) is patched with the configuration of Subscription(child component). "
relationship-count: 9
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
