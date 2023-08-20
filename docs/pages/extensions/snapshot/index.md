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

### Supported Events

The MeshMap Snapshot Action supports all of the events listed in the GitHub documentation: 
For detailed information about each event, including its properties and payloads, refer to the [GitHub Actions documentation](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows).

### Workflow Syntax for Event Filtering

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

Specific events that are relevant to the MeshMap Snapshot Action are:

- **push:** This event is triggered when code is pushed to the repository. It is the most common event used to trigger workflows.
- **pull_request:** This event is triggered when a pull request is opened or updated. It is often used to trigger workflows that run tests or linters on the pull request's code. 

- **workflow_dispatch:** This event is triggered when a workflow is manually triggered. It can be used to trigger workflows on demand, such as for publishing a new release or addressing critical bugs.

For a comprehensive list of events that can be used in GitHub Actions, please refer to the Supported Events section above.

## What Happens to Workflow Customizations on Upgrade?

Customizations to the trigger criteria for the MeshMap Snapshot actions are preserved when upgrading to a new version of the action. However, there may be some cases where customizations are lost, such as when the syntax for specifying the trigger criteria changes in a new version of the action.

Here are some examples of cases where customizations may be lost:

- You currently have a workflow that is triggered on the push event, and the syntax for specifying the push event changes in a new version of the action.

- You have a workflow that is triggered on a custom event, and the custom event is no longer supported in a new version of the action.

It is always a good practice to test your workflows after upgrading to a new version of the MeshMap Snapshot Action to make sure that your customizations are still working as expected.

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




