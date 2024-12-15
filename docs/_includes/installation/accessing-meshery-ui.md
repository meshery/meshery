{% if page.display-title != "false" %}
## Accessing Meshery UI
{% endif %}

#### Accessing Meshery UI

After successfully deploying Meshery, you can access Meshery's web-based user interface. Your default browser will be automatically opened and directed to Meshery UI (default location is [http://localhost:9081](http://localhost:9081)).

You can use the following command to open Meshery UI in your default browser:

{% capture code_content %}$ mesheryctl system dashboard {% endcapture %}
{% include code.html code=code_content %}

If you have installed Meshery on Kubernetes or a remote host, you can access Meshery UI by exposing it as a Kubernetes service or by port forwarding to Meshery UI.

{% capture code_content %}$ mesheryctl system dashboard --port-forward {% endcapture %}
{% include code.html code=code_content %}

Depending upon how you have networking configured in Kubernetes, alternatively, you can use kubectl to port forward to Meshery UI.

{% capture code_content %}$ kubectl port-forward svc/meshery 9081:9081 --namespace meshery{% endcapture %}
{% include code.html code=code_content %}

#### Verify Kubernetes Connection

After installing Meshery, regardless of the installation type, it is important to verify that your kubeconfig file has been uploaded correctly via the UI. 

1. In the Meshery UI, navigate to **Lifecycle** from the menu on the left.
2. Click on Connections.
3. Ensure that your cluster appears in the list of connections and is marked as `Connected`.
4. Click on the cluster name to perform a ping test and confirm that Meshery can communicate with your cluster.

<details>
<summary>Customizing Your Meshery Provider Callback URL</summary>

<p>
  Meshery Server supports customizing your <a href="/extensibility/providers">Meshery Provider</a> authentication flow callback URL. This is helpful when deploying Meshery behind multiple layers of networking infrastructure.
</p>

<p>
  For production deployments, it is recommended to access the Meshery UI by setting up a reverse proxy or using a LoadBalancer. By specifying a custom redirect endpoint, you can ensure that authentication flows complete successfully, even when multiple routing layers are involved.
</p>

<p>
  <b>Note</b>: For production deployments, it is important to preselect the choice of <code>Remote Provider</code> in order to control which identity providers authorized. Learn more about this in the <a href="/extensibility/providers">Extensibility: Providers</a> guide.
</p>

<p>
  Define a custom callback URL by setting up the <code>MESHERY_SERVER_CALLBACK_URL</code> environment variable before installing Meshery.
</p>

<p>
  To customize the authentication flow callback URL, use the following command:
</p>

<pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">$ MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start</div></div>
</pre>

<p>
  Meshery should now be running in your Kubernetes cluster and the Meshery UI should be accessible at the <code>EXTERNAL IP</code> of the <code>meshery</code> service.
</p>

</details>