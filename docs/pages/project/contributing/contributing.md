---
layout: page
title: Contributing
permalink: project/contributing
description: How to contribute to the Meshery project.
language: en
type: project
redirect_from:
  - contributing
  - project/contributing/
category: contributing
---

Please do! Thanks for your help! 🎈 Meshery is community-built and welcomes collaboration. Contributors are expected to adhere to the [CNCF's Code of Conduct](https://github.com/layer5io/meshery/blob/master/CODE_OF_CONDUCT.md).

# Contributing Guides

{% include suggested-reading.html diffName ="true" language = "en" %}

## Not sure where to start?

Follow these steps and you'll be right at home.

1. See the [Newcomers Guide](https://layer5.io/community/newcomers) for how, where, and why to contribute.

2. Sign up for a [_MeshMate_](https://layer5.io/community/meshmates) to find the perfect Mentor to help you explore the Layer5 projects and find your place in the community:

- **Familiarize** yourself with the broader set of community projects (take a look at the [Repository Overview](https://layer5.io/community/handbook/repository-overview): Spend time understanding each of the initiatives through high-level overviews available in the community drive and through discussions with your MeshMate.
- **Identify** your area of interest: Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop.
- **Run** Meshery: Put on your user hat and walk-through all of Meshery’s features and functions as a user.
- **Build** Meshery: Confirm that you have a usable development environment.
- **Discuss** with the community by engaging in the [discussion forum](https://discuss.layer5.io).
- **Contribute** by grabbing any open issue with the [help-wanted label](https://github.com/meshery/meshery/issues/) and jump in. If needed, create a [new issue](https://github.com/meshery/meshery/issues/new/choose). All [pull requests](https://github.com/meshery/meshery/pulls) should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to [automatically close issues in GitHub](https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords).
- **Fill-in** a [community member form](https://layer5.io/newcomers) community member form to gain access to community resources.

## General Contribution Flow

To contribute to Meshery, from creating a fork to creating pull request, please follow the basic fork-and-pull request workflow described [here]({{site.baseurl}}/project/contributing/contributing-gitflow).


<details>

<summary>Signing-off on Commits (Developer Certificate of Origin)</summary>

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

