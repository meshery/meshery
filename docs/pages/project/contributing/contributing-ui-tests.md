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

- A verified account in your choosen provider which integrate with Meshery.
- A compatible browser such as Chromium, Chrome, or Firefox.
- Installations of Golang, NodeJS, and Makefiles for Native OS build (Optional for docker based build).
- Kubernetes clusters (Required for several test cases)
- Already have [Meshery Adapters](https://docs.meshery.io/concepts/architecture/adapters) up and running (Required for several test cases)

## Setting up environment variable

To run the tests successfully, three environment variables must be configured:  
• `REMOTE_PROVIDER_USER_EMAIL`: The email associated with your account within your provider.  
• `REMOTE_PROVIDER_USER_PASSWORD` : The password associated with your account within your provider.  
• `PROVIDER_TOKEN`: You're provider token, that can be generated from your provider account  

In the case you are using Layer5 Cloud as provider, you can generate your token on [Layer5 cloud account token](https://meshery.layer5.io/security/tokens)

During the setup phase, Playwright utilizes these environment variables to log in and store credentials securely in the `playwright/.auth` directory. To protect sensitive data, the `.gitignore` file is configured to exclude the `.env` file and any JSON files within the `/playwright/.auth` directory from the GitHub repository.

There are several tools to help you to working with environment variables locally for each project such as [direnv](https://github.com/direnv/direnv), it can work across multiple shell such as Bash, Powershell, Oh my zsh, Fish, etc

## Starting up Meshery UI and Server

There are a few ways to set up the Meshery UI and server, but for end-to-end testing, we aim to get as close to a production environment as possible. We know developers might need to make some tweaks for UI and Server. Rebuilding the whole project can take time, and we don’t support hot reload because it’s more for development than for end-to-end testing.

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

Alternatively, a Docker-based setup can be utilized, which simplifies the process, ensures consistency across different environments and also more close to production environment as well compare to native solution. 

- Build the docker container locally:

```bash
make docker-testing-env-build
```

- Run the docker container on port 9081

```bash
make docker-testing-env
```

### Meshery CLI

There is also Meshery CLI which can help you run the UI and Server, for more detail, you go to [Meshery CLI documentation](https://docs.meshery.io/project/contributing/contributing-cli-guide#process)

## Setup playwright & Run the test cases

For playwrights always try to use a Native OS as possible, the docker based approach is for unsupported OS only and always not recommended because it runs on top of Ubuntu images, which is very washed if you were already running Ubuntu or Windows.

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

{% include alert.html type="warning" content="Keep in mind this is just for development purposes inside your local system and don’t try to expose your container network to the host system using --network host on production or CI" %}

In the last step, run this command to run the test cases:

```bash
PW_TEST_CONNECT_WS_ENDPOINT=ws://localhost:8080/ npx playwright test
```