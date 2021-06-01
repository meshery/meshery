---
layout: page
title: Contributing to Meshery CLI
permalink: project/contributing-cli
description: How to contribute to Meshery Command Line Interface.
language: en
type: project
category: contributing
---

`mesheryctl` is written in Golang or the Go Programming Language. For development use Go version 1.15+.

{% include alert.html
    type="info"
    title="Meshery CLI Reference Documents"
    content='<ul><li><a href="https://docs.google.com/spreadsheets/d/1q63sIGAuCnIeDs8PeM-0BAkNj8BBgPUXhLbe1Y-318o/edit#gid=0">Meshery Command Tracker</a>: Status of mesheryctl command implementation and platform compatibility.</li>
    <li><a href="https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#">Meshery CLI Commands and Documentation</a>: Detailed documentation of the `mesheryctl` commands.</li>
	<li><a href="https://github.com/layer5io/meshery/labels/component%2Fmesheryctl">mesheryctl open issues and pull requests</a>: Matching the "component/mesheryctl" label.</li></ul>' %}

### Mechanics of Contributing

**Building mesheryctl**

The [`/mesheryctl`](https://github.com/layer5io/meshery/tree/master/mesheryctl) folder contains the complete code for `mesheryctl`. Fork and clone the Meshery repo. `cd mesheryctl` to change directory mesheryctl's source.

After making changes, run `make` in the `mesheryctl` folder to build the binary. You can then use the binary by, say, `./mesheryctl system start`.

**Framework**

`mesheryctl` uses the [Cobra](https://github.com/spf13/cobra) framework. A good first-step towards contributing to `mesheryctl` would be to familiarise yourself with the [Cobra concepts](https://github.com/spf13/cobra#concepts). For manipulating config files, `mesheryctl` uses [Viper](https://github.com/spf13/viper).

**Model and Configuration Data**

A central `struct` is maintained in the `mesheryctl/internal/cli/root/config/config.go` file. These are updated and should be used for getting the Meshery configuration.

**Logging**

For logs, `mesheryctl` uses [Logrus](https://github.com/sirupsen/logrus). Going through the docs and understanding the different [log-levels](https://github.com/sirupsen/logrus#level-logging) will help a lot.

**Linting**

`mesheryctl` uses [golangci-lint](https://github.com/golangci/golangci-lint). See the .github/workflow/ci.yaml for syntax used during Meshery's build process.

**Unit Tests**

Unit test code coverage reports can be found - _tbd_.

# Meshery CLI Style Guide

These guidelines are a collection of principles and conventions that need to be followed while designing mesheryctl commands. `mesheryctl` might be the interface that the users first have with Meshery. As such, `mesheryctl` needs to provide a great UX.

The following principles should be taken in mind while designing `mesheryctl` commands-

## Design Principles

**1. Consistency is quality.**

- _Consistency of interaction drives a quality user experience. Whether that experience is delightful or painful is a related, but separate consideration. Meshery's behavior of user interactions should be consistent even when their user experience is poor._

**2. Intuitive user experiences feel natural to users.**

- _When being designed, each of Meshery's user experiences should be examined first from the user's perspective. Design user experiences that are familiar._

**3. Design for brevity.**

- _Avoid long commands with chained series of flags._

**4. Design with automated testing in mind.**

- _Provide possibility to specify output format as json (-o json) for easy inspection of command response._

Part of delivering a great user experience is providing intuitive interfaces. In the case of `mesheryctl` takes inspiration from and delivers similar user experiences as popular CLIs do in this ecosystem, like `kubectl` and `docker`. Here is relevant `kubectl` information to reference - [Kubectl SIG CLI Community Meeting Minutes](https://docs.google.com/document/u/2/d/1r0YElcXt6G5mOWxwZiXgGu_X6he3F--wKwg-9UBc29I/edit#), [contributing to kubectl](https://github.com/kubernetes/community/blob/master/sig-cli/CONTRIBUTING.md), [code](https://github.com/kubernetes/kubernetes/tree/master/pkg/kubectl/cmd/config).

Command structure and command behavior should be designed in such a way that they are intuitive. Users should ideally be able to understand what a command is used for without having to extensively go through the documentation. For example, `mesheryctl pattern apply -f <pattern name>` requires no further clarification as it is evident that the command will apply the pattern specified.

Consistency is key when designing intuitive interfaces. Although `mesheryctl perf run -f <performance profile name>` may sound more intuitive, users who are experienced in using the CLI will prefer the consistant verb `apply` over `run`. This will also ensure a consistent command language making memorizing easier.

**Flags**

Consistency should also be enforced when chaining commands and using flags. For example, if `mesheryctl pattern` has a `list` and `view` command and has an `-all` and `--output` flag, then, similar commands like `mesheryctl perf` should also support the same commands and flags and provide a consistent user experience.

### Rational defaults overriden with flags

Default behaviour should be optimised for what users will need to do most of the time.

These assumed defaults should be easily overriden by the user with flags.

For example, `mesheryctl system context create <context name>` assumes a default platform for the created context. But this can be easily overriden with the `--platform` flag.

### User Experience: GUI vs CLI

Ideally, all functionaly provided in Meshery UI should be available to users via CLI (in `mesherctl`). Meshery strives for parity of functionality between it's two clients. For example, viewing a performance profile in the GUI and with `mesheryctl system perf view <profile name>` in the CLI should show the same data.

Command line interfaces offer less context to the user, which makes them inherently less intuitive compared to graphical user interfaces. Both of these user interfaces, however, are the same in that they are both clients of Meshery Server. Both clients _are_ a user experience and as such, to be attended to in this way. The following considerations should be accounted for when designing `mesheryctl` experiences:

- Provide only relevant output. Use "debug" logs that can be accessed with `--verbose` flag wherever needed.
- Add headers to output to give context to the user.
- As mentioned [above](#intuition-vs-consistency), similar commands should behave similarly.
- Confirm steps for risky commands. For example, use the `AskForConfirmation` function which will prompt the user to type in "yes" or "no" to continue.
- Anticipate user actions. If the user creates a new context with `mesheryctl system context create` then the next action might be `mesheryctl system start` to start Meshery ot `mesheryctl system context switch` to switch context names.
- Anticipate user errors. For example, if the user types in `mesheryctl system satrt`, using the inbuilt features with the [cobra library](https://github.com/spf13/cobra), we can correct it to `start` and alert the user.

## Designing Commands

If you are working on a new command or adding a new feature on an existing command, it is recommended to setup a design spec so that other contributors can weigh in on the design before you start to code. Broader features should have a design spec made in Google Doc using [this template](https://drive.google.com/drive/folders/1KHtJc4ToklBQ_UUsDgAL2sVZNhOQGzbh). For small changes, communicating over the [issue tracker](https://github.com/layer5io/meshery/issues) or the [discussions](https://github.com/layer5io/meshery/discussions) will be helpful.

When designing for the command line interface, ask and consider the following questions.

##### What the command does

- What makes sense to do from a terminal? What doesn’t?
- What might people want to automate?
- What is the default behavior? What flags might you need to change that behavior?
- What might people try and fail to do and how can you anticipate that?

###### What the command is called

- What should be the command language? (`mesheryctl <command> <subcommand> [args] [flags] [value]`)
- What should be a command vs a flag?
- How can you align the language of the new command with the existing commands?

##### What the command outputs

- How can you make the GUI and the CLI outputs similar?
- What should be outputted normally and what falls into debug logs?

##### How you explain your command

You will need to provide a short and long description of the command for the help pages and also for the Meshery Documentation.

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
