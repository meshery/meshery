---
layout: integration
title: Cilium
subtitle: The easiest way to get production-grade Kubernetes clusters with Cilium up and running
image: /assets/img/integrations/cilium/icons/color/cilium-color.svg
permalink: extensibility/integrations/cilium
docURL: https://docs.meshery.io/extensibility/adapters/cilium
description: 
integrations-category: Cloud Native Network
integrations-subcategory: Service Mesh
registrant: GitHub
components: 
- name: cilium-clusterwide-envoy-config
  colorIcon: assets/img/integrations/cilium/components/cilium-clusterwide-envoy-config/icons/color/cilium-clusterwide-envoy-config-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-clusterwide-envoy-config/icons/white/cilium-clusterwide-envoy-config-white.svg
  description: 
- name: cilium-clusterwide-network-policy
  colorIcon: assets/img/integrations/cilium/components/cilium-clusterwide-network-policy/icons/color/cilium-clusterwide-network-policy-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-clusterwide-network-policy/icons/white/cilium-clusterwide-network-policy-white.svg
  description: 
- name: cilium-egress-gateway-policy
  colorIcon: assets/img/integrations/cilium/components/cilium-egress-gateway-policy/icons/color/cilium-egress-gateway-policy-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-egress-gateway-policy/icons/white/cilium-egress-gateway-policy-white.svg
  description: 
- name: cilium-endpoint
  colorIcon: assets/img/integrations/cilium/components/cilium-endpoint/icons/color/cilium-endpoint-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-endpoint/icons/white/cilium-endpoint-white.svg
  description: 
- name: cilium-envoy-config
  colorIcon: assets/img/integrations/cilium/components/cilium-envoy-config/icons/color/cilium-envoy-config-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-envoy-config/icons/white/cilium-envoy-config-white.svg
  description: 
- name: cilium-external-workload
  colorIcon: assets/img/integrations/cilium/components/cilium-external-workload/icons/color/cilium-external-workload-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-external-workload/icons/white/cilium-external-workload-white.svg
  description: 
- name: cilium-identity
  colorIcon: assets/img/integrations/cilium/components/cilium-identity/icons/color/cilium-identity-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-identity/icons/white/cilium-identity-white.svg
  description: 
- name: cilium-local-redirect-policy
  colorIcon: assets/img/integrations/cilium/components/cilium-local-redirect-policy/icons/color/cilium-local-redirect-policy-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-local-redirect-policy/icons/white/cilium-local-redirect-policy-white.svg
  description: 
- name: cilium-network-policy
  colorIcon: assets/img/integrations/cilium/components/cilium-network-policy/icons/color/cilium-network-policy-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-network-policy/icons/white/cilium-network-policy-white.svg
  description: 
- name: cilium-node-config
  colorIcon: assets/img/integrations/cilium/components/cilium-node-config/icons/color/cilium-node-config-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-node-config/icons/white/cilium-node-config-white.svg
  description: 
- name: cilium-node
  colorIcon: assets/img/integrations/cilium/components/cilium-node/icons/color/cilium-node-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-node/icons/white/cilium-node-white.svg
  description: 
- name: cilium-bgp-advertisement
  colorIcon: assets/img/integrations/cilium/components/cilium-bgp-advertisement/icons/color/cilium-bgp-advertisement-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-bgp-advertisement/icons/white/cilium-bgp-advertisement-white.svg
  description: 
- name: cilium-bgp-cluster-config
  colorIcon: assets/img/integrations/cilium/components/cilium-bgp-cluster-config/icons/color/cilium-bgp-cluster-config-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-bgp-cluster-config/icons/white/cilium-bgp-cluster-config-white.svg
  description: 
- name: cilium-bgp-node-config-override
  colorIcon: assets/img/integrations/cilium/components/cilium-bgp-node-config-override/icons/color/cilium-bgp-node-config-override-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-bgp-node-config-override/icons/white/cilium-bgp-node-config-override-white.svg
  description: 
- name: cilium-bgp-node-config
  colorIcon: assets/img/integrations/cilium/components/cilium-bgp-node-config/icons/color/cilium-bgp-node-config-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-bgp-node-config/icons/white/cilium-bgp-node-config-white.svg
  description: 
- name: cilium-bgp-peer-config
  colorIcon: assets/img/integrations/cilium/components/cilium-bgp-peer-config/icons/color/cilium-bgp-peer-config-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-bgp-peer-config/icons/white/cilium-bgp-peer-config-white.svg
  description: 
- name: cilium-cidr-group
  colorIcon: assets/img/integrations/cilium/components/cilium-cidr-group/icons/color/cilium-cidr-group-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-cidr-group/icons/white/cilium-cidr-group-white.svg
  description: 
- name: cilium-load-balancer-ip-pool
  colorIcon: assets/img/integrations/cilium/components/cilium-load-balancer-ip-pool/icons/color/cilium-load-balancer-ip-pool-color.svg
  whiteIcon: assets/img/integrations/cilium/components/cilium-load-balancer-ip-pool/icons/white/cilium-load-balancer-ip-pool-white.svg
  description: 
components-count: 18
relationships: 
relationship-count: 0
featureList: [
  "Ensure Cilium Service Mesh daemonset operation",
  "Define microservice isolation using eBPF",
  "Visually configure and explore your Cilium Service Mesh topology"
]
howItWorks: "Collaborative Infrastructure as Design"
howItWorksDetails: "Collaboratively manage infrastructure with your coworkers synchronously sharing the same designs."
language: en
list: include
type: extensibility
category: integrations
---
