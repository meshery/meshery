---
layout: default
title: Build & Release (CI)
permalink: project/build-and-release
---

Meshery’s build and release system incorporates many tools, organized into different workflows each triggered by different events. Meshery’s build and release system does not run on a schedule, but is event-driven. GitHub Actions are used to define Meshery’s CI workflows. New builds of Meshery and its various components are automatically generated upon push, release, and other similar events, typically in relation to their respective master branches.

## Artifacts

Today, Meshery and Meshery adapters are released as Docker container images, available on Docker Hub. Meshery adapters are out-of-process adapters (meaning not compiled into the main Meshery binary), and as such, are independent build artifacts.The process of creating Docker images, tagging with the git commit SHA and pushing to Docker Hub is being done automatically using GitHub Actions.

### Artifact Repositories

Artifacts produced in the build processes are published and persisted in different public repositories and in different formats.

| Location      | Project       | Repository    |
| ------------- | ------------- | ------------- |
| Docker Hub    | Meshery       | [https://hub.docker.com/r/layer5/meshery](https://hub.docker.com/r/layer5/meshery) |
| GitHub        | mesheryctl    | [https://github.com/layer5io/meshery/releases](https://github.com/layer5io/meshery/releases) |
| Docker Hub    | Meshery Adapter for \<service-mesh\> | https://hub.docker.com/r/layer5/meshery-\<service-mesh\> |
| Docs          | Meshery Documentation | [https://docs.meshery.io](https://docs.meshery.io) |
| GitHub        | [Service Mesh Performance](https://smp-spec.io) | [https://github.com/layer5io/service-mesh-performance](https://github.com/layer5io/service-mesh-performance) |

## Secrets

Some portions of the workflow require secrets to accomplish their tasks. These secrets are defined within the respective repositories and accessible to workflows during runtime. Currently defined secrets include:

- `DOCKER_USERNAME`: Username of the Docker Hub user with the right privileges to push images
- `DOCKER_PASSWORD`: Password for the Docker Hub user
- `GO_VERSION`: As of December 9th 2020 it is 1.15
- `IMAGE_NAME`: appropriate image name for each of the Docker container images. All are under the `layer5io` org.
- `SLACK_BOT_TOKEN`: Used for notification of new GitHub stars given to the Meshery repo.
- CYPRESS_RECORD_KEY`: Used for integration with the Layer5 account on Cypress.
- `GLOBAL_TOKEN`: Used for securely transmitting performance test results for the None Provider.

The Docker Hub user, `mesheryci`, belongs to the "ciusers" team in Docker Hub and acts as the service account under which these automated builds are being pushed. Every time a new Docker Hub repository is created we have to grant “Admin” (in order to update the README in the Docker Hub repository) permissions to the ciusers team.

## Checks and Tests
Meshery’s CI workflow incorporates several checks (partial list below) during merges and/or commits to any branches and pull requests to master branch to prevent broken code from being merged into master. 

Collectively, Meshery repositories will generally have CI workflow for commits and pull requests that consist of the following actions:

- Lint check (golint)
- Static analysis check (staticcheck)
- Vet (govet)
- Security checks (gosec)
- Unit tests (go tests)
- Build (go build)
- Release binaries through GoReleaser (only for mesheryctl in the Meshery repository)
- Docker build, tag and push

## Automated Builds

All Meshery GitHub repositories are configured with GitHub Actions. Everytime a pull request is submitted against the master branch of any repository, that repository’s GitHub Actions will be invoked (whether the PR is merged or not). Workflows defined in Meshery repository will generally (but not always) perform the following actions:

1. trigger a Docker build to build a Docker container image
1. generate two Docker tags: 
   1. a tag containing the git merge SHA
   1. a tag containing that particular release’s git tag (if one is present)
1. assign each of these two tags to the new container image as well as the latest tag.
1. push the new Docker tags and image to Docker Hub.

### Building `mesheryctl`

As a special case, the meshery repository contains an additional artifact produced during each build. This artifact is mesheryctl which is built as an executable binary. In order to make the job of building mesheryctl easier for a combination of different platform architectures and operating systems, we are using [GoReleaser](https://goreleaser.com). Irrespective of branch, for every git commit and git push to the meshery repository, GoReleaser will execute and generate the OS and arch-specific binaries ( but will NOT publish them to GitHub). Even though mesheryctl binaries are built each time a pull request is merged to master, only stable channel artifacts are published (persisted).

### Releasing `mesheryctl` to GitHub

Only when a git tag containing  a semantic version number is present (is a commit in the master branch) will GoReleaser execute, generate the archives, and also publish the archives to [Meshery’s GitHub releases](https://github.com/layer5io/meshery/releases) automatically. GoReleaser is configured to generate artifacts for the following OS, ARCH combination:

- Darwin - i386, x86_64
- Linux - i386, x86_64
- Windows - i386, x86_64
- FreeBSD - i386, x86_64

The artifacts will be made available as a tar.gz archive for all the operating systems. mesheryctl is bundled into packages for commonly used package managers: homebrew and scoop.

#### Homebrew

GoReleaser facilitates the creation of a brew formula for mesheryctl. The [homebrew-tap](https://github.com/layer5io/homebrew-tap) repository is the location for Layer5’s brew formulas.

#### Scoop

GoReleaser facilitates the creation of a Scoop app for mesheryctl. The [scoop-bucket](https://github.com/layer5io/scoop-bucket) repository is the location of Layer5’s Scoop bucket.

## Release Versioning

We follow the commonly used semantic versioning for Meshery, Meshery Adapter and Performance Benchmark Specification releases. Given a version number MAJOR.MINOR.PATCH.BUILD, increment the:

- MAJOR version - major changes with rare potential for incompatible API changes.
- MINOR version - add functionality in a backwards-compatible manner.
- PATCH version - mostly for bug and security fixes.
- AlPHA/BETA/RC - used to facilitate early testing of an upcoming release.

### Component Versioning

Meshery comprises a number of components including a server, adapters, UI, and CLI. As an application, Meshery is a composition of these different functional components. While all of Meshery’s components generally deploy as a collective unit (together), each component is versioned independently, so as to allow them to be loosely coupled and iterate on functionality independently.  Some of the components must be upgraded simultaneously, while others may be upgraded independently. See [Upgrading Meshery](/guide/upgrade) for more information.
