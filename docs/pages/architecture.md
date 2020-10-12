---
layout: page
title: Meshery Architecture
permalink: architecture
---

## Architecture

<h4>The Meshery architecture can be observed in two perspectives:</h4> 
 
##### 1. <a href="#clients"><b>Clients</b></a>
##### 2. <a href="#providers"><b>Providers</b></a> 

<br /><a href="{{site.baseurl}}/assets/img/architecture/meshery-architecture-diagram.svg"><img src="{{site.baseurl}}/assets/img/architecture/Meshery-architecture-diagram.png" /></a>


<a name="clients">

### 1. <b>Clients</b>

<a href="{{site.baseurl}}/assets/img/architecture/meshery-client-architecture.svg"><img src="{{site.baseurl}}/assets/img/architecture/Meshery-client-architecture.svg" /></a>

<a name="providers">

### 2. <b>Providers</b>

<a href="{{site.baseurl}}/assets/img/architecture/meshery-provider-architecture.svg"><img src="{{site.baseurl}}/assets/img/architecture/Meshery-provider-architecture.svg" /></a>

#### <b>Network Ports</b>

Meshery uses the following list of network ports to interface with its various components:

| Network Application                            | Port             |
| :--------------------------------------------- | :--------------: |
| Meshery REST API                               | 9081/tcp         |
| Learn Layer5 Application                       | 10011            |

#### <b>Adapter Ports</b>

| Service Mesh  | Port          |
| :------------ | ------------: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.port }} |
{% endif -%}
{% endfor %}

See the [**Adapters**](/docs/architecture/adapters) section for more information on the function of an adapter.

#### <b>Statefulness in Meshery components</b>

Some components within Meshery's architecture are concerned with persisting data while others are only
concerned with a long-lived configuration, while others have no state at all.

| Components        | Persistence  | Description                                                           |
| :---------------- | :----------- | :-------------------------------------------------------------------- |
| mesheryctl        | stateless    | command line interface that has a configuration file                  |
| Meshery Adapters  | stateless    | interface with service meshes on a transactional basis                |
| Meshery Server    | caches state | application cache is stored in user's `$HOME/.meshery/` folder        |
| Meshery Providers | stateful     | location of persistent user preferences, environment, tests and so on |
