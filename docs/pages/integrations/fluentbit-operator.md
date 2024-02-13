---
layout: enhanced
title: Fluentbit Operator
subtitle: Collaborative and visual infrastructure as code for Fluentbit Operator
image: /assets/img/integrations/fluentbit-operator/icons/color/fluentbit-operator-color.svg
permalink: extensibility/integrations/fluentbit-operator
docURL: https://docs.meshery.io/extensibility/integrations/fluentbit-operator
description: 
integrations-category: Observability and Analysis
integrations-subcategory: Logging
registrant: Artifact Hub
components: 
- name: filter
  colorIcon: assets/img/integrations/fluentbit-operator/components/filter/icons/color/filter-color.svg
  whiteIcon: assets/img/integrations/fluentbit-operator/components/filter/icons/white/filter-white.svg
  description: 
- name: fluent-bit-config
  colorIcon: assets/img/integrations/fluentbit-operator/components/fluent-bit-config/icons/color/fluent-bit-config-color.svg
  whiteIcon: assets/img/integrations/fluentbit-operator/components/fluent-bit-config/icons/white/fluent-bit-config-white.svg
  description: 
- name: input
  colorIcon: assets/img/integrations/fluentbit-operator/components/input/icons/color/input-color.svg
  whiteIcon: assets/img/integrations/fluentbit-operator/components/input/icons/white/input-white.svg
  description: 
- name: output
  colorIcon: assets/img/integrations/fluentbit-operator/components/output/icons/color/output-color.svg
  whiteIcon: assets/img/integrations/fluentbit-operator/components/output/icons/white/output-white.svg
  description: 
- name: parser
  colorIcon: assets/img/integrations/fluentbit-operator/components/parser/icons/color/parser-color.svg
  whiteIcon: assets/img/integrations/fluentbit-operator/components/parser/icons/white/parser-white.svg
  description: 
- name: fluent-bit
  colorIcon: assets/img/integrations/fluentbit-operator/components/fluent-bit/icons/color/fluent-bit-color.svg
  whiteIcon: assets/img/integrations/fluentbit-operator/components/fluent-bit/icons/white/fluent-bit-white.svg
  description: 
featureList: [
  "Fluent Bit Management: Deploy and destroy Fluent Bit DaemonSet automatically.",
  "Fluentd Management: Deploy and destroy Fluentd StatefulSet automatically.",
  "Custom Configuration: Select input/filter/output plugins via labels."
]
howItWorks: "Fluent Bit will be deployed as a DaemonSet while Fluentd will be deployed as a StatefulSet. "
howItWorksDetails: "Although both Fluent Bit and Fluentd can collect, process(parse and filter) and then forward log to the final destinations, still they have strengths in different aspects.

Fluent Bit is a good choice as a logging agent because of its lightweight and efficiency, while Fluentd is more powerful to perform advanced processing on logs because of its rich plugins.

Fluent Bit only mode: If you just need to collect logs and send logs to the final destinations, all you need is Fluent Bit.
Fluent Bit + Fluentd mode: If you also need to perform some advanced processing on the logs collected or send to more sinks, then you also need Fluentd.
Fluentd only mode: If you need to receive logs through networks like HTTP or Syslog and then process and send the log to the final sinks, you only need Fluentd.
Fluent Operator includes CRDs and controllers for both Fluent Bit and Fluentd which allows you to config your log processing pipelines in the 3 modes mentioned above as you wish."
language: en
list: include
type: extensibility
category: integrations
display-title: "false"
---
<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>

<p>
Fluent Bit is a super fast, lightweight, and highly scalable logging and metrics processor and forwarder. It is the preferred choice for cloud and containerized environments.
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
