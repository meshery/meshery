---
layout: page
title: Contributing to Meshery UI
permalink: project/contributing-ui
description: How to contribute to Meshery UI (web-based user interface).
language: en
type: project
category: contributing
---

### <a name="contributing-ui">UI Contribution Flow</a>

Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

### Architecture

The following is a list of top-level frameworks, libraries, design system used in Meshery UI.

- [NextJS](https://nextjs.org/) - Server and router
- [ReactJS](https://reactjs.org/) - User Interface library
- [Material UI](https://material-ui.com/) - Design System
- [BillboardJS](https://naver.github.io/billboard.js/) - Charting library, used for exposing Grafana and Prometheus-collected metrics

- [MeshMap only] [CytoscapeJS](https://js.cytoscape.org/) - a visualization tool for canvas-based, visual topology (networks)

#### Meshery Server APIs

Go [here](https://docs.meshery.io/extensibility/api#rest) for the docs.

**REST API**

- Meshery provides a REST API available through the default port of `9081/tcp`.
- List of [endpoints](https://docs.meshery.io/reference/rest-apis) (spreadsheet) a simple, static list of REST API endpoints with short description of their purpose.
- Swagger / Open API.
- Collection of sets of REST API docs that Meshery server exposes to clients (like the Meshery UI).

**GraphQL API**

- Meshery provides a GraphQl API available through the default port of `9081/tcp`.
- [Relay](https://relay.dev) is the client used.

### Design

#### Wireframing / Mockups

- Meshery UI in [Figma](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI)

  > Access the [Community Drive](https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA) by completing the community [Member Form](https://layer5.io/newcomer)
  > You need to ask for the access to the above Figma File in [Slack](http://slack.layer5.io/)

#### Design Prologue

Meshery UI is a significant component of the value proposition Meshery offers to individuals and organizations seeking to adopt and operate a service mesh or collection of service meshes.

#### Design Goals

The designs in this specification should result in enabling:

- **User experience should be intuitive**

  This is achieved through sensible defaults, consistency of user interaction paradigms and features that delight the user.

- **Meshery UI should simplify user management of workloads and service meshes**

  Meshery UI needs to be simple, but powerful. This is achieved through intuitive layouts and predefined filters that accommodate common tasks.

- **Extensible**

  Meshery UI should be a first-class component of Meshery, but also facilitate third-party integrations.

#### Design Objectives

The designs in this specification should result in enabling:

- Meshery UI should be event-driven where possible.

### Setting up

#### Install UI dependencies

To install/update the UI dependencies:

```
make setup-ui-libs
```

#### Build and export UI

To build and export the UI code:

```
make build-ui
```

> Now that the UI code is built, Meshery UI will be available at `http://localhost:9081`.

> Changes are not recompiled directly, you will have to run to rebuild the UI to see them


#### Run Meshery
To start running Meshery locally:
```
make run-fast
```

> Now, Meshery will run on the default port `http:localhost:9081`.


#### UI Development Server

If you want to work on the UI, it will be a good idea to use the included UI development server. You can run the UI development server by running the following command:

```
make run-ui-dev
```

> Make sure to have Meshery Server configured, up and running on the default port `http://localhost:9081` before proceeding to access and work on the UI server at `http://localhost:3000`.

> Any UI changes made now will _automatically_ be recompiled and served in the browser.

#### Running Meshery from IDE

If you want to run Meshery from IDE like Goland, VSCode.

- Source these environment variables
  ```
  PROVIDER_BASE_URLS="https://meshery.layer5.io"
  PORT=9081
  DEBUG=true
  ADAPTER_URLS=mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001 mesherylocal.layer5.io:10002 mesherylocal.layer5.io:10003 mesherylocal.layer5.io:10004 mesherylocal.layer5.io:10005 mesherylocal.layer5.io:10006 mesherylocal.layer5.io:10007 mesherylocal.layer5.io:10008 mesherylocal.layer5.io:10009
  ```
- `go tool` argument
  ```shell
  -tags draft
  ```
- Add the below host to `/etc/hosts`
  ```shell
  127.0.0.1 mesherylocal.layer5.io
  ```

## Suggested Reading

{% assign sorted_reading = site.pages | sort: page.title | reverse %}

<ul>
  {% for item in sorted_reading %}
  {% if item.type=="project" and item.category=="contributing" and item.list!="exclude" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>
