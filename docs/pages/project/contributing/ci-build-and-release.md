---
layout: page
title: Build & Release (CI)
permalink: project/contributing/build-and-release
abstract: Details of Meshery's build and release strategy.
language: en
type: project
category: contributing
list: include
training-video: "https://www.youtube.com/watch?v=dlr_nzJV16Q"
---

Meshery’s build and release system incorporates many tools, organized into different workflows each triggered by different events. Meshery’s build and release system does not run on a schedule, but is event-driven. GitHub Actions are used to define Meshery’s CI workflows. New builds of Meshery and its various components are automatically generated upon push, release, and other similar events, typically in relation to their respective master branches.

{% include alert.html type="info" title="Meshery Test Documents" content=" The Meshery Test Strategy (<a href='https://docs.google.com/document/d/11nAxYtz2SUusCYZ0JeNRrOLIxkgmmbUVWz63MBZV2oE/edit'>1</a>,<a href='https://docs.google.com/document/d/14vbwnKafqxrr-cJOmLWEvvj75MFrxFKiUdfPwZRhg64/edit'>2</a>) is a comprehensive document that outlines the approach to testing, types of tests, frameworks used - test strategy - for each of Meshery's architectural components. The <a href='https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit'>Meshery Test Plan</a> is a tactical spreadsheet that *itemizes* specific integration and end-to-end test cases and tracks their status. It includes a list of GitHub workflows and their purpose. Join the <a href='https://docs.google.com/document/d/1GrVdGHZAYeu6wHNLLoiaKNqBtk7enXE9XeDRCvdA4bY/edit#'>Meshery CI meeting</a> and/or the <a href='https://meshery.io/subscribe'>developers mailing list</a> to get involved." %}

## Artifacts

Today, Meshery and Meshery adapters are released as Docker container images, available on Docker Hub. Meshery adapters are out-of-process adapters (meaning not compiled into the main Meshery binary), and as such, are independent build artifacts and Helm charts. The Docker images are created and tagged with the git commit SHA, then pushed to Docker Hub automatically using GitHub Actions. Subsequently, when contributions containing content for the Helm charts of Meshery and Meshery Adapter are linted and merged, they will be pushed and released to [meshery.io](https://github.com/meshery/meshery.io) Github page by GitHub Action automatically.

### Artifact Repositories

Artifacts produced in the build processes are published and persisted in different public repositories and in different formats.

| Location   | Project                                         | Repository                                                                                                           |
| ---------- | ----------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Docker Hub | Meshery                                         | [https://hub.docker.com/r/meshery/meshery](https://hub.docker.com/r/meshery/meshery)                                   |
| GitHub     | mesheryctl                                      | [https://github.com/meshery/meshery/releases](https://github.com/meshery/meshery/releases)                         |
| Docker Hub | Meshery Adapter for \<adapter-name\>            | https://hub.docker.com/r/meshery/meshery-\<adapter-name>\>                                                             |
| Docs       | Meshery Documentation                           | [https://docs.meshery.io](https://docs.meshery.io)                                                                   |
| GitHub     | [Cloud Native Performance](https://smp-spec.io) | [https://github.com/service-mesh-performance](https://github.com/service-mesh-performance/service-mesh-performance)         |
| Github     | Helm charts                                     | [https://github.com/meshery/meshery.io/tree/master/charts](https://github.com/meshery/meshery.io/tree/master/charts) |

## Secrets

Some portions of the workflow require secrets to accomplish their tasks. These secrets are defined within the respective repositories and accessible to workflows during runtime. Currently defined secrets include:

- `DOCKER_USERNAME`: Username of the Docker Hub user with the right privileges to push images
- `DOCKER_PASSWORD`: Password for the Docker Hub user
- `GO_VERSION`: As of December, 2024 is 1.23
- `IMAGE_NAME`: appropriate image name for each of the Docker container images. All are under the `layer5io` org.
- `SLACK_BOT_TOKEN`: Used for notification of new GitHub stars given to the Meshery repo.
- `GLOBAL_TOKEN`: Used for securely transmitting performance test results for the None Provider.
- `NPM_TOKEN`: npm authentication token, used to perform authentication against the npm registry in meshery deployment workflow.
- `GH_ACCESS_TOKEN`: GitHub access token for various operations
- `INTEGRATION_SPREADSHEET_CRED`: Credentials for integration spreadsheet access
- `MAIL_PASSWORD`: Password for email notifications
- `MAIL_USERNAME`: Username for email notifications
- `MESHERY_PROVIDER_TOKEN`: Token for Meshery provider authentication
- `MESHERY_TOKEN`: General Meshery authentication token
- `METAL_AUTH_TOKEN`: Authentication token for metal provider
- `METAL_SERVER1`: Configuration for metal server 1
- `METAL_SERVER2`: Configuration for metal server 2
- `METAL_SERVER3`: Configuration for metal server 3
- `NETLIFY_AUTH_TOKEN`: Authentication token for Netlify
- `NETLIFY_SITE_ID`: Site ID for Netlify deployments
- `PLAYGROUND_CONFIG`: Configuration for playground environments
- `PROVIDER_TOKEN`: General provider authentication token
- `RELEASEDRAFTER_PAT`: Personal access token for Release Drafter
- `RELEASEDRAFTER_PAT`: Personal access token for release notes generation
- `REMOTE_PROVIDER_USER_EMAIL`: Email used for authentication in Playwright tests
- `REMOTE_PROVIDER_USER_PASS`: Password used for authentication in Playwright tests

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
tests in adapters are end-to-end tests and use patternfile. The reusable workflow is present in `.github/workflows` in Meshery repository with the workflow name "Test for Meshery adapters using patternfile"

### Pre-requisite for referencing this workflow

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

### Functionality of Central Workflow

1. Checks out the code of the repository(on the ref of latest commit of branch which made the PR) in which it is referenced.
2. Starts a minikube cluster
3. Builds a docker image of the adapter and sets minikube to use docker's registry.
4. Starts the adapter and Meshery Server (The url to deployment and service yaml of adapter are configurable).
   NOTE: The adapter name has to passed in:

---

      ...
      with:
         adapter_name: < NAME OF THE ADAPTER >

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

### `mesheryctl`

As a special case, the `meshery` repository contains an additional artifact produced during each build. This artifact is `mesheryctl`, which is built as an executable binary. In order to make the job of building `mesheryctl` easier for a combination of different platform architectures and operating systems, we are using [GoReleaser](https://goreleaser.com). Irrespective of branch, for every git commit and git push to the `meshery` repository, GoReleaser will execute and generate the OS and arch-specific binaries (but will NOT persist these artifacts in GitHub). Even though mesheryctl binaries are built each time a pull request is merged to master, only stable channel artifacts are published (persisted).

#### Releasing `mesheryctl` to GitHub

Only when a git tag containing a semantic version number is present (is a commit in the master branch) will GoReleaser execute, generate the archives, and also publish the archives to [Meshery’s GitHub releases](https://github.com/meshery/meshery/releases) automatically. GoReleaser is configured to generate artifacts for the following OS, ARCH combination:

- Darwin - i386, x86_64
- Linux - i386, x86_64
- Windows - i386, x86_64
- FreeBSD - i386, x86_64

The artifacts will be made available as a tar.gz archive for all the operating systems. mesheryctl is bundled into packages for commonly used package managers: homebrew and scoop.

##### Homebrew

GoReleaser facilitates the creation of a brew formula for mesheryctl. The [homebrew-tap](https://github.com/meshery/homebrew-tap) repository is the location for `mesheryctl`'s brew formulas. Releases of mesheryctl are  published in the official homebrew-core tap at https://github.com/Homebrew/homebrew-core/pkgs/container/core%2Fmesheryctl.

##### Scoop

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

Meshery comprises a number of components including a server, adapters, UI, and CLI. As an application, Meshery is a composition of these different functional components. While all of Meshery’s components generally deploy as a collective unit (together), each component is versioned independently, so as to allow them to be loosely coupled and iterate on functionality independently. Some of the components must be upgraded simultaneously, while others may be upgraded independently. See [Upgrading Meshery](/installation/upgrades) for more information.

GitHub release tags will contain a semantic version number. Semantic version numbers will have to be managed manually by tagging a relevant commit in the master branch with a semantic version number (example: v1.2.3).

## Release Process

Documentation of Meshery releases contains a table of releases and release notes and should be updated with each release.

# Meshery Release Documents and Release Lead Responsibilities

## Meshery Release Documents

Meshery uses several types of release documents to standardize the purpose, style, and structure of release information. These documents are crucial for maintaining clear communication about changes, updates, and new features in Meshery releases.

### Types of Release Documents

- **Changelogs**
  - Comprehensive list of all changes since the prior release
  - Generated automatically by tools like ReleaseDrafter
  - Contains detailed, technical information

- **Release Notes**
  - Curated, bulleted list of highlights
  - Summarized and categorized in human-readable format
  - Includes some engineering terminology and issue references
  - Based on ReleaseDrafter output, but human-summarized and refined

- **Release Announcement**
  - Human-written summary highlighting significant items
  - Includes caveats (e.g., incompatibility on upgrade)
  - Provides links to other information sources (upgrade guide, feature blogs, full bug fix list)
  - Contains graphics and links to in-depth documentation
  - Distributed via #announcements Slack channel and public mailing list (for stable releases)

- **User and Upgrade Docs**
  - How-to guides for new features
  - Updates to user-facing documentation
  - Includes Upgrade Guide updates for version-to-version considerations

- **Feature-functionality Blogs**
  - In-depth reviews of new or significantly augmented functionality
  - Explains the what, why, and how of new features

## Meshery Release Lead Responsibilities

The Meshery Release Lead plays a crucial role in coordinating and executing the release process. Their responsibilities span approximately 5 months per release cycle.

### Pre-Release Phase (1 month)

- **Release Planning**
  - Schedule and organize release planning meetings
  - Define and communicate release timelines
  - Coordinate with development, testing, and documentation teams

- **Feature Management**
  - Oversee feature implementation and prioritization
  - Ensure all planned features are completed or properly deferred

- **Quality Assurance**
  - Coordinate with QA team to ensure thorough testing
  - Address and prioritize bug fixes

- **Documentation Preparation**
  - Ensure all new features and changes are properly documented
  - Oversee the creation and updating of release notes

- **Release Build Preparation**
  - Manage the creation of release branches
  - Oversee the release build process
  - Coordinate with DevOps for deployment preparations

### Active Maintenance Phase (6 months)

- **Ongoing Releases**
  - Manage minor releases every 3 weeks
  - Coordinate vulnerability fix integrations
  - Oversee cherry-pick decisions for backports

- **Monitoring and Issue Resolution**
  - Monitor for release blockers
  - Coordinate resolution of critical issues

- **Communication**
  - Lead weekly Meshery build and release meetings
  - Provide regular status updates to the community

- **Documentation Updates**
  - Ensure documentation remains current throughout the maintenance phase
  - Oversee updates to user guides and upgrade instructions

- **End-of-Life Procedures**
  - Manage the process for deprecating the release
  - Ensure proper handover to the next release cycle

### Release Manager Qualifications and Selection

- Must be an active community member for at least 3 months
- Approved by majority vote of current maintainers
- At least one release manager must meet requirements for vulnerability-related access

To volunteer as a Meshery Release Lead, interested individuals should contact a current maintainer or self-nominate.

This structured approach to release documentation and management ensures consistency, clarity, and efficiency in Meshery's release process, facilitating better communication with users and smoother project development.

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

A GitHub Issue labeler bot is configured to automatically assign labels to issues based on which files have changed in which directories. For example, a pull request with changes to files in the “/docs/\*\*” folder will receive the “area/docs” label. Presence of the “area/docs” label is used to trigger documentation builds and Netlify builds of the Meshery Docs. Similar labels are assigned and used to trigger workflows or used as conditional flags in workflows to determine which workflows or which steps in a workflows to run.

## Release Channels

Artifacts of the builds for Meshery and its components are published under two different release channels, so that improved controls may be provided to both Meshery users and Meshery developers. The two release channels are _edge_ and _stable_ release channels.

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

Users can use mesheryctl to switch between release channels, e.g. `mesheryctl system channel [stable|edge]`. Alternatively, users can manually switch between channels by updating the docker image tags in their meshery.yaml / Kubernetes manifest files. This command generates a meshery.yml (a docker-compose file) with release channel-appropriate tags for the different Docker container images.

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
1. On executing `make docs` a `_site` folder is created which has static html files.
1. The `_site` folder is renamed to `v0.X` and is copied into the `docs` folder of the present version.

## Bi-Weekly Meetings

If you are passionate about CI/CD pipelines, DevOps, automated testing, managing deployments, or if you want to learn how to use Meshery and its features, you are invited to join the bi-weekly Build and Release meetings. Find meeting details and agenda in the [community calendar](https://meshery.io/calendar) and the [meeting minutes document](https://docs.google.com/document/d/1GrVdGHZAYeu6wHNLLoiaKNqBtk7enXE9XeDRCvdA4bY/edit#). The meetings are open to everyone and recorded for later viewing. We hope to see you there!


{% include training-video.html %}
