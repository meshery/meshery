---
layout: page
title: Contributing
permalink: project/contributing
description: How to contribute to the Meshery project.
language: en
type: project
category: contributing
---

Please do! Thanks for your help! 🎈

This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [CNCF's Code of Conduct](https://github.com/layer5io/meshery/blob/master/CODE_OF_CONDUCT.md).

## Not sure where to start?

Follow these steps and you'll be right at home.

1. See the [_Community Welcome Guide_](https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit) for how, where, and why to contribute.

2. Sign up for a [_MeshMate_](https://layer5.io/community#meshmate) to find the perfect Mentor to help you explore the Layer5 projects and find your place in the community:

- **Familiarize** yourself with all the Layer5 projects (Take a look at the [Community Drive](https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA) and the [Layer5 Repository Overview](https://docs.google.com/document/d/1brtiJhdzal_O6NBZU_JQXiBff2InNtmgL_G1JgAiZtk/edit#heading=h.uwtb5xf7b5hw): Spend time understanding each of the Layer5 initiatives through high-level overviews available in the community drive and through discussions with your MeshMate.
- **Identify** your area of interest: Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop.
- **Run** Meshery: Put on your user hat and walk-through all of Meshery’s features and functions as a user.
- **Build** Meshery: Confirm that you have a usable development environment.
- **Communicate** with the Layer5 community by joining the [Slack account](http://slack.layer5.io).
- **Contribute** by grabbing any open issue with the [help-wanted label](https://github.com/layer5io/meshery/issues/) and jump in. If needed, create a [new issue](https://github.com/layer5io/meshery/issues/new/choose). All [pull requests](https://github.com/layer5io/meshery/pulls) should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to [automatically close issues in GitHub](https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords).

**Sections**

- [General Contribution Flow]({{site.baseurl}}/project/contributing#general-contribution-flow)
- <a href="#commit-signing">Developer Certificate of Origin</a>
- Meshery Contribution Flow
  - <a href="{{site.baseurl}}/project/contributing-docs">Meshery Documentation</a>
  - <a href="{{site.baseurl}}/project/contributing-cli">Meshery CLI</a>
  - <a href="#contributing-meshery">Meshery Backend</a>
    - <a href="#adapter">Writing a Meshery Adapter</a>
  - <a href="{{site.baseurl}}/project/contributing-ui">Meshery UI</a>
    Relevant coding style guidelines are the Go Code Review Comments and the Formatting and style section of Peter Bourgon's Go: Best Practices for Production Environments.

## <a name="contributing">General Contribution Flow</a>

To contribute to Meshery, please follow the fork-and-pull request workflow described [here]({{site.baseurl}}/project/contributing/CONTRIBUTING-gitflow.md).

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
$ commit -s -m “my commit message w/signoff”
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

#### Build and run Meshery server

To build & run the Meshery server code, run the following command:

```sh
make run-local
```

Any time changes are made to the GO code, you will have to stop the server and run the above command again.
Once the Meshery server is up and running, you should be able to access Meshery on your `localhost` on port `9081` at `http://localhost:9081`. One thing to note, you might NOT see the [Meshery UI](#contributing-ui) until the UI code is built as well.
After running Meshery server, you will need to select your **Cloud Provider** by navigating to `localhost:9081`. Only then you will be able to use the Meshery UI on port `3000`.

#### Building Docker image

To build a Docker image of Meshery, please ensure you have `Docker` installed to be able to build the image. Now, run the following command to build the Docker image:

```sh
make docker
```

# Suggested Reading

{% assign sorted_reading = site.pages | sort: page.title | reverse %}

<ul>
  {% for item in sorted_reading %}
  {% if item.type=="project" and item.category=="contributing" and item.list!="exclude" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>
