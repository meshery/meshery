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

WIP...