---
layout: integration
title: AWS Relational Database Service
subtitle: Collaborative and visual infrastructure as design for AWS Relational Database Service
image: /assets/img/integrations/aws-rds-controller/icons/color/aws-rds-controller-color.svg
permalink: extensibility/integrations/aws-rds-controller
docURL: https://docs.meshery.io/extensibility/integrations/aws-rds-controller
description: 
integrations-category: App Definition and Development
integrations-subcategory: Database
registrant: GitHub
components: 
- name: db-cluster
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-cluster/icons/color/db-cluster-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-cluster/icons/white/db-cluster-white.svg
  description: 
- name: db-cluster-parameter-group
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-parameter-group/icons/color/db-cluster-parameter-group-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-parameter-group/icons/white/db-cluster-parameter-group-white.svg
  description: 
- name: db-instance
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-instance/icons/color/db-instance-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-instance/icons/white/db-instance-white.svg
  description: 
- name: db-parameter-group
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-parameter-group/icons/color/db-parameter-group-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-parameter-group/icons/white/db-parameter-group-white.svg
  description: 
- name: db-proxy
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-proxy/icons/color/db-proxy-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-proxy/icons/white/db-proxy-white.svg
  description: 
- name: db-subnet-group
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-subnet-group/icons/color/db-subnet-group-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-subnet-group/icons/white/db-subnet-group-white.svg
  description: 
- name: global-cluster
  colorIcon: assets/img/integrations/aws-rds-controller/components/global-cluster/icons/color/global-cluster-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/global-cluster/icons/white/global-cluster-white.svg
  description: 
- name: field-export
  colorIcon: assets/img/integrations/aws-rds-controller/components/field-export/icons/color/field-export-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/field-export/icons/white/field-export-white.svg
  description: 
- name: adopted-resource
  colorIcon: assets/img/integrations/aws-rds-controller/components/adopted-resource/icons/color/adopted-resource-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/adopted-resource/icons/white/adopted-resource-white.svg
  description: 
- name: db-snapshot
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-snapshot/icons/color/db-snapshot-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-snapshot/icons/white/db-snapshot-white.svg
  description: 
- name: db-cluster-snapshot
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-snapshot/icons/color/db-cluster-snapshot-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-snapshot/icons/white/db-cluster-snapshot-white.svg
  description: 
- name: db-cluster-snapshot
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-snapshot/icons/color/db-cluster-snapshot-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-snapshot/icons/white/db-cluster-snapshot-white.svg
  description: 
- name: db-cluster-endpoint
  colorIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-endpoint/icons/color/db-cluster-endpoint-color.svg
  whiteIcon: assets/img/integrations/aws-rds-controller/components/db-cluster-endpoint/icons/white/db-cluster-endpoint-white.svg
  description: 
components-count: 13
relationships: 
- type: "Parent"
  kind: "Hierarchical"
  description: ""
- type: "Non Binding"
  kind: "Edge"
  description: ""
- type: "Parent"
  kind: "Hierarchical"
  description: ""
- type: "Binding"
  kind: "Edge"
  description: ""
- type: "Non Binding"
  kind: "Edge"
  description: ""
- type: "Non Binding"
  kind: "Edge"
  description: ""
- type: "Non Binding"
  kind: "Edge"
  description: ""
- type: "Binding"
  kind: "Edge"
  description: ""
- type: "Parent"
  kind: "Hierarchical"
  description: ""
relationship-count: 9
featureList: [
  "Amazon RDS integrates with AWS Config to support compliance and enhance security by recording and auditing changes to the configuration of your DB instance",
  "Amazon Aurora supports quick, efficient cloning operations, where entire multi-terabyte database clusters can be cloned in minutes.",
  "Amazon RDS provides Amazon CloudWatch metrics for your database instances at no additional charge."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
