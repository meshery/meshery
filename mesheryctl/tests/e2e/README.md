# Contributing to mesheryctl's End-to-End tests

End-to-end testing of `mesheryctl` uses the [Bash Automated Testing System](https://github.com/bats-core/bats-core) (BATS) framework to define and execute CLI tests. Each test case is designed to mimic the experience that a Meshery CLI user might have while interacting with `mesheryctl` in their terminal of choice. In this sense, `mesheryctl` tests run end-to-end with each pull request submitted containing changes to either the `/mesheryctl` or the `/server` directories in the `meshery/meshery` repository, ensuring that changes included in those pull requests do not break the existing CLI functionality.


## Prerequisites

Before diving into mesheryctl's testing environment, certain prerequisites are necessary:

- A verified account in your chosen provider which integrate with Meshery.
- bash as shell terminal
- Installations of Golang, NodeJS and Makefiles for Native OS build (Optional for docker based build).
- Kubernetes clusters (Required for connection to Kubernetes test cases)
- Meshery server up and running

### Authtentication

To run the tests successfully, you need be logged. 

### Check log in status

```bash
mesheryctl system check
```

If you need to authenticate, you will see the following message

> Error: !! Authentication token not found. Please supply a valid user token. Login with `mesheryctl system login`

**Accessing Remote Providers**

> In the case you are using Layer5 Cloud as a remote provider, you can <a href="https://cloud.layer5.io/security/tokens">generate a token from your user account</a> > to use while writing and executing tests


### Starting up Server

There are a few ways to set up the Meshery server, but for end-to-end testing, we aim to get as close to a production environment as possible. We know developers might need to make some tweaks for Server. Rebuilding the whole project can take time, and we don’t support hot reload because it’s more for development than for end-to-end testing.


```bash
make server
```

> Be aware that some test cases require the availability of a Kubernetes cluster and one or more  Meshery Adapters. In those cases, please refer to the [installation guides]{{site.baseurl}}/installation) (like that of [installing Meshery on Minikube]({{site.baseurl}}/installation/kubernetes/minikube)). 


### Setup Bats Core

For Bats Core, always try to use a native OS whenever possible. 

You can follow the following steps which is based on the official installation [documentation](https://bats-core.readthedocs.io/en/stable/installation.html)

#### MacOS (homebrew)

```bash
 $ brew install bats-core
```

#### Any OS (npm)

```bash
 $ npm install -g bats
```

#### Windows (from source via bash)

Check out a copy of the Bats repository and install it to $HOME. This will place the bats executable in $HOME/bin, which will need to be added in $PATH.

```bash
$ git clone https://github.com/bats-core/bats-core.git
$ cd bats-core
$ ./install.sh $HOME
```

#### Bats libraries

Some tests could use bats libraires as helpers to create the tests suite. 

## End-to-End Tests

Once all prerequisites have been met, you are setup to run, add, update new test cases.

To keep the development of tests consistents and maintainable, we have put some guidance on implementation process which are describe above.

### Folder structure 

The tests will be available in `mesheryctl/tests/e2e` folder with the following structure which is donedone to follow the code base structure of `mesheryctl`

```
├── aaa-<command>
├── bbb-<command>
├── ...
├── yyy-<command>
├── zzz-<command>
├── helpers
├── README.md
├── run_tests.bash
├── setup_suite.bash
└── teardown_suite.bash
```


### Naming convention

> **Tips**
>
> Bats will run the tests sequentially so we use a prefix code alphanumeric to be able to order tests as required  when necessary
>
> ex: you need to import a model before being able to view the details of that model



#### Folder

- **prefix**: 3 digits alphanumeric code
- **name**: command under test
  
**example**: `011-design`

#### Test file

- **prefix**: 2 digits numeric code
- **name**: subcommand under test
- **extension**: bats

**example**: `05-list.bats`

> For consistency, we will keep the prefix *00-* for the command under test in the folder and subcommands will start at *01-*

Here an example with `mesheryctl model` command

```bash
002-model/
├── 00-model.bats
└── 01-model-list.bats
```

### Run End-to-End (locally)

Move in mesheryctl folder and run 

<!-- 
    TODO: Add make e2e supportwith following changes
    1. move to tests/e2e/helpers
    2. Assert the following lbats libraries are available if not git clone
    - https://github.com/bats-core/bats-file.git
    - https://github.com/bats-core/bats-detik.git
    - https://github.com/bats-core/bats-support.git
    3. back to tests/e2e
    5. run bats *-*/*.bats
-->
```bash
make e2e (Not yet implemented)
```

