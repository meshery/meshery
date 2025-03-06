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

- A verified account in your chosen provider which integrate with Meshery.
- A compatible browser such as Chromium, Chrome, or Firefox.
- Installations of Golang, NodeJS, and Makefiles for Native OS build (Optional for docker based build).
- Kubernetes clusters (Required for connection to Kubernetes test cases)
- Already have [Meshery Adapters](https://docs.meshery.io/concepts/architecture/adapters) up and running (Required for adapters test cases)

## Setting up environment variable

To run the tests successfully, three environment variables must be configured:  
• `REMOTE_PROVIDER_USER_EMAIL` (Required): The email associated with your account within your provider.  
• `REMOTE_PROVIDER_USER_PASSWORD` (Required): The password associated with your account within your provider.  
• `PROVIDER_TOKEN` (Optional): Your provider token, can be generated from an account registered within your provider  

{% include alert.html
    type="info"
    title="Accessing Remote Providers"
    content='In the case you are using Layer5 Cloud as a remote provider, you can <a href="https://cloud.layer5.io/security/tokens">generate a token from your user account</a> to use while writing and executing tests.' %}

During the setup phase, Playwright utilizes these environment variables to log in and store credentials securely in the `playwright/.auth` directory. To protect sensitive data, the `.gitignore` file is configured to exclude the `.env` file and any JSON files within the `/playwright/.auth` directory from the GitHub repository.

There are several tools to help you to working with environment variables locally for each project such as [direnv](https://github.com/direnv/direnv), it can work across multiple shell such as Bash, Powershell, Oh my zsh, Fish, etc

## Starting up Meshery UI and Server

There are a few ways to set up the Meshery UI and server, but for end-to-end testing, we aim to get as close to a production environment as possible. We know developers might need to make some tweaks for UI and Server. Rebuilding the whole project can take time, and we don’t support hot reload because it’s more for development than for end-to-end testing.

{% include alert.html type="warning" title="Several Test may break" content='Some test cases required you to have kubernetes cluster and build meshery adapter as well, be aware of that. Which is out of scope for this documentation<ul><li><a href="https://docs.meshery.io/installation/kubernetes/minikube">Kubernetes Cluster</a>: Installation of kubernetes cluster with Minikube.</li>
<li><a href="https://docs.meshery.io/installation/multiple-adapters">Meshery Adapters</a>: Using Multiple Adapters</li></ul>' %}

### Native OS Build (Recommended)

This approach is very quick to build, but also dependent on your operating system, so you need to have all dependencies necessary to be able compile and running the server.

- Install & Build the NextJS application for both the UI and UI Provider

```bash
make ui-build
```

- Compile the Golang into binary file for Meshery Server

```bash
make build-server
```

- Run the Meshery Server on localhost port 9081

```bash
make server-binary
```

### Meshery CLI

There is also Meshery CLI which can help you run the UI and Server, for more detail, you go to [Meshery CLI documentation](https://docs.meshery.io/project/contributing/contributing-cli-guide#process)

### Docker Based Build

Alternatively, a Docker-based setup can be utilized, simplifying the process, and ensuring consistency across different environments. It is closer to the production environment than the native solution but slower in terms of build time.

- Build the docker container locally:

```bash
make docker-testing-env-build
```

- Run the docker container on port 9081

```bash
make docker-testing-env
```

## Setup Playwright

For Playwrights, always try to use a native OS whenever possible. The Docker-based approach is intended only for unsupported OSes and is generally not recommended because it runs on top of Ubuntu images, which can be redundant if you already using Ubuntu or Windows.

### Playwright on Native OS (Recommended)

Setup playwright:

```bash
make test-setup-ui
```

Run the all project and test cases:

```bash
make test-ui
```

### Playwright server on docker based image

The first step is to pull the docker image from [Azure Container Registry](https://mcr.microsoft.com/en-us/product/playwright/tags) where the playwright stores their image using this command:

```bash
docker pull mcr.microsoft.com/playwright:<version>-<base-image>
```

{% include alert.html type="warning" title="Playwright Versioning" content="Make sure the version you are using matches the version of `@playwright/test` in the `package.json` dev dependencies" %}

Here is the example of pulling playwright v1.44.0 with Ubuntu 22.04 LTS

```bash
docker pull mcr.microsoft.com/playwright:v1.44.0-jammy
```

Starting up playwright docker server:

```bash
docker run --rm --network host --init -it mcr.microsoft.com/playwright:v1.44.0-jammy /bin/sh -c "cd /home/pwuser && npx -y playwright@1.44.0 run-server --port 8080"
```

{% include alert.html type="warning" title="Unsafe Environment" content="Keep in mind this is just for development purposes inside your local system and don’t try to expose your container network to the host system using --network host on production or CI" %}

In the last step go to ui folder, 
```bash
cd ui;
```

## Run the test cases with Playwright CLI

There are several options we can use to run the test cases, in CLI:

To run playwright UI mode using the browser, you can add `--ui` in the cli, for example:
```bash
npx playwright test --ui
```

If you are using playwright from docker, you can use `--ui-port=<playwright-docker-server>`, for example:
```bash
npx playwright test --ui-port=8080
```

To run playwright for specific project only, for example meshery-provider, you can run this command:
```bash
npx playwright test --ui --project=chromium-meshery-provider 
```

To run specific test, you can add the test file location, for example:
```bash
npx playwright test --ui --project=chromium-meshery-provider tests/e2e/service-mesh-performance.spec.js 
```

For more detail, you can read the [Playwright Cli docs](https://playwright.dev/docs/test-cli)


## Testing Meshery & Local Provider 

By default our test cases is running against both Meshery and Local Provider, we are utilizing playwright feature such as:

- StorageState: In meshery [setup auth](https://github.com/meshery/meshery/blob/master/ui/tests/e2e/auth.setup.js), we have 2 storage state, which store a session for Meshery and Local provider.
- Project: After the setup completes, it will run the project-based test depending on which storage state for the Local Provider and one for the Meshery Provider
- Test Parameterize:  In the Local provider we are limiting some features to test against. For the missing features, we leverage this playwright feature to check or even skip the test. If it is not possible to run then you need to specify the `provider` directly from the test, and make sure the test is wrapped using:

```javascript
import { expect, test } from './fixtures/project';

test('Random test', async ({ provider }) => {
  if (provider === "Meshery") {
    // Run this for testing Meshery provider
  }

  if (provider === "None") {
    // Run this for testing Local provider
  }
});
```

## Testing Policy

To maintain consistency across test cases, every new test will be tagged with `@unstable`. This will ensure that it appears with a warning icon rather than a failing icon in the test reporter comments on pull requests. For example:

```javascript
import { expect, test } from './fixtures/project';

test('Random test',  { tag: '@unstable' }, async ({ provider }) => {
  // Test cases here
  // ...
});
```

## Debugging Test on Github Actions

End-to-end test results are stored as artifacts on every PR in Github Actions. In case you need to debug a failed test:

- Visit the PR in question. Go to the bottom of PR directly above the comment.
- Wait until all Github Actions completed, and scroll until you see `Meshery UI and Server / UI end-to-end tests`  
- Click details and it will redirect you to the actions workflow
- Go to summary tab, scroll down until you see artifact, and check the artifact `playwright-report`
- Download the artifact
- Extract the file into a folder
- Go to [Playwright Trace Page](https://trace.playwright.dev/)
- From the test folder pick one folder which represents the test, you want to check
- Upload the trace file


Watch the training session on Playwright testing and trace debugging.

<iframe width="560" height="315" src="https://www.youtube.com/embed/x-W60mvDYuo?si=coN7RpRjkI4a_ndk&amp;start=1524" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
