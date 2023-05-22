---
layout: page
title: Contributing to Meshery UI
permalink: project/contributing/contributing-ui
description: How to contribute to Meshery UI (web-based user interface).
language: en
type: project
category: contributing
---

## <a name="contributing-ui">UI Contribution Flow</a>

Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

## Architecture

The following is a list of top-level frameworks, libraries, design system used in Meshery UI.

- [NextJS](https://nextjs.org/) - Server and router
- [ReactJS](https://reactjs.org/) - User Interface library
- [Material UI](https://material-ui.com/) - Design System
- [BillboardJS](https://naver.github.io/billboard.js/) - Charting library, used for exposing Grafana and Prometheus-collected metrics
- [CytoscapeJS](https://js.cytoscape.org/) - a visualization tool for canvas-based, visual topology (networks)

## Meshery Server APIs

Go [here](https://docs.meshery.io/extensibility/api#rest) for the docs.

### REST API

- Meshery provides a REST API available through the default port of `9081/tcp`.
- List of [endpoints](https://docs.meshery.io/reference/rest-apis) (spreadsheet) a simple, static list of REST API endpoints with short description of their purpose.
- Swagger / Open API.
- Collection of sets of REST API docs that Meshery server exposes to clients (like the Meshery UI).

### GraphQL API

- Meshery provides a GraphQl API available through the default port of `9081/tcp`.
- [Relay](https://relay.dev) is the client used.

## Design

### Wireframing / Mockups

- Meshery UI in [Figma](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI)

> Fill-in a <a href="https://layer5.io/newcomers">community member form</a> to gain access to community resources.
> You need to ask for the access to the above Figma File in [Slack](http://slack.layer5.io/)

### Design Prologue

Meshery UI is a significant component of the value proposition Meshery offers to individuals and organizations seeking to adopt and operate a service mesh or collection of service meshes.

### Design Goals

The designs in this specification should result in enabling:

- **User experience should be intuitive**

  This is achieved through sensible defaults, consistency of user interaction paradigms and features that delight the user.

- **Meshery UI should simplify user management of workloads and service meshes**

  Meshery UI needs to be simple, but powerful. This is achieved through intuitive layouts and predefined filters that accommodate common tasks.

- **Extensible**

  Meshery UI should be a first-class component of Meshery, but also facilitate third-party integrations.

### Design Objectives

The designs in this specification should result in enabling:

- Meshery UI should be event-driven where possible.

## Setup

### Node Version Recommendations 
We recommend to use LTS versions of node:

- Node 16 LTS: https://nodejs.org/download/release/v16.19.0/
- Node 18 LTS: https://nodejs.org/en/

The script supports other node versions as well, including node15 and node17, but is recommended to use node versions above 15 and especially node LTS v18.

### Linting-UI

- When contributing to this project, it is advisable to

  - Use [eslint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) plugin for Visual Studio Code.

  - Disable plugins other than `eslint` for formatting and linting, if any.

### Install UI dependencies

To install/update the UI dependencies:

{% include code.html code="make ui-setup" %}

### Build and export UI

To build and export the UI code:

{% include code.html code="make ui-build" %}

Using this command, changes that you make are not rebuilt automatically. You will have to run this command again to rebuild the UI and see them.

Now that the UI code is built, Meshery UI will be available at `http://localhost:9081` when Meshery Server is running (Read below).

To build and export the UI code _and_ build and run Meshery Server:

{% include code.html code="make ui-server" %}

## Run Meshery Server

To start running Meshery Server locally:

{% include code.html code="make server" %}
Now, Meshery will run on the default port `http://localhost:9081`.


### UI Development Server

If you want to work on the UI, it will be a good idea to use the included UI development server. You can run the UI development server by running the following command:

{% include code.html code="make ui" %}

Refer to [Contributing to Meshery Server](contributing-server), if needed.

> Make sure to have Meshery Server configured, up and running on the default port `http://localhost:9081` and choose a provider to login with (visit `http://localhost:9081`) before proceeding to access and work on the UI server at `http://localhost:3000`.

Any UI changes made now will _automatically_ be rebuilt and served in your browser.

### Running Cypress integration tests

To run cypress integration tests, a convenience make target called `ui-integration-tests` that installs dependencies in `/ui` and `/provider-ui` folders as prerequisite and invokes `ci-test-integration` npm script found in [/ui/package.json](https://github.com/meshery/meshery/blob/master/ui/package.json)
{% include code.html code="make ui-integration-tests" %}
{% include alert.html type="info" title="Above command must be run from Meshery repository's root folder." %}

Refer to [Meshery Cypress Testing](contributing-cypress) for details of how to contribute and benefit from Meshery Cypress (integration & end-to-end) testing.


### Static Files, Icons and Images

The Meshery UI public folder contains static files. Its folder structure looks like this:
```
meshery
└── ui
    └── public
        └── static
            ├── favicon.png
            ├── fonts
            ├── img
            └── style
```

Images and icons used in Meshery UI need to be sourced from the [public directory of images](https://github.com/meshery/meshery/tree/master/ui/public/static/img). The files written inside this directory should only end with the extensions like `.svg`, `.png`, `.jpg` or `.jpeg`. Always use vector-based graphics (`.svg`), unless your have extenuating circumstances.

##### Conventions for SVG files
1. SVGs should be optimized and compressed. 
    1. Use an online, SVG optimizer, like https://www.svgviewer.dev, to compress the file(s) to smaller size.
3. All SVGs should have `height` and `width` properties set to 20px x 20px by default. Ensure that height and width attributes are always set in original SVG.
4. All SVGs should have `height` and `width` included as a style prop in their React component. 
5. Always include this XML header in each SVG image:  
    <pre class="codeblock-pre"><div class="codeblock">
    <div class="clipboardjs">&lt;?xml version="1.0" encoding="UTF-8"?&gt;&lt;!DOCTYPE svg&gt;</div></div>
    </pre>
4. Svg can only fall under two categories, and this categories should be the name of folder \
    1. white: containing white or mono-colored version of that SVG
    2. color: containing colored version of that SVG.
    e.g.: the Meshery logo icon folder structure looks like this:
    ```
    └── img
        └── meshery
            ├── white
            |   └── meshery-white.svg
            └── color
                └── meshery-color.svg
    ```
5. Avoid any kind of duplicity in the versions of icons used.

For accessing the svg file as data-url, the utf8 encoding should be used in place of base64.Use [encodeURIComponent](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/encodeURIComponent) on SVG data URIs. \

{% include code.html code="let svg = 'data:image/svg+xml;utf8,' + encodeURIComponent(svgFile);" %}

<!-- {% include suggested-reading.html %} -->
