---
title: Contributing to Meshery CI/CD Toolchain
description: How to contribute to Meshery's CI/CD toolchain.
categories: [contributing]
---

Meshery's CI/CD toolchain automates and standardizes how the project is built, tested, and released. Its overall build and release process is documented in [Build & Release (CI)]({{< ref "project/contributing/build-and-release.md" >}}).

{{% alert color="info" title="CI/CD Toolchain Reference Documents" %}}<ul>
<li><a href="https://www.gnu.org/software/make/manual/make.html">GNU Make Manual</a>: Reference for Makefile syntax and behavior.</li>
<li><a href="https://docs.docker.com/reference/dockerfile/">Dockerfile Reference</a>: Reference for Dockerfile instructions and syntax.</li>
<li><a href="https://docs.github.com/en/actions">GitHub Actions Documentation</a>: Reference for authoring workflows and actions.</li>
</ul>{{% /alert %}}

## Makefiles

Meshery's Makefile is organized around a set of core recipes that follow a shared standard. Any change to one of these recipes must preserve that standard.

Recipes follow a `noun-verb` naming scheme, in which the noun identifies the component the recipe acts upon and the verb identifies the action performed upon it. A recipe that serves the documentation site is therefore named `site-serve` rather than `serve-site`. Naming recipes in this order keeps those that act upon the same component grouped together.

A recipe resolves directly to the local command it runs rather than delegating to another script. For example, a recipe that serves the documentation site locally invokes the tool directly:

{{< code code=`site:
	hugo server` >}}

Delegating to a further script, such as `npm run hugo`, is done only when that indirection is necessary. Otherwise, Makefile changes follow ordinary Make conventions.

## Dockerfiles

Meshery's Dockerfiles define how its components are packaged into container images. Each Dockerfile follows a shared standard that any change must preserve. The following practices apply to Dockerfile changes:

- Pin base images and installed dependencies to explicit versions so that builds remain reproducible.
- Order instructions from least to most frequently changed to preserve layer caching.
- Use multi-stage builds to keep build-time tooling out of the final image.
- Consolidate related commands to reduce the number of layers.

Every change must produce an image that builds without errors.

## GitHub Actions

Meshery's workflows rely on a number of GitHub Actions that are maintained in their own repositories. Because an action runs only as part of a workflow, contributing a change to one follows a fork-and-test cycle: the action is developed on a branch and exercised through a workflow that references that branch.

### Developing the action on a branch

1. Fork the action's repository.
2. Create a branch on your fork.
3. Make your changes and push them to the branch.

The fork and branch conventions used throughout the project are documented in [Contributing to Meshery using git]({{< ref "project/contributing/contributing-gitflow.md" >}}).

### Testing the action

Testing the action can be accomplished in two ways:

1. Within the fork: add a workflow to your fork that invokes the action and runs on each push, so that every commit exercises the action automatically.
2. From a separate repository: maintain a separate repository containing a workflow that references the action, and trigger it manually to run each test.

### Submitting the change

Once the test workflow has run successfully, include evidence of that run in the pull request. The evidence consists of the workflow logs and a showcase of the successful execution, such as a screenshot of the completed run. This evidence allows reviewers to verify the action's behavior without reproducing the test setup.

{{% alert color="warning" title="Successful runs must be demonstrated" %}}A pull request that modifies a GitHub Action is ready for review only once it includes the workflow logs and a showcase demonstrating that the action ran successfully.{{% /alert %}}