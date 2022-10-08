<p style="text-align:center;" align="center"><a href="https://layer5.io/meshery"><picture align="center">
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/layer5io/layer5/master/src/assets/images/meshery/full-logo/meshery-logo-tag-light-white-text-side.png"  width="70%" align="center" style="margin-bottom:20px;">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/meshery/meshery-logo-tag-light-text-side.png" width="70%" align="center" style="margin-bottom:20px;">
  <img alt="Shows an illustrated light mode meshery logo in light color mode and a dark mode meshery logo dark color mode." src="https://raw.githubusercontent.com/meshery/meshery/master/.github/assets/images/meshery/meshery-logo-tag-light-text-side.png" width="70%" align="center" style="margin-bottom:20px;">
</picture></a><br /><br /></p>

<p align="center">
<a href="https://hub.docker.com/r/layer5/meshery" alt="Docker pulls">
  <img src="https://img.shields.io/docker/pulls/layer5/meshery.svg" /></a>
<a href="https://github.com/issues?utf8=✓&q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Alayer5io+label%3A%22help+wanted%22+" alt="GitHub issues by-label">
  <img src="https://img.shields.io/github/issues/layer5io/meshery/help%20wanted.svg?color=informational" /></a>
<a href="https://github.com/meshery/meshery" alt="LICENSE">
  <img src="https://img.shields.io/github/license/meshery/meshery?color=brightgreen" /></a>
<a href="https://goreportcard.com/report/github.com/meshery/meshery" alt="Go Report Card">
  <img src="https://goreportcard.com/badge/github.com/meshery/meshery" /></a>
<a href="https://github.com/meshery/meshery/actions" alt="Build Status">
  <img src="https://img.shields.io/github/workflow/status/meshery/meshery/Meshery%20Build%20and%20Releaser%20(edge)" /></a>
<a href="https://bestpractices.coreinfrastructure.org/projects/3564" alt="CLI Best Practices">
  <img src="https://bestpractices.coreinfrastructure.org/projects/3564/badge" /></a>
<a href="https://discuss.layer5.io" alt="Discuss Users">
  <img src="https://img.shields.io/discourse/users?label=discuss&logo=discourse&server=https%3A%2F%2Fdiscuss.layer5.io" /></a>
<a href="https://slack.layer5.io" alt="Join Slack">
  <img src="https://img.shields.io/badge/Slack-@layer5.svg?logo=slack"></a>
<a href="https://twitter.com/intent/follow?screen_name=mesheryio" alt="Twitter Follow">
  <img src="https://img.shields.io/twitter/follow/mesheryio.svg?label=Follow+Meshery&style=social" /></a>
</p>

[Meshery](https://meshery.io) is the cloud native management plane offering lifecycle, configuration, and performance management of Kubernetes, service meshes, and your workloads.

<p align="center">
Meshery is a Cloud Native Computing Foundation project.
</p>

# Docker Extension for Meshery

The Docker Extension for Meshery extends Docker Desktop’s position as the cloud native developer’s go-to Kubernetes environment with easy access to the next layer of cloud native infrastructure: service meshes. Service mesh or not, though, Meshery offers a visual topology for designing Docker-Compose applications, operating Kubernetes, service meshes, and their workloads. Meshery brings deep support of 10 different service meshes to the fingertips of Docker Desktop developers in connection with Docker Desktop’s ability to deliver Kubernetes locally.

<h2><a name="running"></a>Get Started with the Docker Extension for Meshery</h2>

<h3>Using Docker Desktop</h3>
<p>Navigate to the Extensions area of Docker Desktop.</p>

<h3>Using <code>docker</code></h3>
<p>Meshery runs as a set of containers inside your Docker Desktop virtual machine.</p>
<pre>docker extension install meshery/docker-extension-meshery</pre>
<p>See the <a href="https://docs.meshery.io/installation/quick-start">quick start</a> guide.</p>
<p style="clear:both;">&nbsp;</p>

## Using the Docker Extension for Meshery

1. Install any service mesh with the check of a box.
1. Import your Docker Compose apps for visual design and deployment to Kubernetes and service meshes.

<p align="center"><a href="https://raw.githubusercontent.com/meshery/meshery/master/install/docker-extension/docs/img/docker-desktop-extension-for-meshery.png"><img src="https://raw.githubusercontent.com/meshery/meshery/master/install/docker-extension/docs/img/docker-desktop-extension-for-meshery.png" width="90%" align="center" /></a></p>

## Docker Extension for Meshery Architecture

The Docker Extension for Meshery deploys Meshery to your local Docker host as a Docker Compose application.

<p align="center"><a href="https://raw.githubusercontent.com/meshery/meshery/master/install/docker-extension/docs/img/docker-extension-for-meshery-architecture.png"><img src="https://raw.githubusercontent.com/meshery/meshery/master/install/docker-extension/docs/img/docker-extension-for-meshery-architecture.png" width="90%" align="center" /></a></p>
Learn more about <a href="https://docs.meshery.io/architecture">Meshery's architecture</a>.

<p style="text-align:center; width:100%;" align="center">
<a href ="https://layer5.io/community"><img alt="MeshMates" src="https://docs.meshery.io/assets/img/readme/layer5-community-sign.png" style="margin-right:10px; margin-bottom:7px;" width="28%" align="center" /></a>
</p>
<p style="text-align:center; width:100%;" align="center">
<h3 style="text-align:center;" align="center"><em>Have questions? Need help?</em> <strong>Ask in the <a href="https://discuss.layer5.io">Community Forum</a></strong>.</h3></p>
