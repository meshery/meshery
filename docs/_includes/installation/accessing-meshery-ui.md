{% if page.display-title != "false" %}
## Accessing Meshery UI
{% endif %}

After successfully deploying Meshery, you can access Meshery's web-based user interface. Your default browser will be automatically opened and directed to Meshery UI (default location is [http://localhost:9081](http://localhost:9081)).

You can use the following command to open Meshery UI in your default browser:

{% capture code_content %} mesheryctl system dashboard {% endcapture %}
{% include code.html code=code_content %}

If you have installed Meshery on Kubernetes or a remote host, you can access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.

{% capture code_content %} mesheryctl system dashboard --port-forward {% endcapture %}
{% include code.html code=code_content %}

Depending upon how you have networking configured in Kubernetes, alternatively, you can use kubectl to port forward to Meshery UI.

{% capture code_content %}kubectl port-forward svc/meshery 9081:9081 --namespace meshery{% endcapture %}
{% include code.html code=code_content %}

<details>
<summary>Customizing Meshery Provider Callback URL</summary>

Customize your Meshery Provider Callback URL. Meshery Server supports customizing authentication flow callback URL, which can be configured in the following way:

{% capture code_content %}$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start{% endcapture %}
{% include code.html code=code_content %}

<br />
Meshery should now be running in your Kubernetes cluster and Meshery UI should be accessible at the `EXTERNAL IP` of `meshery` service.

</details>

Production deployments are recommended to access Meshery UI by setting up a reverse proxy or by using a LoadBalancer.

Log into the [Provider](/extensibility/providers) of your choice.
