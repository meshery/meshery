---
layout: default
title: MeshMap Snapshot
permalink: extensions/snapshot
language: en
abstract: Screenshot service provided via Meshmap to capture a snapshot of their infrastructure at any given time.
display-title: "false"
list: include
type: extensions
---

# <img style="height: 4rem; width: 4rem;" src="{{site.baseurl}}/assets/img/meshmap-icon-color.svg" />  MeshMap Snapshot

Walks in application and takes a shot of your infrastructure using Meshery Extension MeshMap

MeshMap Snapshot is a screenshot service provided via MeshMap for the design or Application user is interested in. It Enables users to visualize the changes being done in the code-base rapidly over each PR and inform the user about any potential changesin their infrastructure. It doesn't need any configuration or setup neither any deployment by the client rather a simple one time setup is able to provide a long time value.

### Guiding Principles

**Step 1: The moving banner board to MeshMap**


The screenshots would be flowing through each PR giving product the eye-balls it needs


**Step 2: Understanding the designs and infrastructure via visual design without downloading or setting up the MeshMap explicitly**


You don't have to explicitly install Meshery or MeshMap on your cluster to see their infrastructure or get benefits of the set of capabilities of MeshMap.


**Step 3: Introducing MeshMap along with Meshery in the GitOps field.**

This sets the foundation of MeshMap in the github lifecycle, PRs, comments and all for future.

## Functional Sequence Diagram
For Github Workflows:
<img src="{{site.baseurl}}/assets/img/meshmap/meshmap-snapshot.png" />

## Example Usage:

### Add Meshmap Snapshot: Github Pull Request
Connect MeshMap to your GitHub repo and see changes pull request-to-pull request. Get snapshots of your infrastructure directly in your PRs.

**Step 1: Log in to your Meshery Dashboard.**

**Step 2: Navigate to Extensions and Enable GitHub Action: MeshMap Snapshot.**

**Step 3: You will be directed to Meshery Cloud.**

**Step 4: Click on the "Let's Go" button in the onboarding modal.**

**Step 5: Now, select the second option, 'GitOps your infrastructure with MeshMap Snapshot'.**

**Step 6: Once that's complete, follow the guided instructions.**

**Steps 7: Voila, you're all set!**



### When Infrastructure is located in the file-system
```yaml
name: 'MeshMap Snapshot With File-located in Fs'
on: # rebuild any PRs and main branch changes
  pull_request:
  push:
    branches:
      - main
      - master
      - 'releases/*'

jobs:
  test: # make sure the action works on a clean machine without building
    runs-on: ubuntu-latest
    steps:
      - name: Set PR number # To comment the final status on the Pull-request opened in any repository
        run: |
          export pull_number=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")
          echo "PULL_NO=$pull_number" >> $GITHUB_ENV
      - uses: actions/checkout@v3 # the repository where your infrastructure is located
      - uses: actions/checkout@v3 #this step would go and would be no longer needed to be written
        with:
          path: action
          repository: layer5labs/meshmap-snapshot
      - id: test_result
        uses: layer5labs/MeshMap-Snapshot@v0.0.4
        with:
          githubToken: ${{ secrets.GITHUB_TOKEN }} # github's personal access token example: "ghp_...."
          providerToken: ${{ secrets.PROVIDER_TOKEN }} # Meshery Cloud Authentication token, signin to meshery-cloud to get one, example: ey.....
          prNumber: ${{ env.PULL_NO }} # auto-filled from the above step
          application_type: "Kubernetes Manifest" # your application type, could be any of three: "Kubernetes Manifest", "Docker Compose", "Helm Chart"
          filePath: "action/__tests__/manifest-test" # relative file-path from the root directory in the github-runner env, you might require to checkout the repository as described in step 2
```
### When Infrastructure is identified via URL
```yaml
name: 'MeshMap Snapshot With URL-Upload'
on: # rebuild any PRs and main branch changes
  pull_request:
  push:
    branches:
      - main
      - master
      - 'releases/*'

jobs:
  test: # make sure the action works on a clean machine without building
    runs-on: ubuntu-latest
    steps:
      - name: Set PR number # To comment the final status on the Pull-request opened in any repository
        run: |
          export pull_number=$(jq --raw-output .pull_request.number "$GITHUB_EVENT_PATH")
          echo "PULL_NO=$pull_number" >> $GITHUB_ENV
      - uses: actions/checkout@v3 #this step would go and would be no longer needed to be written
        with:
          path: action
          repository: layer5labs/meshmap-snapshot
      - id: test_result
        uses: layer5labs/MeshMap-Snapshot@v0.0.4
        with:
          githubToken: ${{ secrets.GITHUB_TOKEN }} # github's personal access token example: "ghp_...."
          providerToken: ${{ secrets.PROVIDER_TOKEN }} # Meshery Cloud Authentication token, signin to meshery-cloud to get one, example: ey.....
          prNumber: ${{ env.PULL_NO }} # auto-filled from the above step
          application_type: "Helm Chart" # your application type, could be any of three: "Kubernetes Manifest", "Docker Compose", "Helm Chart"
          application_url: "https://github.com/meshery/meshery.io/raw/master/charts/meshery-v0.6.88.tgz"
```

#### FileSystem Approach Notes
The filesystem-approach asks for your relative file-path and automatically merges all the yaml files together to bundle up into one. So you might like to give the root directory where all the yamls are located. It doesn't move recursevely in internal folders, so only the first level yaml files are checked.

## List of Input variables supported:
```yaml
designId:  # id of input  #deprecated
  description: "The design uuid, example: 3c116d0a-49ea-4294-addc-d9ab34210662"
  required: false
  default: '{}'
applicationId:  #deprecated
  description: "The application uuid, example: 3c116d0a-49ea-4294-addc-d9ab34210662"
  required: false
githubToken:
  description: "Github PAT token"
  required: true
providerToken:
  description: "Meshery Authentication Provider Token"
  required: true
cypressRecordKey:
  description: "cypress record key"
  required: false
prNumber:
  description: "The Pull request on which comment has to be made"
  required: false
  default: 0
filePath: 
  description: "The relative filepath of the location where the manifests are stored"
  required: false
application_type:
  description: "Application upload type, any of the three, Kubernetes Manifest, Docker Compose, Helm Chart"
  required: true
application_url:
  description: "Application's source url where the manifests or data is stored"
  required: false
```

## Customizing Snapshot Workflow Triggers in MeshMap Snapshot
You can configure your workflows to run when specific activity on GitHub happens, at a scheduled time, or when an event outside of GitHub occurs.

### About events that trigger workflows
GitHub Actions provides a variety of events that can trigger workflows, allowing you to automate your software development process. Each event corresponds to a specific activity, such as creating a pull request, pushing code to a repository, or releasing a new version.

## Supported Events

GitHub Actions supports the following events:

- [`push`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#push-event): Triggers when code is pushed to the repository.

- [`pull_request`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request-event): Triggers when a pull request is opened or updated.

- [`pull_request_target`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request_target-event): Triggers when a pull request is opened against the repository's default branch.

- [`workflow_run`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_run-event): Triggers when a workflow is executed.

- [`check_run`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#check_run-event): Triggers when a check run is requested or completed.

- [`check_suite`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#check_suite-event): Triggers when a check suite is requested or completed.

- [`create`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#create-event): Triggers when a branch or tag is created.

- [`delete`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#delete-event): Triggers when a branch or tag is deleted.

- [`deployment`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#deployment-event): Triggers when a deployment is created or updated.

- [`deployment_status`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#deployment_status-event): Triggers when a deployment's status is updated.

- [`fork`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#fork-event): Triggers when a fork is created.

- [`gollum`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#gollum-event): Triggers when a Wiki page is created or updated.

- [`issue_comment`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#issue_comment-event): Triggers when an issue comment is created or edited.

- [`issues`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#issues-event): Triggers when an issue is opened, edited, or deleted.

- [`label`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#label-event): Triggers when a label is created, edited, or deleted.

- [`milestone`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#milestone-event): Triggers when a milestone is created, edited, or deleted.

- [`page_build`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#page_build-event): Triggers when a GitHub Pages site is built.

- [`project_card`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#project_card-event): Triggers when a project card is created, edited, moved, or deleted.

- [`project_column`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#project_column-event): Triggers when a project column is created, edited, moved, or deleted.

- [`project`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#project-event): Triggers when a project is created, edited, closed, reopened, or deleted.

- [`public`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#public-event): Triggers when a private repository is made public.

- [`pull_request_review`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request_review-event): Triggers when a pull request review is submitted or dismissed.

- [`pull_request_review_comment`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request_review_comment-event): Triggers when a comment on a pull request review is created, edited, or deleted.

- [`registry_package`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#registry_package-event): Triggers when a package is published, updated, or deleted in a GitHub Packages registry.

- [`release`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#release-event): Triggers when a release is published or updated.

- [`status`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#status-event): Triggers when a commit status is created, updated, or deleted.

- [`watch`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#watch-event): Triggers when a user "watches" a repository.

- [`workflow_dispatch`](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#workflow_dispatch-event): Triggers when a workflow is manually triggered using the GitHub UI or API.


For detailed information about each event, including its properties and payloads, refer to the [GitHub Actions documentation](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows).

## Workflow Syntax for Event Filtering

When defining workflows, you can use the `on` keyword to specify which events trigger the workflow. You can further filter the triggering conditions by using the `types`, `branches`, `tags`, and other options. For example:

```yaml
on:
  push:
    branches:
      - main
  pull_request:
    types:
      - opened
      - synchronize
```

## Usage:

After testing you can [create a v1 tag](https://github.com/actions/toolkit/blob/master/docs/action-versioning.md) to reference the stable and latest V1 action

# General Upgrade Guide
[MeshMap Snapshot Release Page](https://github.com/layer5labs/meshmap-snapshot/releases)

```
 - id: test_result
        uses: layer5labs/MeshMap-Snapshot@v0.0.5 # <-- Update the version to latest from the MeshMap-Snapshot release page
        with:
        ...
```

## Upgrade/Migrate Guide - For Meshery
1. Given changes done in `action.yml` in MeshMap Snapshot, updating the workflows is required.
2. Given changes done other than in `action.yml` in MeshMap Snapshot, the update in the `.github/worflows` is not a hard requirement, but doesnt hurt.



