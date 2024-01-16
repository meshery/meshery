---
layout: enhanced
title: Azure Application Gateway 
subtitle: Collaborative and visual infrastructure as code for Azure Application Gateway 
image: /assets/img/integrations/azure-application-gateway-.svg
permalink: extensibility/integrations/azure-application-gateway-
docURL: https://docs.meshery.io/extensibility/integrations/aks-appgw-fe
description: 
category: Cloud Native Network
subcategory: Service Proxy
registrant: artifacthub
components: 
	-	name: AzureAssignedIdentity
		colorIcon: assets/img/integrations/components/AzureAssignedIdentity-color.svg
		whiteIcon: assets/img/integrations/components/AzureAssignedIdentity-white.svg
		description: 
	-	name: AzureIdentity
		colorIcon: assets/img/integrations/components/AzureIdentity-color.svg
		whiteIcon: assets/img/integrations/components/AzureIdentity-white.svg
		description: 
	-	name: AzureIdentityBinding
		colorIcon: assets/img/integrations/components/AzureIdentityBinding-color.svg
		whiteIcon: assets/img/integrations/components/AzureIdentityBinding-white.svg
		description: 
	-	name: AzurePodIdentityException
		colorIcon: assets/img/integrations/components/AzurePodIdentityException-color.svg
		whiteIcon: assets/img/integrations/components/AzurePodIdentityException-white.svg
		description: 
featureList: [
  "URL routing and cookie-based affinity
",
  "Support for public, private, and hybrid web sites and 
integrated web application firewall",
  "Secure Sockets Layer (SSL) termination and End-to-end SSL"
]
howItWorks: Collaborative Infrastructure as Code
howItWorksDetails: Application Gateway Ingress Controller runs in its own pod on the customer’s AKS. Ingress Controller monitors a subset of Kubernetes’ resources for changes. The state of the AKS cluster is translated to Application Gateway specific configuration and applied to the Azure Resource Manager. The continuous re-configuration of Application Gateway ensures uninterrupted flow of traffic to AKS’ services. The diagram below illustrates the flow of state and configuration changes from the Kubernetes API, via Application Gateway Ingress Controller, to Resource Manager and then Application Gateway.
language: en
list: include
---
<p>
Create generic frontend web application in AKS with SSL termination on an Application Gateway
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
