# Port Forwarding

> Access Meshery UI using port forwarding when LoadBalancer is not available.

---

LLMS index: [llms.txt](/pr-preview/pr-22666/llms.txt)

---

# Accessing Meshery UI Using Port Forwarding

Port forwarding is a universal method to access Meshery UI when a `LoadBalancer` service is not available or practical. This is common in local development environments (Minikube, KinD) or when you don't want to expose Meshery externally.

This guide is applicable to all cloud providers including AWS (EKS), Azure (AKS), GCP (GKE), and local cluster services (Minikube, KinD).

## When to Use Port Forwarding

Use port forwarding in the following scenarios:

- **Local development clusters** (Minikube, KinD) where LoadBalancer controllers are not installed
- **Testing environments** where you don't want to expose Meshery UI publicly
- **Temporary access** to Meshery UI without configuring ingress or load balancers
- **Debugging** Meshery deployments in restricted networks

## Prerequisites

- Meshery is installed and running in your Kubernetes cluster
- `kubectl` is configured to access your cluster
- Port `9081` is available on your local machine (or specify an alternative port)

## Using mesheryctl (Recommended)

The easiest way to access Meshery UI with port forwarding is using `mesheryctl`:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system dashboard --port-forward</code>
	</div>
</pre>


This command:
- Automatically sets up port forwarding from port `9081` on your local machine to the Meshery service in the cluster
- Opens Meshery UI in your default browser
- Handles service discovery automatically

### Custom Port

If port `9081` is already in use, specify an alternative local port:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system dashboard --port-forward -p 9082</code>
	</div>
</pre>


### Skip Browser Opening

To set up port forwarding without automatically opening the browser:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">mesheryctl system dashboard --port-forward --skip-browser</code>
	</div>
</pre>


## Using kubectl

For more control or when `mesheryctl` is not available, use `kubectl` directly:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl port-forward svc/meshery 9081:9081 --namespace meshery</code>
	</div>
</pre>


This command forwards:
- Local port `9081` → Service port `9081` in the `meshery` namespace

### Custom Namespace or Service Name

If you installed Meshery in a different namespace or with a different service name:



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl port-forward svc/meshery 9081:9081 --namespace &lt;your-namespace&gt;</code>
	</div>
</pre>


### Port Forward to a Specific Pod

To port forward directly to a Meshery pod (useful for debugging):



<pre class="codeblock-pre">
	<div class="codeblock">
		<code class="clipboardjs">kubectl port-forward deployment/meshery 9081:8080 --namespace meshery</code>
	</div>
</pre>


## Verify the Connection

Once port forwarding is established:

1. Open your browser and navigate to `http://localhost:9081`
2. You should see the Meshery UI login screen
3. Log in with your Meshery credentials

## Troubleshooting

### Port Already in Use

If port `9081` is already in use on your machine:

- Use an alternative port: `mesheryctl system dashboard --port-forward -p 9082`
- Or stop the process using port `9081`

### Connection Refused

If you get a connection refused error:

1. Verify Meshery is running: `kubectl get pods --namespace meshery`
2. Check the service exists: `kubectl get svc --namespace meshery`
3. Ensure the service name is `meshery` (or adjust the port-forward command accordingly)

### Timeout Errors

If port forwarding times out:

1. Check your kubeconfig context: `kubectl config current-context`
2. Verify network connectivity to the cluster
3. Check if firewall rules are blocking the connection

## LoadBalancer vs Port Forwarding

| Feature | LoadBalancer | Port Forwarding |
|---------|--------------|-----------------|
| **External Access** | Yes (cloud provider assigns external IP) | No (local access only) |
| **Production Use** | Recommended | Not recommended |
| **Setup Complexity** | Requires cloud load balancer | Simple, one command |
| **Use Case** | Production, multi-user access | Development, testing, temporary access |

For production deployments, consider using:
- **LoadBalancer** service type (cloud providers provision external IP automatically)
- **Ingress controllers** (NGINX, Traefik) for advanced routing
- **Reverse proxies** for additional security layers

## Related Documentation

- [Helm Installation](/pr-preview/pr-22666/installation/kubernetes/helm/) - Install Meshery using Helm
- [Install Meshery on Kubernetes](/pr-preview/pr-22666/installation/kubernetes/) - General Kubernetes installation guide
- [mesheryctl system dashboard](/pr-preview/pr-22666/reference/mesheryctl/system/dashboard/) - CLI reference for dashboard command
- [Accessing Meshery UI](/pr-preview/pr-22666/installation/accessing-meshery-ui/) - General guide for accessing Meshery UI

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/new-contributor-intro-vedit-agrawal-bughunterx2101-looking-to-pair-with-a-meshmate-for-meshery-ui-ai-adapter/8172" target="_blank" rel="noopener noreferrer">
            New Contributor Intro: Vedit Agrawal (@BugHunterX2101) | Looking to pair with a MeshMate for Meshery UI &amp; AI Adapter
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-seeking-an-experienced-meshery-contributor-as-a-guide/8166" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate: Seeking an Experienced Meshery Contributor as a Guide
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-to-reactivate-my-slack-account/8168" target="_blank" rel="noopener noreferrer">
            Request to Reactivate My Slack Account
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-frontend-backend-cloud-native-development/8161" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate — Frontend, Backend &amp; Cloud-Native Development
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-for-a-meshmate-active-meshery-contributor/8158" target="_blank" rel="noopener noreferrer">
            Request for a MeshMate — Active Meshery Contributor
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/bug-new-catalog-publish-request-email-renders-with-oversized-social-icons-and-squashed-content-on-mobile/8138" target="_blank" rel="noopener noreferrer">
            [Bug] &#34;New Catalog publish request&#34; email renders with oversized social icons and squashed content on mobile
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/how-to-present-an-issue-pr-in-a-development-meeting-feedback-wanted/8107" target="_blank" rel="noopener noreferrer">
            How to Present an Issue/PR in a Development Meeting: feedback wanted
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newbie-here-excited-to-join/8052" target="_blank" rel="noopener noreferrer">
            Newbie here -- Excited to join
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
