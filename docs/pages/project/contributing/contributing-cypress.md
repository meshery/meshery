---
layout: page
title: Contributing to Meshery's End-to-End Tests using Cypress
permalink: project/contributing/contributing-cypress
description: How to contribute to End-to-End Tests using Cypress.
language: en
type: project
category: contributing
list: include
---

## Introduction

To automate functional integration and end-to-end testing through Meshery UI, [Cypress](https://www.cypress.io/) is leveraged as it allows for both UI Integration & End-to-end test scripting with javascript through its [modern features](https://docs.cypress.io/guides/overview/why-cypress#Features) and [supported test types](https://docs.cypress.io/guides/overview/why-cypress#Who-uses-Cypress).

## Understanding the test framework directories

Clone the `meshery/meshery` repo and navigate to the [/ui/cypress/](https://github.com/meshery/meshery/tree/master/ui/cypress) directory.

```
.
├── actionHelpers
│   └── service-mesh-configuration-management.js
├── fixtures
│   ├── clusterVersion.json
│   ├── configuration
│   ├── example.json
│   ├── getMeshAdapters.json
│   ├── postMeshManage.json
│   ├── stats.json
│   └── sync.json
├── integration
│   ├── e2e
│   │   ├── lifecyclecheck_spec.js
│   │   ├── service_mesh_configuration_management_spec.js
│   │   ├── settings_spec.js
│   │   └── userpreference_spec.js
│   ├── integration
│   │   ├── configuration_filters_spec.js
│   │   ├── discoverCluster_spec.js
│   │   ├── indexui_spec.js
│   │   ├── settings_spec.js
│   │   └── userpreference_spec.js
│   └── sample_spec.js
├── plugins
│   └── index.js
├── support
│   ├── commands.js
│   └── index.js

```

Let's walk-through each of these directories and comprehend their purpose.

### Directory: `./actionHelpers/` ([code](https://github.com/meshery/meshery/tree/master/ui/cypress/actionHelpers))

Helpers to provide common UI or API level actions across our different cypress integration and end-to-end tests.

### Directory: `./fixtures/` ([code](https://github.com/meshery/meshery/tree/master/ui/cypress/fixtures))

Our [Fixture Files](https://docs.cypress.io/guides/core-concepts/writing-and-organizing-tests#Fixture-Files) which are used by our tests as:

- external pieces of static data to [Stub](https://docs.cypress.io/guides/guides/network-requests#Stubbing) response data in integration tests **(i.e. [/integration/integration/configuration_filters_spec.js](https://github.com/meshery/meshery/blob/master/ui/cypress/integration/integration/configuration_filters_spec.js))**
- or reuse data as test input in end-to-end tests **(i.e. [/integration/e2e/service_mesh_configuration_management_spec.js](https://github.com/meshery/meshery/blob/master/ui/cypress/integration/e2e/service_mesh_configuration_management_spec.js))**.

### Directory: `./integration/integration/` ([code](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/integration))

Integration tests for Meshery UI that stub server requests to:

1. Prevent state mutation across tests.
1. Build the way we want the data to be structured without contract of server being available.
1. Test critical edge cases without server, etc.

Follow [this guidance](https://docs.cypress.io/guides/getting-started/testing-your-app#Stubbing-the-server) regarding when it's a good idea to stub the server versus allowing the frontend to reach out the actual server and its underlying resources.

### Directory: `./integration/e2e/` ([code](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/e2e))

End-to-end tests for both Meshery UI and Meshery Server where its usually necessary to [seed data](https://docs.cypress.io/guides/getting-started/testing-your-app#Seeding-data), occasionally [bypass our UI](https://docs.cypress.io/guides/getting-started/testing-your-app#Bypassing-your-UI), use [actual server responses](https://docs.cypress.io/guides/guides/network-requests#Use-Server-Responses) and define cypress [routes](https://docs.cypress.io/guides/guides/network-requests#Routing) to [wait](https://docs.cypress.io/guides/guides/network-requests#Waiting) and [assert](https://docs.cypress.io/guides/guides/network-requests#Assertions) on requests and/or their responses.

### Directory: `./plugins/` ([code](https://github.com/meshery/meshery/tree/master/ui/cypress/plugins))

Define Cypress plugins to leverage as "Seams" for Meshery's workflows to run the project's own custom code to execute during particular stages of Cypress lifecycle.

### Directory: `./support/` ([code](https://github.com/meshery/meshery/tree/master/ui/cypress/support))

This is where Meshery's Cypress supportFile resides ([./support/index.js](https://github.com/meshery/meshery/blob/master/ui/cypress/support/index.js)). It's processed and loaded automatically before tests run and it imports our [./support/commands.js](https://github.com/meshery/meshery/blob/master/ui/cypress/support/commands.js) file which allows us to sparingly define our Cypress [Custom Commands](https://docs.cypress.io/api/cypress-api/custom-commands) to reuse functions needed across most or all test suites.

## How to manually run end-to-end tests

Steps to start Cypress depend on whether your Meshery installation is built from source code or from a deployed release. The following steps try to simplify the former which should be the most frequently needed scenario:

### Run Meshery UI dev server and Cypress

If its the first time you're opening cypress:
{% include code.html code="cd ui && npm i && npm run cy:dev:open" %}

Else, just run:
{% include code.html code="npm run cy:dev:open" %}

{%include alert.html content="Note the difference between local dependencies for Integration vs End-to-End Tests" %}
keep in mind that if running integration tests (tests in [./integration/integration/](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/integration) folder) the server doesn't need to be running but for full blown end-to-end tests (tests in [./integration/e2e/](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/e2e) folder) its recommended to run `make server` OR make sure a Meshery user build (see [Getting Started](/installation/quick-start.md)) is installed and running locally before choosing one of those tests.

## References

- [Testing Your App](https://docs.cypress.io/guides/getting-started/testing-your-app)
- [Writing and Organizing Tests](https://docs.cypress.io/guides/core-concepts/writing-and-organizing-tests)
- [Best Practices](https://docs.cypress.io/guides/references/best-practices)
- [Learn Cypress](https://learn.cypress.io)
