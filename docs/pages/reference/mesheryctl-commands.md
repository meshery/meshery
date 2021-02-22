---
layout: default
title: mesheryctl Command Reference
abstract: "A guide to Meshery's CLI: mesheryctl"
permalink: reference/mesheryctl
redirect_from: 
 - guides/mesheryctl
 - guides/mesheryctl-commands
type: Reference
language: en
lang: en
categories: en
---
## Categories and Command Structure

Meshery CLI commands are categorized by function, which are:

- `mesheryctl` - Global flags
- `mesheryctl system` - Meshery Lifecycle and Troubleshooting
- `mesheryctl mesh` - Service Mesh Lifecycle & Configuration Management
- `mesheryctl perf` -  Service Mesh Performance Management

## Global Commands and Flags

| command    |     flag      |   function  |   Usage   |
| :--------- | :-----------: | :---------- | :-------- |{% for command_hash in site.data.mesheryctlcommands.global.commands %}{% assign command = command_hash[1] %}
| {{ command.name }} |       | {{ command.description }} | {{command.usage}} |{% endfor %}{% for subcommand_hash in site.data.mesheryctlcommands.global.subcommands %}{% assign subcommand = subcommand_hash[1] %}
|          | {{ subcommand.name }}  | {{ subcommand.description }} | {{ subcommand.usage }} |{% endfor %}

## Meshery Lifecycle Management

Installation, troubleshooting and debugging of Meshery and its adapters.

| command    |    arg   |      flag         | function  |  Usage  |
| :--------- | :------: | :---------------: | :-------- | :-----: |{% for command_hash in site.data.mesheryctlcommands.lifecycle.commands %}{% assign command = command_hash[1] %}
| {{ command.name }}  |     |  {{ command.flag }}  | {{ command.description }} | {{ command.usage }} |{% endfor %}{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.subcommands %}{% assign subcommand = subcommand_hash[1] %}
|         | {{ subcommand.name }}  |    | {{ subcommand.description }} | {{ subcommand.usage }} |{% endfor %}{% for flag_hash in site.data.mesheryctlcommands.lifecycle.flags %}{% assign flag = flag_hash[1] %}
| {{ flag.name }}  |  {{flag.arg}}    |  {{flag.flag}}  | {{ flag.description }} | {{ flag.usage }} |{% endfor %}

## Performance Management

| command    |      flag         | function  | Usage |
| :--------- | :-----------------: | :---------- | :----------- |{% for command_hash in site.data.mesheryctlcommands.performance.commands %}{% assign command = command_hash[1] %}
| {{ command.name }} |       | {{ command.description }} | {{ command.usage }}  |{% endfor %}{% for flag_hash in site.data.mesheryctlcommands.performance.flags %}{% assign flag = flag_hash[1] %}
|          | {{ flag.name }}  | {{ flag.description }} | {{ flag.usage }} |{% endfor %}

## Service Mesh Lifecycle Management

| command    |    arg   |      flag     | function  |   Usage    |
| :--------- | :-------:| :-----------: | :-------- | :--------- |{% for command_hash in site.data.mesheryctlcommands.meshes.commands %}{% assign command = command_hash[1] %}
| {{ command.name }} |       |      | {{ command.description }} |     |{% endfor %}{% for subcommand_hash in site.data.mesheryctlcommands.meshes.subcommands %}{% assign subcommand = subcommand_hash[1] %}
|        | {{ subcommand.name }}  |    | {{ subcommand.description }} | {{ subcommand.usage }} |{% endfor %}{% for flag_hash in site.data.mesheryctlcommands.meshes.flags %}{% assign flag = flag_hash[1] %}
|        |        | {{ flag.name }}  | {{ flag.description }} | {{ flag.usage }} |{% endfor %}
