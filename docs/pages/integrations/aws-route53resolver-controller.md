---
layout: enhanced
title: AWS Route 53 Resolver
subtitle: Collaborative and visual infrastructure as code for AWS Route 53 Resolver
image: /assets/img/integrations/aws-route-53-resolver.svg
permalink: extensibility/integrations/aws-route-53-resolver
docURL: https://docs.meshery.io/extensibility/integrations/aws-route53resolver-controller
description: 
category: Cloud Native Network
subcategory: Networking Content Delivery
registrant: aws
components: 
	-	name: FieldExport
		colorIcon: assets/img/integrations/components/FieldExport-color.svg
		whiteIcon: assets/img/integrations/components/FieldExport-white.svg
		description: 
	-	name: AdoptedResource
		colorIcon: assets/img/integrations/components/AdoptedResource-color.svg
		whiteIcon: assets/img/integrations/components/AdoptedResource-white.svg
		description: 
	-	name: ResolverEndpoint
		colorIcon: assets/img/integrations/components/ResolverEndpoint-color.svg
		whiteIcon: assets/img/integrations/components/ResolverEndpoint-white.svg
		description: 
	-	name: ResolverRule
		colorIcon: assets/img/integrations/components/ResolverRule-color.svg
		whiteIcon: assets/img/integrations/components/ResolverRule-white.svg
		description: 
featureList: [
  "Local VPC domain names for EC2 instances (for example, ec2-192-0-2-44.compute-1.amazonaws.com).
",
  "Records in private hosted zones (for example, acme.example.com).
",
  "For public domain names, Route 53 Resolver performs recursive lookups against public name servers on the internet.
"
]
howItWorks: Collaborative Infrastructure as Code
howItWorksDetails: Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs.
language: en
list: include
---
<p>
Amazon Route 53 Resolver responds recursively to DNS queries from AWS resources for public records, Amazon VPC-specific DNS names, and Amazon Route 53 private hosted zones, and is available by default in all VPCs.
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
