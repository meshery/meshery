---
layout: default
title: Install Meshery CLI on Windows
permalink: installation/windows
type: installation
category: mesheryctl
redirect_from:
- installation/platforms/windows
display-title: "true"
language: en
list: include
image: /assets/img/platforms/wsl2.png
abstract: Guide for installing mesheryctl and setting up a Meshery development environment on Windows.
---

On Windows systems, `mesheryctl` can be installed via Scoop or [downloaded directly](https://github.com/meshery/meshery/releases/latest).

{% include mesheryctl/installation-scoop.md %}

## Install `mesheryctl` as a direct download

Follow the installation steps to install the `mesheryctl` CLI. Then, execute:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">./mesheryctl system start</div></div>
</pre>

If you are installing Meshery on Docker, execute the following command:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">./mesheryctl system start -p docker</div></div>
</pre>

Optionally, move the `mesheryctl` binary to a directory in your `PATH`.

## Local Development Setup

For contributors looking to build Meshery from source on Windows, follow these steps to bypass common environment issues:

### 1. Run the Server (Go)
From the root directory of the repository, use the specific path to the main entry point:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">go run server/cmd/main.go</div></div>
</pre>

### 2. Run the UI (Next.js)
Navigate to the `ui` folder. If the default `npm run dev` fails due to `rm` command errors (Linux-specific commands), use the following workflow:
<pre class="codeblock-pre">
<div class="codeblock"><div class="clipboardjs">cd ui
npm install
npx next dev</div></div>
</pre>

**Note:** If you need to clean the build cache on Windows, install `rimraf` via npm (`npm install -g rimraf`) and use `rimraf .next out` instead of the standard `make clean` command.

# Related Reading

## Meshery CLI Guides

Guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "name" %}

<ul>
  {% for item in sorted_guides %}
  {% if item.type=="guides" and item.category=="mesheryctl" and item.list!="exclude" and item.language=="en" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
    <li><a href="{{ site.baseurl }}/installation/upgrades#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
</ul>

{% include related-discussions.html tag="mesheryctl" %}

{:toc}