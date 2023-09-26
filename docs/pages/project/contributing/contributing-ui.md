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

Meshery UI is built on React and Next.js. In nearly all cases, while contributing to Meshery UI, you will need to build and run Meshery Server as well. Meshery Server is written in `Go` (Golang) and leverages Go Modules. To make building of each component easier a `Makefile` is included in the main repository folder. Use `make` to build any and all components.

## Architecture

The following is a list of top-level frameworks, libraries, design system used in Meshery UI.

- [NextJS](https://nextjs.org/) - Server and router
- [ReactJS](https://reactjs.org/) - User Interface library
- [Material UI](https://material-ui.com/) - Design System
- [BillboardJS](https://naver.github.io/billboard.js/) - Charting library, used for exposing Grafana and Prometheus-collected metrics
- [CytoscapeJS](https://js.cytoscape.org/) - a visualization tool for canvas-based, visual topology (networks)

## Meshery Server APIs

The [API overview](https://docs.meshery.io/extensibility/api) in the Extensibility guide offers high-level insight, while each API reference (below) offers details pertainting to each API endpoints' behavior and use.

### REST API
Meshery provides a REST API available through the default port of `9081/tcp` at `<hostname>:<port>/api/`. See the [REST API Reference](https://docs.meshery.io/reference/rest-apis) for a complete list of endpoints available with short description of their purpose and example code.

### GraphQL API

Meshery provides a GraphQl API available through the default port of `9081/tcp` at `<hostname>:<port>/api/graphql/query`. See the [GraphQL API Reference](https://docs.meshery.io/reference/graphql-apis) [Relay](https://relay.dev) is the client used.

## Design

Meshery UI is a significant component of the value proposition Meshery offers to individuals and organizations seeking to adopt and operate a service mesh or collection of service meshes.

### User Flows / Wireframing / Mockups

Figma is used as the user flow and UI mockup design tool. The [Meshery UI design file](https://www.figma.com/file/SMP3zxOjZztdOLtgN4dS2W/Meshery-UI) can be viewed by using this [open invitation](https://www.figma.com/team_invite/redeem/qJy1c95qirjgWQODApilR9) for view access. Fill-in a [community member form](https://layer5.io/newcomers) to gain access to additional community resources.

### Design Goals

The designs in this specification should result in enabling:

- **User experience should be intuitive** - Meshery's UX, which includes Meshery UI as a component of the overall UX, needs to be simple, but powerful. This is achieved through intuitive layouts and predefined filters that accommodate common tasks. This is achieved through sensible defaults, consistency of user interaction paradigms and features that delight the user.

- **Extensible** - Meshery UI should be a first-class component of Meshery, but also facilitate third-party integrations.

- **Real-time** - Meshery UI should be event-driven where possible.

## Setup

### Node Version Recommendations 
We recommend using Node 18 LTS. Meshery UI's build script supports other Node versions as well, including node16 and node17, but is recommended to use Node 18 LTS.

### Linting-UI

When contributing to Meshery UI, it is advisable to:

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

**Please note**: If you see "Meshery Development Incompatible" while trying to sign into Meshery Server, then follow these steps:

<a href="{{ site.baseurl }}/assets/img/meshery-development-incompatible-error.png">
  <img style= "width: 600px;" src="{{ site.baseurl }}/assets/img/meshery-development-incompatible-error.png" />
</a>

Potential Solution: 

-  Go to your meshery folder in your local-system where you’ve cloned it.
Execute: 

- `git remote add upstream https://github.com/meshery/meshery`
- `git fetch upstream`
- Restart the meshery server
- Additionally, before restarting the server, if you like to pull the latest changes, you can do: `git pull upstream master`


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

{% include suggested-reading.html %}
