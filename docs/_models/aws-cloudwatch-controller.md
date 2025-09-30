---
layout: integration
title: AWS CloudWatch
subtitle: Collaborative and visual infrastructure as design for AWS CloudWatch
image: /assets/img/integrations/aws-cloudwatch-controller/icons/color/aws-cloudwatch-controller-color.svg
permalink: extensibility/integrations/aws-cloudwatch-controller
docURL: https://docs.meshery.io/extensibility/integrations/aws-cloudwatch-controller
description: 
integrations-category: Observability and Analysis
integrations-subcategory: Management Governance
registrant: GitHub
components: 
- name: field-export
  colorIcon: assets/img/integrations/aws-cloudwatch-controller/components/field-export/icons/color/field-export-color.svg
  whiteIcon: assets/img/integrations/aws-cloudwatch-controller/components/field-export/icons/white/field-export-white.svg
  description: 
- name: adopted-resource
  colorIcon: assets/img/integrations/aws-cloudwatch-controller/components/adopted-resource/icons/color/adopted-resource-color.svg
  whiteIcon: assets/img/integrations/aws-cloudwatch-controller/components/adopted-resource/icons/white/adopted-resource-white.svg
  description: 
- name: metric-alarm
  colorIcon: assets/img/integrations/aws-cloudwatch-controller/components/metric-alarm/icons/color/metric-alarm-color.svg
  whiteIcon: assets/img/integrations/aws-cloudwatch-controller/components/metric-alarm/icons/white/metric-alarm-white.svg
  description: 
- name: metric-stream
  colorIcon: assets/img/integrations/aws-cloudwatch-controller/components/metric-stream/icons/color/metric-stream-color.svg
  whiteIcon: assets/img/integrations/aws-cloudwatch-controller/components/metric-stream/icons/white/metric-stream-white.svg
  description: 
components-count: 4
relationships: 
- type: "Binding"
  kind: "Edge"
  description: "An edge relationship between MetricAlarm and Function"
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between AdoptedResource and MetricAlarm"
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between MetricStream and Instance"
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between AdoptedResource and MetricStream"
- type: "Non Binding"
  kind: "Edge"
  description: "An edge relationship between MetricAlarm and Instance"
relationship-count: 5
featureList: [
  "Provides you with data and actionable insights to monitor your applications, respond to system-wide performance changes, and optimize resource utilization.",
  "Collects monitoring and operational data in the form of logs, metrics, and traces.",
  "Get a unified view of operational health and gain complete visibility of your AWS resources, applications, and services running on AWS and on-premises."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
