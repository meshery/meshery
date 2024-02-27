---
layout: enhanced
title: Azure Application Gateway 
subtitle: Collaborative and visual infrastructure as code for Azure Application Gateway 
image: /assets/img/integrations/aks-appgw-fe/icons/color/aks-appgw-fe-color.svg
permalink: extensibility/integrations/aks-appgw-fe
docURL: https://docs.meshery.io/extensibility/integrations/aks-appgw-fe
description: 
integrations-category: Cloud Native Network
integrations-subcategory: Service Proxy
registrant: Artifact Hub
components: 
- name: azure-assigned-identity
  colorIcon: assets/img/integrations/aks-appgw-fe/components/azure-assigned-identity/icons/color/azure-assigned-identity-color.svg
  whiteIcon: assets/img/integrations/aks-appgw-fe/components/azure-assigned-identity/icons/white/azure-assigned-identity-white.svg
  description: 
- name: azure-identity
  colorIcon: assets/img/integrations/aks-appgw-fe/components/azure-identity/icons/color/azure-identity-color.svg
  whiteIcon: assets/img/integrations/aks-appgw-fe/components/azure-identity/icons/white/azure-identity-white.svg
  description: 
- name: azure-identity-binding
  colorIcon: assets/img/integrations/aks-appgw-fe/components/azure-identity-binding/icons/color/azure-identity-binding-color.svg
  whiteIcon: assets/img/integrations/aks-appgw-fe/components/azure-identity-binding/icons/white/azure-identity-binding-white.svg
  description: 
- name: azure-pod-identity-exception
  colorIcon: assets/img/integrations/aks-appgw-fe/components/azure-pod-identity-exception/icons/color/azure-pod-identity-exception-color.svg
  whiteIcon: assets/img/integrations/aks-appgw-fe/components/azure-pod-identity-exception/icons/white/azure-pod-identity-exception-white.svg
  description: 
featureList: [
  "URL routing and cookie-based affinity
",
  "Support for public, private, and hybrid web sites and 
integrated web application firewall",
  "Secure Sockets Layer (SSL) termination and End-to-end SSL"
]
howItWorks: "Collaborative Infrastructure as Code"
howItWorksDetails: "Application Gateway Ingress Controller runs in its own pod on the customer’s AKS. Ingress Controller monitors a subset of Kubernetes’ resources for changes. The state of the AKS cluster is translated to Application Gateway specific configuration and applied to the Azure Resource Manager. The continuous re-configuration of Application Gateway ensures uninterrupted flow of traffic to AKS’ services. The diagram below illustrates the flow of state and configuration changes from the Kubernetes API, via Application Gateway Ingress Controller, to Resource Manager and then Application Gateway."
language: en
list: include
type: extensibility
category: integrations
display-title: "false"
---
<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>

<p>
Create generic frontend web application in AKS with SSL termination on an Application Gateway
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
