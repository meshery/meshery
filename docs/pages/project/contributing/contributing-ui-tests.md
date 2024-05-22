---
layout: page
title: Contributing to Meshery's End-to-End Tests
permalink: project/contributing/contributing-ui-tests
abstract: How to contribute to End-to-End Tests using Playwright.
language: en
type: project
category: contributing
list: include
---

To automate functional integration and end-to-end testing Meshery uses [Playwright](https://playwright.dev) as one of the tools to automate browser testing. End-to-end tests run with each pull request to ensure that the changes do not break the existing functionality.

## Executing Tests

Clone the `meshery/meshery` repo. Tests are written in JavaScript. Navigate to the [/ui/tests/e2e/](https://github.com/meshery/meshery/tree/master/ui/tests/e2e) directory.

Install the dependencies by running the following command:

```bash

npm install

npm playwright install --with-deps

```

To run the tests, you can use the following commands:

```bash

// for running the whole test suite with all browsers
npm run test:e2e

// for running only on chromium
npm run test:e2e:chromium

// for only running fast tests
npm run test:e2e:fast

```

<!-- ## Writing and Organizing Tests 


### Best Practices

-->
