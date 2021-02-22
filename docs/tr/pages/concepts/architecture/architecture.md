---
layout: page
title: Meshery Mimarisi
permalink: /tr/concepts/architecture
abstract: "Meshery mimarisinin farklı bileşenlerine ve bunların bir sistem olarak nasıl etkileşim kurduğuna genel bakış"
language: tr
lang: tr
categories: tr
type: concepts
list: include
---

# Mimari

#### Meshery'nin mimarisi iki açıdan görülebilir:
 
##### 1. [**Müşteriler**](#1-client-architecture)
##### 2. [**Sağlayıcılar**](#2-provider-architecture)

![Meshery architecture](/assets/img/architecture/Meshery-architecture-diagram.png)

### 1. **İstemci Mimarisi**

![Client architecture](/assets/img/architecture/Meshery-client-architecture.svg)

### 2. **Satıcı Mimarisi**

![Provider architecture](/assets/img/architecture/Meshery-provider-architecture.svg)

#### **Ağ bağlantı noktaları**

Meshery, birden çok bileşenine bağlanmak için aşağıdaki ağ bağlantı noktaları listesini kullanır:

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
