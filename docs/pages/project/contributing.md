---
layout: page
title: Contributing
permalink: project/contributing
---


Please do! Thanks for your help! :balloon:

This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [CNCF's Code of Conduct](https://github.com/layer5io/meshery/blob/master/CODE_OF_CONDUCT.md). 
## Not sure where to start? 

Follow these steps and you'll be right at home.

1. See the [*Community Welcome Guide*](https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit) for how, where, and why to contribute.

2. Sign up for a [*MeshMate*](https://layer5.io/community#meshmate) to find the perfect Mentor to help you explore the Layer5 projects and find your place in the community:
  - **Familiarize** yourself with all the Layer5 projects (Take a look at the [Community Drive](https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA) and the [Layer5 Repository Overview](https://docs.google.com/document/d/1brtiJhdzal_O6NBZU_JQXiBff2InNtmgL_G1JgAiZtk/edit#heading=h.uwtb5xf7b5hw): Spend time understanding each of the Layer5 initiatives through high-level overviews available in the community drive and through discussions with your MeshMate.
  - **Identify** your area of interest: Use the time with your MeshMate to familiarize yourself with the architecture and technologies used in the projects. Inform your MeshMate of your current skills and what skills you aim to develop.
  - **Run** Meshery: Put on your user hat and walk-through all of Meshery’s features and functions as a user.
  - **Build** Meshery: Confirm that you have a usable development environment.
  - **Communicate** with the Layer5 community by joining the [Slack account](http://slack.layer5.io).
  - **Contribute** by grabbing any open issue with the [help-wanted label](https://github.com/layer5io/meshery/issues/) and jump in. If needed, create a [new issue](https://github.com/layer5io/meshery/issues/new/choose). All [pull requests](https://github.com/layer5io/meshery/pulls) should reference an open issue. Include keywords in your pull request descriptions, as well as commit messages, to [automatically close issues in GitHub](https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords).

**Sections**
- General Contribution Flow
- <a href="#commit-signing">Developer Certificate of Origin</a>
- Meshery Contribution Flow
  - <a href="#contributing-docs">Meshery Documentation</a>
  - <a href="#contributing-meshery">Meshery Backend</a>
    - <a href="#adapter">Writing a Meshery Adapter</a>
  - <a href="#contributing-ui">Meshery UI</a>
Relevant coding style guidelines are the Go Code Review Comments and the Formatting and style section of Peter Bourgon's Go: Best Practices for Production Environments.

## <a name="contributing">General Contribution Flow</a>

To contribute to Meshery, please follow the fork-and-pull request workflow described [here](CONTRIBUTING-gitflow.md).

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

To ensure all your commits are signed, you may choose to add this alias to your global ```.gitconfig```:

*~/.gitconfig*
```
[alias]
  amend = commit -s --amend
  cm = commit -s -m
  commit = commit -s
```
Or you may configure your IDE, for example, Visual Studio Code to automatically sign-off commits for you:

<a href="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" ><img src="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" width="50%"><a>
  
### <a name="contributing-docs">Documentation Contribution Flow</a>
Please contribute! Meshery documentation uses GitHub Pages to host the docs site. Learn more about [Meshery's documentation framework](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit?usp=sharing). The process of contributing follows this flow:

1. Create a fork, if you have not already, by following the steps described [here](CONTRIBUTING-gitflow.md)
1. In the local copy of your fork, navigate to the docs folder.
`cd docs`
1. Create and checkout a new branch to make changes within
`git checkout -b <my-changes>`
1. Edit/add documentation.
`vi <specific page>.md`
1. Run site locally to preview changes.
`make site`
1. Commit, [sign-off](#commit-signing), and push changes to your remote branch.
`git push origin <my-changes>`
1. Open a pull request (in your web browser) against our main repo: https://github.com/layer5io/meshery.


### <a name="contributing-meshery">Meshery Contribution Flow</a>
Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

Relevant coding style guidelines are the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) and the _Formatting and style_ section of Peter Bourgon's [Go: Best 
Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style).

__Please note__: All `make` commands should be run in a terminal from within the Meshery's main folder.

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

#### <a name="adapter">Writing a Meshery Adapter</a>
Meshery uses adapters to provision and interact with different service meshes. Follow these instructions to create a new adapter or modify and existing adapter.

1. Get the proto buf spec file from Meshery repo: 
```wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto```
1. Generate code
    1. Using Go as an example, do the following:
        - adding GOPATH to PATH: `export PATH=$PATH:$GOPATH/bin`
        - install grpc: `go get -u google.golang.org/grpc`
        - install protoc plugin for go: `go get -u github.com/golang/protobuf/protoc-gen-go`
        - Generate Go code: `protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/`
    1. For other languages, please refer to gRPC.io for language-specific guides.
1. Implement the service methods and expose the gRPC server on a port of your choice (e.g. 10000). 

_Tip:_ The [Meshery Adapter for Istio](https://github.com/layer5io/meshery-istio) is a good reference adapter to use as an example of a Meshery Adapter written in Go.

### <a name="contributing-mesheryctl">mesheryctl Contribution Flow</a>

`mesheryctl` is written in Golang or the Go Programming Language. For development use Go version 1.15+.

The [`/mesheryctl`](https://github.com/layer5io/meshery/tree/master/mesheryctl) folder contains the complete code for `mesheryctl`.

After making changes, run `make` in the `mesheryctl` folder to build the binary. You can then use the binary by, say, `./mesheryctl system start`.

#### mesheryctl command reference

Detailed documentation of the `mesheryctl` commands is available in the [Meshery Docs](https://docs.meshery.io/reference/mesheryctl).

#### Guidelines and resources for contributing to mesheryctl

`mesheryctl` might be the interface that the users first have with Meshery. As such, `mesheryctl` needs to provide a great UX.

The following principles should be taken in mind while designing `mesheryctl` commands-

1. Provide user experiences that are familiar.
2. Make the commands and their behavior intuitive.
3. Avoid long commands with chained series of flags.
4. Design with automated testing in mind, e.g. provide possibility to specify output format as json (-o json) for easy inspection of command response.

Part of delivering a great user experience is providing intuitive interfaces. In the case of `mesheryctl`, we should take inspiration from and deliver similar user experiences as popular CLIs do in this ecosystem, like `kubectl` and `docker`. Here is relevant `kubectl` information to reference - [Kubectl SIG CLI Community Meeting Minutes](https://docs.google.com/document/u/2/d/1r0YElcXt6G5mOWxwZiXgGu_X6he3F--wKwg-9UBc29I/edit#), [contributing to kubectl](https://github.com/kubernetes/community/blob/master/sig-cli/CONTRIBUTING.md), [code](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubectl/cmd/config).

`mesheryctl` uses the [Cobra](https://github.com/spf13/cobra) framework. A good first-step towards contributing to `mesheryctl` would be to familiarise yourself with the [Cobra concepts](https://github.com/spf13/cobra#concepts).

For manipulating config files, `mesheryctl` uses [Viper](https://github.com/spf13/viper).

A central `struct` is maintained in the `mesheryctl/internal/cli/root/config/config.go` file. These are updated and should be used for getting the Meshery configuration.

For logs, `mesheryctl` uses [Logrus](https://github.com/sirupsen/logrus). Going through the docs and understanding the different [log-levels](https://github.com/sirupsen/logrus#level-logging) will help a lot.

`mesheryctl` uses [golangci-lint](https://github.com/golangci/golangci-lint). Refer it for lint checks.

### <a name="contributing-ui">UI Contribution Flow</a>
Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

#### Install UI dependencies
To install/update the UI dependencies:
```
make setup-ui-libs
```

#### Build and export UI
To build and export the UI code:
```
make build-ui
```

Now that the UI code is built, Meshery UI will be available at `http://localhost:9081`.
Any time changes are made to the UI code, the above code will have to run to rebuild the UI.

#### UI Development Server
If you want to work on the UI, it will be a good idea to use the included UI development server. You can run the UI development server by running the following command:
```
make run-ui-dev
```

Make sure to have Meshery server configured, up and running on the default port `http://localhost:9081` before proceeding to access and work on the UI server at `http://localhost:3000`.
Any UI changes made now will automatically be recompiled and served in the browser.

#### Running Meshery from IDE
If you want to run Meshery from IDE like Goland, VSCode. set below environment variable
```
PROVIDER_BASE_URLS="https://meshery.layer5.io"
PORT=9081
DEBUG=true
ADAPTER_URLS=mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001 mesherylocal.layer5.io:10002 mesherylocal.layer5.io:10003 mesherylocal.layer5.io:10004 mesherylocal.layer5.io:10005 mesherylocal.layer5.io:10006 mesherylocal.layer5.io:10007 mesherylocal.layer5.io:10008 mesherylocal.layer5.io:10009
```
go tool argument
```shell
-tags draft
```
update /etc/hosts
```shell
127.0.0.1 mesherylocal.layer5.io
```
