---
layout: enhanced
title: Fluentbit Operator
subtitle: Collaborative and visual infrastructure as code for Fluentbit Operator
image: /assets/img/integrations/fluentbit-operator.svg
permalink: extensibility/integrations/fluentbit-operator
docURL: https://docs.meshery.io/extensibility/integrations/fluentbit-operator
description: 
category: Observability and Analysis
subcategory: Logging
registrant: artifacthub
components: 
	-	name: Filter
		colorIcon: assets/img/integrations/components/Filter-color.svg
		whiteIcon: assets/img/integrations/components/Filter-white.svg
		description: 
	-	name: FluentBitConfig
		colorIcon: assets/img/integrations/components/FluentBitConfig-color.svg
		whiteIcon: assets/img/integrations/components/FluentBitConfig-white.svg
		description: 
	-	name: Input
		colorIcon: assets/img/integrations/components/Input-color.svg
		whiteIcon: assets/img/integrations/components/Input-white.svg
		description: 
	-	name: Output
		colorIcon: assets/img/integrations/components/Output-color.svg
		whiteIcon: assets/img/integrations/components/Output-white.svg
		description: 
	-	name: Parser
		colorIcon: assets/img/integrations/components/Parser-color.svg
		whiteIcon: assets/img/integrations/components/Parser-white.svg
		description: 
	-	name: FluentBit
		colorIcon: assets/img/integrations/components/FluentBit-color.svg
		whiteIcon: assets/img/integrations/components/FluentBit-white.svg
		description: 
featureList: [
  "Fluent Bit Management: Deploy and destroy Fluent Bit DaemonSet automatically.",
  "Fluentd Management: Deploy and destroy Fluentd StatefulSet automatically.",
  "Custom Configuration: Select input/filter/output plugins via labels."
]
howItWorks: Fluent Bit will be deployed as a DaemonSet while Fluentd will be deployed as a StatefulSet. 
howItWorksDetails: Although both Fluent Bit and Fluentd can collect, process(parse and filter) and then forward log to the final destinations, still they have strengths in different aspects.

Fluent Bit is a good choice as a logging agent because of its lightweight and efficiency, while Fluentd is more powerful to perform advanced processing on logs because of its rich plugins.

Fluent Bit only mode: If you just need to collect logs and send logs to the final destinations, all you need is Fluent Bit.
Fluent Bit + Fluentd mode: If you also need to perform some advanced processing on the logs collected or send to more sinks, then you also need Fluentd.
Fluentd only mode: If you need to receive logs through networks like HTTP or Syslog and then process and send the log to the final sinks, you only need Fluentd.
Fluent Operator includes CRDs and controllers for both Fluent Bit and Fluentd which allows you to config your log processing pipelines in the 3 modes mentioned above as you wish.
language: en
list: include
---
<p>
Fluent Bit is a super fast, lightweight, and highly scalable logging and metrics processor and forwarder. It is the preferred choice for cloud and containerized environments.
</p>
<p>
    Collaboratively and visually diagram your cloud native infrastructure with GitOps-style pipeline integration. Design, test, and manage configuration your Kubernetes-based, containerized applications as a visual topology.
</p>
<p>
    Looking for best practice cloud native design and deployment best practices? Choose from thousands of pre-built components in MeshMap. Choose from hundreds of ready-made design patterns by importing templates from Meshery Catalog or use our low code designer, MeshMap, to create and deploy your own cloud native infrastructure designs.
</p>
