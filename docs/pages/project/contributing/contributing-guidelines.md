---
layout: default
title: "Meshery CLI Contributing Guidelines"
abstract: "Design principles and code conventions."
permalink: project/contributing/contributing-cli-guide
type: project
language: en
category: contributing
list: include
---

These guidelines are a collection of principles and conventions that need to be followed while designing mesheryctl commands.

## Design Principles

**1. Consistency is quality.**

- _Consistency of interaction drives a *quality* user experience. Whether that experience is delightful or painful is a related, but separate consideration. Meshery's behavior of user interactions should be consistent even when their user experience is poor._

**2. Intuitive user experiences feel natural to users.**

- _When being designed, each of Meshery's user experiences should be examined first from the user's perspective._

Commands should be designed in such a way that they are intuitive. Users should ideally be able to understand what a command is used for without having to extensively go through the documentation.

For example, `mesheryctl pattern apply -f <pattern name>` requires no further clarification as it is evident that the command will apply the pattern specified.

Consistency is key when designing intuitive interfaces. Although `mesheryctl perf run -f <performance profile name>` may sound more intuitive, users who are experienced in using the CLI will prefer the consistant verb `apply` over `run`. This will also ensure a consistent command language making memorizing easier. Consistency should also be enforced when chaining commands and using flags.

For example, if `mesheryctl pattern` has a `list` and `view` command and has an `-all` and `--output` flag, then, similar commands like `mesheryctl perf` should also support the same commands and flags and provide a consistent user experience.

### Rational defaults overriden with flags

Default behaviour should be optimised for what users will need to do most of the time.

These assumed defaults should be easily overriden by the user with flags.

For example, `mesheryctl system context create <context name>` assumes a default platform for the created context. But this can be easily overriden with the `--platform` flag.

### User Experience: GUI vs CLI

Ideally, all functionality provided in Meshery UI should be available to users via CLI (in `mesherctl`). Meshery strives for parity of functionality between it's two clients. For example, viewing a performance profile in the GUI and with `mesheryctl system perf view <profile name>` in the CLI should show the same data.

Command line interfaces offer less context to the user, which makes them inherently less intuitive compared to graphical user interfaces. Both of these user interfaces, however, are the same in that they are both clients of Meshery Server. Both clients _are_ a user experience and as such, to be attended to in this way. The following considerations should be accounted for when designing `mesheryctl` experiences:

- Provide only relevant output. Use "debug" logs that can be accessed with `--verbose` flag wherever needed.
- Add headers to output to give context to the user.
- As mentioned [above](#intuition-vs-consistency), similar commands should behave similarly.
- Confirm steps for risky commands. For example, use the `AskForConfirmation` function which will prompt the user to type in "yes" or "no" to continue.
- Anticipate user actions. If the user creates a new context with `mesheryctl system context create` then the next action might be `mesheryctl system start` to start Meshery ot `mesheryctl system context switch` to switch context names.
- Anticipate user errors. For example, if the user types in `mesheryctl system satrt`, using the inbuilt features with the [cobra library](https://github.com/spf13/cobra), we can correct it to `start` and alert the user.

## Process

### Contributor-facing Design Artifacts

The following table outlines a set of design artifacts produced in the process of creating and implementing functional specifications.

| **Document** | **Purpose** |
| [Meshery CLI Commands & Documentation](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#) | Overarching index of all mesheryctl commands |
| [Meshery CLI Design Specification](https://docs.google.com/document/d/1Iw88bZEL_fWeajxHh0BgmhuXGKl455P4CF-vm1yCkU8/edit#) - See [Creating a Functional Specification](https://docs.google.com/document/d/1RP3IWLc-MiQS-QYasqCoVuCH7--G87p5ezE5f_nOzB8/edit) | Breakout of individual command purpose, syntax, flags, and behavior. Links to these separate, breakout documents should be pasted into the Meshery CLI Commands & Documentation doc. |
| [Mesheryctl Command Matrix](https://docs.google.com/spreadsheets/d/1q63sIGAuCnIeDs8PeM-0BAkNj8BBgPUXhLbe1Y-318o/edit#gid=0) | A spreadsheet of all commands, subcommands and flags with an indication of their implementation status would be most helpful for both documentation and planning purposes. |

### User-facing Documentation Artifacts

The following table outlines a set of design artifacts produced in the process of creating and implementing functional specifications.

| **Document** | **Purpose** |
| [mesheryctl Command Reference](https://docs.meshery.io/reference/mesheryctl) | Overarching index in the form of succinct, categorized, and tabularly displayed reference of all commands and their syntax with simple example usage. |
| `mesheryctl <name of command>` | Breakout of individual command purpose, syntax, flags, and behavior. Links to these separate, breakout documents should be available from the mesheryctl Command Reference. |

## When designing for the CLI, consider...

### What the command does

What makes sense to do from a terminal? What doesn’t?
What might people want to automate?
What is the default behavior? What flags might you need to change that behavior?
What might people try and fail to do and how can you anticipate that?

### What the command is called

What should be the command language? (`mesheryctl <command> <subcommand> [args] [flags] [value]`)
What should be a command vs a flag?
How can you align the language of the new command with the existing commands?

### What the command outputs

How can you make the GUI and the CLI outputs similar?
What should be outputted normally and what falls into debug logs?

### How you explain your command

You will need to provide a short and long description of the command for the help pages and also for the Meshery Documentation.

## Proposing new features / creating a design specification

If you are working on a new command or adding a new feature on an existing command, it is recommended to setup a design spec so that other contributors can weigh in on the design before you start to code.

Broader features should have a design specification created in the form of a Google Doc. Refer to "[Process: Creating a Functional Specification or other Document](https://docs.google.com/document/d/1RP3IWLc-MiQS-QYasqCoVuCH7--G87p5ezE5f_nOzB8/edit?usp=sharing)", which has detailed information on creating a design specification.

For small changes, communicating over the [issue tracker](https://github.com/meshery/meshery/issues) or the [discussions](https://github.com/meshery/meshery/discussions) will be helpful.


