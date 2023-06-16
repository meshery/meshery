---
layout: default
title: MeshMap Snapshot
permalink: extensions/snapshot
language: en
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

## Usage:

After testing you can [create a v1 tag](https://github.com/actions/toolkit/blob/master/docs/action-versioning.md) to reference the stable and latest V1 action

