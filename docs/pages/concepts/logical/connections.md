---
layout: default
title: Connections
permalink: concepts/logical/connections
type: concepts
abstract: "Meshery Connections are managed and unmanaged resources that either through discovery or manual entry are managed by a state machine and used within one or more Environments."
language: en
list: include
redirect_from:
- concepts/connections
---
Meshery Connections are managed and unmanaged resources that either through discovery or manual entry are tracked by Meshery. Connections can be assigned as resources to an Environment. 

{% include alert.html type="info" title="Connections as resources in Environments" content="Meshery Environments allow you to logically group related <a href='/concepts/logical/connections'>Connections</a> and their associated <a href='/concepts/logical/credentials'>Credentials</a>. Environments make it easier for you to manage, share, and work with a collection of resources as a group, instead of dealing with all your Connections and Credentials on an individual basis." %}

{% include alert.html type="dark" title="Managed vs Unmanaged Connections" content="Managed Connections are those that are discovered by MeshSync and are managed by Meshery. Unmanaged Connections are those that are manually added by the user and are not managed by Meshery." %}

## States and the Lifecycle of Connections

Meshery tracks the status of each connections throughout the connection's lifecycle. Meshery is intentional about the currently assigned state and which state a connection may or may not transition to and from. To better understand connection states and their meaning, let's consider an example in which you a `Kubernetes` cluster with `Prometheus` installed.

![]({{site.baseurl}}/assets/img/lifecycle-management/states-for-kubernetes-cluster-connections.svg)

### State: Discovered

All resources discovered by [MeshSync's](meshsync.md) multi-tier discovery or provided as part of config, and if Meshery can integrate, a connection with state as `Discovered` will be created. Though, the connection/resources are not tested for its reachability/usability i.e. Meshery has not made an attempt to connect or manage the connection.

When a connection has been discovered, it will be listed in the MeshSync browser / Connections table in Meshery UI. You can self transition a particular connection to [Register](#state-registered) / [Ignore](#state-ignored) state.

> Example: MeshSync discovers Prometheus components and inform Meshery Server about available Prometheus connection, but Meshery is yet to [connect](#state-connected) and start scraping metrics.

### State: Registered

The connection in this state have been verified for its use and reachability but not yet being used. Almost all reachable connections will auto transition to Registered state from [Discovered](#state-discovered) state and it is upto the user what to do with this connection (i.e. User needs to administratively process the connection). It can be transitioned to [Connected](#state-connected), [Maintenance](#state-maintenance) and [Not Found](#state-not-found).

> Example: User manually selects the registered Prometheus connection and transition to the [connected](#state-connected) state (i.e. User administratively processes the connection).

### State: Connected

The connection in this state is administratively processed and being actively managed by Meshery. User can interface and invoke set of actions with the connection.<br>
From this state the transition can happen to either [Maintenance](#state-maintenance) or [Ignore](#state-ignored) state.<br>
Auto transition to [Disconnected](#state-disconnected) state will occur if Meshery can no longer communicate with the connection, which can occur due to connectivity issue/AuthN-AuthZ/connection was deleted outside Meshery or any other issue.

> Example: Meshery is communicating with Prometheus APIs to scrape metrics and present it in the UI.

_Certain connections can auto-transition to connected state._

### State: Ignored

The connection is administratively processed to be ignored from Meshery's view of management. Meshery will not re-discover this connection even when current user session gets expired.

> Example: Meshery server will stop/not scrape metrics from Prometheus. Though, the previous data (if connected previously) will continue to exist and user needs to manually delete.

{% include alert.html type="info" title="Ignored versus Disconnected" content="You might intentionally choose to have Meshery ignore a given Prometheus connection, explicitly leaving in Meshery’s field of view, but identifying it as a connection not to manage. This is distinctly different than a Prometheus that Meshery was managing, but has been turned off/uninstalled and now Meshery is disconnected from the Prometheus." %}

### State: Maintenance

The connection is administratively processed to be offline for maintenance tasks. This is different from being [Disconnected](#state-disconnected)/[Ignored](#state-ignored).

### State: Disconnected

The connection was previously [discovered](#state-discovered)/[registered](#state-registered)/[connected](#state-connected) but is not available currently. This could happen due to connectivity issue/AuthN-AuthZ/connection was deleted outside meshery/administratively disconnected.

> Example: Prometheus crashed/API token provided at time of registration is revoked.

{% include alert.html type="info" title="Disconnected vs Deleted" content="The connection was previously connected but is unreachable due to connectivity issue/AuthN-AuthZ/connection was **deleted outside Meshery** i.e. Connection was deleted beyond the Meshery's view of management." %}

### State: Deleted

The connection is administratively processed to be deleted and removed from Meshery's view of management. All the available/collected data will also be deleted.

> Example: Prometheus metrics will no longer be accessible to you from the Meshery UI.

### State: Not Found

User tried registering the connection **manually** but Meshery could not connect to it or if the connection is unavailable now. User can delete the connection or try re-registering.

{% include alert.html type="info" title="Not Found vs Disconnected" content="You might attempt to transition to Connected state but the connection is unavaialble now due to being deleted/some other reason. This is distinctly different than a cluster with Prometheuses installed for `application monitoring` which was connected previously but is now unreachable from Meshery's view of management due to change in API token/similar issue." %}

_Connections like **Registration of Meshery server with remote provider** (and few other connection types) can self transtion to the valid states._

## Registering Connections with Remote Providers

To register a connection with a remote provider, you need to follow these steps:

1. Obtain the necessary credentials or access tokens from the remote provider.
2. Open the Meshery UI and navigate to the Connections page.
3. Click on the "Add Connection" button.
4. Fill in the required information, such as the provider type, name, and credentials.
5. Click on the "Register" button to register the connection.

Once the connection is registered, Meshery will verify its reachability and usability. If successful, the connection will transition to the "Registered" state. From there, you can choose to administratively process the connection and transition it to the "Connected" state.

Note that some connections, such as the registration of Meshery server with remote providers, can self-transition to valid states.

For more information on the different states and the lifecycle of connections, refer to the documentation above.

![]({{site.baseurl}}/assets/img/architecture/meshery-server-registration-with-remote-providers.svg)
