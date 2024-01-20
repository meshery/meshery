---
layout: enhanced
title: Kong API Gateway
subtitle: Collaborative and visual infrastructure as code for Kong API Gateway
image: /assets/img/integrations/kong/icons/color/kong-color.svg
permalink: extensibility/integrations/kong
docURL: https://docs.meshery.io/extensibility/integrations/kong
description: 
integrations-category: Cloud Native Network
integrations-subcategory: API Gateway
registrant: artifacthub
components: 
- name: ingress-class-parameters
  colorIcon: assets/img/integrations/kong/components/ingress-class-parameters/icons/color/ingress-class-parameters-color.svg
  whiteIcon: assets/img/integrations/kong/components/ingress-class-parameters/icons/white/ingress-class-parameters-white.svg
  description: 
- name: kong-cluster-plugin
  colorIcon: assets/img/integrations/kong/components/kong-cluster-plugin/icons/color/kong-cluster-plugin-color.svg
  whiteIcon: assets/img/integrations/kong/components/kong-cluster-plugin/icons/white/kong-cluster-plugin-white.svg
  description: 
- name: kong-consumer
  colorIcon: assets/img/integrations/kong/components/kong-consumer/icons/color/kong-consumer-color.svg
  whiteIcon: assets/img/integrations/kong/components/kong-consumer/icons/white/kong-consumer-white.svg
  description: 
- name: kong-ingress
  colorIcon: assets/img/integrations/kong/components/kong-ingress/icons/color/kong-ingress-color.svg
  whiteIcon: assets/img/integrations/kong/components/kong-ingress/icons/white/kong-ingress-white.svg
  description: 
- name: kong-plugin
  colorIcon: assets/img/integrations/kong/components/kong-plugin/icons/color/kong-plugin-color.svg
  whiteIcon: assets/img/integrations/kong/components/kong-plugin/icons/white/kong-plugin-white.svg
  description: 
- name: tcp-ingress
  colorIcon: assets/img/integrations/kong/components/tcp-ingress/icons/color/tcp-ingress-color.svg
  whiteIcon: assets/img/integrations/kong/components/tcp-ingress/icons/white/tcp-ingress-white.svg
  description: 
- name: udp-ingress
  colorIcon: assets/img/integrations/kong/components/udp-ingress/icons/color/udp-ingress-color.svg
  whiteIcon: assets/img/integrations/kong/components/udp-ingress/icons/white/udp-ingress-white.svg
  description: 
- name: kong-upstream-policy
  colorIcon: assets/img/integrations/kong/components/kong-upstream-policy/icons/color/kong-upstream-policy-color.svg
  whiteIcon: assets/img/integrations/kong/components/kong-upstream-policy/icons/white/kong-upstream-policy-white.svg
  description: 
featureList: [
  "Advanced routing, load balancing, health checking - all configurable via a RESTful admin API or declarative configuration.",
  "Authentication and authorization for APIs using methods like JWT, basic auth, OAuth, ACLs and more.",
  "Sophisticated deployment models like Declarative Databaseless Deployment and Hybrid Deployment (control plane/data plane separation) without any vendor lock-in."
]
howItWorks: "Collaborative Infrastructure as Code"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
display-title: "false"
---
<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>

<p>
Kong API Gateway is a cloud-native, platform-agnostic, scalable API Gateway distinguished for its high performance and extensibility via plugins.
</p>
<p>
By providing functionality for proxying, routing, load balancing, health checking, authentication (and more), Kong serves as the central layer for orchestrating microservices or conventional API traffic with ease.
</p>
<p>
Kong runs natively on Kubernetes thanks to its official Kubernetes Ingress Controller.</p>
