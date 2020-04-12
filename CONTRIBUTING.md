# <a name="contributing">Contributing Overview</a>
Please do! Thanks for your help improving the project! :balloon:

All contributors are welcome. Please see the [newcomers welcome guide](https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit) for how, where and why to contribute. This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

Not sure where to start? First, see the [newcomers welcome guide](https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit). Grab an open issue with the [help-wanted label](../../labels/help%20wanted) and jump in. Join the [Slack account](http://slack.layer5.io) and engage in conversation. Create a [new issue](/../../issues/new/choose) if needed.  All [pull requests](/../../pulls) should reference an open [issue](/../../issues). Include keywords in your pull request descriptions, as well as commit messages, to [automatically close issues in GitHub](https://help.github.com/en/github/managing-your-work-on-github/closing-issues-using-keywords).

**Sections**
- <a name="contributing">General Contribution Flow</a>
  - <a href="#commit-signing">Developer Certificate of Origin</a>
- Meshery Contribution Flow
  - <a href="#contributing-docs">Meshery Documentation</a>
  - <a href="#contributing-meshery">Meshery Backend</a>
    - <a href="#adapter">Writing a Meshery Adapter</a>
  - <a href="#contributing-ui">Meshery UI</a>
Relevant coding style guidelines are the Go Code Review Comments and the Formatting and style section of Peter Bourgon's Go: Best Practices for Production Environments.
# <a name="contributing">General Contribution Flow</a>

In order to contribute to Meshery, please follow the fork-and-pull request workflow described [here](./git-workflow.md).

## <a name="commit-signing">Signing-off on Commits (Developer Certificate of Origin)</a>

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
  
## <a name="contributing-docs">Documentation Contribution Flow</a>
Please contribute! Meshery documentation uses GitHub Pages to host the docs site. Learn more about [Meshery's documentation framework](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit?usp=sharing). The process of contributing follows this flow:

1. Create a fork, if you have not already, by following the steps described [here](./git-workflow.md)
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


## <a name="contributing-meshery">Meshery Contribution Flow</a>
Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

Relevant coding style guidelines are the [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) and the _Formatting and style_ section of Peter Bourgon's [Go: Best 
Practices for Production Environments](https://peter.bourgon.org/go-in-production/#formatting-and-style).

__Please note__: All `make` commands should be run in a terminal from within the Meshery's main folder.

### Prerequisites for building Meshery in your development environment:
1. `Go` version 1.11+ installed if you want to build and/or make changes to the existing code.
1. `GOPATH` environment variable should be configured appropriately
1. `npm` and `node` should be installed your machine, preferably the latest versions.
1. Fork this repository (`git clone https://github.com/layer5io/meshery.git`), clone your forked version of Meshery to your local, preferably outside `GOPATH`. If you happen to checkout Meshery inside your `GOPATH` and you have version of `Go` prior to version 1.13, please set an environment variable `GO111MODULE=on` to enable GO Modules.

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

_Tip:_ The [Meshery adapter for Istio](https://github.com/layer5io/meshery-istio) is a good reference adapter to use as an example of a Meshery adapter written in Go.

## <a name="contributing-ui">UI Contribution Flow</a>
Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

### Install UI dependencies
To install/update the UI dependencies:
```
make setup-ui-libs
```

### Build and export UI
To build and export the UI code:
```
make build-ui
```

Now that the UI code is built, Meshery UI will be available at `http://localhost:9081`.
Any time changes are made to the UI code, the above code will have to run to rebuild the UI.

### UI Development Server
If you want to work on the UI, it will be a good idea to use the included UI development server. You can run the UI development server by running the following command:
```
make run-ui-dev
```

Once you have the server up and running, you will be able to access the Meshery UI at `http://localhost:3000`. One thing to note is that for the UI dev server to work, you need Meshery server running on the default port of `9081`.
Any UI changes made now will automatically be recompiled and served in the browser.

### Running Meshery from IDE
If you want to run Meshery from IDE like Goland, VSCode. set below environment variable
```
SAAS_BASE_URL=https://meshery.layer5.io
PORT=9081
DEBUG=true
ADAPTER_URLS=mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001 mesherylocal.layer5.io:10002 mesherylocal.layer5.io:10003 mesherylocal.layer5.io:10004
```
go tool argument
```shell
-tags draft
```
update /etc/hosts
```shell
127.0.0.1 mesherylocal.layer5.io
```

# <a name="maintaining"> Reviews</a>
All contributors are invited to review pull requests. See this short video on [how to review a pull request](https://www.youtube.com/watch?v=isLfo7jfE6g&feature=youtu.be).

# New to Git?
Resources: https://lab.github.com and https://try.github.com/

### License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
