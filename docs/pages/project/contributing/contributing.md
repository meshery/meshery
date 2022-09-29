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

1. See the [_Community Welcome Guide_](https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit) for how, where, and why to contribute.

2. Sign up for a [_MeshMate_](https://layer5.io/community/meshmates) to find the perfect Mentor to help you explore the Layer5 projects and find your place in the community:

- **Familiarize** yourself with the broader set of community projects (take a look at the [Repository Overview](https://docs.google.com/document/d/1brtiJhdzal_O6NBZU_JQXiBff2InNtmgL_G1JgAiZtk/edit#heading=h.uwtb5xf7b5hw): Spend time understanding each of the initiatives through high-level overviews available in the community drive and through discussions with your MeshMate.
- **Identify** your area of interest: Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop.
- **Run** Meshery: Put on your user hat and walk-through all of Meshery’s features and functions as a user.
- **Build** Meshery: Confirm that you have a usable development environment.
- **Discuss** with the community by engaging in the [discussion forum](https://discuss.layer5.io).
- **Contribute** by grabbing any open issue with the [help-wanted label](https://github.com/meshery/meshery/issues/) and jump in. If needed, create a [new issue](https://github.com/meshery/meshery/issues/new/choose). All [pull requests](https://github.com/meshery/meshery/pulls) should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to [automatically close issues in GitHub](https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords).
- **Fill-in** a [community member form](https://layer5.io/newcomers) community member form to gain access to community resources.

## <a name="contributing">General Contribution Flow</a>

To contribute to Meshery, from creating a fork to creating pull request, please follow the basic fork-and-pull request workflow described [here]({{site.baseurl}}/project/contributing/contributing-gitflow).

### <a name="commit-signing">Signing-off on Commits (Developer Certificate of Origin)</a>

To contribute to this project, you must agree to the Developer Certificate of
Origin (DCO) for each commit you make. The DCO is a simple statement that you,
as a contributor, have the legal right to make the contribution.

See the [DCO](https://developercertificate.org) file for the full text of what you must agree to
and how it works [here](https://github.com/probot/dco#how-it-works).
To signify that you agree to the DCO for contributions, you simply add a line to each of your
git commit messages:

```
Signed-off-by: Jane Smith <jane.smith@example.com>
```

In most cases, you can add this signoff to your commit automatically with the
`-s` or `--signoff` flag to `git commit`. You must use your real name and a reachable email
address (sorry, no pseudonyms or anonymous contributions). An example of signing off on a commit:

```
$ git commit -s -m “my commit message w/signoff”
```

To ensure all your commits are signed, you may choose to add this alias to your global `.gitconfig`:

_~/.gitconfig_

```
[alias]
  amend = commit -s --amend
  cm = commit -s -m
  commit = commit -s
```

Or you may configure your IDE, for example, Visual Studio Code to automatically sign-off commits for you:

<a href="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" ><img src="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" width="50%"><a>

### <a name="contributing-meshery">Meshery Contribution Flow</a>

Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

Relevant coding style guidelines are the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) and the _Formatting and style_ section of Peter Bourgon's [Go: Best
Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style).

**Please note**: All `make` commands should be run in a terminal from within the Meshery's main folder.

#### Prerequisites for building Meshery in your development environment:

1. `Go` version 1.11+ installed if you want to build and/or make changes to the existing code.
1. `GOPATH` environment variable should be configured appropriately
1. `npm` and `node` should be installed on your machine, preferably the latest versions.
1. Fork this repository (`git clone https://github.com/layer5io/meshery.git`), clone your forked version of Meshery to your local, preferably outside `GOPATH`. If you happen to checkout Meshery inside your `GOPATH` and you have a version of `Go` prior to version 1.13, please set an environment variable `GO111MODULE=on` to enable GO Modules.
