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

Meshery’s build and release system incorporates many tools, organized into different workflows each triggered by different events. Meshery’s build and release system does not run on a schedule, bu[...]

{% include alert.html type="info" title="Meshery Test Documents" content=" The Meshery Test Strategy (<a href='https://docs.google.com/document/d/11nAxYtz2SUusCYZ0JeNRrOLIxkgmmbUVWz63MBZV2oE/edit'>1</[...]

## Artifacts

Today, Meshery and Meshery adapters are released as Docker container images, available on Docker Hub. Meshery adapters are out-of-process adapters (meaning not compiled into the main Meshery binary), [...]

All repositories under the `github.com/meshery` and `github.com/meshery-extensions` organizations use immutable releases.
<img width="954" height="166" alt="immutable-releases-setting" src="https://github.com/user-attachments/assets/4435086f-db09-449e-a154-70979b8b01d1" />


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

Some portions of the workflow require secrets to accomplish their tasks. These secrets are defined within the respective repositories and accessible to workflows during runtime. Currently defined secr[...]

- `DOCKER_USERNAME`: Username of the Docker Hub user with the right privileges to push images
- `DOCKER_PASSWORD`: Password for the Docker Hub user
- `IMAGE_NAME`: appropriate image name for each of the Docker container images. All are under the `meshery` org.
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
- `NETLIFY_AUTH_TOKEN`: Authentication token for Netlify
- `NETLIFY_SITE_ID`: Site ID for Netlify deployments
- `PLAYGROUND_CONFIG`: Configuration for playground environments
- `PROVIDER_TOKEN`: General provider authentication token
- `RELEASEDRAFTER_PAT`: Personal access token for Release Drafter
- `RELEASEDRAFTER_PAT`: Personal access token for release notes generation
- `REMOTE_PROVIDER_USER_EMAIL`: Email used for authentication in Playwright tests
- `REMOTE_PROVIDER_USER_PASS`: Password used for authentication in Playwright tests

The Docker Hub user, `mesheryci`, belongs to the "ciusers" team in Docker Hub and acts as the service account under which these automated builds are being pushed. Every time a new Docker Hub repositor[...]

## Secure Coding Practices

Meshery's community follows a set of secure coding practices to reduce the likelihood of introducing vulnerabilities and to make remediation faster and more consistent when issues are discovered. These practices are applied by contributors, reviewers, and CI/CD workflows.

Key practices include:

- Security-first code review
  - Every pull request undergoes human review with attention to security-relevant changes (authentication/authorization, input parsing, network boundaries, use of secrets, cryptography).
  - Reviewers use a security checklist when reviewing changes that touch sensitive areas.
  - Changes that affect the security posture (e.g., bypassing validation, changing default configurations, adding new dependencies) must include rationale and mitigation steps.

- Automated static analysis and linters
  - CI runs static analysis tools such as staticcheck and go vet to detect common bugs and suspicious constructs.
  - gosec (or equivalent tools) is used to identify common security issues in Go code (e.g., SQL injection risk patterns, weak crypto usage, insecure file handling).
  - Frontend code is scanned with appropriate linters and security-oriented tooling.

- Dependency and supply-chain security
  - Dependabot and SCA tooling
    - Dependabot is enabled for this repository and integrated into CI and pre-release checks. Dependabot alerts and automated dependency update PRs are used to surface and remediate vulnerable dependencies; they are triaged and addressed according to severity and release impact.
    - See the repository's SECURITY-INSIGHTS.yml for the declared SCA tooling and integration details: https://github.com/meshery/meshery/blob/cfbf38de49c0863059fcbab93274b6de24cba4b3/SECURITY-INSIGHTS.yml
  - Dependabot, other SCA scanners, and GitHub's security alerts are monitored and acted upon. When introducing third-party packages, contributors must prefer well-maintained libraries and justify the addition in the PR description.
  - Meshery's Software Bill of Materials (SBOM) is produced as a build artifact and is available as part of the repository's BOM workflow; it is also referenced in SECURITY-INSIGHTS.yml.

- Secret management and least privilege
  - Secrets must never be committed to source code. GitHub secret scanning and repository-level secret scanning are used to detect accidental exposures.
  - CI secrets are stored in GitHub Secrets (or equivalent) and scoped with least privilege. Tokens used by workflows are restricted and rotated regularly.
  - Access to production signing keys and high-privilege tokens is limited to the minimum set of accounts and workflows necessary.

- Container and image security
  - Images are built from minimal, well-maintained base images.
  - Container images are scanned for known vulnerabilities before release (e.g., with Trivy or other image scanners).
  - Image tagging, immutable releases, and reproducible build practices help reduce supply-chain risks.
  - An SBOM (Software Bill of Materials) is created for releases where practical to improve traceability of included components.

- Secure defaults and configuration
  - New features and configurations should ship with secure defaults (e.g., secure cipher suites, reasonable timeouts, least-privilege access).
  - Configuration values that weaken security must be clearly documented and gated (feature flags or explicit opt-in).

- Runtime and integration security testing
  - Integration and E2E tests include assertions that critical services are not exposed unintentionally and that authentication/authorization flows work as expected.
  - Tests avoid storing credentials in plaintext and use test accounts with limited scope where possible.

- Incident handling and responsible disclosure
  - Vulnerabilities reported to the project are triaged promptly and handled according to priority and severity.
  - The Meshery community coordinates disclosure and remediation and communicates fixes and mitigation steps through release notes and advisories.
  - Contributors are encouraged to follow responsible disclosure channels (contacting maintainers or security@meshery.io where available).

- Developer training and documentation
  - Security best practices are documented for contributors and included in the project's contributing guidelines where appropriate.
  - New contributors are encouraged to familiarize themselves with the security checklist and the project's secure development practices.

These practices are enforced via a combination of automated CI checks and human review. If you are contributing a change that impacts security-sensitive areas, call this out in your PR description so reviewers and automated checks can give it the appropriate scrutiny.

## Checks and Tests

Meshery’s CI workflow incorporates several checks (partial list below) during merges and/or commits to any branches and pull requests to master branch to prevent broken code from being merged into m[...]

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
tests in adapters are end-to-end tests and use patternfile. The reusable workflow is present in `.github/workflows` in Meshery repository with the workflow name "Test for Meshery adapters using patter[...]

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

All Meshery GitHub repositories are configured with GitHub Actions. Everytime a pull request is submitted against the master branch of any repository, that repository’s GitHub Actions will be invoke[...]

1. trigger a Docker build to build a Docker container image
1. generate two Docker tags:
   1. a tag containing the git merge SHA
   1. a tag containing that particular release’s git tag (if one is present)
1. assign each of these two tags to the new container image as well as the latest tag.
1. push the new Docker tags and image to Docker Hub.

### `mesheryctl`

As a special case, the `meshery` repository contains an additional artifact produced during each build. This artifact is `mesheryctl`, which is built as an executable binary. In order to make the job [...]

#### Releasing `mesheryctl` to GitHub

Only when a git tag containing a semantic version number is present (is a commit in the master branch) will GoReleaser execute, generate the archives, and also publish the archives to [Meshery’s Git[...]

- Darwin - i386, x86_64
- Linux - i386, x86_64
- Windows - i386, x86_64
- FreeBSD - i386, x86_64

The artifacts will be made available as a tar.gz archive for all the operating systems. mesheryctl is bundled into packages for commonly used package managers: homebrew and scoop.

##### Homebrew

GoReleaser facilitates the creation of a brew formula for mesheryctl. The [homebrew-tap](https://github.com/meshery/homebrew-tap) repository is the location for `mesheryctl`'s brew formulas. Releases [...]

##### Scoop

GoReleaser facilitates the creation of a Scoop app for mesheryctl. The [scoop-bucket](https://github.com/layer5io/scoop-bucket) repository is the location of Meshery’s Scoop bucket.

## Helm Charts Lint Check, Build, and Release

The charts lint check, charts build, and charts release workflows are all triggered by GitHub events. Sometimes this event is the opening, updating, or merging of a branch, or sometimes a manual invoc[...]

### Check Helm Charts

Every PR which includes changes to the files under `install/kubernetes/` directory in the `meshery/meshery` will trigger a Github Action to check for any mistakes in Helm charts using the `helm lint` [...]

### Release Helm Charts to Github and Artifact Hub

New Meshery Helm charts are published upon trigger of a release event in the `meshery/meshery` repo. New versions of Meshery's Helm charts are published to [Meshery's Helm charts release page](https:/[...]

## Release Versioning

Meshery and its components follow the commonly used, semantic versioning for its release numbering scheme. Given a version number MAJOR.MINOR.PATCH.BUILD, increment the:

- MAJOR version - major changes with rare potential for incompatible API changes.
- MINOR version - add functionality in a backwards-compatible manner.
- PATCH version - mostly for bug and security fixes.
- AlPHA/BETA/RC - used to facilitate early testing of an upcoming release.

### Component Versioning

Meshery comprises a number of components including a server, adapters, UI, and CLI. As an application, Meshery is a composition of these different functional components. While all of Meshery’s compo[...]

GitHub release tags will contain a semantic version number. Semantic version numbers will have to be managed manually by tagging a relevant commit in the master branch with a semantic version number ([...]

## Release Process

Documentation of Meshery releases contains a table of releases and release notes and should be updated with each release.

# Meshery Release Documents and Release Lead Responsibilities

## Meshery Release Documents

Meshery uses several types of release documents to standardize the purpose, style, and structure of release information. These documents are crucial for maintaining clear communication about changes, [...]

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

This structured approach to release documentation and management ensures consistency, clarity, and efficiency in Meshery's release process, facilitating better communication with users and smoother pr[...]

### Automated Releases

Releases are manually triggered by a member of the release team publishing a release. Release names and release tags need to be assigned by the publishing team member. GitHub Action workflows will tri[...]

### Workflow Triggers

The following events will trigger one or more workflows:

1. Tagged Release
1. Commit pushed to the master branch
1. PR opened or commit pushed to PR branch
1. PR merged to the master branch

### Release Notes

While use of GitHub Actions facilitates automated builds, ReleaseDrafter is helping with facilitating automated release notes and versioning.

### Generating Release Notes

ReleaseDrafter generates a GitHub tag and release draft. ReleaseDrafter action will trigger and will automatically draft release notes according to the configuration set-up. ReleaseDrafter drafts rele[...]

#### Automated Release Notes Publishing

The publishing of release notes to Meshery Docs is automated. Triggered by a release event, a workflow will checkout the Meshery repo, copy the auto-drafted release notes into a Jekyll collection in M[...]

#### Automated Release Notes Sending

The sending of release notes is now automated as a step in the stable release channel workflow. The release notes are automatically sent to the [developers@meshery.io mailing list](https://groups.goog[...]

#### Automated Pull Request Labeler

A GitHub Issue labeler bot is configured to automatically assign labels to issues based on which files have changed in which directories. For example, a pull request with changes to files in the “/d[...]

## Release Channels

Artifacts of the builds for Meshery and its components are published under two different release channels, so that improved controls may be provided to both Meshery users and Meshery developers. The t[...]

Relative to stable releases, edge releases occur much more frequently. Edge releases are made with each merge to master, unless that merge to master is for a stable release. Stable releases are made w[...]

### Stable Channel

The following is an example of the release channels and the docker tags used to differentiate them. The latest tag will be applied only to images in the stable release channel. Here are two releases w[...]

**Latest Stable Image**

- meshery/meshery:stable-latest
- meshery/meshery:stable-v0.4.1
- meshery/meshery:stable-324vdgb (sha)

**Older Stable Image**

- meshery/meshery:stable-v0.4.0
- meshery/meshery:stable-289d02 (sha)

Every docker image built receives either the edge tags or the stable tags. Which set of image tags assigned is determined by whether a release tag is present or not. In other words, stable channel doc[...]

### Edge Channel

The Edge release channel provides early access to the latest features, allowing experimentation and feedback. Recent updates introduce dynamic versioning for Edge releases and improve compatibility wi[...]

Stable and edge releases are both published to the same Docker Hub repository. Docker Hub repositories differentiate release channels by image tag. The following Docker images tagging convention is fo[...]

**Latest Edge Image**

- meshery/meshery:edge-latest
- meshery/meshery:edge-289d02 (sha)

**Older Edge Image**

- meshery/meshery:edge-324vdgb (sha)

### Switching Between Meshery Release Channels

Users are empowered to switch between release channels at their leisure.

#### Switching Release Channels Using mesheryctl

Users can use mesheryctl to switch between release channels, e.g. `mesheryctl system channel [stable|edge]`. Alternatively, users can manually switch between channels by updating the docker image tags[...]

#### Viewing Release Channel and Version Information in Meshery UI

Users are shown their Meshery deployment’s release channel subscription enient new setting in the Preferences area of the Meshery UI, so that people can alternatively use the UI to switch between ch[...]

## Release Cadence

Minor releases of the Meshery project are release frequently (on a monthly basis on average) with patch releases made on-demand in-between those times. The project does not have long term releases tha[...]

### Release Support

General community support is available. Separately, third parties and partners may offer longer-term support solutions.

#### Pre v1.0

Project focuses on functionality, quality and adoption, while retaining the flexibility for shifts in architecture.

#### Post v1.0

Once a 1.0 release has been made, Around once a month or so, the project maintainers will take one of these daily builds and run it through a number of additional qualification tests and tag the build[...]

The different types (Daily, Stable, LTS) represent different product quality levels and different levels of support from the Meshery team. In this context, support means that we will produce patch rel[...]

## Versioning Documentation

### For new major release

The structure which the docs follow right now is, The main `docs` folder has the most recent version of documentation, while there are sub-folders for previous versions, v0.x (x being the last major r[...]
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

If you are passionate about CI/CD pipelines, DevOps, automated testing, managing deployments, or if you want to learn how to use Meshery and its features, you are invited to join the bi-weekly Build a[...]


{% include training-video.html %}

```
