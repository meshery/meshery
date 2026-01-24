---
layout: page
title: Contributing
permalink: project/contributing
abstract: How to contribute to the Meshery project.
# redirect_from:
#  - project/contributing/
display-title: "false"
language: en
type: project
category: contributing
list: exclude
display-suggested-reading: false
abstract: How to contribute to the Meshery project and any of its components.
---

{% include sdd-alert.md %}

# Contributing

Please do! Thanks for your help! 🎈 Meshery is community-built and welcomes collaboration. Contributors are expected to adhere to the [CNCF's Code of Conduct](https://github.com/meshery/meshery/blob/master/CODE_OF_CONDUCT.md).

## General Contribution Flow

Meshery and it's various architectural components are written in different languages, including Golang, Javascript (React.js and Next.js) To make building, testing, and the experience of contributing consistent across all Meshery components, a `Makefile` is included in the every repository. These `make` targets are what you will use to build, run, test, and document.

To contribute to Meshery, please follow this basic fork-and-pull request [gitflow]({{site.baseurl}}/project/contributing/contributing-gitflow).

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

<pre><code>$ git commit -s -m “my commit message w/signoff”</code></pre>
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
  <li><strong><a href="https://play.meshery.io">Play with Meshery</a></strong>: Put on your user hat and walk-through all of Meshery’s features and functions as a user.</li>
  <li><strong>Build Meshery Server and UI</strong>: Confirm that you have a usable development environment. See <a href="#contributing-guides">Guides</a> below.</li>
  <li><strong>Discuss</strong> with the community by engaging in the <a href="https://meshery.io/community#discussion-forums">discussion forum</a>.</li>
  <li><strong>Contribute</strong> by grabbing any open issue with the <a href="https://github.com/issues?q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Ameshery+org%3Ameshery-extensions+org%3Aservice-mesh-performance+org%3Aservice-mesh-patterns+label%3A%22help+wanted%22+">help-wanted label</a> and jump in. If needed, create a <a href="https://github.com/meshery/meshery/issues/new/choose">new issue</a>. All <a href="https://github.com/meshery/meshery/pulls">pull requests</a> should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to <a href="https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords">automatically close issues in GitHub</a>.</li>
  <li><strong>Fill in</strong> a <a href="https://meshery.io/newcomers">community member form</a> to gain access to community resources.</li>
</ol>
</details>

## Specific Contribution Guides

Here is a complete list of all of Meshery's contributing guides from Server to UI to CLI to Extensions and so on.

<!-- Contributing & Community -->
  <div class="section">
      <!-- CONTRIBUTING -->
      <ul class="section-title">
       {% assign contributing = site.pages | where: "category","contributing" %}
          {% for item in contributing %}
          {% if item.category=="contributing" and item.language=="en" -%}
            <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
            {% if item.abstract != " " %}
              - {{ item.abstract }}
            {% endif %}
            </li>
            {% endif %}
          {% endfor %}
      </ul>
  </div>
