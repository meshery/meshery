---
layout: page
title: Build & Release (CI)
permalink: project/contributing/build-and-release
description: Details of Meshery's build and release strategy.
language: en
type: project
category: contributing
---

Meshery’s build and release system incorporates many tools, organized into different workflows each triggered by different events. Meshery’s build and release system does not run on a schedule, but is event-driven. GitHub Actions are used to define Meshery’s CI workflows. New builds of Meshery and its various components are automatically generated upon push, release, and other similar events, typically in relation to their respective master branches.

## Artifacts

Today, Meshery and Meshery adapters are released as Docker container images, available on Docker Hub. Meshery adapters are out-of-process adapters (meaning not compiled into the main Meshery binary), and as such, are independent build artifacts and Helm charts.The process of creating Docker images, tagging with the git commit SHA and pushing to Docker Hub is being done automatically using GitHub Actions. And When the contribution includes content of Helm chart of Meshery and Meshery Adapter was lint and merged, it will be pushing and release to [meshery.io](https://github.com/meshery/meshery.io) Github page by GitHub Action automatically.

### Artifact Repositories

Artifacts produced in the build processes are published and persisted in different public repositories and in different formats.

| Location      | Project       | Repository    |
| ------------- | ------------- | ------------- |
| Docker Hub    | Meshery       | [https://hub.docker.com/r/layer5/meshery](https://hub.docker.com/r/layer5/meshery) |
| GitHub        | mesheryctl    | [https://github.com/layer5io/meshery/releases](https://github.com/layer5io/meshery/releases) |
| Docker Hub    | Meshery Adapter for \<service-mesh\> | https://hub.docker.com/r/layer5/meshery-\<service-mesh\> |
| Docs          | Meshery Documentation | [https://docs.meshery.io](https://docs.meshery.io) |
| GitHub        | [Service Mesh Performance](https://smp-spec.io) | [https://github.com/layer5io/service-mesh-performance](https://github.com/layer5io/service-mesh-performance) |
| Github        | Helm charts   | [https://github.com/meshery/meshery.io/tree/master/charts](https://github.com/meshery/meshery.io/tree/master/charts) |

## Secrets

Some portions of the workflow require secrets to accomplish their tasks. These secrets are defined within the respective repositories and accessible to workflows during runtime. Currently defined secrets include:

- `DOCKER_USERNAME`: Username of the Docker Hub user with the right privileges to push images
- `DOCKER_PASSWORD`: Password for the Docker Hub user
- `GO_VERSION`: As of July 21st 2021 it is 1.16
- `IMAGE_NAME`: appropriate image name for each of the Docker container images. All are under the `layer5io` org.
- `SLACK_BOT_TOKEN`: Used for notification of new GitHub stars given to the Meshery repo.
- `CYPRESS_RECORD_KEY`: Used for integration with the Layer5 account on Cypress.
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
- Helm charts lint (helm)
- Helm charts release, tag and push(stefanprodan/helm-gh-pages@master)

## Tests for adapters
All Meshery adapters use a central workflow that is referenced in each of their test workflows which get triggered every time a pull request is made. These
tests in adapters are end-to-end tests and use patternfile. The reusable workflow is present in .github/workflows in Meshery repository with the workflow name "Test for Meshery adapters using patternfile"



### The pre-requisite of referencing this workflow is -
1. Using actions/upload-artifact@v2 a patternfile has to be uploaded as an artifact with the name as "patternfile".
2. The name of the uploaded patterfile should be passed in

---
      ...
      with:
          patternfile_name: < name of the patternfile which is uploaded >


3. Note: This Job is pre-run to the actual test. This is done in order to create patternfiles dynamically and use them. Therefore name of this jobs has to be passed as

---
      ...
      needs: 
         < Name of the pre-requisite job >


4. There should be an infinite token passed in: (Or else local provider will be used)

---
      ...
      secrets:
        token: ${{ secrets.PROVIDER_TOKEN }}


### The central workflow functionally does -
1. Checks out the code of the repository(on the ref of latest commit of branch which made the PR) in which it is referenced.
2. Starts a minikube cluster
3. Builds a docker image of the adapter and sets minikube to use docker's registry.
4. Starts the adapter and meshery server (The url to deployment and service yaml of adapter are configurable).
 NOTE: The service mesh name( whose adapter we are testing ) has to passed in:

 ---
      ...
      with:
         adapter_name: < NAME OF THE SERVICE MESH >

5. The uploaded patternfile is deployed.
6. Workflow sleeps for some time.
7. Then the assertion is made that the pods passed in-

---
      ...
      with:
         expected_pods: < pod1,pod2,pod3 >  #comma separated pod names that will be expected to be present after patternfile is deployed
8. And these pods are present in their respective namespaces passed in-

--- 
      ...
      with:
         expected_pods_namespaces: < pod1ns, pod2ns , pod3ns >


### Expected inputs of this workflow 

---
    inputs:
      expected_pods:
        required: true
        type: string
      expected_pods_namespaces:
        required: true
        type: string
      service_url:
        required: true
        type: string
      deployment_url:
        required: true
        type: string
      adapter_name:
        required: true
        type: string   
      patternfile_name:
        required: true
        type: string   
      provider: 
        required: false
        type: string   
    secrets:
      token:

### Expected outputs of this workflow
 The pods passed in “expected_pods” are running in the subsequent namespaces passed in “expected_pods_namespaces”. If not, the workflow fails
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

## Helm Charts Lint Check, Build, and Release

The charts lint check, charts build, and charts release workflows are all triggered by GitHub events. Sometimes this event is the opening, updating, or merging of a branch, or sometimes a manual invocation, or a GitHub release event.
### Check Helm Charts

Every PR which includes changes to the files under `install/kubernetes/` directory in the `meshery/meshery` will trigger a Github Action to check for any mistakes in Helm charts using the `helm lint` command.

### Release Helm Charts to Github and Artifact Hub

New Meshery Helm charts are published upon trigger of a release event in the `meshery/meshery` repo. New versions of Meshery's Helm charts are published to [Meshery's Helm charts release page](https://github.com/meshery/meshery.io/tree/master/charts). [Artifact Hub] (https://artifacthub.io/packages/helm/meshery/meshery) syncs with these updated Meshery Helm charts.

## Release Versioning

Meshery and its components follow the commonly used, semantic versioning for its release numbering scheme. Given a version number MAJOR.MINOR.PATCH.BUILD, increment the:

- MAJOR version - major changes with rare potential for incompatible API changes.
- MINOR version - add functionality in a backwards-compatible manner.
- PATCH version - mostly for bug and security fixes.
- AlPHA/BETA/RC - used to facilitate early testing of an upcoming release.

### Component Versioning

Meshery comprises a number of components including a server, adapters, UI, and CLI. As an application, Meshery is a composition of these different functional components. While all of Meshery’s components generally deploy as a collective unit (together), each component is versioned independently, so as to allow them to be loosely coupled and iterate on functionality independently.  Some of the components must be upgraded simultaneously, while others may be upgraded independently. See [Upgrading Meshery](/guide/upgrade) for more information.

GitHub release tags will contain a semantic version number. Semantic version numbers will have to be managed manually by tagging a relevant commit in the master branch with a semantic version number (example: v1.2.3). 

## Release Process

Documentation of Meshery releases contains a table of releases and release notes and should be updated with each release.

### Automated Releases

Releases are manually triggered by a member of the release team publishing a release. Release names and release tags need to be assigned by the publishing team member. GitHub Action workflows will trigger and take care of running the required steps and publishing all artifacts (e.g., binary and docker images).

### Workflow Triggers

The following events will trigger one or more workflows:

1. Tagged Release
1. Commit pushed to the master branch
1. PR opened or commit pushed to PR branch
1. PR merged to the master branch

### Release Notes

While use of GitHub Actions facilitates automated builds, ReleaseDrafter is helping with facilitating automated release notes and versioning. 
### Generating Release Notes

ReleaseDrafter generates a GitHub tag and release draft. ReleaseDrafter action will trigger and will automatically draft release notes according to the configuration set-up. ReleaseDrafter drafts releases as soon as a commit is made into master after the previous release. The GitHub Action, ReleaseDrafter, is compatible with semantic releases and is used to auto-increment the semantic version number by looking at the previous release version. 

#### Automated Release Notes Publishing

The publishing of release notes to Meshery Docs is automated. Triggered by a release event, a workflow will checkout the Meshery repo, copy the auto-drafted release notes into a Jekyll collection in Meshery Docs, and generate a pull request.

#### Automated Release Notes Sending

The sending of release notes is now automated as a step in the stable release channel workflow. The release notes are automatically sent to the [developers@meshery.io mailing list](https://groups.google.com/a/meshery.io/g/developers).

#### Automated Pull Request Labeler

A GitHub Issue labeler bot is configured to automatically assign labels to issues based on which files have changed in which directories. For example, a pull request with changes to files in the “/docs/**” folder will receive the “area/docs” label. Presence of the “area/docs” label is used to trigger documentation builds and Netlify builds of the Meshery Docs. Similar labels are assigned and used to trigger workflows or used as conditional flags in workflows to determine which workflows or which steps in a workflows to run. 

## Release Channels

Artifacts of the builds for Meshery and its components are published under two different release channels, so that improved controls may be provided to both Meshery users and Meshery developers. The two release channels are *edge* and *stable* release channels.

Relative to stable releases, edge releases occur much more frequently. Edge releases are made with each merge to master, unless that merge to master is for a stable release. Stable releases are made with each merge to master when a GitHub release tag is also present in the workflow.

### Stable Channel

The following is an example of the release channels and the docker tags used to differentiate them. The latest tag will be applied only to images in the stable release channel. Here are two releases with two different images.

**Latest Stable Image**

- layer5/meshery:stable-latest
- layer5/meshery:stable-v0.4.1
- layer5/meshery:stable-324vdgb (sha)

**Older Stable Image**

- layer5/meshery:stable-v0.4.0
- layer5/meshery:stable-289d02 (sha)

Every docker image built receives either the edge tags or the stable tags. Which set of image tags assigned is determined by whether a release tag is present or not. In other words, stable channel docker images get the “stable” tags only in the presence of a release tag (e.g. v0.4.1).

### Edge Channel

The Edge release channel provides early access to the latest features, allowing experimentation and feedback. Recent updates introduce dynamic versioning for Edge releases and improve compatibility with remote providers, ensuring a seamless experience. This makes the Edge channel valuable for users seeking cutting-edge features while contributing to the Meshery project. Docker image tags remain consistent for differentiation.

Stable and edge releases are both published to the same Docker Hub repository. Docker Hub repositories differentiate release channels by image tag. The following Docker images tagging convention is followed:

**Latest Edge Image**

- layer5/meshery:edge-latest
- layer5/meshery:edge-289d02 (sha)

**Older Edge Image**

- layer5/meshery:edge-324vdgb (sha)


### Switching Between Meshery Release Channels

Users are empowered to switch between release channels at their leisure.

#### Switching Release Channels Using mesheryctl

Users can use mesheryctl to switch between release channels, e.g. `mesheryctl system channel [stable|edge]`.  Alternatively, users can manually switch between channels by updating the docker image tags in their meshery.yaml / Kubernetes manifest files. This command generates a meshery.yml (a docker-compose file) with release channel-appropriate tags for the different Docker container images.

#### Viewing Release Channel and Version Information in Meshery UI

Users are shown their Meshery deployment’s release channel subscription enient new setting in the Preferences area of the Meshery UI, so that people can alternatively use the UI to switch between channels if they like. Version numbers for Meshery adapters are also shown in the UI.

## Release Cadence

Minor releases of the Meshery project are release frequently (on a monthly basis on average) with patch releases made on-demand in-between those times. The project does not have long term releases that are sustained with bug fixes, yet. Bug fixes and patches will be released as needed on the latest release version.

### Release Support

General community support and commercial support from Layer5 is available. Separately, third parties and partners may offer longer-term support solutions.

#### Pre v1.0

Project focuses on functionality, quality and adoption, while retaining the flexibility for shifts in architecture. 

#### Post v1.0

Once a 1.0 release has been made, Around once a month or so, the project maintainers will take one of these daily builds and run it through a number of additional qualification tests and tag the build as a Stable release. Around once a quarter or so, the project maintainers take one of these Stable releases, run through a bunch more tests and tag the build as a Long Term Support (LTS) release. Finally, if we find something wrong with an LTS release, we issue patches.

The different types (Daily, Stable, LTS) represent different product quality levels and different levels of support from the Meshery team. In this context, support means that we will produce patch releases for critical issues and offer technical assistance. 

## Versioning Documentation 
### For new major release

The structure which the docs follow right now is, The main `docs` folder has the most recent version of documentation, while there are sub-folders for previous versions, v0.x (x being the last major release).
On release of a new major version, the static html files for the most recent version is generated and is renamed as the release version (v0.x). 

##### Steps:
After cloning the Meshery repository
1. `cd docs` > `bundle install` > `make docs`
1. On executing `make docs` a `_site` folder is created which has static html files. 
1. The `_site` folder is renamed to `v0.x`. 
1. This `v0.x` folder is now the latest version of docs. 

##### _In the `v0.x` folder_
1. Search and replace all the instances where there is a direct path is defined to include the version name in the path, i.e, all paths to intra-page links and images should start with `/v0.x/`.
- Look for `href="/` and replace with `href="/0.x/`
- Look for `src="/`and replace with `src="/0.x/` <br/><br/>
<a href="{{ site.baseurl }}/assets/img/versioning-guide/search-and-replace.png">
  <img src="{{ site.baseurl }}/assets/img/versioning-guide/search-and-replace.png" />
</a>


### For old release

For older releases we have to travel back in time. Using the `Tags` in github we go to a previous release, `v0.X.x`, the `.x` here should be the latest version of the archived docs. 

##### Steps: 
1. Copy the commit ID for that release. <br/><br/>
<a href="{{ site.baseurl }}/assets/img/versioning-guide/commit-ID.png">
  <img src="{{ site.baseurl }}/assets/img/versioning-guide/commit-ID.png" />
</a>

1. `git checkout <commit ID>` > `cd docs` > `bundle install` > `make docs`
1.  On executing `make docs` a `_site` folder is created which has static html files. 
1.  The `_site` folder is renamed to `v0.X` and is copied into the `docs` folder of the present version. 

## Bi-Weekly Meetings

If you are passionate about CI/CD pipelines, DevOps, automated testing, managing deployments, or if you want to learn how to use Meshery and its features, you are invited to join the bi-weekly Build and Release meetings. Find meeting details and agenda in the [community calendar](https://meshery.io/calendar) and the [meeting minutes document](https://docs.google.com/document/d/1GrVdGHZAYeu6wHNLLoiaKNqBtk7enXE9XeDRCvdA4bY/edit#). The meetings are open to everyone and recorded for later viewing. We hope to see you there!

These [steps]({{site.baseurl}}/project/build-and-release#in-the-v0x-folder) for replacing all the instances of direct path are to be followed. 

