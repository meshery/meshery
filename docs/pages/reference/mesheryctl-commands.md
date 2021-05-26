---
layout: default
title: mesheryctl Command Reference
abstract: "A guide to Meshery's CLI: mesheryctl"
permalink: reference/mesheryctl
redirect_from: reference/mesheryctl/commands/
type: Reference
---
## Categories and Command Structure

Meshery CLI commands are categorized by function, which are:

- `mesheryctl` - Global flags and CLI configuration
- `mesheryctl system` - Meshery Lifecycle and Troubleshooting
- `mesheryctl mesh` - Service Mesh Lifecycle & Configuration Management: provisioning and configuration best practices
- `mesheryctl perf` -  Service Mesh Performance Management: Workload and service mesh performance characterization
- `mesheryctl pattern` - Service Mesh Pattern Configuration & Management: Service mesh patterns and Open Application Model integration
- `mesheryctl fitler` - Data Plane Intelligence: Registry and configuration of WebAssembly filters for Envoy

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
      <td rowspan=6><a href="{{ site.baseurl }}/reference/mesheryctl/mesheryctl">{{ command.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command.description }}</td>
    </tr>
    {% for subcommand_hash in site.data.mesheryctlcommands.global.subcommands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
         <td><a href="{{ site.baseurl }}/reference/mesheryctl/mesheryctl/version">{{ subcommand.name }}</a></td>
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
      <td rowspan=16><a href="{{ site.baseurl }}/reference/mesheryctl/system">{{ command.name }}</a></td>
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
        <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/start">{{ command.name }}</a></td>
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
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/stop">{{ subcommand.name }}</a></td>
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
  {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.completion.command %}{% assign subcommand = subcommand_hash[1] %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/completion">{{ subcommand.name }}</a></td>
      <td></td>
      <td>{{ subcommand.description }}</td>
    </tr>
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.completion.flag %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.flag }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
  {% endfor %}
  {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.update.command %}{% assign subcommand = subcommand_hash[1] %}
    <tr>
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/update">{{ subcommand.name }}</a></td>
      <td></td>
      <td>{{ subcommand.description }}</td>
    </tr>
    {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.update.flag %}{% assign flag = flag_hash[1] %}
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
          <td rowspan=5><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel/">{{ command.name }}</a></td>
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
        <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel/view">{{ subcommand.name }}</a></td>
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
        <td rowspan=10><a href="{{ site.baseurl }}/reference/mesheryctl/system/context">{{ command.name }}</a></td>
        <td></td>
        <td></td>
        <td>{{ command.description }}</td>
      </tr>
      {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.create.command %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/create">{{ subcommand.name }}</a></td>
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
      {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.delete.command %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/delete">{{ subcommand.name }}</a></td>
          <td></td>
          <td>{{ subcommand.description }}</td>
        </tr>
      {% endfor %}
      {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.view.command %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/view">{{ subcommand.name }}</a></td>
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
      {% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system-context.switch.command %}{% assign subcommand = subcommand_hash[1] %}
        <tr>
          <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/switch">{{ subcommand.name }}</a></td>
          <td></td>
          <td>{{ subcommand.description }}</td>
        </tr>
      {% endfor %}
      {% for flag_hash in site.data.mesheryctlcommands.lifecycle.system-context.switch.flags %}{% assign flag = flag_hash[1] %}
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
      <td rowspan=11><a href="{{ site.baseurl }}/reference/mesheryctl/perf">{{ command.name }}</a></td>
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
      <td rowspan=10><a href="{{ site.baseurl }}/reference/mesheryctl/mesh">{{ command.name }}</a></td>
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
      <td rowspan=9><a href="{{ site.baseurl }}/reference/mesheryctl/pattern">{{ command.name }}</a></td>
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
    {% endfor %}
    {% for subcommand_hash in site.data.mesheryctlcommands.pattern.list.commands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/pattern/list">{{ subcommand.name }}</a></td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
      {% for flag_hash in site.data.mesheryctlcommands.pattern.list.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
    {% endfor %}
    {% for subcommand_hash in site.data.mesheryctlcommands.pattern.apply.commands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/pattern/apply">{{ subcommand.name }}</a></td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
      {% for flag_hash in site.data.mesheryctlcommands.pattern.apply.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
    {% endfor %}
    {% for subcommand_hash in site.data.mesheryctlcommands.pattern.view.commands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/pattern/view">{{ subcommand.name }}</a></td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
      {% for flag_hash in site.data.mesheryctlcommands.pattern.view.flags %}{% assign flag = flag_hash[1] %}
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
