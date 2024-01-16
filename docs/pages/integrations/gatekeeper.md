---
layout: enhanced
title: OPA Gatekeeper
subtitle: Collaborative and visual infrastructure as code for Frinx Machine
image: /assets/img/integrations/opa-gatekeeper.svg
permalink: extensibility/integrations/opa-gatekeeper
docURL: https://docs.meshery.io/extensibility/integrations/gatekeeper
description: 
category: Security & Compliance
subcategory: Security & Compliance
registrant: artifacthub
components: 
	-	name: Assign
		colorIcon: assets/img/integrations/components/Assign-color.svg
		whiteIcon: assets/img/integrations/components/Assign-white.svg
		description: 
	-	name: AssignMetadata
		colorIcon: assets/img/integrations/components/AssignMetadata-color.svg
		whiteIcon: assets/img/integrations/components/AssignMetadata-white.svg
		description: 
	-	name: ConstraintTemplate
		colorIcon: assets/img/integrations/components/ConstraintTemplate-color.svg
		whiteIcon: assets/img/integrations/components/ConstraintTemplate-white.svg
		description: 
	-	name: ExpansionTemplate
		colorIcon: assets/img/integrations/components/ExpansionTemplate-color.svg
		whiteIcon: assets/img/integrations/components/ExpansionTemplate-white.svg
		description: 
	-	name: ModifySet
		colorIcon: assets/img/integrations/components/ModifySet-color.svg
		whiteIcon: assets/img/integrations/components/ModifySet-white.svg
		description: 
	-	name: MutatorPodStatus
		colorIcon: assets/img/integrations/components/MutatorPodStatus-color.svg
		whiteIcon: assets/img/integrations/components/MutatorPodStatus-white.svg
		description: 
	-	name: Provider
		colorIcon: assets/img/integrations/components/Provider-color.svg
		whiteIcon: assets/img/integrations/components/Provider-white.svg
		description: 
	-	name: AssignImage
		colorIcon: assets/img/integrations/components/AssignImage-color.svg
		whiteIcon: assets/img/integrations/components/AssignImage-white.svg
		description: 
	-	name: ExpansionTemplatePodStatus
		colorIcon: assets/img/integrations/components/ExpansionTemplatePodStatus-color.svg
		whiteIcon: assets/img/integrations/components/ExpansionTemplatePodStatus-white.svg
		description: 
featureList: [
  "Native Kubernetes CRDs for instantiating the policy library (aka constraints)",
  "An extensible, parameterized policy library",
  "Ongoing synchronization of Kubernetes configuration and changes across any number of clusters."
]
howItWorks: Collaborative Infrastructure as Code
howItWorksDetails: Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs.
language: en
list: include
---
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
