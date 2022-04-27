---
layout: page
title: Contributing to Meshery Cypress Tests
permalink: project/contributing/contributing-cypress
description: How to contribute to Meshery Cypress (integration & end-to-end) Tests.
language: en
type: project
category: contributing
---

## <a name="contributing-cypress-intro">Introduction</a>

To automate functional integration and end-to-end testing through Meshery UI, [Cypress](https://www.cypress.io/) is leveraged as it allows for both UI Integration & End-to-end test scripting with javascript through its [modern features](https://docs.cypress.io/guides/overview/why-cypress#Features) and [supported test types](https://docs.cypress.io/guides/overview/why-cypress#Who-uses-Cypress).

## Test Framework Architecture

The following is the foldering structure and description of what each files are intended for:

- [/ui/cypress](https://github.com/meshery/meshery/tree/master/ui/cypress) - All cypress test-related javascript files.
- [/ui/cypress/actionHelpers](https://github.com/meshery/meshery/tree/master/ui/cypress/actionHelpers) - Helpers to provide common UI or API level actions across our different cypress integration and end-to-end tests.
- [/ui/cypress/fixtures](https://github.com/meshery/meshery/tree/master/ui/cypress/fixtures) - Our [Fixture Files](https://docs.cypress.io/guides/core-concepts/writing-and-organizing-tests#Fixture-Files) which are used by our tests as external pieces of static data to [Stub](https://docs.cypress.io/guides/guides/network-requests#Stubbing) response data in integration tests **(i.e. [/ui/cypress/integration/integration/configuration_filters_spec.js](https://github.com/meshery/meshery/blob/master/ui/cypress/integration/integration/configuration_filters_spec.js))** or reuse data as test input in end-to-end tests **(i.e. [/ui/cypress/integration/e2e/service_mesh_configuration_management_spec.js](https://github.com/meshery/meshery/blob/master/ui/cypress/integration/e2e/service_mesh_configuration_management_spec.js))**.
- [/ui/cypress/integration](https://github.com/meshery/meshery/tree/master/ui/cypress/integration) - 

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

### Linting-UI

- When contributing to this project, it is advisable to

  - Use [eslint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) plugin for Visual Studio Code.

  - Disable plugins other than `eslint` for formatting and linting, if any.

### Install UI dependencies

To install/update the UI dependencies:

```
make setup-ui-libs
```

### Build and export UI

To build and export the UI code:

```
make build-ui
```

> Now that the UI code is built, Meshery UI will be available at `http://localhost:9081` when Meshery Server is running (Read below).

> Changes are not recompiled directly, you will have to run to rebuild the UI to see them

## Run Meshery Server

To start running Meshery Server locally:

```
make server
```
> Now, Meshery will run on the default port `http://localhost:9081`.

{% include alert.html type="warning" title="Usage of 'make run-fast' is deprecated!" %}

### UI Development Server

If you want to work on the UI, it will be a good idea to use the included UI development server. You can run the UI development server by running the following command:

```
make run-ui-dev
```

Refer to [Contributing to Meshery Server](contributing-server), as needed.

> Make sure to have Meshery Server configured, up and running on the default port `http://localhost:9081` before proceeding to access and work on the UI server at `http://localhost:3000`.

> Any UI changes made now will _automatically_ be recompiled and served in the browser.

### Running Meshery Cypress Tests



Refer to [Contributing to Cypress UI Tests](contributing-cypress).

### Running Meshery from IDE

All of the above steps would get the Meshery's development server running for you to work on in any IDE of your choice.

{% include suggested-reading.html %}
