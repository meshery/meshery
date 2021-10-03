---
layout: default
title: Database
permalink: /es/concepts/architecture/database
type: concepts
redirect_from: architecture/database
abstract: "Meshery ofrece soporte para el almacenamiento en caché interno con la ayuda de bases de datos de archivos. Esto se ha implementado con varias bibliotecas que admiten diferentes tipos de formatos de datos."
language: es
list: include
---

## ¿Qué son las bases de datos de Meshery??

Las base de datos de Meshery funcionan como repositorios para [MeshSync](/concepts/architecture/meshsync), preferencias del usuario y configuración del sistema. Ambas bases de datos se consideran efímeras y deben tratarse como cachés. La retención de datos está vinculada a la vida útil de su instancia de Meshery Server. [Remote Providers](/extensibility/providers) puede ofrecer persistencia de datos a largo plazo. Las API de Meshery ofrecen mecanismos para los clientes, como [`mesheryctl`](/reference/mesheryctl) y la interfaz de usuario de Meshery para recuperar datos.

Consulte la figura siguiente para obtener detalles adicionales sobre los formatos de datos admitidos y el tipo de datos almacenados.

[![Architecture Diagram]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)

### Componentes

Meshery Database tiene varios tipos de implementaciones de bases de datos para admitir varios casos de uso. Se enumeran a continuación:
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Componente      | Libraría                               |
| :------------- | :------------------------------------ |
| Bitcask        | git.mills.io/prologic/bitcask         |
| SQLite         | gorm.io/gorm, gorm.io/driver/sqlite   |
