---
title: Command Line Reference
description: "A guide to Meshery's CLI: mesheryctl"
aliases:
- /reference/mesheryctl/commands/
- /reference/mesheryctl/commands
- /reference/mesheryctl/
---
## Categories and Command Structure
Meshery CLI commands are categorized by function, which are:
- [`mesheryctl`](#global-commands-and-flags) - Global flags and CLI configuration
- [`mesheryctl system`](#meshery-lifecycle-management-and-troubleshooting) - Meshery Lifecycle and Troubleshooting
- [`mesheryctl adapter`](#cloud-native-lifecycle-and-configuration-management) - Lifecycle & Configuration Management: provisioning and configuration best practices
- [`mesheryctl perf`](#cloud-native-performance-management) - Performance Management: Workload and cloud native performance characterization
- [`mesheryctl design`](#cloud-native-design-configuration-and-management) - Design Patterns: Cloud native patterns and best practices
- [`mesheryctl filter`](#data-plane-intelligence) - Data Plane Intelligence: Registry and configuration of WebAssembly filters for Envoy
- [`mesheryctl model`](#meshery-models) - A unit of packaging to define managed infrastructure and their relationships, and details specifics of how to manage them.
- [`mesheryctl component`](#meshery-components) - Fundamental building block used to represent and define the infrastructure under management
- [`mesheryctl registry`](#meshery-registry-management) - Model Database: Manage the state and contents of Meshery's internal registry of capabilities.
- [`mesheryctl environment`](#meshery-environment) - Logical group of connections and their associated credentials.
- [`mesheryctl connection`](#meshery-connection) - Managed or unmanaged resources that either through discovery or manual entry are tracked by Meshery.
- [`mesheryctl exp`](#experimental-featuresexp) - Experimental features
## Global Commands and Flags

{{< mesheryctl-command-table command="global" >}}

## Meshery Lifecycle Management and Troubleshooting
Installation, troubleshooting and debugging of Meshery and its adapters

{{< mesheryctl-command-table command="system" >}}

## Cloud Native Performance Management

{{< mesheryctl-command-table command="perf" >}}

## Cloud Native Lifecycle and Configuration Management

{{< mesheryctl-command-table command="adapter" >}}

## Cloud Native Design Configuration and Management

{{< mesheryctl-command-table command="design" >}}

## Data Plane Intelligence

{{< mesheryctl-command-table command="filter" >}}
## Meshery Registry Management

{{< mesheryctl-command-table command="registry" >}}

## Meshery Models

{{< mesheryctl-command-table command="model" >}}
## Meshery Components

{{< mesheryctl-command-table command="component" >}}

## Meshery Environment

{{< mesheryctl-command-table command="environment" >}}
## Meshery Connection

{{< mesheryctl-command-table command="connection" >}}

## Experimental Features(exp)

{{< mesheryctl-command-table command="exp" >}}
## Frequently Asked Questions for Meshery CLI 
Refer the following <a href='/guides/mesheryctl/working-with-mesheryctl'>Frequently asked questions</a> related to Meshery CLI.

## See Also 