---
layout: enhanced
title: Fluentd Operator
subtitle: Collaborative and visual infrastructure as code for Fluentd Operator
image: /assets/img/integrations/fluentd-operator/icons/color/fluentd-operator-color.svg
permalink: extensibility/integrations/fluentd-operator
docURL: https://docs.meshery.io/extensibility/integrations/fluentd-operator
description: 
integrations-category: Observability and Analysis
integrations-subcategory: Logging
registrant: Artifact Hub
components: 
- name: fluentd-config
  colorIcon: assets/img/integrations/fluentd-operator/components/fluentd-config/icons/color/fluentd-config-color.svg
  whiteIcon: assets/img/integrations/fluentd-operator/components/fluentd-config/icons/white/fluentd-config-white.svg
  description: 
featureList: [
  "Fluent Operator provides great flexibility in building a logging layer based on Fluent Bit and Fluentd.",
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
Fluentd: Unified Logging Layer (project under CNCF)
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
