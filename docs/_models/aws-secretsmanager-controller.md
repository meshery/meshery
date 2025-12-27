---
layout: integration
title: AWS Secrets Manager
subtitle: Collaborative and visual infrastructure as design for AWS Secrets Manager
image: /assets/img/integrations/aws-secretsmanager-controller/icons/color/aws-secretsmanager-controller-color.svg
permalink: extensibility/integrations/aws-secretsmanager-controller
docURL: https://docs.meshery.io/extensibility/integrations/aws-secretsmanager-controller
description: 
integrations-category: Security & Compliance
integrations-subcategory: Security Identity Compliance
registrant: GitHub
components: 
- name: field-export
  colorIcon: assets/img/integrations/aws-secretsmanager-controller/components/field-export/icons/color/field-export-color.svg
  whiteIcon: assets/img/integrations/aws-secretsmanager-controller/components/field-export/icons/white/field-export-white.svg
  description: 
- name: adopted-resource
  colorIcon: assets/img/integrations/aws-secretsmanager-controller/components/adopted-resource/icons/color/adopted-resource-color.svg
  whiteIcon: assets/img/integrations/aws-secretsmanager-controller/components/adopted-resource/icons/white/adopted-resource-white.svg
  description: 
- name: secret
  colorIcon: assets/img/integrations/aws-secretsmanager-controller/components/secret/icons/color/secret-color.svg
  whiteIcon: assets/img/integrations/aws-secretsmanager-controller/components/secret/icons/white/secret-white.svg
  description: 
- name: iam-role-selector
  colorIcon: assets/img/integrations/aws-secretsmanager-controller/components/iam-role-selector/icons/color/iam-role-selector-color.svg
  whiteIcon: assets/img/integrations/aws-secretsmanager-controller/components/iam-role-selector/icons/white/iam-role-selector-white.svg
  description: 
components-count: 4
relationships: 
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship RDS stores master credentials in Secrets Manager"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship Lambda retrieves secrets from Secrets Manager"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship ECS tasks fetch secrets at runtime"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship EKS uses Secrets Store CSI driver to mount secrets"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between adoptedresourcex and Secret "
relationship-count: 5
featureList: [
  "Centrally store and manage credentials, API keys, and other secrets.",
  "Use AWS Identity and Access Management (IAM) permissions policies to manage access to your secrets.",
  "Rotate secrets on demand or on a schedule, without redeploying or disrupting active applications."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
