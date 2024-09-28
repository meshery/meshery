---
layout: default
title: Command Line Reference
abstract: "A guide to Meshery's CLI: mesheryctl"
permalink: reference/mesheryctl
redirect_from:
  - reference/mesheryctl/commands/
  - reference/mesheryctl/commands
  - reference/mesheryctl/
type: Reference
language: en
abstract: "A guide to Meshery's CLI: mesheryctl"
---

## Categories and Command Structure

Meshery CLI commands are categorized by function, which are:

- `mesheryctl` - Global flags and CLI configuration
- `mesheryctl system` - Meshery Lifecycle and Troubleshooting
- `mesheryctl adapter` - Lifecycle & Configuration Management: provisioning and configuration best practices
- `mesheryctl perf` - Performance Management: Workload and cloud native performance characterization
- `mesheryctl design` - Design Patterns: Cloud native patterns and best practices
- `mesheryctl filter` - Data Plane Intelligence: Registry and configuration of WebAssembly filters for Envoy
- `mesheryctl model` - A unit of packaging to define managed infrastructure and their relationships, and details specifics of how to manage them.
- `mesheryctl components` - Fundamental building block used to represent and define the infrastructure under management
- `mesheryctl registry` - Model Database: Manage the state and contents of Meshery's internal registry of capabilities.
- `mesheryctl exp` - Experimental features


## Global Commands and Flags
<table>
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command1 = site.data.mesheryctlcommands.cmds.global %}
    <tr>
      <td rowspan=6><a href="{{ site.baseurl }}/reference/mesheryctl/main">{{ command1.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command1.description }}</td>
    </tr>
  {% for flag_hash in command1.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td></td>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
  {% endfor %}
  {% assign subcommand1 = command1.subcommands.version %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/version">{{ subcommand1.name }}</a></td>
      <td></td>
      <td>{{ subcommand1.description }}</td>
    </tr>
  {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
  {% endfor %}
  {% assign subcommand2 = command1.subcommands.completion %}
  <tr>
    <td><a href="{{ site.baseurl }}/reference/mesheryctl/completion">{{ subcommand2.name }}</a></td>
    <td></td>
    <td>{{ subcommand2.description }}</td>
  </tr>
  {% for flag_hash in subcommand2.flag %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
  {% endfor %}
</thead>
</table>

## Meshery Lifecycle Management and Troubleshooting

Installation, troubleshooting and debugging of Meshery and its adapters
<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Arguments</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command2 = site.data.mesheryctlcommands.cmds.system %}
    <tr>
      <td rowspan=34><a href="{{ site.baseurl }}/reference/mesheryctl/system">{{ command2.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command2.description }}</td>
    </tr>
    {% for flag_hash in command2.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand1 = command2.subcommands.start %}
      <tr>
        <td rowspan=5><a href="{{ site.baseurl }}/reference/mesheryctl/system/start">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
    {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand2 = command2.subcommands.stop %}
    <tr>
      <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/stop">{{ subcommand2.name }}</a></td>
      <td></td>
      <td>{{ subcommand2.description }}</td>
    </tr>
    {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand4 = command2.subcommands.update %}
    <tr>
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/update">{{ subcommand4.name }}</a></td>
      <td></td>
      <td>{{ subcommand4.description }}</td>
    </tr>
    {% for flag_hash in subcommand4.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand5 = command2.subcommands.config %}
    <tr>
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/config">{{ subcommand5.name }}</a></td>
      <td></td>
      <td>{{ subcommand5.description }}</td>
    </tr>
    {% for flag_hash in subcommand5.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand7 = command2.subcommands.reset %}
    <tr>
      <td rowspan="1"><a href="{{ site.baseurl }}/reference/mesheryctl/system/reset">{{ subcommand7.name }}</a></td>
      <td></td>
      <td>{{ subcommand7.description }}</td>
    </tr>
    {% for flag_hash in subcommand7.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand7 = command2.subcommands.logs %}
    <tr>
      <td rowspan="2"><a href="{{ site.baseurl }}/reference/mesheryctl/system/logs">{{ subcommand7.name }}</a></td>
      <td></td>
      <td>{{ subcommand7.description }}</td>
    </tr>
    {% for flag_hash in subcommand7.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand8 = command2.subcommands.restart %}
    <tr>
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/restart">{{ subcommand8.name }}</a></td>
      <td></td>
      <td>{{ subcommand8.description }}</td>
    </tr>
    {% for flag_hash in subcommand8.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand9 = command2.subcommands.status %}
    <tr>
      <td rowspan="2"><a href="{{ site.baseurl }}/reference/mesheryctl/system/status">{{ subcommand9.name }}</a></td>
      <td></td>
      <td>{{ subcommand9.description }}</td>
    </tr>
    {% for flag_hash in subcommand9.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand13 = command2.subcommands.dashboard %}
    <tr>
      <td rowspan="3"><a href="{{ site.baseurl }}/reference/mesheryctl/system/dashboard">{{ subcommand13.name }}</a></td>
      <td></td>
      <td>{{ subcommand13.description }}</td>
    </tr>
    {% for flag_hash in subcommand13.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand10 = command2.subcommands.login %}
    <tr>
      <td rowspan="2"><a href="{{ site.baseurl }}/reference/mesheryctl/system/login">{{ subcommand10.name }}</a></td>
      <td></td>
      <td>{{ subcommand10.description }}</td>
    </tr>
    {% for flag_hash in subcommand10.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand11 = command2.subcommands.logout %}
    <tr>
      <td rowspan="1"><a href="{{ site.baseurl }}/reference/mesheryctl/system/logout">{{ subcommand11.name }}</a></td>
      <td></td>
      <td>{{ subcommand11.description }}</td>
    </tr>
    {% for flag_hash in subcommand11.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign subcommand12 = command2.subcommands.check %}
    <tr>
      <td rowspan=6><a href="{{ site.baseurl }}/reference/mesheryctl/system/check">{{ subcommand12.name }}</a></td>
      <td></td>
      <td>{{ subcommand12.description }}</td>
    </tr>
    {% for flag_hash in subcommand12.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign command3 = site.data.mesheryctlcommands.cmds.system-channel %}
    <tr>
      <td rowspan=5><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel">{{ command3.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command3.description }}</td>
    </tr>
    {% assign subcommand1 = command3.subcommands.set %}
    <tr>
      <td rowspan=1><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel/set">{{ subcommand1.name }}</a></td>
      <td></td>
      <td>{{ subcommand1.description }}</td>
    </tr>
    {% assign subcommand2 = command3.subcommands.switch %}
    <tr>
      <td rowspan=1><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel/switch">{{ subcommand2.name }}</a></td>
      <td></td>
      <td>{{ subcommand2.description }}</td>
    </tr>
    {% assign subcommand3 = command3.subcommands.view %}
    <tr>
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel/view">{{ subcommand3.name }}</a></td>
      <td></td>
      <td>{{ subcommand3.description }}</td>
    </tr>
    {% for flag_hash in subcommand3.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign command4 = site.data.mesheryctlcommands.cmds.system-context %}
    <tr>
      <td rowspan=12><a href="{{ site.baseurl }}/reference/mesheryctl/system/context">{{ command4.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command4.description }}</td>
    </tr>
    {% assign subcommand1 = command4.subcommands.create %}
    <tr>
      <td rowspan=5><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/create">{{ subcommand1.name }}</a></td>
      <td></td>
      <td>{{ subcommand1.description }}</td>
    </tr>
    {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign subcommand2 = command4.subcommands.delete %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/delete">{{ subcommand2.name }}</a></td>
      <td></td>
      <td>{{ subcommand2.description }}</td>
    </tr>
    {% assign subcommand3 = command4.subcommands.view %}
    <tr>
      <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/view">{{ subcommand3.name }}</a></td>
      <td></td>
      <td>{{ subcommand3.description }}</td>
    </tr>
    {% for flag_hash in subcommand3.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign subcommand4 = command4.subcommands.switch %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/switch">{{ subcommand4.name }}</a></td>
      <td></td>
      <td>{{ subcommand4.description }}</td>
    </tr>
    {% for flag_hash in subcommand4.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign subcommand5 = command4.subcommands.list %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/list">{{ subcommand5.name }}</a></td>
      <td></td>
      <td>{{ subcommand5.description }}</td>
    </tr>
    {% for flag_hash in subcommand5.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign command5 = site.data.mesheryctlcommands.cmds.system-provider %}
    <tr>
      <td rowspan=8><a href="{{ site.baseurl }}/reference/mesheryctl/system/provider">{{ command5.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command5.description }}</td>
    </tr>
    {% assign subcommand2 = command5.subcommands.list %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/provider/list">{{ subcommand2.name }}</a></td>
      <td></td>
      <td>{{ subcommand2.description }}</td>
    </tr>
    {% assign subcommand3 = command5.subcommands.reset %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/provider/reset">{{ subcommand3.name }}</a></td>
      <td></td>
      <td>{{ subcommand3.description }}</td>
    </tr>
    {% assign subcommand4 = command5.subcommands.switch %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/provider/switch">{{ subcommand4.name }}</a></td>
      <td></td>
      <td>{{ subcommand4.description }}</td>
    </tr>
    {% assign subcommand5 = command5.subcommands.view %}
    <tr>
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/provider/view">{{ subcommand5.name }}</a></td>
      <td></td>
      <td>{{ subcommand5.description }}</td>
    </tr>
    {% for flag_hash in subcommand5.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
    {% assign subcommand1 = command5.subcommands.set %}
    <tr>
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/provider/set">{{ subcommand1.name }}</a></td>
      <td></td>
      <td>{{ subcommand1.description }}</td>
    </tr>
    {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
    <tr>
      <td>{{ flag.name }}</td>
      <td>{{ flag.description }}</td>
    </tr>
    {% endfor %}
</thead>
</table>


## Cloud Native Performance Management
<table class="table-wrapper">
<thead>
  <tr>
    <th>Main Command</th>
    <th>Arguments</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command5 = site.data.mesheryctlcommands.cmds.perf %}
    <tr>
      <td rowspan=20><a href="{{ site.baseurl }}/reference/mesheryctl/perf">{{ command5.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command5.description }}</td>
    </tr>
    {% for flag_hash in command5.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
  {% assign subcommand1 = command5.subcommands.apply %}
      <tr>
        <td rowspan=11><a href="{{ site.baseurl }}/reference/mesheryctl/perf/apply">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
  {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
  {% endfor %}
  {% assign subcommand2 = command5.subcommands.profile %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/perf/profile">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
  {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
  {% endfor %}
  {% assign subcommand3 = command5.subcommands.result %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/perf/result">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
  {% for flag_hash in subcommand3.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
  {% endfor %}
</thead>
</table>

## Cloud Native Lifecycle and Configuration Management
<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command7 = site.data.mesheryctlcommands.cmds.adapter %}
    {% assign subcommand_flag_count = 0 %}
    {% for subcommand_hash in command7.subcommands %}
      {% assign subcommand = subcommand_hash[1] %}
      {% assign subcommand_flag_count = subcommand_flag_count | plus: subcommand.flags.size %}
    {% endfor %}
    {% assign total_rowspan = command7.subcommands.size | plus: subcommand_flag_count | plus: command7.flags.size | plus: 1 %}
    <tr>
      <td rowspan={{ total_rowspan }}><a href="{{ site.baseurl }}/reference/mesheryctl/{{ command7.name }}">{{ command7.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command7.description }}</td>
    </tr>
    {% for subcommand_hash in command7.subcommands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td rowspan={{ subcommand.flags.size | plus: 1 }} ><a href="{{ site.baseurl }}/reference/mesheryctl/{{ command7.name }}/{{ subcommand.name }}">{{ subcommand.name }}</a></td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
      {% for flag_hash in subcommand.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% endfor %}
</thead>
</table>

## Cloud Native Design Configuration and Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command7 = site.data.mesheryctlcommands.cmds.design %}
    <tr>
      <td rowspan=23><a href="{{ site.baseurl }}/reference/mesheryctl/design">{{ command7.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command7.description }}</td>
    </tr>
    {% assign subcommand1 = command7.subcommands.apply %}
      <tr>
        <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/design/apply">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
      {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand2 = command7.subcommands.view %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/design/view">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
      {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand3 = command7.subcommands.list %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/design/list">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
      {% for flag_hash in subcommand3.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand4 = command7.subcommands.delete %}
      <tr>
        <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/design/delete">{{ subcommand4.name }}</a></td>
        <td></td>
        <td>{{ subcommand4.description }}</td>
      </tr>
      {% for flag_hash in subcommand4.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand5 = command7.subcommands.import %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/design/import">{{ subcommand5.name }}</a></td>
        <td></td>
        <td>{{ subcommand5.description }}</td>
      </tr>
      {% for flag_hash in subcommand5.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand6 = command7.subcommands.onboard %}
      <tr>
        <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/design/onboard">{{ subcommand6.name }}</a></td>
        <td></td>
        <td>{{ subcommand6.description }}</td>
      </tr>
      {% for flag_hash in subcommand6.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
      {% assign subcommand6 = command7.subcommands.export %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/design/export">{{ subcommand6.name }}</a></td>
        <td></td>
        <td>{{ subcommand6.description }}</td>
      </tr>
      {% for flag_hash in subcommand6.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand7 = command7.subcommands.offboard %}
      <tr>
        <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/design/offboard">{{ subcommand7.name }}</a></td>
        <td></td>
        <td>{{ subcommand7.description }}</td>
      </tr>
      {% for flag_hash in subcommand7.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
</thead>
</table>
 
## Data Plane Intelligence
<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command9 = site.data.mesheryctlcommands.cmds.filter %}
    <tr>
      <td rowspan=10><a href="{{ site.baseurl }}/reference/mesheryctl/filter">{{ command9.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command9.description }}</td>
    </tr>
    {% assign subcommand1 = command9.subcommands.import %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/filter/import">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
      {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand2 = command9.subcommands.delete %}
      <tr>
        <td rowspan=1><a href="{{ site.baseurl }}/reference/mesheryctl/filter/delete">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
      {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand3 = command9.subcommands.list %}
      <tr>
        <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/filter/list">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
      {% for flag_hash in subcommand3.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand4 = command9.subcommands.view %}
      <tr>
        <td rowspan=3><a href="{{ site.baseurl }}/reference/mesheryctl/filter/view">{{ subcommand4.name }}</a></td>
        <td></td>
        <td>{{ subcommand4.description }}</td>
      </tr>
      {% for flag_hash in subcommand4.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
</thead>
</table>

## Meshery Registry Management
<table>
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command10 = site.data.mesheryctlcommands.cmds.registry %}
    <tr>
      <td rowspan=14><a href="{{ site.baseurl }}/reference/mesheryctl/registry">{{ command10.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command10.description }}</td>
    </tr>
    {% for flag_hash in command10.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
{% assign subcommand1 = command10.subcommands.publish %}
      <tr>
        <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/registry/publish">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
      {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
{% assign subcommand2 = command10.subcommands.update %}
<tr>
        <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/registry/update">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
      {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
{% assign subcommand3 = command10.subcommands.generate %}
<tr>
        <td rowspan=6><a href="{{ site.baseurl }}/reference/mesheryctl/registry/generate">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
      {% for flag_hash in subcommand3.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
</thead>
</table>

## Meshery Models
<table>
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command12 = site.data.mesheryctlcommands.cmds.model %}
    <tr>
      <td rowspan=10><a href="{{ site.baseurl }}/reference/mesheryctl/{{ command12.name }}">{{ command12.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command12.description }}</td>
    </tr>
    {% for flag_hash in command12.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% for subcommand_hash in command12.subcommands %}{% assign subcomand = subcommand_hash[1] %}
      <tr>
        <td rowspan={{ subcomand.flags.size | plus:1 }} ><a href="{{ site.baseurl }}/reference/mesheryctl/{{ command12.name }}/{{ subcomand.name }}">{{ subcomand.name }}</a></td>
        <td></td>
        <td>{{ subcomand.description }}</td>
      </tr>
      {% for flag_hash in subcomand.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% endfor %}
</thead>
</table>

## Meshery Components
<table>
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command13 = site.data.mesheryctlcommands.cmds.components %}
    <tr>
      <td rowspan=9><a href="{{ site.baseurl }}/reference/mesheryctl/{{ command13.name }}">{{ command13.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command13.description }}</td>
    </tr>
    {% for flag_hash in command13.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% for subcommand_hash in command13.subcommands %}{% assign subcomand = subcommand_hash[1] %}
      <tr>
        <td rowspan={{ subcomand.flags.size | plus:1 }} ><a href="{{ site.baseurl }}/reference/mesheryctl/{{ command13.name }}/{{ subcomand.name }}">{{ subcomand.name }}</a></td>
        <td></td>
        <td>{{ subcomand.description }}</td>
      </tr>
      {% for flag_hash in subcomand.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% endfor %}
</thead>
</table>

## Experimental Features(exp)
<table>
<thead>
  <tr>
    <th>Command</th>
    <th>Subcommand</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command14 = site.data.mesheryctlcommands.cmds.exp %}
  {% for cmd_hash in command14 %}
    {% assign cmd = cmd_hash[1] %}
    {% assign subcommand_flag_count = 0 %}
    {% for subcommand_hash in cmd.subcommands %}
      {% assign subcommand = subcommand_hash[1] %}
      {% assign subcommand_flag_count = subcommand_flag_count | plus: subcommand.flags.size %}
    {% endfor %}
    {% assign total_rowspan = cmd.subcommands.size | plus: subcommand_flag_count | plus: cmd.flags.size | plus: 1 %}
    <tr>
      <td rowspan="{{ total_rowspan }}"><a href="{{ site.baseurl }}/reference/mesheryctl/exp/{{ cmd.name }}">{{ cmd.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ cmd.description }}</td>
    </tr>
    {% for flag_hash in cmd.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% for subcommand_hash in cmd.subcommands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
        <td rowspan="{{ subcommand.flags.size | plus: 1 }}"><a href="{{ site.baseurl }}/reference/mesheryctl/exp/{{ cmd.name }}/{{ subcommand.name }}">{{ subcommand.name }}</a></td>
        <td></td>
        <td>{{ subcommand.description }}</td>
      </tr>
      {% for flag_hash in subcommand.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% endfor %}
  {% endfor %}
</thead>
</table>

{% include related-discussions.html tag="mesheryctl" %}

