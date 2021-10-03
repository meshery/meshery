---
layout: default
title: Referencia del Comando mesheryctl
abstract: "Una Guía para la CLI de Meshery: mesheryctl"
permalink: es/reference/mesheryctl
language: es
lang: es
categories: es
redirect_from: 
 - guides/mesheryctl
 - guides/mesheryctl-commands
type: Reference
---

## Comandos Globales y Banderas

| comando    |     bandera      |   función  |   Uso   |
| :--------- | :-----------: | :---------- | :-------- |{% for command_hash in site.data.mesheryctlcommands.global.commands %}{% assign command = command_hash[1] %}
| {{ command.name }} |       | {{ command.description }} | {{command.usage}} |{% endfor %}{% for subcommand_hash in site.data.mesheryctlcommands.global.subcommands %}{% assign subcommand = subcommand_hash[1] %}
|          | {{ subcommand.name }}  | {{ subcommand.description }} | {{ subcommand.usage }} |{% endfor %}


## Gestión del Ciclo de Vida de Meshery

Instalación, resolución de problemas y depuración de Meshery y sus adaptadores.


| comando    |    arg   |      bandera         | función  |  Uso  |
| :--------- | :------: | :---------------: | :-------- | :-----: |{% for command_hash in site.data.mesheryctlcommands.lifecycle.commands %}{% assign command = command_hash[1] %}
| {{ command.name }}  |     |  {{ command.flag }}  | {{ command.description }} | {{ command.usage }} |{% endfor %}{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.subcommands %}{% assign subcommand = subcommand_hash[1] %}
|         | {{ subcommand.name }}  |    | {{ subcommand.description }} | {{ subcommand.usage }} |{% endfor %}{% for flag_hash in site.data.mesheryctlcommands.lifecycle.flags %}{% assign flag = flag_hash[1] %}
| {{ flag.name }}  |  {{flag.arg}}    |  {{flag.flag}}  | {{ flag.description }} | {{ flag.usage }} |{% endfor %}


## Gestión del Desempeño

| comando    |      bandera         | función  | Uso |
| :--------- | :-----------------: | :---------- | :----------- |{% for command_hash in site.data.mesheryctlcommands.performance.commands %}{% assign command = command_hash[1] %}
| {{ command.name }} |       | {{ command.description }} | {{ command.usage }}  |{% endfor %}{% for flag_hash in site.data.mesheryctlcommands.performance.flags %}{% assign flag = flag_hash[1] %}
|          | {{ flag.name }}  | {{ flag.description }} | {{ flag.usage }} |{% endfor %}

## Gestión del Ciclo de Vida de la Malla de Servicios

| comando    |    arg   |      bandera     | función  |   Uso    |
| :--------- | :-------:| :-----------: | :-------- | :--------- |{% for command_hash in site.data.mesheryctlcommands.meshes.commands %}{% assign command = command_hash[1] %}
| {{ command.name }} |       |      | {{ command.description }} |     |{% endfor %}{% for subcommand_hash in site.data.mesheryctlcommands.meshes.subcommands %}{% assign subcommand = subcommand_hash[1] %}
|        | {{ subcommand.name }}  |    | {{ subcommand.description }} | {{ subcommand.usage }} |{% endfor %}{% for flag_hash in site.data.mesheryctlcommands.meshes.flags %}{% assign flag = flag_hash[1] %}
|        |        | {{ flag.name }}  | {{ flag.description }} | {{ flag.usage }} |{% endfor %}


