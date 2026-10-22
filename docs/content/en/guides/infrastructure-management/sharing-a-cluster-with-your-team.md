---
title: Sharing a Cluster with Your Team
categories: [infrastructure]
weight: -2
description: Give your teammates access to a Kubernetes cluster you manage in Meshery by grouping it into an Environment, attaching the Environment to a Workspace, and assigning teams to the Workspace.
---

A Kubernetes cluster you connect to Meshery is private to you until you share it. Use this guide to let your teammates see the cluster, deploy designs to it, and work with its resources, without handing out your kubeconfig. You share a cluster in three steps:

1. Put the cluster's [Connection]({{< ref "concepts/logical/connections/index.md" >}}) in an [Environment]({{< ref "concepts/logical/environments.md" >}}).
2. Attach the Environment to a [Workspace]({{< ref "concepts/logical/workspaces.md" >}}).
3. Assign your teammates' teams to the Workspace.

## Before you begin

- A Meshery Server signed in with a Remote Provider, such as Meshery Cloud. Organizations, teams, and other users come from the Remote Provider. With the Local provider you are the only user, although you can still use Environments to group your own clusters.
- The cluster registered as a Kubernetes Connection. See [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}).
- Teammates who belong to a team in your organization. Your provider manages teams and invitations.
- Permission to create Environments and Workspaces and to assign resources to them. If a button is disabled, ask your organization admin.

## Create an Environment and add the cluster

{{< tabpane text=true >}}
{{% tab header="Meshery UI" %}}
1. Open **Lifecycle → Environments** and click **Create**.
2. Enter a name, such as `staging`, and an optional description, then save.
3. On the new Environment's card, click **Assigned Connections**. The **`<environment>` Resources** dialog opens.
4. Move the cluster's Connection from the available list to the assigned list, then click **Save**.

You can also assign Environments from **Lifecycle → Connections**: use the **Environments** column on the cluster's row.
{{% /tab %}}
{{% tab header="mesheryctl" %}}
Find your organization ID, then create the Environment:

```bash
mesheryctl organization list
mesheryctl environment create --orgId <org-id> --name staging --description "Staging clusters"
mesheryctl environment list --orgId <org-id>
```

`mesheryctl` does not assign Connections to Environments. Use Meshery UI for step 3 and step 4.
{{% /tab %}}
{{< /tabpane >}}

An Environment can hold several clusters, and a cluster can belong to several Environments. Group clusters the way your team deploys to them, for example one Environment per stage.

## Create a Workspace and attach the Environment

{{< tabpane text=true >}}
{{% tab header="Meshery UI" %}}
1. Open **Lifecycle → Workspaces** and click **Create**.
2. Enter a name, such as `payments-team`, and an optional description, then save.
3. On the Workspace's card, open its **Environments** control. The **Assign Environments to `<workspace>`** dialog opens.
4. Move your Environment to the assigned list and save.
{{% /tab %}}
{{% tab header="mesheryctl" %}}
```bash
mesheryctl workspace create --orgId <org-id> --name payments-team --description "Payments team"
mesheryctl workspace list --orgId <org-id>
```

`mesheryctl` does not attach Environments to Workspaces. Use Meshery UI for step 3 and step 4.
{{% /tab %}}
{{< /tabpane >}}

To share designs as well, open the Workspace's **Designs** control and assign them in the **Assign Designs and Views to `<workspace>`** dialog.

## Assign teams to the Workspace

1. In **Lifecycle → Workspaces**, open the Workspace's actions and click **Assign Teams**.
2. In the **Assign Teams to `<workspace>`** dialog, move your teammates' teams to the assigned list and save.

Members of those teams now see the Workspace, its Environments, and the clusters in them, and can choose those clusters when they [deploy a design]({{< ref "guides/infrastructure-management/deploying-a-design.md" >}}). What each person can do is set by their role in your organization. See [Remote Provider Permissions]({{< ref "reference/extensibility/authorization/index.md" >}}).

## Make sure teammates can reach the cluster

Access in Meshery is only half of the setup. The Meshery Server your teammate uses also needs a network path to your cluster's API server:

- **Shared Meshery Server.** Your teammates sign in to the same Meshery Server that you connected the cluster to. That server already reaches the cluster, so you only need to make the server itself reachable to them, for example over your VPN.
- **Their own Meshery Server.** Your teammates run their own Meshery Server and sign in with the same Remote Provider. They see the shared Connection, but their server must reach your cluster's API server over the network.

## Stop sharing

Remove access at whichever level fits:

- To remove a team's access to everything in the Workspace, unassign the team from the Workspace.
- To stop sharing one cluster but keep the rest of the Environment shared, unassign the Connection from the Environment.
- To stop sharing every cluster in an Environment, unassign the Environment from the Workspace.

## Related

- [Environments]({{< ref "concepts/logical/environments.md" >}}) and [Workspaces]({{< ref "concepts/logical/workspaces.md" >}})
- [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}})
- [Deploying a Design]({{< ref "guides/infrastructure-management/deploying-a-design.md" >}})

{{< discuss >}}
