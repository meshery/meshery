---
layout: default
title: mesheryctl Command Reference
abstract: "A guide to Meshery's CLI: mesheryctl"
permalink: reference/mesheryctl
redirect_from: 
#  - guides/mesheryctl
#  - guides/mesheryctl-commands
type: Reference
---
## Categories and Command Structure

Meshery CLI commands are categorized by function, which are:

- `mesheryctl` - Global flags
- `mesheryctl system` - Meshery Lifecycle and Troubleshooting
- `mesheryctl mesh` - Service Mesh Lifecycle & Configuration Management
- `mesheryctl perf` -  Service Mesh Performance Management
- `mesheryctl pattern` - Service Mesh Pattern Configuration & Management

## Global Commands and Flags

<table>
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% for command_hash in site.data.mesheryctlcommands.global.commands %}{% assign command = command_hash[1] %}
    <tr>
      <td rowspan=6>{{ command.name }}</td>
      <td></td>
      <td></td>
      <td>{{ command.description }}</td>
    </tr>
    {% for subcommand_hash in site.data.mesheryctlcommands.global.subcommands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
         <td>{{ subcommand.name }}</td>
         <td></td>
         <td>{{ subcommand.description }}</td>
      </tr>
      {% for flag_hash in site.data.mesheryctlcommands.global.flags %}{% assign flag = flag_hash[1] %}
        <tr>
         <td></td>
         <td>{{ flag.name }}</td>
         <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% endfor %}
  {% endfor %}
</thead>
</table>

## Meshery Lifecycle Management and Troubleshooting

Installation, troubleshooting and debugging of Meshery and its adapters.

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Arguments</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% for command_hash in site.data.mesheryctlcommands.lifecycle.system.commands %}{% assign command = command_hash[1] %}
    <tr>
      <td rowspan=15>{{ command.name }}</td>
      <td></td>
      <td></td>
      <td>{{ command.description }}</td>
    </tr>
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.flag }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% for command_hash in site.data.mesheryctlcommands.lifecycle.system.start.command %}{% assign command = command_hash[1] %}
    <tr>
      <td rowspan=4>{{ command.name }}</td>
      <td></td>
      <td>{{ command.description }}</td>
    </tr>
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.start.flag %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.flag }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
  {% endfor %}
  {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.stop.command %}{% assign subcommand = subcommand_hash[1] %}
    <tr>
      <td rowspan=2>{{ subcommand.name }}</td>
      <td></td>
      <td>{{ subcommand.description }}</td>
    </tr>
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.stop.flag %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.flag }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
  {% endfor %}
    {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.subcommands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td>{{ subcommand.name }}</td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
    {% endfor %}
  {% endfor %}
  {% for command_hash in site.data.mesheryctlcommands.lifecycle.system-channel.commands %}{% assign command = command_hash[1] %}
        <tr>
          <td rowspan=5>{{ command.name }}</td>
          <td></td>
          <td></td>
          <td>{{ command.description }}</td>
        </tr>
        {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.subcommands %}{% assign subcommand = subcommand_hash[1] %}
          <tr>
            <td>{{ subcommand.name }}</td>
            <td></td>
            <td>{{ subcommand.description }}</td>
          </tr>
        {% endfor %}
          {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-channel.flags %}{% assign flag = flag_hash[1] %}
            <tr>
              <td>{{ flag.arg }}</td>
              <td>{{ flag.flag }}</td>
              <td>{{ flag.description }}</td>
            </tr>
          {% endfor %}
    {% endfor %}
    {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-channel.view.command %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td rowspan=2>{{ subcommand.name }}</td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
    {% endfor %}
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-channel.view.flag %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.flag }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% for command_hash in site.data.mesheryctlcommands.lifecycle.system-context.commands %}{% assign command = command_hash[1] %}
      <tr>
        <td rowspan=10>{{ command.name }}</td>
        <td></td>
        <td></td>
        <td>{{ command.description }}</td>
      </tr>
      {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.create.command %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td rowspan=4>{{ subcommand.name }}</td>
          <td></td>
          <td>{{ subcommand.description }}</td>
        </tr>
      {% endfor %}
      {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.create.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.flag }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
      {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.subcommands %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td>{{ subcommand.name }}</td>
          <td></td>
          <td>{{ subcommand.description }}</td>
        </tr>
      {% endfor %}
      {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.arg }}</td>
          <td>{{ flag.flag }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
      {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.view.command %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td rowspan=3>{{ subcommand.name }}</td>
          <td></td>
          <td>{{ subcommand.description }}</td>
        </tr>
      {% endfor %}
      {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.view.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.flag }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
  {% endfor %}
</thead>
</table>

## Service Mesh Performance Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% for command_hash in site.data.mesheryctlcommands.performance.commands %}{% assign command = command_hash[1] %}
    <tr>
      <td rowspan=11>{{ command.name }}</td>
      <td></td>
      <td>{{ command.description }}</td>
    </tr>
    {% for flag_hash in site.data.mesheryctlcommands.performance.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
  {% endfor %}
</thead>
</table>


## Service Mesh Lifecycle and Configuration Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% for command_hash in site.data.mesheryctlcommands.meshes.commands %}{% assign command = command_hash[1] %}
    <tr>
      <td rowspan=10>{{ command.name }}</td>
      <td></td>
      <td></td>
      <td>{{ command.description }}</td>
    </tr>
  {% endfor %}
    {% for subcommand_hash in site.data.mesheryctlcommands.meshes.init.commands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td rowspan=4>{{ subcommand.name }}</td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
    {% endfor %}
      {% for flag_hash in site.data.mesheryctlcommands.meshes.init.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% for subcommand_hash in site.data.mesheryctlcommands.meshes.validate.commands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td rowspan=5>{{ subcommand.name }}</td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
    {% endfor %}
      {% for flag_hash in site.data.mesheryctlcommands.meshes.validate.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
</thead>
</table>

## Service Mesh Pattern Configuration and Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% for command_hash in site.data.mesheryctlcommands.pattern.commands %}{% assign command = command_hash[1] %}
    <tr>
      <td rowspan=4>{{ command.name }}</td>
      <td></td>
      <td></td>
      <td>{{ command.description }}</td>
    </tr>
    {% for flag_hash in site.data.mesheryctlcommands.pattern.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
      {% for subcommand_hash in site.data.mesheryctlcommands.pattern.subcommands %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td>{{ subcommand.name }}</td>
          <td></td>
          <td>{{ subcommand.description }}</td>
        </tr>
      {% endfor %}
    {% endfor %}
  {% endfor %}
</thead>
</table>