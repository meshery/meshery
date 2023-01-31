---
layout: default
title: Command Reference
abstract: "A guide to Meshery's CLI: mesheryctl"
permalink: reference/mesheryctl
redirect_from:
  - reference/mesheryctl/commands/
  - reference/mesheryctl/commands
  - reference/mesheryctl/
type: Reference
---

## Categories and Command Structure

Meshery CLI commands are categorized by function, which are:

- `mesheryctl` - Global flags and CLI configuration
- `mesheryctl system` - Meshery Lifecycle and Troubleshooting
- `mesheryctl mesh` - Service Mesh Lifecycle & Configuration Management: provisioning and configuration best practices
- `mesheryctl perf` - Service Mesh Performance Management: Workload and service mesh performance characterization
- `mesheryctl pattern` - Service Mesh Pattern Configuration & Management: Service mesh patterns and Open Application Model integration
- `mesheryctl app` - Service Mesh Application Management
- `mesheryctl filter` - Data Plane Intelligence: Registry and configuration of WebAssembly filters for Envoy (experimental feature)

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

Installation, troubleshooting and debugging of Meshery and its adapters.

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
      <td rowspan=25><a href="{{ site.baseurl }}/reference/mesheryctl/system">{{ command2.name }}</a></td>
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
        <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/start">{{ subcommand1.name }}</a></td>
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
      <td rowspan=2><a href="{{ site.baseurl }}/reference/mesheryctl/system/stop">{{ subcommand2.name }}</a></td>
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
    {% assign subcommand6 = command2.subcommands.reset %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/reset">{{ subcommand6.name }}</a></td>
      <td></td>
      <td>{{ subcommand6.description }}</td>
    </tr>
    {% for flag_hash in subcommand6.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand7 = command2.subcommands.logs %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/logs">{{ subcommand7.name }}</a></td>
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
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/status">{{ subcommand9.name }}</a></td>
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
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/dashboard">{{ subcommand13.name }}</a></td>
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
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/login">{{ subcommand10.name }}</a></td>
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
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/logout">{{ subcommand11.name }}</a></td>
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
      <td rowspan=5><a href="{{ site.baseurl }}/reference/mesheryctl/system/check">{{ subcommand12.name }}</a></td>
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
            <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel/set">{{ subcommand1.name }}</a></td>
            <td></td>
            <td>{{ subcommand1.description }}</td>
          </tr>
        {% assign subcommand2 = command3.subcommands.switch %}
          <tr>
            <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel/switch">{{ subcommand2.name }}</a></td>
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
        <td rowspan=10><a href="{{ site.baseurl }}/reference/mesheryctl/system/context">{{ command4.name }}</a></td>
        <td></td>
        <td></td>
        <td>{{ command4.description }}</td>
      </tr>
      {% assign subcommand1 = command4.subcommands.create %}
        <tr>
          <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/create">{{ subcommand1.name }}</a></td>
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
          <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/context/switch">{{ subcommand4.name }}</a></td>
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

## Service Mesh Performance Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Arguments</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command5 = site.data.mesheryctlcommands.cmds.perf %}
    <tr>
      <td rowspan=19><a href="{{ site.baseurl }}/reference/mesheryctl/perf">{{ command5.name }}</a></td>
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
        <td rowspan=9><a href="{{ site.baseurl }}/reference/mesheryctl/perf/apply">{{ subcommand1.name }}</a></td>
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
        <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/perf/result">{{ subcommand3.name }}</a></td>
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

## Service Mesh Lifecycle and Configuration Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command6 = site.data.mesheryctlcommands.cmds.mesh %}
    <tr>
      <td rowspan=10><a href="{{ site.baseurl }}/reference/mesheryctl/mesh">{{ command6.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command6.description }}</td>
    </tr>
    {% assign subcommand1 = command6.subcommands.validate %}
      <tr>
        <td rowspan=5><a href="{{ site.baseurl }}/reference/mesheryctl/mesh/validate">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
      {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% assign subcommand2 = command6.subcommands.deploy %}
      <tr>
        <td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/mesh/deploy">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
      {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
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
  {% assign command7 = site.data.mesheryctlcommands.cmds.pattern %}
    <tr>
      <td rowspan=12><a href="{{ site.baseurl }}/reference/mesheryctl/pattern">{{ command7.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command7.description }}</td>
    </tr>
    {% for flag_hash in command7.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand1 = command7.subcommands.list %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/pattern/list">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
      {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
    {% assign subcommand2 = command7.subcommands.apply %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/pattern/apply">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
      {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
    {% assign subcommand3 = command7.subcommands.view %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/pattern/view">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
      {% for flag_hash in subcommand3.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
    {% assign subcommand4 = command7.subcommands.delete %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/pattern/delete">{{ subcommand4.name }}</a></td>
        <td></td>
        <td>{{ subcommand4.description }}</td>
      </tr>
      {% for flag_hash in subcommand4.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
</thead>
</table>

## Service Mesh Application Management

<table>
<thead>
  <tr>
    <th>Main Command</th>
    <th>Command</th>
    <th>Flag</th>
    <th>Function</th>
  </tr>
  {% assign command8 = site.data.mesheryctlcommands.cmds.app %}
    <tr>
      <td rowspan=12><a href="{{ site.baseurl }}/reference/mesheryctl/app">{{ command8.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command8.description }}</td>
    </tr>
    {% for flag_hash in command8.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand1 = command8.subcommands.onboard %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/app/onboard">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
    {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand2 = command8.subcommands.offboard %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/app/offboard">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
    {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand3 = command8.subcommands.list %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/app/list">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
    {% assign subcommand4 = command8.subcommands.view %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/app/view">{{ subcommand4.name }}</a></td>
        <td></td>
        <td>{{ subcommand4.description }}</td>
      </tr>
      {% for flag_hash in subcommand4.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
</thead>
</table>

## Data Plane Intelligence (experimental feature)

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
      <td rowspan=12><a href="{{ site.baseurl }}/reference/mesheryctl/exp/filter">{{ command9.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command9.description }}</td>
    </tr>
    {% for flag_hash in command9.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand1 = command9.subcommands.apply %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/exp/filter/apply">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
    {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand2 = command9.subcommands.delete %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/exp/filter/delete">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
    {% for flag_hash in subcommand2.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td></td>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
    {% endfor %}
    {% assign subcommand3 = command9.subcommands.list %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/exp/filter/list">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
    {% assign subcommand4 = command9.subcommands.view %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/exp/filter/view">{{ subcommand4.name }}</a></td>
        <td></td>
        <td>{{ subcommand4.description }}</td>
      </tr>
      {% for flag_hash in subcommand4.flags %}{% assign flag = flag_hash[1] %}
        <tr>
          <td></td>
          <td>{{ flag.name }}</td>
          <td>{{ flag.description }}</td>
      </tr>
      {% endfor %}
</thead>
</table>
{% include related-discussions.html tag="mesheryctl" %}

