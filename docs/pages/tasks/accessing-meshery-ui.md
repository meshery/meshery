---
layout: default
title: Accessing Meshery UI
permalink: tasks/accessing-meshery-ui
type: tasks
language: en
list: include
abstract: "How to find where Meshery UI is exposed from your Kubernetes cluster."
---

To access Meshery’s UI via port-forwarding, please refer to the port-forwarding guide for detailed instructions.

Use `mesheryctl system dashboard` to open your default browser to Meshery UI, [click here](/reference/mesheryctl/system/dashboard) to see the reference.


## Docker Desktop

Your default browser will be opened and directed to Meshery's web-based user interface typically found at `http://localhost:9081`.


## Kubernetes

Access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.

#### [Optional] Port Forward to Meshery UI

{% capture code_content %}kubectl port-forward svc/meshery 9081:9081 --namespace meshery{% endcapture %}
{% include code.html code=code_content %}
