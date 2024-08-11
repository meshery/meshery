---
layout: integration
title: Azure Active Directory (AAD)
subtitle: Collaborative and visual infrastructure as design for Azure Active Directory (AAD)
image: /assets/img/integrations/aad-pod-identity/icons/color/aad-pod-identity-color.svg
permalink: extensibility/integrations/aad-pod-identity
docURL: https://docs.meshery.io/extensibility/integrations/aad-pod-identity
description: 
integrations-category: Provisioning
integrations-subcategory: Security & Compliance
registrant: Artifact Hub
components: 
- name: azure-assigned-identity
  colorIcon: assets/img/integrations/aad-pod-identity/components/azure-assigned-identity/icons/color/azure-assigned-identity-color.svg
  whiteIcon: assets/img/integrations/aad-pod-identity/components/azure-assigned-identity/icons/white/azure-assigned-identity-white.svg
  description: 
- name: azure-identity
  colorIcon: assets/img/integrations/aad-pod-identity/components/azure-identity/icons/color/azure-identity-color.svg
  whiteIcon: assets/img/integrations/aad-pod-identity/components/azure-identity/icons/white/azure-identity-white.svg
  description: 
- name: azure-identity-binding
  colorIcon: assets/img/integrations/aad-pod-identity/components/azure-identity-binding/icons/color/azure-identity-binding-color.svg
  whiteIcon: assets/img/integrations/aad-pod-identity/components/azure-identity-binding/icons/white/azure-identity-binding-white.svg
  description: 
- name: azure-pod-identity-exception
  colorIcon: assets/img/integrations/aad-pod-identity/components/azure-pod-identity-exception/icons/color/azure-pod-identity-exception-color.svg
  whiteIcon: assets/img/integrations/aad-pod-identity/components/azure-pod-identity-exception/icons/white/azure-pod-identity-exception-white.svg
  description: 
featureList: [
  "Use Azure Active Directory pod-managed identities in Azure Kubernetes Service.",
  "A maximum of 200 pod-managed identities are allowed for a cluster.",
  "Pod-managed identities are available on Linux node pools only."
]
howItWorks: "Federate with any external identity provider"
howItWorksDetails: "Azure AD workload identity authentication method replaces pod-managed identity, which integrates with the Kubernetes native capabilities to federate with any external identity providers on behalf of the application."
language: en
list: include
type: extensibility
category: integrations
---
