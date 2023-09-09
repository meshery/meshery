---
layout: default
title: Database
permalink: concepts/architecture/database
type: components
redirect_from: architecture/database
abstract: "Meshery offers support for internal caching with the help of file databases. This has been implemented with several libraries that supports different kinds of data formats."
language: en
list: include
---

## What are the Meshery Databases?

Meshery Databases function as repositories for [MeshSync](/concepts/architecture/meshsync), user preferences and system settings. Both databases are considered ephemeral and should be treated as caches. Data retention is tethered to the lifetime of their Meshery Server instance. [Remote Providers](/extensibility/providers) may offer long-term data persistence. Meshery's APIs offer mechanisms for clients, like [`mesheryctl`](/reference/mesheryctl) and Meshery UI to retrieve data.

See the figure below for additional details of the data formats supported and type of data stored.

[![Architecture Diagram]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)

### Components

Meshery Database has several kinds of database implementations to support various usecases. They are listed below:
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Component      | Library                               |
| :------------- | :------------------------------------ |
| Bitcask        | git.mills.io/prologic/bitcask         |
| SQLite         | gorm.io/gorm, gorm.io/driver/sqlite   |
