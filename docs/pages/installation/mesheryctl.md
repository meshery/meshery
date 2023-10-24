---
layout: default
title: Install mesheryctl
permalink: installation/mesheryctl
type: installation
redirect_from:
 - installation/mesheryctl
 - installation/mesheryctl/
display-title: "true"
language: en
list: include
# image: /assets/img/platforms/brew.png
---

Meshery's command line client is `mesheryctl` and is the recommended tool for configuring and deploying one or more Meshery deployments. To install `mesheryctl` on your system, you may choose from any of the following supported methods.

`mesheryctl` can be installed via `bash`, [Homebrew]({{site.baseurl}}/installation/linux-mac/brew), [Scoop]({{site.baseurl}}/installation/windows/scoop) or [directly downloaded](https://github.com/meshery/meshery/releases/latest).

{% include mesheryctl/installation-brew.md %}

{% include mesheryctl/installation-bash.md %}

{% include mesheryctl/installation-scoop.md %}

Continue deploying Meshery onto one of the [Supported Platforms]({{ site.baseurl }}/installation/platforms).
