---
title: Managing Connections
aliases: 
- /tasks/lifecycle-management
- /guides/infrastructure-management/managing-connections
categories: [infrastructure]
description: Find, inspect, transition, reconfigure, and remove the Connections Meshery manages, from Meshery UI or mesheryctl.
---

<a name="lifecycle-management"></a>

Every cluster, metrics server, and source that Meshery manages is a [Connection]({{< ref "concepts/logical/connections/index.md" >}}). After you [register a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}), use this guide to keep it in the state you want: find it, check its status, change its state, reconfigure it, and remove it when you are done.

[![States of a Kubernetes cluster connection](./images/states-for-kubernetes-cluster-connections.svg)](./images/states-for-kubernetes-cluster-connections.svg)

## Find your Connections

{{< tabpane text=true >}}
{{% tab header="Meshery UI" %}}
Open **Lifecycle → Connections**. The table lists each Connection with its **Name**, **Environments**, **Kind**, **Category**, **Status**, and timestamps. Search by name, or filter the table by kind and status. Click a Kubernetes Connection's row to expand it and see its server, its Operator, MeshSync, and Broker status, and its [Diagnostics]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md#diagnostics" >}}).
{{% /tab %}}
{{% tab header="mesheryctl" %}}
```bash
mesheryctl connection list
mesheryctl connection list --kind kubernetes --status connected
mesheryctl connection list --count
mesheryctl connection view <connection-name|connection-id> --output-format json
```

`list` shows 10 Connections per page. Use `--page` and `--pagesize` to see more. Repeat `--kind` or `--status` to match several values.
{{% /tab %}}
{{< /tabpane >}}

## Change a Connection's state

The **Status** column holds a selector for each Connection. Choose the state you want, then confirm the transition. The selector only offers transitions that the Connection's kind allows, and you need permission to change Connection state. Common transitions:

| Goal | Choose | Result |
|---|---|---|
| Start managing a discovered or registered Connection | **Connected** | Meshery connects and, for a Kubernetes cluster, starts MeshSync. |
| Stop managing it but keep it | **Disconnected** | Meshery disconnects. Reconnect it later from the same selector. |
| Stop Meshery from rediscovering it | **Ignored** | Meshery keeps the Connection on record and leaves it alone. |
| Remove it | **Deleted** | Meshery deletes the Connection. You cannot transition it again. |

The **Status** filter also lists **Maintenance**, but no current Connection kind offers a transition into it.

For what each state means, see [States and the Lifecycle of Connections]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}}).

## Reconfigure a Connection

Open the Connection's **Actions** menu:

- **Configure** opens the [Connection Wizard]({{< ref "guides/infrastructure-management/registering-a-connection.md#updating-a-connection" >}}) in configure mode. For a Kubernetes Connection, this is where you rename it, change its status, switch the MeshSync deployment mode, and flush MeshSync data.
- **Configure Controllers** (Kubernetes only) tunes the Operator, MeshSync, and Broker for that cluster. See [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}).
- **Copy link** copies a direct link to the Connection.

## Group Connections into Environments

Use the **Environments** column to assign a Connection to one or more [Environments]({{< ref "concepts/logical/environments.md" >}}). Environments are how you target a deployment at a set of clusters and how you share Connections with your team through [Workspaces]({{< ref "concepts/logical/workspaces.md" >}}).

## Remove a Connection

{{< tabpane text=true >}}
{{% tab header="Meshery UI" %}}
Set the Connection's **Status** to **Deleted** and confirm.
{{% /tab %}}
{{% tab header="mesheryctl" %}}
```bash
mesheryctl connection delete <connection-id>
```
{{% /tab %}}
{{< /tabpane >}}

Deleting a Kubernetes Connection removes it from Meshery. It does not delete anything that Meshery deployed into the cluster. To remove the Operator, MeshSync, and Broker from the cluster, uninstall them from the Connection first, as described in [Operating a connection]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md#operating-a-connection" >}}).

## Lifecycle FAQs

<details>
<summary>
<strong>Question:</strong> Meshery is a collaboration platform. When I bring my infrastructure under Meshery's management, will my Kubernetes clusters be available to all other users? Can other people access my cluster?
</summary><strong>Answer:</strong> <p>Only if you explicitly allow it. These controls apply:</p>
<p><b>1. Ownership:</b> The person who provides the Kubernetes context creates and owns the Connection to that cluster. That person may share the Connection with others on their teams.</p>
<p><b>2. Permission:</b> To grant team members access, create an Environment, assign the Kubernetes Connection to it, then create a Workspace and assign the Environment to it. Members of the teams you share the Workspace with can then access the cluster.</p>
<p><b>3. Connectivity:</b> Those people also need a network path to the cluster. There are two common setups:</p>
<p>3.a) The other user runs their own Meshery Server (or signs into a shared instance such as the Playground) and sees the Connection when they sign in. Their Meshery Server must be able to reach your Kubernetes API server over the network.</p>
<p>3.b) The other user signs into your Meshery Server, which already has network access to your cluster. You then need to make your Meshery Server reachable to them, for example over a VPN.</p>
<p>For more information, see <a href="{{< ref "reference/extensibility/authorization/index.md" >}}">Remote Provider Permissions</a>.</p>
</details>

{{< discuss >}}
