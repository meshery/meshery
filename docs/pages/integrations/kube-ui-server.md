---
layout: enhanced
title: Kubernetes UI Server
subtitle: Collaborative and visual infrastructure as code for Kubernetes UI Server
image: /assets/img/integrations/kubernetes-ui-server.svg
permalink: extensibility/integrations/kubernetes-ui-server
docURL: https://docs.meshery.io/extensibility/integrations/kube-ui-server
description: 
category: Provisioning
subcategory: Automation & Configuration
registrant: artifacthub
components: 
	-	name: AppBinding
		colorIcon: assets/img/integrations/components/AppBinding-color.svg
		whiteIcon: assets/img/integrations/components/AppBinding-white.svg
		description: 
	-	name: HelmRelease
		colorIcon: assets/img/integrations/components/HelmRelease-color.svg
		whiteIcon: assets/img/integrations/components/HelmRelease-white.svg
		description: 
	-	name: HelmRepository
		colorIcon: assets/img/integrations/components/HelmRepository-color.svg
		whiteIcon: assets/img/integrations/components/HelmRepository-white.svg
		description: 
	-	name: Feature
		colorIcon: assets/img/integrations/components/Feature-color.svg
		whiteIcon: assets/img/integrations/components/Feature-white.svg
		description: 
	-	name: FeatureSet
		colorIcon: assets/img/integrations/components/FeatureSet-color.svg
		whiteIcon: assets/img/integrations/components/FeatureSet-white.svg
		description: 
	-	name: ResourceDashboard
		colorIcon: assets/img/integrations/components/ResourceDashboard-color.svg
		whiteIcon: assets/img/integrations/components/ResourceDashboard-white.svg
		description: 
	-	name: ResourceEditor
		colorIcon: assets/img/integrations/components/ResourceEditor-color.svg
		whiteIcon: assets/img/integrations/components/ResourceEditor-white.svg
		description: 
	-	name: ChartPreset
		colorIcon: assets/img/integrations/components/ChartPreset-color.svg
		whiteIcon: assets/img/integrations/components/ChartPreset-white.svg
		description: 
	-	name: ClusterChartPreset
		colorIcon: assets/img/integrations/components/ClusterChartPreset-color.svg
		whiteIcon: assets/img/integrations/components/ClusterChartPreset-white.svg
		description: 
	-	name: ProjectQuota
		colorIcon: assets/img/integrations/components/ProjectQuota-color.svg
		whiteIcon: assets/img/integrations/components/ProjectQuota-white.svg
		description: 
featureList: [
  "WhoAmI service returns the user info of the user making the api call.",
  "PodView resource exposes actual resource usage by a Pod. The resource usage information is read from Prometheus.",
  "Identity Server is a Kubernetes extended apiserver (EAS). As an EAS, it has access to the user who is making an api call to the whoami server."
]
howItWorks: Collaborative Infrastructure as Code
howItWorksDetails: Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs.
language: en
list: include
---
<p>
Kubernetes UI Server is an extended api server for Kubernetes. This exposes a number of apis for a Kubernetes cluster, such as:
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
