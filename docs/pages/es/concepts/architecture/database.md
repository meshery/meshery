---
layout: default
title: Database
permalink: es/concepts/architecture/database
type: concepts
redirect_from: architecture/database
abstract: "Meshery ofrece soporte para el almacenamiento en caché interno con la ayuda de base de datos de archivos. Esto se ha implementado con varias bibliotecas que admiten diferentes tipos de formatos de datos."
language: es
list: include
---

## ¿Qué son las Meshery Databases?

Las Meshery Databases funcionan como repositorios para [MeshSync](/concepts/architecture/meshsync), preferencias del usuario y configuraciones del sistema. Ambas bases de datos se consideran efímeras y deben ser tratadas como cachés. La retención de datos está atada a la vida útil de su instancia de Meshery Server. Los [Remote Providers](/extensibility/providers) pueden ofrecer persistencia de datos a largo plazo. Las API de Meshery ofrecen mecanismos para clientes, como [`mesheryctl`](/reference/mesheryctl) y UI de Meshery para recuperar datos.

Ver la figura a continuación para obtener detalles adicionales de los formatos de datos soportados y el tipo de datos almacenados.

[![Diagrama de Arquitectura]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)

### Componentes

Meshery Database tiene varios tipos de implementaciones de base de datos para soportar varios casos de uso. Estos son listados a continuación:
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Componente     | Biblioteca                            |
| :------------- | :------------------------------------ |
| Bitcask        | git.mills.io/prologic/bitcask         |
| SQLite         | gorm.io/gorm, gorm.io/driver/sqlite   |
