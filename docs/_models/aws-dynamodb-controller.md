---
layout: integration
title: AWS DynamoDB
subtitle: Collaborative and visual infrastructure as design for AWS DynamoDB
image: /assets/img/integrations/aws-dynamodb-controller/icons/color/aws-dynamodb-controller-color.svg
permalink: extensibility/integrations/aws-dynamodb-controller
docURL: https://docs.meshery.io/extensibility/integrations/aws-dynamodb-controller
description: 
integrations-category: App Definition and Development
integrations-subcategory: Database
registrant: GitHub
components: 
- name: field-export
  colorIcon: assets/img/integrations/aws-dynamodb-controller/components/field-export/icons/color/field-export-color.svg
  whiteIcon: assets/img/integrations/aws-dynamodb-controller/components/field-export/icons/white/field-export-white.svg
  description: 
- name: adopted-resource
  colorIcon: assets/img/integrations/aws-dynamodb-controller/components/adopted-resource/icons/color/adopted-resource-color.svg
  whiteIcon: assets/img/integrations/aws-dynamodb-controller/components/adopted-resource/icons/white/adopted-resource-white.svg
  description: 
- name: backup
  colorIcon: assets/img/integrations/aws-dynamodb-controller/components/backup/icons/color/backup-color.svg
  whiteIcon: assets/img/integrations/aws-dynamodb-controller/components/backup/icons/white/backup-white.svg
  description: 
- name: table
  colorIcon: assets/img/integrations/aws-dynamodb-controller/components/table/icons/color/table-color.svg
  whiteIcon: assets/img/integrations/aws-dynamodb-controller/components/table/icons/white/table-white.svg
  description: 
- name: global-table
  colorIcon: assets/img/integrations/aws-dynamodb-controller/components/global-table/icons/color/global-table-color.svg
  whiteIcon: assets/img/integrations/aws-dynamodb-controller/components/global-table/icons/white/global-table-white.svg
  description: 
- name: iam-role-selector
  colorIcon: assets/img/integrations/aws-dynamodb-controller/components/iam-role-selector/icons/color/iam-role-selector-color.svg
  whiteIcon: assets/img/integrations/aws-dynamodb-controller/components/iam-role-selector/icons/white/iam-role-selector-white.svg
  description: 
components-count: 6
relationships: 
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship Lambda triggered by DynamoDB Streams for real-time processing"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between table and VPCEndpoint"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship EC2 instances use DynamoDB SDK for NoSQL operations"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship ECS services use DynamoDB for scalable NoSQL storage"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between AdoptedResource and GlobalTable"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between AdoptedResource and GlobalTable"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between AdoptedResource and Table"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between Backup and Table"
- type: "parent"
  kind: "hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of (parent component) is patched with the configuration of (child component). "
relationship-count: 9
featureList: [
  "Handle more than 10 trillion requests per day and can support peaks of more than 20 million requests per second.
",
  "Secure your data with encryption at rest, automatic backup and restore, and guaranteed reliability with an SLA of up to 99.999% availability.",
  "Fast and flexible NoSQL database service for any scale"
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
