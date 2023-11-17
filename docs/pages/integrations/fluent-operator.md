---
layout: default
title: Fluent Operator
permalink: integrations/fluent-operator
type: installation
category: integrations
display-title: "false"
language: en
list: include
image: /assets/img/integrations/fluent-operator.svg
---

<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>


<!-- This needs replaced with the Category property, not the sub-category.
 #### Category: fluent-operator -->

### Overview & Features:
1. Fluent Operator provides great flexibility in building a logging layer based on Fluent Bit and Fluentd.

2. Collaborative and visual infrastructure as code for Fluent Operator

4. 
    Connect GitHub with Meshery and import selectively import your existing Helm Charts, Docker Compose applications, and Kubernetes manifests.Visually configure and customize your cloud native infrastructure.
    Save and share your design patterns to GitHub using either public or private repositories.



    Learn more about <a href="/blog/service-mesh-specifications/pipelining-service-mesh-specifications">pipelining service mesh specifications</a> and using Service Mesh Interface and Service Mesh Performance specs on your CI/CD pipelines with Meshery's GitHub Actions.



5. Fluent Bit Management: Deploy and destroy Fluent Bit DaemonSet automatically.

6. Fluentd Management: Deploy and destroy Fluentd StatefulSet automatically.

7. Custom Configuration: Select input/filter/output plugins via labels.

8. Fluent Bit will be deployed as a DaemonSet while Fluentd will be deployed as a StatefulSet. 

9. Although both Fluent Bit and Fluentd can collect, process(parse and filter) and then forward log to the final destinations, still they have strengths in different aspects.

Fluent Bit is a good choice as a logging agent because of its lightweight and efficiency, while Fluentd is more powerful to perform advanced processing on logs because of its rich plugins.

Fluent Bit only mode: If you just need to collect logs and send logs to the final destinations, all you need is Fluent Bit.
Fluent Bit + Fluentd mode: If you also need to perform some advanced processing on the logs collected or send to more sinks, then you also need Fluentd.
Fluentd only mode: If you need to receive logs through networks like HTTP or Syslog and then process and send the log to the final sinks, you only need Fluentd.
Fluent Operator includes CRDs and controllers for both Fluent Bit and Fluentd which allows you to config your log processing pipelines in the 3 modes mentioned above as you wish.

