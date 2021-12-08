---
layout: default
title: Referencia del Comando mesheryctl
abstract: "Una Guía para la CLI de Meshery: mesheryctl"
permalink: es/reference/mesheryctl
language: es
redirect_from:
  - es/guides/mesheryctl
  - es/guides/mesheryctl-commands
type: Reference
---

## Categorías y estructura de comando

Los comandos de la CLI de Meshery se clasifican por función, que son:

- `mesheryctl` - Indicadores globales y configuración de CLI
- `mesheryctl system` - Ciclo de vida y resolución de problemas de Meshery
- `mesheryctl mesh` - Service Mesh Lifecycle & Configuration Management: mejores prácticas de aprovisionamiento y configuración
- `mesheryctl perf` - Gestión del rendimiento de la malla de servicios: caracterización del rendimiento de la malla de servicios y cargas de trabajo
- `mesheryctl pattern` - Configuración y gestión de patrones de malla de servicios: Patrones de malla de servicios e integración de modelos de aplicaciones abiertas
- `mesheryctl app` - Gestión de aplicaciones Service Mesh
- `mesheryctl filter` - Data Plane Intelligence: Registro y configuración de filtros de WebAssembly para Envoy (¡próximamente!)

## Banderas y comandos globales

<table>
<thead>
  <tr>
    <th>comando</th>
    <th>Subcomando</th>
    <th>Bandera</th>
    <th>Función</th>
  </tr>
  {% assign command1 = site.data.mesheryctlcommands.es_cmds.global %}
    <tr>
      <td rowspan=6><a href="{{ site.baseurl }}/reference/mesheryctl/mesheryctl">{{ command1.name }}</a></td>
      <td></td>
      <td></td>
      <td>{{ command1.description }}</td>
    </tr>
    {% for subcommand_hash in command1.subcommands %}{% assign subcommand = subcommand_hash[1] %}
      <tr>
         <td><a href="{{ site.baseurl }}/reference/mesheryctl/mesheryctl/version">{{ subcommand.name }}</a></td>
         <td></td>
         <td>{{ subcommand.description }}</td>
      </tr>
      {% for flag_hash in command1.flags %}{% assign flag = flag_hash[1] %}
        <tr>
         <td></td>
         <td>{{ flag.name }}</td>
         <td>{{ flag.description }}</td>
        </tr>
      {% endfor %}
    {% endfor %}
</thead>
</table>

## Gestión y resolución de problemas del ciclo de vida de Meshery

Instalación, resolución de problemas y depuración de Meshery y sus adaptadores.

<table>
<thead>
  <tr>
    <th>Comando principal</th>
    <th>Argumentos</th>
    <th>Bandera</th>
    <th>Función</th>
  </tr>
  {% assign command2 = site.data.mesheryctlcommands.es_cmds.system %}
    <tr>
      <td rowspan=20><a href="{{ site.baseurl }}/reference/mesheryctl/system">{{ command2.name }}</a></td>
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
    {% assign subcommand3 = command2.subcommands.completion %}
    <tr>
      <td><a href="{{ site.baseurl }}/reference/mesheryctl/system/completion">{{ subcommand3.name }}</a></td>
      <td></td>
      <td>{{ subcommand3.description }}</td>
    </tr>
    {% for flag_hash in subcommand3.flag %}{% assign flag = flag_hash[1] %}
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

{% assign command3 = site.data.mesheryctlcommands.es_cmds.system-channel %}

<tr>
<td rowspan=4><a href="{{ site.baseurl }}/reference/mesheryctl/system/channel">{{ command3.name }}</a></td>
<td></td>
<td></td>
<td>{{ command3.description }}</td>
</tr>
{% for subcommand_hash in command3.subcommands %}{% assign subcommand = subcommand_hash[1] %}
<tr>
<td>{{ subcommand.name }}</td>
<td></td>
<td>{{ subcommand.description }}</td>
</tr>
{% endfor %}
{% assign subcommand1 = command3.subcommands.view %}
{% for flag_hash in subcommand1.flag %}{% assign flag = flag_hash[1] %}
<tr>
<td>{{ flag.name }}</td>
<td>{{ flag.description }}</td>
</tr>
{% endfor %}
{% assign command4 = site.data.mesheryctlcommands.es_cmds.system-context %}
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

## Gestión del rendimiento de la malla de servicios

<table>
<thead>
  <tr>
    <th>Comando principal</th>
    <th>Argumentos</th>
    <th>Bandera</th>
    <th>Función</th>
  </tr>
  {% assign command5 = site.data.mesheryctlcommands.es_cmds.perf %}
    <tr>
      <td rowspan=18><a href="{{ site.baseurl }}/reference/mesheryctl/perf">{{ command5.name }}</a></td>
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
        <td rowspan=10><a href="{{ site.baseurl }}/reference/mesheryctl/perf/apply">{{ subcommand1.name }}</a></td>
        <td></td>
        <td>{{ subcommand1.description }}</td>
      </tr>
  {% for flag_hash in subcommand1.flags %}{% assign flag = flag_hash[1] %}
      <tr>
        <td>{{ flag.name }}</td>
        <td>{{ flag.description }}</td>
      </tr>
  {% endfor %}
  {% assign subcommand2 = command5.subcommands.list %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/perf/list">{{ subcommand2.name }}</a></td>
        <td></td>
        <td>{{ subcommand2.description }}</td>
      </tr>
  {% assign subcommand3 = command5.subcommands.view %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/perf/view">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
</thead>
</table>

## Gestión del ciclo de vida y la configuración de Service Mesh

<table>
<thead>
  <tr>
    <th>Comando principal</th>
    <th>Argumentos</th>
    <th>Bandera</th>
    <th>Función</th>
  </tr>
  {% assign command6 = site.data.mesheryctlcommands.es_cmds.mesh %}
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

## Configuración y gestión de patrones de malla de servicios

<table>
<thead>
  <tr>
    <th>Comando principal</th>
    <th>Argumentos</th>
    <th>Bandera</th>
    <th>Función</th>
  </tr>
  {% assign command7 = site.data.mesheryctlcommands.es_cmds.pattern %}
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

## Gestión de aplicaciones Service Mesh

<table>
<thead>
  <tr>
    <th>Comando principal</th>
    <th>Argumentos</th>
    <th>Bandera</th>
    <th>Función</th>
  </tr>
  {% assign command8 = site.data.mesheryctlcommands.es_cmds.app %}
    <tr>
      <td rowspan=12><a href="{{ site.baseurl }}/reference/mesheryctl/apps">{{ command8.name }}</a></td>
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
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/apps/onboard">{{ subcommand1.name }}</a></td>
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
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/apps/offboard">{{ subcommand2.name }}</a></td>
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
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/apps/list">{{ subcommand3.name }}</a></td>
        <td></td>
        <td>{{ subcommand3.description }}</td>
      </tr>
    {% assign subcommand4 = command8.subcommands.view %}
      <tr>
        <td><a href="{{ site.baseurl }}/reference/mesheryctl/apps/view">{{ subcommand4.name }}</a></td>
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
