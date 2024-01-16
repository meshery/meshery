---
layout: enhanced
title: Contour
subtitle: Collaborative and visual infrastructure as code for Contour
image: /assets/img/integrations/contour.svg
permalink: extensibility/integrations/contour
docURL: https://docs.meshery.io/extensibility/integrations/contour-operator
description: 
category: Cloud Native Network
subcategory: Service Proxy
registrant: artifacthub
components: 
	-	name: Contour
		colorIcon: assets/img/integrations/components/Contour-color.svg
		whiteIcon: assets/img/integrations/components/Contour-white.svg
		description: 
	-	name: ContourConfiguration
		colorIcon: assets/img/integrations/components/ContourConfiguration-color.svg
		whiteIcon: assets/img/integrations/components/ContourConfiguration-white.svg
		description: 
	-	name: ContourDeployment
		colorIcon: assets/img/integrations/components/ContourDeployment-color.svg
		whiteIcon: assets/img/integrations/components/ContourDeployment-white.svg
		description: 
	-	name: ExtensionService
		colorIcon: assets/img/integrations/components/ExtensionService-color.svg
		whiteIcon: assets/img/integrations/components/ExtensionService-white.svg
		description: 
	-	name: Gateway
		colorIcon: assets/img/integrations/components/Gateway-color.svg
		whiteIcon: assets/img/integrations/components/Gateway-white.svg
		description: 
	-	name: GatewayClass
		colorIcon: assets/img/integrations/components/GatewayClass-color.svg
		whiteIcon: assets/img/integrations/components/GatewayClass-white.svg
		description: 
	-	name: HTTPProxy
		colorIcon: assets/img/integrations/components/HTTPProxy-color.svg
		whiteIcon: assets/img/integrations/components/HTTPProxy-white.svg
		description: 
	-	name: HTTPRoute
		colorIcon: assets/img/integrations/components/HTTPRoute-color.svg
		whiteIcon: assets/img/integrations/components/HTTPRoute-white.svg
		description: 
	-	name: TCPRoute
		colorIcon: assets/img/integrations/components/TCPRoute-color.svg
		whiteIcon: assets/img/integrations/components/TCPRoute-white.svg
		description: 
	-	name: TLSCertificateDelegation
		colorIcon: assets/img/integrations/components/TLSCertificateDelegation-color.svg
		whiteIcon: assets/img/integrations/components/TLSCertificateDelegation-white.svg
		description: 
	-	name: TLSRoute
		colorIcon: assets/img/integrations/components/TLSRoute-color.svg
		whiteIcon: assets/img/integrations/components/TLSRoute-white.svg
		description: 
	-	name: UDPRoute
		colorIcon: assets/img/integrations/components/UDPRoute-color.svg
		whiteIcon: assets/img/integrations/components/UDPRoute-white.svg
		description: 
featureList: [
  "Supports dynamic configuration updates out of the box while maintaining a lightweight profile.",
  "Introduces a new ingress API (HTTPProxy) which is implemented via a Custom Resource Definition (CRD).",
  "Solves shortcomings in the original design."
]
howItWorks: Collaborative Infrastructure as Code
howItWorksDetails: Builds upon the basic Kubernetes resource and controller concepts, but includes domain-specific knowledge to automate the entire lifecycle of Contour.
language: en
list: include
---
<p>
Contour is a Kubernetes ingress controller using Envoy proxy.
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
