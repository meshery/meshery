## Accessing Meshery UI

After successfully deploying Meshery, you can access Meshery's web-based user interface. Your default browser will be automatically opened and directed to Meshery UI (default location is [http://localhost:9081](http://localhost:9081)).

If you are unable to access Meshery UI, you can use the following command to open Meshery UI in your default browser:

{% capture code_content %} mesheryctl system dashboard {% endcapture %}
{% include code.html code=code_content %}

If you have installed Meshery on Kubernetes or a remote host, you can access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.

{% capture code_content %} mesheryctl system dashboard --port-forward {% endcapture %}
{% include code.html code=code_content %}

Alternatively, you can use kubectl to port forward to Meshery UI. Example

{% capture code_content %}kubectl port-forward svc/meshery 9081:9081 --namespace meshery{% endcapture %}
{% include code.html code=code_content %}

Log into the [Provider](/extensibility/providers) of your choice.
