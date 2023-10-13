---
layout: default
title: Port-Forwarding
permalink: services/port-forward
type: services
display-title: "true"
language: en
list: exclude
---
<div class="prereqs"><p><strong style="font-size: 20px;">Prerequisites</strong> </p> 
  <ol>
    <li>Make sure <a href="https://kubernetes.io/docs/tasks/tools/" class="meshery-light">kubectl</a> is installed on your local machine.</li>
    <li>You should have access to the cluster/platform where Meshery is deployed.</li>
    <li>Ensure that the kubeconfig file has the correct current context/cluster configuration.</li>
  </ol>
</div>

## Port Forward to Meshery UI

{% capture code_content %}kubectl port-forward svc/meshery --namespace meshery 9081:9081{% endcapture %}
{% include code.html code=code_content %}

Meshery UI should now be accessible at [http://localhost:9081](http://localhost:9081).