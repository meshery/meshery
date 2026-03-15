---
title: AWS Route 53 Resolver
subtitle: Collaborative and visual infrastructure as design for AWS Route 53 Resolver
image: /extensions/models/aws-route53resolver-controller/icons/color/aws-route53resolver-controller-color.svg
docURL: https://docs.meshery.io/extensibility/integrations/aws-route53resolver-controller
description: 
integrations-category: Cloud Native Network
integrations-subcategory: Networking Content Delivery
registrant: GitHub
components: 
- name: field-export
  colorIcon: /extensions/models/aws-route53resolver-controller/components/field-export/icons/color/field-export-color.svg
  whiteIcon: /extensions/models/aws-route53resolver-controller/components/field-export/icons/white/field-export-white.svg
  description: 
- name: adopted-resource
  colorIcon: /extensions/models/aws-route53resolver-controller/components/adopted-resource/icons/color/adopted-resource-color.svg
  whiteIcon: /extensions/models/aws-route53resolver-controller/components/adopted-resource/icons/white/adopted-resource-white.svg
  description: 
- name: resolver-endpoint
  colorIcon: /extensions/models/aws-route53resolver-controller/components/resolver-endpoint/icons/color/resolver-endpoint-color.svg
  whiteIcon: /extensions/models/aws-route53resolver-controller/components/resolver-endpoint/icons/white/resolver-endpoint-white.svg
  description: 
- name: resolver-rule
  colorIcon: /extensions/models/aws-route53resolver-controller/components/resolver-rule/icons/color/resolver-rule-color.svg
  whiteIcon: /extensions/models/aws-route53resolver-controller/components/resolver-rule/icons/white/resolver-rule-white.svg
  description: 
- name: iam-role-selector
  colorIcon: /extensions/models/aws-route53resolver-controller/components/iam-role-selector/icons/color/iam-role-selector-color.svg
  whiteIcon: /extensions/models/aws-route53resolver-controller/components/iam-role-selector/icons/white/iam-role-selector-white.svg
  description: 
components-count: 5
relationships: 
- type: "non-binding"
  kind: "edge"
  description: ""
- type: "non-binding"
  kind: "edge"
  description: ""
- type: "non-binding"
  kind: "edge"
  description: ""
- type: "non-binding"
  kind: "edge"
  description: ""
- type: "non-binding"
  kind: "edge"
  description: ""
relationship-count: 5
featureList: [
  "Local VPC domain names for EC2 instances (for example, ec2-192-0-2-44.compute-1.amazonaws.com).
",
  "Records in private hosted zones (for example, acme.example.com).
",
  "For public domain names, Route 53 Resolver performs recursive lookups against public name servers on the internet.
"
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
categories: [integrations]
aliases:
- /extensibility/integrations/aws-route53resolver-controller
---
