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

The following is the foldering structure and description of what each files are intended for under our [/ui/cypress/](https://github.com/meshery/meshery/tree/master/ui/cypress) folder:

### [./actionHelpers/](https://github.com/meshery/meshery/tree/master/ui/cypress/actionHelpers)

Helpers to provide common UI or API level actions across our different cypress integration and end-to-end tests.

### [./fixtures/](https://github.com/meshery/meshery/tree/master/ui/cypress/fixtures)

Our [Fixture Files](https://docs.cypress.io/guides/core-concepts/writing-and-organizing-tests#Fixture-Files) which are used by our tests as: 
- external pieces of static data to [Stub](https://docs.cypress.io/guides/guides/network-requests#Stubbing) response data in integration tests **(i.e. [/integration/integration/configuration_filters_spec.js](https://github.com/meshery/meshery/blob/master/ui/cypress/integration/integration/configuration_filters_spec.js))**
- or reuse data as test input in end-to-end tests **(i.e. [/integration/e2e/service_mesh_configuration_management_spec.js](https://github.com/meshery/meshery/blob/master/ui/cypress/integration/e2e/service_mesh_configuration_management_spec.js))**.

### [./integration/integration/](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/integration)

Integration tests for Meshery UI that stub server requests to:
1. Prevent state mutation across tests.
2. Build the way we want the data to be structured without contract of server being available.
3. Test critical edge cases without server, etc. 

Refer [here](https://docs.cypress.io/guides/getting-started/testing-your-app#Stubbing-the-server) for details about when it's a good idea to stub the server versus allowing the frontend to reach out the actual server and its underlying resources.

### [./integration/e2e/](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/e2e)

End-to-end tests for both Meshery UI and Meshery Server where its usually necessary to [seed data](https://docs.cypress.io/guides/getting-started/testing-your-app#Seeding-data), ocassionally [bypass our UI](https://docs.cypress.io/guides/getting-started/testing-your-app#Bypassing-your-UI), use [actual server responses](https://docs.cypress.io/guides/guides/network-requests#Use-Server-Responses) and define cypress [routes](https://docs.cypress.io/guides/guides/network-requests#Routing) to [wait](https://docs.cypress.io/guides/guides/network-requests#Waiting) and [assert](https://docs.cypress.io/guides/guides/network-requests#Assertions) on requests and/or their responses.

### [./plugins/](https://github.com/meshery/meshery/tree/master/ui/cypress/plugins)

Define cypress [plugins](https://docs.cypress.io/guides/tooling/plugins-guide) to leverage as "Seams" for us to run our own custom code to execute during particular stages of Cypress lifecycle.

### [./support/](https://github.com/meshery/meshery/tree/master/ui/cypress/support)

This is where our [supportFile](https://docs.cypress.io/guides/references/configuration#Folders-Files) resides ([./support/index.js](https://github.com/meshery/meshery/blob/master/ui/cypress/support/index.js)). Its processed and loaded automatically before tests run and it imports our [./support/commands.js](https://github.com/meshery/meshery/blob/master/ui/cypress/support/commands.js) file which allows us to sparingly define our Cypress [Custom Commands](https://docs.cypress.io/api/cypress-api/custom-commands) to reuse functions needed across most or all test suites.

## How to Open Cypress

Steps to start cypress depends on if Meshery installation that's being targeted for testing is a user build or a dev build (from source code) but the following steps try to simplify the latter which should be the most frequently needed scenario:

### Run Meshery UI Dev Server & Open Cypress
If its the first time you're opening cypress:
`cd ui && npm i && npm run cy:dev:open`

Else, just run:
`npm run cy:dev:open`

Note: keep in mind that if running integration tests (tests in [./integration/integration/](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/integration) folder) the server doesn't need to be running but for full blown end-to-end tests (tests in [./integration/e2e/](https://github.com/meshery/meshery/tree/master/ui/cypress/integration/e2e) folder) its recommended to run `make server` OR make sure a Meshery user build is installed and running locally before choosing one of those tests.

## References
- [Writing Your First Test](https://docs.cypress.io/guides/getting-started/writing-your-first-test)
- [Testing Your App](https://docs.cypress.io/guides/getting-started/testing-your-app)
- [Introduction to Cypress](https://docs.cypress.io/guides/core-concepts/introduction-to-cypress)
- [Writing and Organizing Tests](https://docs.cypress.io/guides/core-concepts/writing-and-organizing-tests)
- [Best Practices](https://docs.cypress.io/guides/references/best-practices)
- [Learn Cypress](https://learn.cypress.io)
