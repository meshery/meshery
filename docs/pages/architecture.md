---
layout: page
title: Meshery Architecture
permalink: architecture
---

## Architecture

#### The Meshery architecture can be observed in two perspectives:
 
##### 1. [**Clients**](#clients)
##### 2. [**Providers**](#providers)


![Meshery architecture](/docs/assets/img/architecture/Meshery-architecture-diagram.png)


<a name="clients">

### 1. **Clients**

![Client architecture](/docs/assets/img/architecture/Meshery-client-architecture.svg)


<a name="providers">

### 2. **Providers**

![Provider architecture](/docs/assets/img/architecture/Meshery-provider-architecture.svg)

#### **Network Ports**

Meshery uses the following list of network ports to interface with its various components:

| Network Application                            | Port             |
| :--------------------------------------------- | :--------------: |
| Meshery REST API                               | 9081/tcp         |
| Learn Layer5 Application                       | 10011            |

#### **Adapter Ports**

| Service Mesh  | Port          |
| :------------ | ------------: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.port }} |
{% endif -%}
{% endfor %}

See the [**Adapters**](/docs/architecture/adapters) section for more information on the function of an adapter.

#### **Statefulness in Meshery components**

Some components within Meshery's architecture are concerned with persisting data while others are only
concerned with a long-lived configuration, while others have no state at all.

| Components        | Persistence  | Description                                                           |
| :---------------- | :----------- | :-------------------------------------------------------------------- |
| mesheryctl        | stateless    | command line interface that has a configuration file                  |
| Meshery Adapters  | stateless    | interface with service meshes on a transactional basis                |
| Meshery Server    | caches state | application cache is stored in user's `$HOME/.meshery/` folder        |
| Meshery Providers | stateful     | location of persistent user preferences, environment, tests and so on |
