# mesheryctl

`mesheryctl` is the CLI client for Meshery.

# Contributing

Please refer the [Meshery Contributing Guidelines](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md) for setting up your development environment.

Refer the [mesheryctl- Command Reference and Tracker](https://docs.google.com/spreadsheets/d/1q63sIGAuCnIeDs8PeM-0BAkNj8BBgPUXhLbe1Y-318o/edit#gid=0) for current status of `mesheryctl`.

For a quick introduction to `mesheryctl`, checkout [Beginner's guide to contributing to Meshery and mesheryctl](https://youtu.be/hh_kFLZx3G4).

## Building and running `mesheryctl`

The [`/mesheryctl`](https://github.com/meshery/meshery/tree/master/mesheryctl) folder contains the complete code for `mesheryctl`.

`mesheryctl` is written in Golang or the Go Programming Language. For development use Go version 1.15+.

After making changes, run `make` in the `mesheryctl` folder to build the binary. You can then use the binary by, say, `./mesheryctl system start`.

Refer the [Meshery CLI Commands and Documentation](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#heading=h.5fucij4hc5wt) for a complete reference of `mesheryctl`.

## `mesheryctl` command reference

Detailed documentation of the `mesheryctl` commands is available in the [Meshery Docs](https://docs.meshery.io/reference/mesheryctl).

## General guidelines and resources

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

All contributors are invited to review [pull requests](https://github.com/meshery/meshery/pulls) on `mesheryctl` as on other Layer5 projects.
