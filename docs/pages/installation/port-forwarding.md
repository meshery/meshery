---
layout: default
title: Port-Forwarding
permalink: services/port-forward
type: services
display-title: "true"
language: en
list: exclude
---

## Prerequisites:
1. Make sure [kubectl](https://kubernetes.io/docs/tasks/tools/) is installed on your local machine.
2. You should have access to the location where Meshery is deployed.
3. Ensure that the kubeconfig file has the correct current context/cluster configuration.
<br />

## Port Forward to Meshery UI

{% capture code_content %}kubectl port-forward svc/meshery --namespace meshery 9081:9081{% endcapture %}
{% include code.html code=code_content %}

Meshery UI should now be accessible at [http://localhost:9081](http://localhost:9081).