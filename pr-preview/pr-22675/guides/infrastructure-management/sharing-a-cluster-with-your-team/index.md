# Sharing a Cluster with Your Team

> Give your teammates access to a Kubernetes cluster you manage in Meshery by grouping it into an Environment, attaching the Environment to a Workspace, and assigning teams to the Workspace.

---

LLMS index: [llms.txt](/pr-preview/pr-22675/llms.txt)

---

A Kubernetes cluster you connect to Meshery is private to you until you share it. Use this guide to let your teammates see the cluster, deploy designs to it, and work with its resources, without handing out your kubeconfig. You share a cluster in three steps:

1. Put the cluster's [Connection](/pr-preview/pr-22675/concepts/logical/connections/) in an [Environment](/pr-preview/pr-22675/concepts/logical/environments/).
2. Attach the Environment to a [Workspace](/pr-preview/pr-22675/concepts/logical/workspaces/).
3. Assign your teammates' teams to the Workspace.

## Before you begin

- A Meshery Server signed in with a Remote Provider, such as Meshery Cloud. Organizations, teams, and other users come from the Remote Provider. With the Local provider you are the only user, although you can still use Environments to group your own clusters.
- The cluster registered as a Kubernetes Connection. See [Registering a Connection](/pr-preview/pr-22675/guides/infrastructure-management/registering-a-connection/).
- Teammates who belong to a team in your organization. Your provider manages teams and invitations.
- Permission to create Environments and Workspaces and to assign resources to them. If a button is disabled, ask your organization admin.

## Create an Environment and add the cluster




<ul class="nav nav-tabs" id="tabs-4" role="tablist">
  <li class="nav-item">
      <button class="nav-link active"
          id="tabs-04-00-tab" data-bs-toggle="tab" data-bs-target="#tabs-04-00" role="tab"
          data-td-tp-persist="meshery ui" aria-controls="tabs-04-00" aria-selected="true">
        Meshery UI
      </button>
    </li><li class="nav-item">
      <button class="nav-link"
          id="tabs-04-01-tab" data-bs-toggle="tab" data-bs-target="#tabs-04-01" role="tab"
          data-td-tp-persist="mesheryctl" aria-controls="tabs-04-01" aria-selected="false">
        mesheryctl
      </button>
    </li>
</ul>

<div class="tab-content" id="tabs-4-content">
    <div class="tab-body tab-pane fade show active"
        id="tabs-04-00" role="tabpanel" aria-labelled-by="tabs-04-00-tab" tabindex="4">
        <ol>
<li>Open <strong>Lifecycle → Environments</strong> and click <strong>Create</strong>.</li>
<li>Enter a name, such as <code>staging</code>, and an optional description, then save.</li>
<li>On the new Environment&rsquo;s card, click <strong>Assigned Connections</strong>. The <strong><code>&lt;environment&gt;</code> Resources</strong> dialog opens.</li>
<li>Move the cluster&rsquo;s Connection from the available list to the assigned list, then click <strong>Save</strong>.</li>
</ol>
<p>You can also assign Environments from <strong>Lifecycle → Connections</strong>: use the <strong>Environments</strong> column on the cluster&rsquo;s row.</p>

    </div>
    <div class="tab-body tab-pane fade"
        id="tabs-04-01" role="tabpanel" aria-labelled-by="tabs-04-01-tab" tabindex="4">
        <p>Find your organization ID, then create the Environment:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl organization list
</span></span><span class="line"><span class="cl">mesheryctl environment create --orgId &lt;org-id&gt; --name staging --description <span class="s2">&#34;Staging clusters&#34;</span>
</span></span><span class="line"><span class="cl">mesheryctl environment list --orgId &lt;org-id&gt;
</span></span></code></pre></div><p><code>mesheryctl</code> does not assign Connections to Environments. Use Meshery UI for step 3 and step 4.</p>

    </div>
</div>


An Environment can hold several clusters, and a cluster can belong to several Environments. Group clusters the way your team deploys to them, for example one Environment per stage.

## Create a Workspace and attach the Environment




<ul class="nav nav-tabs" id="tabs-5" role="tablist">
  <li class="nav-item">
      <button class="nav-link active"
          id="tabs-05-00-tab" data-bs-toggle="tab" data-bs-target="#tabs-05-00" role="tab"
          data-td-tp-persist="meshery ui" aria-controls="tabs-05-00" aria-selected="true">
        Meshery UI
      </button>
    </li><li class="nav-item">
      <button class="nav-link"
          id="tabs-05-01-tab" data-bs-toggle="tab" data-bs-target="#tabs-05-01" role="tab"
          data-td-tp-persist="mesheryctl" aria-controls="tabs-05-01" aria-selected="false">
        mesheryctl
      </button>
    </li>
</ul>

<div class="tab-content" id="tabs-5-content">
    <div class="tab-body tab-pane fade show active"
        id="tabs-05-00" role="tabpanel" aria-labelled-by="tabs-05-00-tab" tabindex="5">
        <ol>
<li>Open <strong>Lifecycle → Workspaces</strong> and click <strong>Create</strong>.</li>
<li>Enter a name, such as <code>payments-team</code>, and an optional description, then save.</li>
<li>On the Workspace&rsquo;s card, open its <strong>Environments</strong> control. The <strong>Assign Environments to <code>&lt;workspace&gt;</code></strong> dialog opens.</li>
<li>Move your Environment to the assigned list and save.</li>
</ol>

    </div>
    <div class="tab-body tab-pane fade"
        id="tabs-05-01" role="tabpanel" aria-labelled-by="tabs-05-01-tab" tabindex="5">
        <div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl workspace create --orgId &lt;org-id&gt; --name payments-team --description <span class="s2">&#34;Payments team&#34;</span>
</span></span><span class="line"><span class="cl">mesheryctl workspace list --orgId &lt;org-id&gt;
</span></span></code></pre></div><p><code>mesheryctl</code> does not attach Environments to Workspaces. Use Meshery UI for step 3 and step 4.</p>

    </div>
</div>


To share designs as well, open the Workspace's **Designs** control and assign them in the **Assign Designs and Views to `<workspace>`** dialog.

## Assign teams to the Workspace

1. In **Lifecycle → Workspaces**, open the Workspace's actions and click **Assign Teams**.
2. In the **Assign Teams to `<workspace>`** dialog, move your teammates' teams to the assigned list and save.

Members of those teams now see the Workspace, its Environments, and the clusters in them, and can choose those clusters when they [deploy a design](/pr-preview/pr-22675/guides/infrastructure-management/deploying-a-design/). What each person can do is set by their role in your organization. See [Remote Provider Permissions](/pr-preview/pr-22675/reference/extensibility/authorization/).

## Make sure teammates can reach the cluster

Access in Meshery is only half of the setup. The Meshery Server your teammate uses also needs a network path to your cluster's API server:

- **Shared Meshery Server.** Your teammates sign in to the same Meshery Server that you connected the cluster to. That server already reaches the cluster, so you only need to make the server itself reachable to them, for example over your VPN.
- **Their own Meshery Server.** Your teammates run their own Meshery Server and sign in with the same Remote Provider. They see the shared Connection, but their server must reach your cluster's API server over the network.

## Stop sharing

Remove access at whichever level fits:

- To remove a team's access to everything in the Workspace, unassign the team from the Workspace.
- To stop sharing one cluster but keep the rest of the Environment shared, unassign the Connection from the Environment.
- To stop sharing every cluster in an Environment, unassign the Environment from the Workspace.

A cluster can belong to several Environments, and an Environment can be attached to several Workspaces. To revoke a team's access to a cluster completely, remove every path between them: check each Environment that holds the cluster and each Workspace those Environments are attached to.

## Related

- [Environments](/pr-preview/pr-22675/concepts/logical/environments/) and [Workspaces](/pr-preview/pr-22675/concepts/logical/workspaces/)
- [Managing Connections](/pr-preview/pr-22675/guides/infrastructure-management/lifecycle-management/)
- [Deploying a Design](/pr-preview/pr-22675/guides/infrastructure-management/deploying-a-design/)

<div class="alert alert-dark" role="alert">
  <h4 class="alert-heading">Discussion Forum</h4>
  <p>Don't find an answer to your question here? Ask on the <a href="https://discuss.meshery.io/">Discussion Forum</a>.</p>
</div>
