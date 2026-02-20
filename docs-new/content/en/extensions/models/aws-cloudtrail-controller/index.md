---
title: AWS CloudTrail
subtitle: Collaborative and visual infrastructure as design for AWS CloudTrail
image: /extensions/models/aws-cloudtrail-controller/icons/color/aws-cloudtrail-controller-color.svg
docURL: https://docs.meshery.io/extensibility/integrations/aws-cloudtrail-controller
description: 
integrations-category: Observability and Analysis
integrations-subcategory: Logging
registrant: GitHub
components: 
- name: field-export
  colorIcon: /extensions/models/aws-cloudtrail-controller/components/field-export/icons/color/field-export-color.svg
  whiteIcon: /extensions/models/aws-cloudtrail-controller/components/field-export/icons/white/field-export-white.svg
  description: 
- name: adopted-resource
  colorIcon: /extensions/models/aws-cloudtrail-controller/components/adopted-resource/icons/color/adopted-resource-color.svg
  whiteIcon: /extensions/models/aws-cloudtrail-controller/components/adopted-resource/icons/white/adopted-resource-white.svg
  description: 
- name: event-data-store
  colorIcon: /extensions/models/aws-cloudtrail-controller/components/event-data-store/icons/color/event-data-store-color.svg
  whiteIcon: /extensions/models/aws-cloudtrail-controller/components/event-data-store/icons/white/event-data-store-white.svg
  description: 
- name: trail
  colorIcon: /extensions/models/aws-cloudtrail-controller/components/trail/icons/color/trail-color.svg
  whiteIcon: /extensions/models/aws-cloudtrail-controller/components/trail/icons/white/trail-white.svg
  description: 
- name: iam-role-selector
  colorIcon: /extensions/models/aws-cloudtrail-controller/components/iam-role-selector/icons/color/iam-role-selector-color.svg
  whiteIcon: /extensions/models/aws-cloudtrail-controller/components/iam-role-selector/icons/white/iam-role-selector-white.svg
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
- type: "non-binding"
  kind: "edge"
  description: ""
relationship-count: 6
featureList: [
  "Data events that capture data plane actions within a resource, such as reading or writing an Amazon S3 object.",
  "Configuration items from AWS Config that capture resource configuration history and resource compliance history as evaluated by AWS Config rules.",
  "Audit evidence from AWS Audit Manager that contains the information needed to demonstrate compliance with the requirements as specified by Audit Manager controls."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
categories: [integrations]
aliases:
- /extensibility/integrations/aws-cloudtrail-controller
---
