---
layout: integration
title: Azure Cache
subtitle: Collaborative and visual infrastructure as design for Azure Cache
image: /assets/img/integrations/azure-cache/icons/color/azure-cache-color.svg
permalink: extensibility/integrations/azure-cache
docURL: https://docs.meshery.io/extensibility/integrations/azure-cache
description: 
integrations-category: Cloud Native Storage
integrations-subcategory: Content Delivery Network
registrant: GitHub
components: 
- name: redis
  colorIcon: assets/img/integrations/azure-cache/components/redis/icons/color/redis-color.svg
  whiteIcon: assets/img/integrations/azure-cache/components/redis/icons/white/redis-white.svg
  description: 
- name: redis-enterprise-database
  colorIcon: assets/img/integrations/azure-cache/components/redis-enterprise-database/icons/color/redis-enterprise-database-color.svg
  whiteIcon: assets/img/integrations/azure-cache/components/redis-enterprise-database/icons/white/redis-enterprise-database-white.svg
  description: 
- name: redis-enterprise
  colorIcon: assets/img/integrations/azure-cache/components/redis-enterprise/icons/color/redis-enterprise-color.svg
  whiteIcon: assets/img/integrations/azure-cache/components/redis-enterprise/icons/white/redis-enterprise-white.svg
  description: 
- name: redis-firewall-rule
  colorIcon: assets/img/integrations/azure-cache/components/redis-firewall-rule/icons/color/redis-firewall-rule-color.svg
  whiteIcon: assets/img/integrations/azure-cache/components/redis-firewall-rule/icons/white/redis-firewall-rule-white.svg
  description: 
- name: redis-linked-server
  colorIcon: assets/img/integrations/azure-cache/components/redis-linked-server/icons/color/redis-linked-server-color.svg
  whiteIcon: assets/img/integrations/azure-cache/components/redis-linked-server/icons/white/redis-linked-server-white.svg
  description: 
- name: redis-patch-schedule
  colorIcon: assets/img/integrations/azure-cache/components/redis-patch-schedule/icons/color/redis-patch-schedule-color.svg
  whiteIcon: assets/img/integrations/azure-cache/components/redis-patch-schedule/icons/white/redis-patch-schedule-white.svg
  description: 
components-count: 6
relationships: 
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of Redis Cache(parent component) is patched with the configuration of RedisFirewallRule(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of Redis Cache(parent component) is patched with the configuration of RedisLinkedServer(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of Redis Enterprise Cache(parent component) is patched with the configuration of RedisEnterpriseDatabase(child component). "
- type: "Parent"
  kind: "Hierarchical"
  description: "A hierarchical inventory relationship in which the configuration of Redis Cache(parent component) is patched with the configuration of RedisPatchSchedule(child component). "
relationship-count: 4
featureList: [
  "Drag-n-drop cloud native infrastructure designer to configure, model, and deploy your workloads.",
  "Invite anyone to review and make changes to your private designs.",
  "Ongoing synchronization of Kubernetes configuration and changes across any number of clusters."
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
