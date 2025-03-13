# Contributing to mesheryctl's End-to-End tests

End-to-end testing of `mesheryctl` uses the [Bash Automated Testing System](https://github.com/bats-core/bats-core) (BATS) framework to define and execute CLI tests. Each test case is designed to mimic the experience that a Meshery CLI user might have while interacting with `mesheryctl` in their terminal of choice. In this sense, `mesheryctl` tests run end-to-end with each pull request submitted containing changes to either the `/mesheryctl` or the `/server` directories in the `meshery/meshery` repository, ensuring that changes included in those pull requests do not break the existing CLI functionality.


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


```bash
make server
```

> Be aware that some test cases require the availability of a Kubernetes cluster and one or more  Meshery Adapters. In those cases, please refer to the [installation guides]{{site.baseurl}}/installation) (like that of [installing Meshery on Minikube]({{site.baseurl}}/installation/kubernetes/minikube)). 


## Setup Bats Core

For Bats Core, always try to use a native OS whenever possible. The Docker-based approach is intended only for unsupported OSes and is generally not recommended because it runs on top of Ubuntu images, which can be redundant if you already using Ubuntu or Windows.


