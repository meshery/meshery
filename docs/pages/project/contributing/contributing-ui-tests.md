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

To automate functional integration and end-to-end testing Meshery uses [Playwright](https://playwright.dev/) as one of the tools to automate browser testing. End-to-end tests run with each pull request to ensure that the changes do not break the existing functionality.

## Prerequisites:

Before diving into Meshery's testing environment, certain prerequisites are necessary:

- A verified Layer5 cloud account.
- A compatible browser such as Chromium, Chrome, or Firefox.
- Installations of Golang, NodeJS, and Makefiles for local OS setups (Optional for docker based build).
- Kubernetes clusters (Optional but several test cases will fail)
- Already setting up Meshery adapters for several test cases (Optional but several test cases will fail)

## Setting up environment variable

To run the tests successfully, three environment variables must be configured:
• `REMOTE_PROVIDER_USER_EMAIL`: The email associated with your Layer5 cloud account.
• `REMOTE_PROVIDER_USER_PASSWORD` : The password for your Layer5 cloud account.
• `PROVIDER_TOKEN`: You're provider token, that can be generated from [Layer5 cloud account](https://meshery.layer5.io/security/tokens)

During the setup phase, Playwright utilizes these environment variables to log in and store credentials securely in the `playwright/.auth` directory. To protect sensitive data, the `.gitignore` file is configured to exclude the `.env` file and any JSON files within the `/playwright/.auth` directory from the GitHub repository.

There are several tools to help you to working with environment variables locally for each project such as [direnv](https://github.com/direnv/direnv), it can work across multiple shell such as Bash, Powershell, Oh my zsh, Fish, etc` dev dependencies

## Starting up Meshery UI and Server

There are several methods to set up the Meshery UI and server, but for e2e testing scenario we aim for environment which close to production as possible and also we aware there some minimum changes developer also need to make sometimes. So rebuilding the entire project will take time and we are not supporting hot reload because it's optimizing for development not e2e testing

{% include alert.html type="warning" content="Some test cases required you to have kubernetes cluster and build meshery adapter as well, be aware of that. Which is out of scope for this documentation" %}

### Native OS

This approach is very quick to build, but also dependent on your operating system, so you need to have all dependencies necessary to be able compile and running the server.

- Install & Build the NextJS static site generator application for both the UI and UI Provider

```bash
make ui-build
```

- Compile golang server

```bash
make build-server
```

- Run the server locally on port 8080

```bash
make server-binary
```

### Docker Based

Alternatively, a Docker-based setup can be utilized, which simplifies the process and ensures consistency across different environments.

- Build the docker container locally:

```bash
make docker-testing-env-build
```

- Run the docker container on port 9081

```bash
make docker-testing-env
```

## Setup playwright & Run the test cases

### Native OS

Setup playwright:

```bash
make test-setup-ui
```

Run the test cases:

```bash
make test-ui
```

### Docker based

The first step is to pull the docker image from [Azure Container Registry](https://mcr.microsoft.com/en-us/product/playwright/tags) where the playwright stores their image using this command:

```bash
docker pull mcr.microsoft.com/playwright:<version>-<base-image>
```

{% include alert.html type="warning" content="Make sure the version you are using matches the version of `@playwright/test` in the `package.json` dev dependencies" %}

Here is the example of pulling playwright v1.44.0 with Ubuntu 22.04 LTS

```bash
docker pull mcr.microsoft.com/playwright:v1.44.0-jammy
```

Starting up playwright docker server:

```bash
docker run --rm --network host --init -it mcr.microsoft.com/playwright:v1.44.0-jammy /bin/sh -c "cd /home/pwuser && npx -y playwright@1.44.0 run-server --port 8080"
```

In the last step, run this command to run the test cases:

```bash
PW_TEST_CONNECT_WS_ENDPOINT=ws://localhost:8080/ npx playwright test
```

<!-- ## Writing and Organizing Tests


### Best Practices

-->

