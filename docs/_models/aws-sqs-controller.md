---
layout: integration
title: AWS Simple Queuing Service
subtitle: Collaborative and visual infrastructure as design for AWS Simple Queuing Service
image: /assets/img/integrations/aws-sqs-controller/icons/color/aws-sqs-controller-color.svg
permalink: extensibility/integrations/aws-sqs-controller
docURL: https://docs.meshery.io/extensibility/integrations/aws-sqs-controller
description: 
integrations-category: App Definition and Development
integrations-subcategory: Streaming & Messaging
registrant: GitHub
components: 
- name: field-export
  colorIcon: assets/img/integrations/aws-sqs-controller/components/field-export/icons/color/field-export-color.svg
  whiteIcon: assets/img/integrations/aws-sqs-controller/components/field-export/icons/white/field-export-white.svg
  description: 
- name: adopted-resource
  colorIcon: assets/img/integrations/aws-sqs-controller/components/adopted-resource/icons/color/adopted-resource-color.svg
  whiteIcon: assets/img/integrations/aws-sqs-controller/components/adopted-resource/icons/white/adopted-resource-white.svg
  description: 
- name: queue
  colorIcon: assets/img/integrations/aws-sqs-controller/components/queue/icons/color/queue-color.svg
  whiteIcon: assets/img/integrations/aws-sqs-controller/components/queue/icons/white/queue-white.svg
  description: 
- name: iam-role-selector
  colorIcon: assets/img/integrations/aws-sqs-controller/components/iam-role-selector/icons/color/iam-role-selector-color.svg
  whiteIcon: assets/img/integrations/aws-sqs-controller/components/iam-role-selector/icons/white/iam-role-selector-white.svg
  description: 
components-count: 4
relationships: 
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship S3 sends object event notifications to SQS"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship SNS topics fan out messages to SQS queues"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship EventBridge sends events to SQS for queuing"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship EC2 applications send/receive messages from SQS queues"
- type: "non-binding"
  kind: "edge"
  description: "An edge relationship between adoptedresource and queue"
relationship-count: 5
featureList: [
  "Scalable message queuing",
  "Reliable message delivery",
  "Easy integration with other AWS services"
]
howItWorks: "Integrates with AWS SQS"
howItWorksDetails: "Simplified message queuing and decoupling of microservices on AWS"
language: en
list: include
type: extensibility
category: integrations
---
