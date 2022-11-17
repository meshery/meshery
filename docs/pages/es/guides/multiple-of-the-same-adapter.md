---
layout: default
title: Usar varios adaptadores
permalink: es/guides/multiple-adapters
language: es
type: Guides
---

## Configuración Avanzada

Meshery es capaz de ejecutar cero o más adaptadores de malla de servicio. Sin ningún adaptador de malla de servicio, algunas, pero no todas, las características de Meshery seguirán funcionando (por ejemplo, pruebas de rendimiento de cargas de trabajo que no se ejecutan en una malla de servicio).

### Modificación de la configuración de implementación del adaptador predeterminada

La cantidad de adaptadores, el tipo de adaptadores, dónde se implementan, cómo se nombran y en qué puerto están expuestos son todas opciones de implementación configurables. Para modificar la configuración predeterminada, busque _~/.meshery/meshery.yaml_ en su sistema. _~/.meshery/meshery.yaml_ es un archivo Docker Compose.

#### Configuración: Ejecutar menos Adaptadores Meshery

En el archivo de configuración _~/.meshery/meshery.yaml_, remover la linea(s) del adaptador(es) que quiere eliminar de su sistema.

#### Configuración: Ejecutar más de una instancia del mismo adaptador Meshery

La configuración predeterminada de una implementación de Meshery incluye una instancia de cada uno de los adaptadores de Meshery (que han alcanzado un estado de versión estable). Puede optar por ejecutar varias instancias del mismo tipo de adaptador Meshery; p.ej. dos instancias del adaptador _meshery-istio_. Para hacerlo, modifique _~/.meshery/meshery.yaml_ para incluir múltiples copias del adaptador dado.

Demostración de Meshery administrando implementaciones de malla de servicios en múltiples clústeres:

<iframe class="container" width="560" height="315" src="https://www.youtube.com/embed/yWPu3vq4vEs?start=5041" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

Vea en YouTube: [Cloud Native Austin Virtual Meetup: April 2020](https://youtu.be/yWPu3vq4vEs?t=5041&list=PL3A-A6hPO2IOpTbdH89qR-4AE0ON13Zie)
