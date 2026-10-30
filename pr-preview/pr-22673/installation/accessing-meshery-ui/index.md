# Accessing Meshery UI

> Accessing Meshery's web-based user interface

---

LLMS index: [llms.txt](/pr-preview/pr-22673/llms.txt)

---

<h2>Accessing Meshery UI</h2>


<p>Once Meshery is running, <code>mesheryctl system start</code> opens Meshery UI in your default browser. On Docker, Meshery UI listens at <a href="http://localhost:9081">http://localhost:9081</a> by default. On Kubernetes, <code>mesheryctl</code> opens the address of the <code>meshery</code> service in the <code>meshery</code> namespace, which the Helm chart creates as a <code>LoadBalancer</code> service on port 9081.</p>

<p>To open Meshery UI in your default browser at any time, run:</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">mesheryctl system dashboard</code>
</div>
</pre>

<p>If your cluster does not give the <code>meshery</code> service a reachable address, forward a local port to Meshery instead. By default, <code>mesheryctl</code> listens on local port 9081; use <code>-p</code> to choose another port.</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">mesheryctl system dashboard --port-forward</code>
</div>
</pre>

<p>You can also port-forward with <code>kubectl</code>:</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">kubectl port-forward svc/meshery 9081:9081 --namespace meshery</code>
</div>
</pre>

<h4>Verify Kubernetes Connection</h4>

Whichever way you installed Meshery, verify that Meshery is connected to your cluster:

<ol>
<li>In the Meshery UI, navigate to <strong>Lifecycle</strong> from the menu on the left.</li>
<li>Select <strong>Connections</strong>.</li>
<li>Ensure that your cluster appears in the list of connections and is marked as <code>Connected</code>.</li>
<li>Click the cluster name to run a connectivity test and confirm that Meshery can reach your cluster.</li>
</ol>

<details>
<summary>Customizing Your Meshery Provider Callback URL</summary>

<p>
  You can customize the callback URL that your <a href="/pr-preview/pr-22673/reference/extensibility/providers/">Meshery Provider</a> redirects to after sign-in. This helps when Meshery runs behind several layers of networking infrastructure.
</p>

<p>
  For production deployments, we recommend that you expose Meshery UI through a reverse proxy or a load balancer. A custom callback URL makes sure that sign-in completes even when requests pass through several routing layers.
</p>

<p>
  <b>Note</b>: For production deployments, use a Remote Provider so that you control which identity providers can sign users in. Learn more about this in the <a href="/pr-preview/pr-22673/reference/extensibility/providers/">Extensibility: Providers</a> guide.
</p>

<p>
  To set a custom callback URL, set the <code>MESHERY_SERVER_CALLBACK_URL</code> environment variable when you start Meshery:
</p>

<pre class="codeblock-pre" style="font-size: 1.75rem !important;">
<div class="codeblock">
<code class="clipboardjs">MESHERY_SERVER_CALLBACK_URL=https://custom-host mesheryctl system start</code>
</div>
</pre>

<p>
  On Kubernetes, Meshery UI is then available at the <code>EXTERNAL-IP</code> of the <code>meshery</code> service. Run <code>kubectl get svc meshery -n meshery</code> to find it.
</p>

</details>
