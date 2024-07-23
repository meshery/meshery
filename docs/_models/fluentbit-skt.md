---
layout: integration
title: Fluentbit Skt
subtitle: Collaborative and visual infrastructure as design for Fluentbit Skt
image: /assets/img/integrations/fluentbit-skt/icons/color/fluentbit-skt-color.svg
permalink: extensibility/integrations/fluentbit-skt
docURL: https://docs.meshery.io/extensibility/integrations/fluentbit-skt
description: 
integrations-category: Observability and Analysis
integrations-subcategory: Logging
registrant: Artifact Hub
components: 
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
---
