# <a name="contributing">Contributing</a>
Please do! Contributions, updates, [issues](/../../issues) and [pull requests](/../../pulls) are welcome. This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

Not sure where to start? See the [newcomers welcome guide](https://docs.google.com/document/d/14Fofs9BysojB5igihXBI_SsFWoSUu-QRsGnnFqUvR0M/edit) for how, where and why to contribute. Or grab an open issue with the [help-wanted label](../../labels/help%20wanted) and jump in.

# <a name="contributing">General Contribution Flow</a>
Please contribute! The Meshery documentation uses GitHub Pages to host the docs site. Learn more about [Meshery's documentation framework](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit?usp=sharing). The process of contributing to documentation follows this flow:

1. Get a local copy of the documentation.
`git clone https://github.com/layer5io/meshery`
1. Navigate to the docs folder.
`cd docs`
1. Create and checkout a new branch to make changes within
`git checkout -b <my-changes>`
1. Edit/add documentation.
`vi <specific page>.md`
1. Run site locally to preview changes.
`make site`
1. Commit and push changes to your remote branch.
`git push origin <my-changes>`
1. Open a pull request (in your web browser) against the master branch on https://github.com/layer5io/meshery.

# <a name="adapter">Writing an Adapter</a>
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

# <a name="building">Building Meshery</a>
Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

__Please note__: All `make` commands should be run in a terminal from within the Meshery's main folder.

## Prerequisites for building Meshery in your development environment:
1. `Go` version 1.11+ installed if you want to build and/or make changes to the existing code.
1. `GOPATH` environment variable should be configured appropriately
1. `npm` and `node` should be installed your machine, preferrably the latest versions.
1. Clone this repository (`git clone https://github.com/layer5io/meshery.git`), preferrably outside `GOPATH`. If you happen to checkout Meshery inside your `GOPATH`, please set an environment variable `GO111MODULE=on` to enable GO Modules.


## Building and packaging

### Build and run Meshery server
To build & run the Meshery server code, run the following command:
```
make run-local
```

Any time changes are made to the GO code, you will have to stop the server and run the above command again.
Once the Meshery server is up and running, you should be able to access Meshery on your `localhost` on port `9081` at `http://localhost:9081`. One thing to note, you might NOT see the Meshery UI until the UI code is built as well.

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

### Building Docker image
To build a Docker image of Meshery please ensure you have `Docker` installed to be able to build the image. Now, run the following command to build the Docker image:
```
make docker
```
### License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
