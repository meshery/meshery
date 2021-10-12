---
layout: default
title: Adaptadores
permalink: es/concepts/architecture/adapters
type: concepts
redirect_from: es/architecture/adapters
abstract: "Meshery ofrece soporte para más adaptadores que cualquier otro proyecto o producto en el mundo. Meshery utiliza adaptadores para gestionar las distintas mallas de servicio."
language: es
list: include
---

Como administrador de múltiples mallas, Meshery ofrece soporte para más adaptadores que cualquier otro proyecto o producto en el mundo. Meshery utiliza adaptadores para gestionar las distintas mallas de servicio.

## que son el adaptador mashery?

Los adaptadores permiten a Meshery interactuar con las diferentes mallas de servicio, exponiendo su valor diferenciado a los usuarios.

Meshery tiene adaptadores para administrar las siguientes mallas de servicio.
{% assign sorted = site.adapters | sort: "project_status" | reverse %}

| Estado del adaptador | Malla de servicio | Versión de Service Mesh | Puerto |
| :------------------: | :---------------- | :---------------------: | :----: |
{% for adapter in sorted -%}
{% if adapter.project_status -%}
| {{ adapter.project_status }} | <img src="{{ adapter.image }}" style="width:20px" /> [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.version }} | {{ adapter.port }} |
{% endif -%}
{% endfor %}

### Preguntas frecuentes sobre adaptadores

#### Todos los adaptadores de malla de servicio son iguales?

No, se escriben diferentes adaptadores de malla de servicios para exponer el valor único de cada malla de servicios. En consecuencia, no son igualmente capaces ya que cada malla de servicio no es igualmente capaz que la otra.

Los adaptadores tienen un conjunto de operaciones que se agrupan en función de tipos de operaciones predefinidos. Ver el [extensibilidad]({{site.baseurl}}/extensibility) página para obtener más detalles sobre las operaciones del adaptador.

#### Cómo puedo crear un nuevo adaptador?

Ver el [extensibilidad]({{site.baseurl}}/extensibility) documentación para obtener detalles sobre cómo se fabrican los nuevos adaptadores Meshery.

#### Puedo ejecutar más de una instancia del mismo adaptador Meshery?

La configuración predeterminada de una implementación de Meshery incluye una instancia de cada uno de los adaptadores de Meshery (que han alcanzado un estado de versión estable). Puede optar por ejecutar varias instancias del mismo tipo de adaptador Meshery; p.ej. dos instancias del `meshery-istio` adaptador. Para hacerlo, modifique ~/.meshery/meshery.yaml para incluir varias copias del adaptador dado.

Ver el "[Adaptadores múltiples]({{site.baseurl}}/guides/multiple-adapters)" guía para obtener más información.
