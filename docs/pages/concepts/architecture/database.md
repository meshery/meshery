---
layout: default
title: Database
permalink: concepts/architecture/database
type: concepts
redirect_from: architecture/database
abstract: "Meshery offers support for internal caching with the help of file databases. This has been implemented with several libraries that supports different kinds of data formats."
language: en
list: include
---

## What is the Meshery Database?

The Meshery Database offers support for internal caching with the help of file databases. This has been implemented with several libraries that supports different kinds of data formats. 

## Components Involved:

Meshery Database has several kinds of database implementations to support various usecases. They are listed below:
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Component      | Library                               |
| :------------- | :------------------------------------ |
| Bitcask        | github.com/prologic/bitcask           |
| SQLite         | gorm.io/gorm, gorm.io/driver/sqlite   |
