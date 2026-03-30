# Kanvas Snapshot

Give your pipelines super powers with Kanvas Snapshot GitHub Actions.

## Kanvas Snapshot GitHub Action
1. See your deployment before you merge.
1. Connect Kanvas to your GitHub repo and see changes pull request-to-pull request.
1. Get snapshots of your infrastructure directly in your PRs.

## Overview

The **Kanvas Snapshot GitHub Action** allows you to generate and snapshot a visual diagram of your infrastructure as code whether that be Kubernetes manifests, Helm charts, or Docker compose files. This action automates the process of creating Kanvas Snapshots. 

### Features

1. **Snapshot Generation:** Create visual snapshots of infrastructure as code, complete with associated resources.
2. **Delivery via Email or PR Comment:** Snapshots can be delivered via email, image in your design catalog, or via comment in a GitHub pull request.
3. **Plugins:**
  1.  See [Helm Kanvas Snapshot](https://github.com/meshery/helm-kanvas-snapshot) plugin
  2.  See [Kubectl Kanvas Snapshot](https://github.com/meshery/kubectl-kanvas-snapshot) plugin

## Standard Installation and Use

See [Configuring Kanvas Snapshots](https://docs.layer5.io/cloud/tutorials/gitops-snapshots/) in Layer5 Docs.


<!--## Functional/Sequence Diagram
For Github Workflows:
![sequence-diag](.github/readme/images/meshmap-snapshot.png)-->

## Customized Use

<details>
  <summary><h4>Example: Importing design from filesystem</h4>
</summary>
<pre><code>
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
        uses: layer5labs/MeshMap-Snapshot@v0.2.6
        with:
          githubToken: ${{ secrets.GITHUB_TOKEN }} # github's personal access token example: "ghp_...."
          mesheryToken: ${{ secrets.MESHERY_TOKEN }} # Meshery Cloud Authentication token, signin to meshery-cloud to get one, example: ey.....
          prNumber: ${{ env.PULL_NO }} # auto-filled from the above step
          application_type: "Kubernetes Manifest" # your application type, could be any of three: "Kubernetes Manifest", "Docker Compose", "Helm Chart"
          filePath: "action/__tests__/manifest-test" # relative file-path from the root directory in the github-runner env, you might require to checkout the repository as described in step 2
</code></pre>
<h4>Notes</h4>
The filesystem-approach asks for your relative file-path and automatically merges all the yaml files together to bundle up into one. So you might like to give the root directory where all the yamls are located. It doesn't move recursevely in internal folders, so only the first level yaml files are checked.
</details>
<details>
  <summary><h4>Example: Importing design from URL</h4>
</summary>
<pre><code>
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
          mesheryToken: ${{ secrets.MESHERY_TOKEN }} # Meshery Cloud Authentication token, signin to meshery-cloud to get one, example: ey.....
          prNumber: ${{ env.PULL_NO }} # auto-filled from the above step
          application_type: "Helm Chart" # your application type, could be any of three: "Kubernetes Manifest", "Docker Compose", "Helm Chart"
          application_url: "https://github.com/meshery/meshery.io/raw/master/charts/meshery-v0.6.88.tgz"
</code></pre>
</details>
<details>
  <summary><h4>Customizing your workflow</h4></summary>

The following is a list of Action Inputs.

<pre><code lang="yaml">
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
</code></pre>
</details>


<div>&nbsp;</div>
<h1>Join the Layer5 Community!</h1>
<p style="clear:both;">
<h2><a name="contributing"></a><a name="community"></a> <a href="https://layer5.io/community">Community</a> and <a href="https://layer5.io/community/handbook">Contributions</a></h2>
<a href="https://slack.layer5.io"><img src="/.github/readme/images/Layer5-MeshMentors.png" width="150px" /></a>
<p> We warmly welcome all contributors! Our projects are community-built and each welcomes open collaboration. As you get started, please review this project's <a href="https://github.com/layer5io/layer5/blob/master/CONTRIBUTING.md">contributing guidelines</a>. Whether you are a user or code contributor and whether you're opening an <a href="/../../issues">issue</a> or a <a href="/../../pulls">pull request</a>, know that any form of your engagement is considered contribution and is appreciated. Contributors are expected to adhere to the <a href="https://github.com/cncf/foundation/blob/master/code-of-conduct.md">CNCF Code of Conduct</a>.
</p>
<p>
  Join us in the <a href="https://discuss.layer5.io">discussion forum</a> and on <a href="https://slack.layer5.io"><img src=".github/readme/images/community.svg" height="16px" align="bottom" /> Slack</a> to learn more about Layer5 and its community! Make sure you see the <a href="https://layer5.io/community/newcomers">newcomer's guide</a> for a tour of resources available to you.
</p>
<p>
<a href="https://slack.meshery.io">

<picture align="right">
  <source media="(prefers-color-scheme: dark)" srcset=".github/assets/images/buttons/slack-dark-128.webp"  width="110px" align="right" style="margin-left:10px;margin-top:10px;">
  <source media="(prefers-color-scheme: light)" srcset=".github/assets/images/buttons/slack-128.webp" width="110px" align="right" style="margin-left:10px;padding-top:5px;">
  <img alt="Shows an illustrated light mode meshery logo in light color mode and a dark mode meshery logo dark color mode." src=".github/assets/images/buttons/slack-128.webp" width="110px" align="right" style="margin-left:10px;padding-top:13px;">
</picture>
</a>
<a href="https://layer5.io/community"><img alt="Layer5 Cloud Native Community" src=".github/readme/images/community.svg" style="margin-right:10px;" width="130px" align="left" /></a>
✔️ <em><strong>Join</strong></em> any or all of the weekly meetings on the <a href="https://meet.layer5.io">community calendar</a>.<br />
✔️ <em><strong>Watch</strong></em> community <a href="http://youtube.com/Layer5io?sub_confirmation=1">meeting recordings</a>.<br />
✔️ <em>Fill-in</em> a <a href="https://layer5.io/newcomers">community member form</a> to gain access to community resources.
<br />
✔️ <em><strong>Discuss</strong></em> in the <a href="https://discuss.layer5.io">Community Forum</a>.<br />
✔️ <em><strong>Explore more</strong></em> in the <a href="https://layer5.io/community/handbook">Community Handbook</a>.<br />
</p>
<div align="center"><i>Not sure where to start?</i> Grab an open issue with the <a href="https://github.com/issues?q=is%3Aopen+is%3Aissue+archived%3Afalse+(org%3Alayer5io+OR+org%3Ameshery+OR+org%3Alayer5labs+OR+org%3Aservice-mesh-performance+OR+org%3Aservice-mesh-patterns+OR+org%3Ameshery-extensions)+label%3A%22help+wanted%22">help-wanted label</a>.
</div>
<br />
<ul>
  <li>Find us on Twitter: <a href="https://twitter.com/layer5">@layer5</a>, <a href="https://twitter.com/mesheryio">@mesheryio</a>, and <a href="https://x.com/kanvas_new">@kanvas_new</a>.</li>
  <li>Visit us on LinkedIn: <a href="https://www.linkedin.com/company/layer5">Layer5</a>, <a href="https://www.linkedin.com/showcase/meshery/">Meshery</a>, and <a href="https://www.linkedin.com/showcase/service-mesh-performance">Cloud Native Performance</a>.</li>
  <li>Subscribe on <a href="http://youtube.com/Layer5io?sub_confirmation=1">Youtube</a>.</li>
</ul>

### License

All of Layer5's projects are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
