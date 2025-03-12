# Contributing to mesheryctl's End-to-End tests

To automate ne-to-end testing mesehryctl uses [bats](https://github.com/bats-core/bats-core) framework to automate cli testing. End-to-end run to ensure that the changes do not break te existing functionnality.


## Prerequisites:

Before diving into mesheryctl's testing environment, certain prerequisites are necessary:

- A verified account in your chosen provider which integrate with Meshery.
- bash as shell terminal
- Installations of Golang, NodeJS and Makefiles for Native OS build (Optional for docker based build).
- Kubernetes clusters (Required for connection to Kubernetes test cases)
- Meshery server up and running

## Authtentication

To run the tests successfully, you need be logged. 

### Check log in status

```bash
mesheryctl system check
```

If you need to authenticate, you will see the following message

> Error: !! Authentication token not found. Please supply a valid user token. Login with `mesheryctl system login`

**Accessing Remote Providers**

> In the case you are using Layer5 Cloud as a remote provider, you can <a href="https://cloud.layer5.io/security/tokens">generate a token from your user account</a> > to use while writing and executing tests


## Starting up Server

There are a few ways to set up the Meshery server, but for end-to-end testing, we aim to get as close to a production environment as possible. We know developers might need to make some tweaks for Server. Rebuilding the whole project can take time, and we don’t support hot reload because it’s more for development than for end-to-end testing.


> Some test cases required you to have kubernetes cluster and build meshery adapter as well, be aware of that. Which is out of scope for this documentation<ul><li><a href="https://docs.meshery.io/installation/kubernetes/minikube">Kubernetes Cluster</a>: Installation of kubernetes cluster with Minikube.</li>
> <li><a href="https://docs.meshery.io/installation/multiple-adapters">Meshery Adapters</a>: Using Multiple Adapters</li></ul>

### Native OS Build (Recommended)

This approach is very quick to build, but also dependent on your operating system, so you need to have all dependencies necessary to be able compile and running the server.

- Compile the Golang into binary file for Meshery Server

```bash
make build-server
```

- Run the Meshery Server on localhost port 9081

```bash
make server-binary
```

### Docker Based Build

Alternatively, a Docker-based setup can be utilized, simplifying the process, and ensuring consistency across different environments. It is closer to the production environment than the native solution but slower in terms of build time.

- Build the docker container locally:

```bash
make docker-testing-env-build
```

- Run the docker container on port 9081

```bash
make docker-testing-env
```

## Setup Bats Core

For Bats Core, always try to use a native OS whenever possible. The Docker-based approach is intended only for unsupported OSes and is generally not recommended because it runs on top of Ubuntu images, which can be redundant if you already using Ubuntu or Windows.

