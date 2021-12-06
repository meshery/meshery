---
layout: default
title: Arquitectura
permalink: es/concepts/architecture
redirect_from: architecture
type: concepts
abstract: descripción general de diferentes componentes individuales de la arquitectura de Meshery y cómo interactuan como un sistema.
language: es
list: include
---

## Componentes, su propósito y Lenguajes

Meshery y sus componentes están escritos utilizando los siguientes lenguajes de programación y tecnologías.

| Componentes                                                          | Lenguajes y Tecnologías                                                                                  |
| :------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------- |
| Meshery Server                                                       | Golang, gRPC, GraphQL, [SMP](https://smp-spec.io)                                                        |
| [Meshery Adapters](/concepts/architecture/adapters)                  | Golang, gRPC, [CloudEvents](https://cloudevents.io/), [SMI](https://smi-spec.io), [OAM](https://oam.dev) |
| [Meshery WASM Filters](https://github.com/layer5io/wasm-filters)     | Rust y C++                                                                                               |
| Meshery UI                                                           | ReactJS, NextJS, BillboardJS                                                                             |
| Meshery Provider UI                                                  | ReactJS, NextJS                                                                                          |
| [Meshery Remote Providers](/extensibility/providers)                 | _cualquiera_ - debe adherirse a los [Extension Points]({{site.baseurl}}/extensibility) de Meshery        |
| [Meshery Operator](/concepts/architecture/operator)                  | Golang                                                                                                   |
| &nbsp;&nbsp;&nbsp;&nbsp; [MeshSync](/concepts/architecture/meshsync) | Golang                                                                                                   |
| &nbsp;&nbsp;&nbsp;&nbsp; [Broker](/concepts/architecture/broker)     | Golang, NATS                                                                                             |
| [Meshery Database](/concepts/architecture/database)                  | Golang, SQLlite                                                                                          |

## Despliegues

Meshery se despliega como un conjunto de contenedores. Los contenedores de Meshery pueden ser desplegados en Docker o Kubernetes. Los componentes de Meshery se conectan entre sí a través de peticiones gRPC. Meshery Server almacena la ubicación de los otros componentes y se conecta con estos componentes según sea necesario. Típicamente, una conexión de Meshery Server a Meshery Adapters es iniciada desde un petición de cliente (generalmente `mesheryctl` o Meshery UI) para recopilar información del Adapter o invocar una operación del Adapter.

### Adapters

En Meshery v0.6.0, los Adapters se registrarán con Meshery Server a través de HTTP POST. Si Meshery Server no está disponible, los Meshery Adapters se apartan y vuelven a reintentar conectarse a Meshery Server perpetuamente.

[![Arquitectura de Meshery]({{ site.baseurl }}/assets/img/architecture/meshery-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/meshery-architecture.svg)

_Figure: Despliegue de Meshery dentro o fuera de un clúster de Kubernetes_

#### Adapters y Registro de Capacidades

Cada Meshery Adapter ofrece su propia funcionalidad específica de un servicio único. Como tal, en el momento del despliegue, el Meshery Adapter registrará sus capacidades específicas de service mesh (sus operaciones) con el registro de capacidades de Meshery Server.

[![Registro de Operación de Meshery Adapter]({{ site.baseurl }}/assets/img/adapters/meshery-adapter-operation-registration.svg
)]({{ site.baseurl }}/assets/img/adapters/meshery-adapter-operation-registration.svg)

_Figure: Registro de Operación de Service Mesh Adapter_

### Clientes

La REST API de Meshery puede ser consumida por cualquier número de clientes. Los clientes necesitan presentar un token JWT válido.

[![Arquitectura de cliente]({{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/Meshery-client-architecture.svg)

_Figure: Los clientes usan [REST API](extensibility/api#rest), [GraphQL API](extensibility/api#graphql) de Meshery o una combinación de ambos._

### Providers

Como un punto de extensión, Meshery soporta dos tipos de providers: _Local_ y _Remoto_.

[![Arquitectura de Provider]({{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.svg)]({{ site.baseurl }}/assets/img/architecture/Meshery-provider-architecture.svg)

## Modelo de Objetos

Este diagrama describe las construcciones lógicas dentro de Meshery y sus relaciones.

[![Modelo de Objetos]({{ site.baseurl }}/assets/img/architecture/meshery_logical_object_model.svg)]({{ site.baseurl }}/assets/img/architecture/meshery_logical_object_model.svg)

## Meshery Operator y MeshSync

Meshery Operator es el operador multi-service mesh (un controlador de Kubernetes personalizado) que administra MeshSync y su agente de mensajes.

[![Meshery Operator y MeshSync]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg
)]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg)

_Ver la sección [**Operator**]({{ site.baseurl }}/concepts/architecture/operator) para obtener más información sobre la función de un operador y la sección [**MeshSync**]({{ site.baseurl }}/concepts/architecture/meshsync) para obtener más infromación sobre la función de MeshSync._

## Base de datos

La base de datos de Meshery es responsable de recopilar y centralizar el estado de todos los elementos bajo administración, incluyendo la infraestructura, la aplicación y los propios componentes de Meshery. La base de datos de Meshery, a la vez que persiste para archivar, es tratada como caché.

[![Meshery Database]({{ site.baseurl }}/assets/img/architecture/meshery-database.svg)]({{ site.baseurl }}/concepts/architecture/database)

_Ver la sección [**Database**]({{ site.baseurl }}/concepts/architecture/database) para obtener más información sobre la función de la base de datos._

### **Statefulness en los componentes de Meshery**

Algunos componentes dentro de la arquitectura de Meshery se ocupan de la persistencia de datos, mientras que otros solo se relacionan con una configuración de larga duración, mientras que otros no tienen ningún estado en absoluto.

| Componente        | Persistencia | Descripción                                                                               |
| :---------------- | :----------- | :---------------------------------------------------------------------------------------- |
| mesheryctl        | stateless    | interfaz de línea de comandos que tiene un archivo de configuración                       |
| Meshery Adapters  | stateless    | interfaz con service meshes sobre una base transaccional                                  |
| Meshery Server    | caches state | la memoria caché de la aplicación se almacena en la carpeta `$HOME/.meshery/` del usuario |
| Meshery Providers | stateful     | ubicación de las preferencias de usuario persistentes, ambiente, pruebas y etc.           |
| Meshery Operator  | stateless    | operador de controladores personalizados de Meshery, espacialmente MeshSync               |
| MeshSync          | stateless    | Controlador personalizado de Kubernetes, continuamente ejecutando el descubrimiento       |

### **Puertos de Red**

Meshery usa la siguiente lista de puertos de red para interactuar con sus diversos componenetes:

| Componente               |                           Puerto                           |
| :----------------------- | :--------------------------------------------------------: |
| Meshery REST API         |                          9081/tcp                          |
| Meshery GraphQL          |                          9081/tcp                          |
| Meshery Broker           | 4222/tcp, 6222/tcp, 8222/tcp, 7777/tcp, 7422/tcp, 7522/tcp |
| Learn Layer5 Application |                         10011/tcp                          |
| Meshery Adapters         |                         10000+/tcp                         |
| Meshery Remote Providers |                          443/tcp                           |

### **Puertos de Adapter**

| Service Mesh | Port |
| :----------- | ---: |
{% for adapter in site.adapters -%}
{% if adapter.port -%}
| <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) |&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; {{ adapter.port }} |
{% endif -%}
{% endfor %}

_Ver la sección [**Adapters**]({{ site.baseurl }}/concepts/architecture/adapters) para obtener más información sobre la función del adapter._

