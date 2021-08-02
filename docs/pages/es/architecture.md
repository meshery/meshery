---
layout: page
title: Arquitectura de Meshery 
permalink: es/architecture
language: es
---

# Arquitectura

#### La arquitectura de Meshery puede ser vista desde dos perspectivas:
 
##### 1. [**Clientes**](#1-client-architecture)
##### 2. [**Proveedores**](#2-provider-architecture)

![Meshery architecture](/assets/img/architecture/Meshery-architecture-diagram.png)

### 1. **Arquitectura del Client**

![Client architecture](/assets/img/architecture/Meshery-client-architecture.svg)

### 2. **Arquitectura del Proveedor**

![Provider architecture](/assets/img/architecture/Meshery-provider-architecture.svg)

#### **Puertos de Red**

Meshery utiliza la siguiente lista de puertos de red para conectarse con sus múltiples componentes:

| Applicación de Red                             | Puerto           |
| :--------------------------------------------- | :--------------: |
| Meshery REST API                               | 9081/tcp         |
| Learn Layer5 Application                       | 10011            |

#### **Puertos del Adaptador**

| Service Mesh  | Puerto          |
| :------------ | ------------: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.port }} |
{% endif -%}
{% endfor %}

Véase la sección de [**Adaptadores**](/docs/architecture/adapters) para más información sobre el funcionamiento de un adaptador.

#### **Estado en los componentes de Meshery **

Algunos componentes dentro la arquitectura de Meshery se preocupan por la persistencia de datos mientras que otros se preocupan por una configuración duradera, mientras que otros no lo consideran.

| Componentes             | Persistencia   | Descripción                                                                      |
| :---------------------- | :------------- | :------------------------------------------------------------------------------- |
| mesheryctl              | stateless      | interface de línea de comando con un archivo de configuración                    |
| Adaptadores de Meshery  | stateless      | interface con service meshes con una base transaccional                          |
| Servidor de Meshery     | caché de estado| caché de aplicación guardada en la carpeta de usuario $HOME/.meshery/`           |
| Proveedores de Meshery  | stateful       | ubicación de preferencias persitentes de usuario, ambiente, pruebas, entre otros |
