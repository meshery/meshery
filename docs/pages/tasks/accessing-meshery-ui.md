---
layout: default
title: Accessing Meshery UI
permalink: tasks/accessing-meshery-ui
type: tasks
language: en
list: include
---

To access Meshery’s UI via port-forwarding, please refer to the port-forwarding guide for detailed instructions.

Use `mesheryctl system dashboard` to open your default browser to Meshery UI, [click here](/reference/mesheryctl/system/dashboard) to see the reference.


## Docker Desktop

Your default browser will be opened and directed to Meshery's web-based user interface typically found at `http://localhost:9081`.


## Kubernetes

Access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.

#### [Optional] Port Forward to Meshery UI

{% capture code_content %}export POD_NAME=$(kubectl get pods --namespace meshery -l "app.kubernetes.io/name=meshery,app.kubernetes.io/instance=meshery" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace meshery port-forward $POD_NAME 9081:8080{% endcapture %}
{% include code.html code=code_content %}
