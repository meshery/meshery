---
title: Database
description: "Meshery offers support for internal caching with the help of file databases. This has been implemented with several libraries that supports different kinds of data formats."
aliases:
- /architecture/database/
---

## What are the Meshery Databases?

Meshery Databases function as repositories for [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}), user preferences and system settings. Both databases are considered ephemeral and should be treated as caches. Data retention is tethered to the lifetime of their Meshery Server instance. [Remote Providers]({{< ref "reference/extensibility/providers/index.md" >}}) may offer long-term data persistence. Meshery's APIs offer mechanisms for clients, like [`mesheryctl`]({{< ref "reference/mesheryctl/_index.md" >}}) and Meshery UI to retrieve data.

See the figure below for additional details of the data formats supported and type of data stored.

[![Architecture Diagram](./images/meshery-database.webp)](./images/meshery-database.webp)

### Components

Meshery Server stores its data in a SQLite database file, `mesherydb.sql`, in the directory set by the [`USER_DATA_FOLDER`]({{< ref "installation/advanced/environment-variables.md" >}}) environment variable (default: `~/.meshery/config`). It is implemented with the libraries listed below:

| Component      | Library                               |
| :------------- | :------------------------------------ |
| SQLite         | gorm.io/gorm, gorm.io/driver/sqlite   |
