---
layout: page
title: Contributing to Meshery CLI End-to-End Tests
permalink: project/contributing/contributing-cli-tests
abstract: How to contribute to Meshery Command Line Interface end-to-end testing with BATS.
language: en
type: project
category: contributing
list: include
display-title: false
---

{% include alert.html
    type="danger"
    title="DOCUMENT INCOMPLETE"
    content="This document is incomplete and is still under improvement. Help wanted! 😃" %}


Meshery CLI is the command line interface for Meshery. Meshery CLI, otherwise known as `mesheryctl`, is a client of Meshery Server's [REST API]({{site.baseurl}}/extensibility/api). It provides a way to interact with Meshery and perform various operations such as installing, configuring, and managing cloud native infrastructure.

This document is intended to help you contribute to the end-to-end tests for `mesheryctl`, the Meshery CLI. It is designed to be a guide for developers who are new to the project and want to contribute to the testing of `mesheryctl`.

The end-to-end tests for `mesheryctl` are designed to ensure that the CLI is working as expected and that it is compatible with the various cloud native infrastructure and public cloud services that Meshery supports. These tests are run automatically on every pull request to ensure that any changes made to the code do not break the existing functionality of the CLI.

{% include alert.html
    type="info"
    title="Meshery CLI Reference Documents"
    content='<ul><li><a href="/project/contributing/contributing-cli">Contributing to Meshery CLI</a></li><li><a href="https://docs.google.com/spreadsheets/d/1q63sIGAuCnIeDs8PeM-0BAkNj8BBgPUXhLbe1Y-318o/edit#gid=0">Meshery Command Tracker</a>: Status of mesheryctl command implementation and platform compatibility.</li>
    <li><a href="https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#">Meshery CLI Commands and Documentation</a>: Detailed documentation of the `mesheryctl` commands.</li>
	<li><a href="https://github.com/meshery/meshery/labels/component%2Fmesheryctl">mesheryctl open issues and pull requests</a>: Matching the "component/mesheryctl" label.</li></ul>' %}

## Familiarizing with Meshery CLI End-to-End Testing

End-to-end testing of `mesheryctl` uses the [Bash Automated Testing System](https://github.com/bats-core/bats-core) (BATS) framework to define and execute CLI tests. Each test case is designed to mimic the experience that a Meshery CLI user might have while interacting with `mesheryctl` in their terminal of choice. In this sense, `mesheryctl` tests run end-to-end with each pull request submitted containing changes to either the `/mesheryctl` or the `/server` directories in the `meshery/meshery` repository, ensuring that changes included in those pull requests do not break the existing CLI functionality.

### Prerequisites

Before diving into `mesheryctl`'s testing environment, certain prerequisites must be met to ensure a smooth testing experience. These prerequisites include:

- A working installation of Meshery CLI and Meshery Server.
  - A verified account in your chosen provider which integrate with Meshery.
- A working installation of the BATS testing framework.
  - `bash` as shell terminal.
- [Optional] A working installation of a Kubernetes cluster (Minikube, Kind, etc.) for testing Kubernetes-related functionality.
- `jq` and `yq`, tools for processing JSON and YAML inputs, respectively.

### Setup Bats Core

For Bats Core, always try to use a BATS-native OS whenever possible. This is because BATS Core does not support Windows. If you are using Windows, you can use WSL (Windows Subsystem for Linux) to run BATS Core. See the official [BATS installation documentation](https://bats-core.readthedocs.io/en/stable/installation.html) for more information on how to install BATS Core on your system. Here are quick start steps.

#### MacOS (homebrew)

```bash
brew install bats-core
```

#### Any OS (npm)

```bash
npm install -g bats
```

#### Windows (from source via bash)

Check out a copy of the Bats repository and install it to $HOME. This will place the bats executable in $HOME/bin, which will need to be added in $PATH.

```bash
git clone https://github.com/bats-core/bats-core.git
cd bats-core
./install.sh $HOME
```

Some tests could use bats libraires as helpers to create the tests suite.

#### Setup Dependencies

There are needed dependencies to test whether the server is up and running. Resolve them by navigating to `meshery/mesheryctl` directory and then executing:

```bash
  make e2e-libs
```

### Starting Meshery Server

There are a few ways to set up the Meshery server, but for end-to-end testing, we aim to get as close to a production environment as possible. We know developers might need to make some tweaks for Server. Rebuilding the whole project can take time, and we don’t support hot reload because it’s more for development than for end-to-end testing.

```bash
make server
```

Be aware that some test cases require the availability of a Kubernetes cluster and one or more  Meshery Adapters. In those cases, please refer to the [installation guides]({{site.baseurl}}/installation)) (like that of [installing Meshery on Minikube]({{site.baseurl}}/installation/kubernetes/minikube)). 

### Authentication

To run the tests successfully, you need be logged in to your Meshery account. This is required to access the Meshery Server and perform operations on it. Whether using the local or a remote provider, you will need to generate a token from your user account to use while writing and executing tests.

**Remote Provider Example**

If you are using Layer5 Cloud as your remote provider, you can [generate and download an API token](https://cloud.layer5.io/security/tokens) from your user account for use while writing and executing tests.

### Verify your API Token

```bash
mesheryctl system check
```

If you see this error message - `Error: !! Authentication token not found. Please supply a valid user token. Login with mesheryctl system login`, you will need to authenticate using the command:

```bash
mesheryctl system login
```


## Writing End-to-End Test Cases

Once all prerequisites have been met, you are setup to run, add, update new test cases.To keep the development of tests consistents and maintainable, we have put some guidance on implementation process which are describe above.

### Folder structure

The tests will be available in[`mesheryctl/tests/e2e`](https://github.com/meshery/meshery/tree/master/mesheryctl/tests/e2e) folder with the following structure which is done to follow the code base structure of `mesheryctl`.

```shell
├── aaa-<command>
├── bbb-<command>
├── ...
├── yyy-<command>
├── zzz-<command>
├── helpers
├── README.md
├── run_tests.bash
├── setup_suite.bash
└── teardown_suite.bash
```

### Test Case Naming Convention

Bats will run the tests sequentially so we use a prefix code alphanumeric to be able to order tests as required  when necessary. The prefix code is used to order the tests in a way that makes sense for the command under test. The prefix code is used to group the tests by command and subcommand. The prefix code is also used to order the tests within the command and subcommand. For example, you need to import a [Meshery Model]({{site.baseurl}}/concepts/logical/models) before being able to view the details of that model

**Test Folder**
- **prefix**: 3 digits alphanumeric code
- **name**: command under test
  
Example: `011-design`

**Test File**

- **prefix**: 2 digits numeric code
- **name**: subcommand under test
- **extension**: bats

Example: `05-list.bats`

For consistency, we will keep the prefix *00-* for the command under test in the folder and subcommands will start at *01-*. Here an example with `mesheryctl model` command:

```bash
002-model/
├── 00-model.bats
└── 01-model-list.bats
```

## Run End-to-End (locally)

<!-- 
    TODO: Add make e2e support with following changes
    1. move to tests/e2e/helpers
    2. Assert the following lbats libraries are available if not git clone
    - https://github.com/bats-core/bats-file.git
    - https://github.com/bats-core/bats-detik.git
    - https://github.com/bats-core/bats-support.git
    3. back to tests/e2e
    5. run bats *-*/*.bats

```bash
make cli-tests _not yet implemented_
```
-->

Make sure you are in `meshery/mesheryctl` directory

**Run all tests** 

```bash
make e2e
```
<!-- TODO: https://github.com/meshery/meshery/issues/14105
**Run a specific test file**

Switch to the directory containing the test file and execute:

```bash
MESHERYCTL_BIN=<path to mesheryctl binary> bats <file_name>.bats
```
-->

**More on running tests locally**

Breaking down the execution of `make e2e`, two commands are executed before subsequent commands are run; 
- `make`: which builds the binary and
- `e2e-libs`: which gets needed dependencies for running the tests.

These steps can become redundant if you're not running the end-to-end tests for the first time.

It is important to point out that there are other ways to run end-to-end tests locally.
To use these, ensure you are in the `meshery/mesheryctl/tests/e2e` directory.

**Run tests with already built binary**

This excludes the need to build the binary everytime there is an attempt to run all the tests.

```bash
  bash run_tests_local.sh
```

NB: This works if there is an existing `mesheryctl` binary. If there isn't, the binary will be built. 

**Enforce rebuilding the  binary**

This involves parsing a flag for the binary to be built whether it exists or not. This comes in handy when you have local changes and possibly will like to test.

```bash
  bash run_tests_local.sh -b
```

### Find Tests here
Refer to [Meshery Test Plan](https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit?usp=sharing) for test scenarios.

To filter and view only CLI-related test cases using the Sheet Views feature:
1. In the top menu bar, click Data → Change view
2. Choose the pre-defined view labeled "CLI"

![Meshery Test Plan Screenshot](/assets/img/contributing/meshery-test-plan-v0.8.0-ui.png)

## Developement

**End-to-End Tests:**

* **Purpose:** Validate the complete flow of an application, from the user interface (in this case, the `mesheryctl` CLI) down to the underlying systems (Meshery server, Kubernetes, etc.).
* **Scope:** Test multiple components working together to achieve a specific user scenario.
* **Speed:** Generally slower to execute compared to unit and integration tests as they involve setting up and interacting with a real or near-real environment.
* **Cost:** Can be more expensive to set up and maintain due to the complexity of the environment.
* **Frequency:** Typically run less frequently than unit and integration tests, often as part of continuous integration pipelines or release processes.

For `mesheryctl`, E2E tests will focus on verifying that CLI commands perform their intended actions against a running Meshery instance **in context of how users experience the CLI**, this means we are focusing on UX. This might involve deploying applications, managing connections, interacting with Meshery features, and verifying the expected outcomes.


### Implementation

We will exclusively use the Bats Core framework and its built-in functionalities for writing E2E tests. This ensures consistency and leverages a well-established testing tool for shell scripts.

**Key Principles:**

* **Pure Bats Core:** Avoid relying on external custom scripts or libraries beyond what Bats Core provides, while there might be occasional need to deviate from the library. Take in consideration that doing so all increases the possibility for bugs as well as our sustaining costs.
* **Focus on `mesheryctl`:** The tests should primarily interact with the `mesheryctl` CLI.
* **Clear Assertions:** Use Bats Core's assertion functions (`assert`, `assert_success`, `assert_failure`, `assert_output`, etc.) to verify expected outcomes.
* **Setup and Teardown:** Utilize `setup()` and `teardown()` functions to prepare the testing environment and clean up afterwards.
* **Helper Scripts:** If a custom script or function is absolutely necessary to facilitate testing (and cannot be achieved with standard Bats Core), it **must** be created as a `.bash` file within the `helpers` folder. Each helper script/function should have a clear description of its purpose within the file itself. Avoid inline custom scripting within the test files.
* **Consistency**: 

#### Test Naming Convention

It must follow this naming convention

```
<mesheryctl command> [subcommand] <execution context> <expected result>
```

**Example:**

```bash
@test "mesheryctl model view without model name should display an error" {
  ... test implementation ...
}
```

#### Test Data

If a command requries a specific id, name or any predefined value ensure that the data is created by your test or another test beforehand. Do not rely on external or uncontrolled data as it will lead to unexpected results.

**Example:**

In the following example, we must have create a model with the name `model-test` before creating or running the following test

```bash
@test "mesheryctl model view providing a model name should display model information" {
  run MESHERYCTL_BIN model view model-test
  ... ...
}
```

#### Writing E2E Tests with Bats Core

Official documentation is available at [https://bats-core.readthedocs.io/en/stable/](https://bats-core.readthedocs.io/en/stable/)

Github organization [https://github.com/bats-core](https://github.com/bats-core) contains bats-core repository and also bats libraries repositoires  

1.  **Basic Test Structure:** A Bats test file consists of one or more test cases defined using the `@test` keyword.

    ```bash
    #!/usr/bin/env bats

    @test "Ensure mesheryctl version command works" {
      run mesheryctl version
      assert_success
      assert_output --partial "mesheryctl version"
    }
    ```

2.  **Interacting with `mesheryctl`:** Execute `mesheryctl` commands within your test cases using the `run` command. Capture the output and exit status for assertions.

    ```bash
    @test "Deploy a sample application" {
      run mesheryctl design apply -f samples/apps/nginx.yaml
      assert_success
      assert_output --partial "Successfully applied"
    }
    ```

3.  **Assertions:** Use Bats Core's assertion functions to validate the results of your `mesheryctl` commands. Refer to the Bats Core documentation for a complete list of assertions.

    * `assert_success`: Checks if the command exited with a status code of 0.
    * `assert_failure`: Checks if the command exited with a non-zero status code.
    * `assert_output`: Checks if the command's output matches a given string.
    * `assert_output --partial`: Checks if the command's output contains a given substring.
    * `assert_equal`: Checks if two strings are equal.
    * `assert_not_equal`: Checks if two strings are not equal.
    * `assert_file_exists`: Checks if a file exists.
    * `assert_file_contains`: Checks if a file contains a given string.

4.  **Setup and Teardown:**

    * `setup()`: This function is executed before each test case within a file. Use it to set up the necessary environment for your tests (e.g., ensure Meshery is running, configure connections).

    * `teardown()`: This function is executed after each test case within a file. Use it to clean up any resources created during the test (e.g., undeploy applications, reset configurations).

5.  **Test Data:** If your tests require specific input files (e.g., Kubernetes manifests), store them in a relevant directory (e.g., `tests/e2e/fixtures`).

6.  **Error Handling:** Consider how your tests will handle errors. Use appropriate assertions to check for expected failures and provide informative error messages.

### Bug Reporting During Test Implementation

While implementing new E2E tests, you might discover bugs in `mesheryctl`. It is crucial to report these findings properly:

1.  **Create a New Issue:** Navigate to the [meshery/meshery](https://github.com/meshery/meshery/issues) repository on GitHub and click "New issue."

2.  **Use the `mesheryctl Bug Report` Template:** Look for and select the "mesheryctl Bug Report" issue template. This template provides a structured format for reporting bugs related to the CLI.

3.  **Provide Detailed Information:** Fill out the sections of the bug report template with as much detail as possible, including:

    * **Steps to reproduce:** Clearly outline the exact `mesheryctl` commands and environment setup that led to the bug.
    * **Expected behavior:** Describe what you anticipated `mesheryctl` to do.
    * **Actual behavior:** Describe what `mesheryctl` actually did. Include any error messages or unexpected output.
    * **Environment details:** Provide information about your operating system, `mesheryctl` version (obtained using `mesheryctl version`), and the Meshery environment (e.g., local Docker, remote Kubernetes).
    * **Relevant logs:** Include any relevant logs from `mesheryctl` or the Meshery server that might help diagnose the issue.
    * **Context:** Explain the context in which you encountered the bug (e.g., while testing a specific feature).

4.  **Link to the Test:** If the bug was discovered while writing a specific test, mention the test file and test case in the issue.
