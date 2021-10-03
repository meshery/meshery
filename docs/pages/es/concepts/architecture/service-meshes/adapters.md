---
layout: default
title: Adaptadores
permalink: es/concepts/architecture/adapters
redirect_from: es/architecture/adapters
abstract: "Meshery ofrece soporte para más adaptadores que cualquier otro proyecto o producto en el mundo. Meshery utiliza adaptadores para gestionar las distintas mallas de servicio."
language: es
lang: es
categories: es
type: concepts
list: include
---

Como administrador de múltiples mallas, Meshery ofrece soporte para más adaptadores que cualquier otro proyecto o producto en el mundo. Meshery utiliza adaptadores para gestionar las distintas mallas de servicio.

## ¿Qué son los Adaptadores de Meshery?

Los adaptadores permiten a Meshery interactuar con las diferentes mallas de servicio, exponiendo su valor diferenciado a los usuarios.

Meshery tiene adaptadores para gestionar las siguientes mallas de servicio.
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Adapter Status |  Service Mesh  | Service Mesh Version | Port          |
| :------------: | :------------ | :------------: | :------------ |
{% for adapter in sorted -%}
{% if adapter.project_status -%}
| {{ adapter.project_status }} | <img src="{{ adapter.image }}" style="width:20px" /> [Adaptador Meshery para {{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.version }} | {{ adapter.port }} |
{% endif -%}
{% endfor %}

### FAQs sobre Adaptadores

#### ¿Se hace igual cada adaptador de malla de servicio?
No, se escriben diferentes adaptadores de malla de servicios para exponer el valor único de cada malla de servicios. En consecuencia, no son igualmente capaces ya que cada malla de servicio no es igualmente capaz que la otra.

Los adaptadores tienen un conjunto de operaciones que se agrupan en función de tipos de operaciones predefinidos. Vea la página [extensibility]({{site.baseurl}}/extensibility) para más detalles de las operaciones de adaptadores.

#### ¿Cómo puedo crear un nuevo adaptador?

Vea la documentación [extensibility]({{site.baseurl}}/extensibility) para los detalles sobre como son hechos los nuevos adaptadores Meshery.

#### ¿Puedo ejecutar más de una instancia del mismo adaptador Meshery?
La configuración predeterminada de una implementación de Meshery incluye una instancia de cada uno de los adaptadores de Meshery (que han alcanzado un estado de versión estable). Puede optar por ejecutar varias instancias del mismo tipo de adaptador Meshery; p.ej. dos instancias del adaptador `meshery-istio`. Para hacerlo, modifique ~/.meshery/meshery.yaml para incluir múltiples copias del adaptador dado.

Vea la guía "[Multiple Adapters]({{site.baseurl}}/guides/multiple-adapters)" para más información.
