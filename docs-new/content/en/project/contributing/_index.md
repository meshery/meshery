---
title: Contributing
description: "How to contribute to the Meshery project and its components"
weight: 20
aliases:
  - /project/contributing
display_title: false
---

# Contributing 

Please do! Thanks for your help! 🎈 Meshery is community-built and welcomes collaboration. Contributors are expected to adhere to the [CNCF's Code of Conduct](https://github.com/meshery/meshery/blob/master/CODE_OF_CONDUCT.md).

## General Contribution Flow

Meshery and it's various architectural components are written in different languages, including Golang, Javascript (React.js and Next.js) To make building, testing, and the experience of contributing consistent across all Meshery components, a `Makefile` is included in the every repository. These `make` targets are what you will use to build, run, test, and document.

To contribute to Meshery, please follow this basic fork-and-pull request [gitflow](/project/contributing/contributing-gitflow).

<details>
<summary>Adding your sign-off on commits (Developer Certificate of Origin)</summary>
<ul>
<li>
To contribute to this project, you must agree to the Developer Certificate of
Origin (DCO) for each commit you make. The DCO is a simple statement that you,
as a contributor, have the legal right to make the contribution.
</li>

<li>
See the <a href="https://developercertificate.org"> DCO </a> file for the full text of what you must agree to
and how it works <a href="https://github.com/probot/dco#how-it-works">here</a>.
To signify that you agree to the DCO for contributions, you simply add a line to each of your
git commit messages:

<pre><code>
Signed-off-by: Jane Smith <jane.smith@example.com>
</code></pre></li>

<li>
In most cases, you can add this signoff to your commit automatically with the
<code>-s</code> or <code>--signoff</code> flag to <code>git commit</code>. You must use your real name and a reachable email
address (sorry, no pseudonyms or anonymous contributions). An example of signing off on a commit:

<pre><code>$ git commit -s -m "my commit message w/signoff"</code></pre>
</li>

<li>
To ensure all your commits are signed, you may choose to add this alias to your global <code>.gitconfig</code>:

~/.gitconfig

<pre><code>
[alias]
  amend = commit -s --amend
  cm = commit -s -m
  commit = commit -s
</code></pre>

Or you may configure your IDE, for example, VSCode to automatically sign-off commits for you:<a href="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" ><img src="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" width="50%"/></a>

</li>
</ul>
</details>

## Not sure where to start?

<details>
<summary>Follow these steps and you'll be right at home.</summary>
<ol>
  <li>See the <a href="https://meshery.io/community">Newcomers Guide</a> for how, where, and why to contribute.</li>
  <li>Sign up for a <a href="https://meshery.io/community#meshmates"><em>MeshMate</em></a> to find the perfect Mentor to help you explore the projects and find your place in the community.</li>
  <li><strong>Familiarize</strong> yourself with the broader set of projects in Meshery's ecosystem, including the <a href="https://github.com/meshery-extensions">meshery-extensions</a> repositories (this <a href="https://layer5.io/community/handbook/repository-overview">Repository Overview</a> is a helpful resource): Spend time understanding each of the initiatives through high-level overviews available in the community drive and through discussions with your MeshMate.</li>
  <li><strong>Identify</strong> your area of interest: Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop.</li>
  <li><strong><a href="https://play.meshery.io">Play with Meshery</a></strong>: Put on your user hat and walk-through all of Meshery's features and functions as a user.</li>
  <li><strong>Build Meshery Server and UI</strong>: Confirm that you have a usable development environment. See <a href="#contributing-guides">Guides</a> below.</li>
  <li><strong>Discuss</strong> with the community by engaging in the <a href="https://meshery.io/community#discussion-forums">discussion forum</a>.</li>
  <li><strong>Contribute</strong> by grabbing any open issue with the <a href="https://github.com/issues?q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Ameshery+org%3Ameshery-extensions+org%3Aservice-mesh-performance+org%3Aservice-mesh-patterns+label%3A%22help+wanted%22+">help-wanted label</a> and jump in. If needed, create a <a href="https://github.com/meshery/meshery/issues/new/choose">new issue</a>. All <a href="https://github.com/meshery/meshery/pulls">pull requests</a> should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to <a href="https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords">automatically close issues in GitHub</a>.</li>
  <li><strong>Fill in</strong> a <a href="https://meshery.io/newcomers">community member form</a> to gain access to community resources.</li>
</ol>
</details>

## Specific Contribution Guides

Here is a complete list of all of Meshery's contributing guides from Server to UI to CLI to Extensions and so on.

- [Build & Release (CI)](/project/contributing/build-and-release) - Details of Meshery's build and release strategy.
- [Contributing to Meshery Adapters](/project/contributing/contributing-adapters) - How to contribute to Meshery Adapters
- [Contributing to Meshery CLI End-to-End Tests](/project/contributing/contributing-cli-tests) - How to contribute to Meshery Command Line Interface end-to-end testing with BATS.
- [Contributing to Meshery CLI](/project/contributing/contributing-cli) - How to contribute to Meshery Command Line Interface.
- [Contributing to Meshery Docker Extension](/project/contributing/contributing-docker-extension) - How to contribute to Meshery Docker Extension
- [Meshery Documentation Structure and Organization](/project/contributing/contributing-docs/contributing-docs-structure) - Audience, high-level outline & information architecture for Meshery Documentation
- [Contributing to Meshery Docs](/project/contributing/contributing-docs) - How to contribute to Meshery Docs.
- [How to write MeshKit compatible errors](/project/contributing/contributing-error) - How to declare errors in Meshery components.
- [Contributing to Meshery using git](/project/contributing/contributing-gitflow) - How to contribute to Meshery using git
- [Meshery CLI Style Guide](/project/contributing/contributing-cli/contributing-cli-guide) - Design principles and code conventions.
- [Contributing to Model Components](/project/contributing/contributing-models/contributing-components) - How to contribute to Meshery Model Components
- [Contributing to Models Quickstart](/project/contributing/contributing-models/contributing-models-quick-start) - A quick guide to creating your own Meshery Models quickly.
- [Contributing to Models](/project/contributing/contributing-models) - How to contribute to Meshery Models, Components, Relationships, Policies...
- [Contributing to Meshery Policies](/project/contributing/contributing-policies) - How to contribute to Meshery Policies
- [Contributing to Model Relationships](/project/contributing/contributing-models/contributing-relationships) - How to contribute to Meshery Models Relationships, Policies...
- [Contributing to Meshery Schemas](/project/contributing/contributing-schemas) - How to contribute to Meshery Schemas
- [Contributing to Meshery Server Events](/project/contributing/contributing-server/contributing-server-events) - Guide to help backend contributors send server events using Golang.
- [Contributing to Meshery Server](/project/contributing/contributing-server) - How to contribute to Meshery Server
- [Contributing to Meshery Ui - Notification Center](/project/contributing/contributing-ui/contributing-ui-notification-center) - How to contribute to the Notification Center in Meshery's web-based UI.
- [Schema-Driven UI Development in Meshery](/project/contributing/contributing-ui/contributing-ui-schemas) - How to contribute to Meshery Schemas for UI
- [Contributing to Meshery UI - Sistent](/project/contributing/contributing-ui/contributing-ui-sistent) - How to contribute to the Meshery's web-based UI using sistent design system.
- [Contributing to Meshery UI End-to-End Tests](/project/contributing/contributing-ui/contributing-ui-tests) - How to contribute to end-to-end testing in Meshery UI using Playwright.
- [Contributing to Meshery UI - Dashboard Widgets](/project/contributing/contributing-ui/contributing-ui-widgets) - Guide to extending Meshery dashboards with custom widgets.
- [Contributing to Meshery UI](/project/contributing/contributing-ui) - How to contribute to Meshery UI
- [Dev Setup on Windows](/project/contributing/meshery-windows) - Development environment setup guide for Windows
