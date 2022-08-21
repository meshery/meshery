---
layout: page
title: Capacidad de cumplimiento de SMI
permalink: es/functionality/smi-conformance
type: functionality
language: es
---

## Conformidad

Definiendo “Conformidad” - Es importante reconocer que la conformidad consiste en capacidades y estado de conformidad. Nosotros definimos la conformidad como una combinación de estos 2 conceptos.

1. La conformidad de SMI reconoce que
   ...algunos service mesh participantes nunca implementan completamente funciones a conciencia (especificaciones SMI)...

2. La conformidad SMI identifica
   ...una diferencia entre la completa implementación de una especificación y el cumplimiento de las partes que implementa...

## Capacidad

Dado que algunas implementaciones del service mesh nunca tienen la intención de implementar completamente las especificaciones de SMI, por cada prueba individual, existen tres posibles designaciones de capacidad.

- Completa - los service mesh tienen esta capacidad.

- Parcial - los service mesh tienen implementada una porción de esta capacidad (puede o no tener esta capacidad completamente en el futuro).

- Ninguna - actualmente los service mesh no tienen esta capacidad (puede o no tener esta capacidad en el futuro).

## Enfoque de las pruebas de conformidad

Cada versión de Kubernetes, versión de service mesh y categoría de SMI serán sometidas a varias pruebas. Cada prueba se llevará a cabo de forma automatizada y concurrente, principalmente invocando a Meshery para ejecutar las pruebas de conformidad.

#### Estos son los pasos

- Configurar una versión específica de un service mesh
- [Opcional] Configurar Prometheus
- Desplegar una aplicación de muestra elegida
- Desplegar el operador SMI necesario, como smi-metrics
  la prueba se ejecutará, lo que puede implicar la llamada a un punto final de API
- Validar la respuesta

No todas las pruebas pueden ser validadas simplemente usuando la respuesta, en estos casos podemos también conectarnos a una instancia de Prometheus, el cual está preconfigurado para recolectar todas las métricas para la prueba, y usar estas métricas para definir expectativas. Trabajaremos en ejemplos específicos.

Finalmente, los resultados de las ejecuciones de pruebas de Meshery son persistentes (como repositorios de GitHub) y publicados en la página web de conformidad. Eventualmente, podemos construir un sistema el cual nos permitirá ejecutar de forma granular pruebas en demanda para una versión elegida de Kubernetes, service mesh y operador de SMI.

## Definiciones de las pruebas de conformidad

Las pruebas de conformidad son clasificadas por tipo de especificación SMI. Un conjunto de pruebas son definidas para cada especificación SMI. Dentro de cada conjunto de pruebas, se definen dos tipos de pruebas de aserción: aserción de presencia y aserción de capacidad.

## Validación de conformidad

La conformidad con las especificaciones de SMI se realizarán a través del aprovisionamiento automatizado de service mesh individuales y el despliegue de una carga de trabajo en común. Se utiliza una aplicación de muestra sencilla e instrumentada como carga de trabajo para probar.

## Definición de conformidad

La conformidad con las especificaciones de SMI es definida como una serie de pruebas de aserción. Una prueba de aserción es una condición que debe probarse para confirmar la conformidad con un requisito. Una prueba de aserción es una condición que desde la perspectiva de las pruebas de validación, determinan la conformidad requerida para que se prueben varias condiciones. La colección de pruebas de aserción clasificadas por la especificación de SMI define colectivamente el conjunto de pruebas de conformidad de SMI. Meshery es el mejor arnés de prueba utilizado para ajustar las pruebas de conformidad de SMI para diferentes service mesh y diferentes cargas de trabajo.

## Pasos para ejecutar pruebas de rendimiento

### Precondiciones

- La capacidad de un determinado service mesh para adherirse a una especificación de SMI es validada ejecutando una carga de trabajo sobre el service mesh.
- Los depliegues de cargas de trabajo son configurados de forma específica para los requisitos de incorporación del service mesh dado.
- Las pruebas son definidas para validar la conformidad para cada tipo de especificación de SMI (por ejemplo, métricas, acceso, tráfico...).

### Invocación

- Las pruebas de aserción son definidas de acuerdo a una carga de trabajo específica y son desplegadas con las cargas de trabajo que se están probando (las pruebas de aserción están empaquetadas).
- Un resultado de prueba es recolectado con la evaluación de cada aserción.
- Los resultados de futuras pruebas se transmitirán individualmente a Meshery después de que se evalúe cada aserción.
- Una vez que todas las aserciones son evaluadas, los resultados de las pruebas son retornadas a Meshery para su representación visual.

## Informe de conformidad

### Procedencia de los resultados de pruebas

A cada proyecto de service mesh participante se le pedirá que incorpore la herramienta de conformidad, Meshery, en sus pipelines de CI o, alternativamente, que ejecute manualmente el conjunto de pruebas de conformidad cuando se realice un despliegue del proyecto de service mesh. La herramienta de conformidad ejecutará el conjunto de pruebas y automáticamente actualizará el panel de conformidad cuando la prueba sea completada.
Para garantizar la procedencia de los resultados de las pruebas que representan un service mesh determinado, se pedirá a cada proyecto que identifique uno o más cuentas de GitHub que se utilizarán para publicar las pruebas. Idealmente, esta cuenta de GitHub es una cuenta servicerobot que se usa dentro del pipeline de CI del proyecto.

Este método de proporcionar la verification de resultados es usado de manera similar para esos mismos proyectos de service mesh que también usan Meshery, para proporcionar sus resultados de pruebas de rendimiento.

Cada proyecto de service mesh necesita identificar su cuenta de servicerobot, actualizando esta lista httpsmeshery.iosmi-conformancesm-service-accounts. Identifique un usuario proveedor de Meshery determinado y designe su "cuenta de servicio de CI".

## Publicación de resultados de pruebas

Un informe de cara al público mostrará el estado actual e histórico de la capacidad del service mesh individual con cada una de las especificaciones de SMI. El informe será de natuaraleza visual, pero también está disponible como yaml. El panel se publicará públicamente aquí httpsmeshery.iosmi-conformance (actualmente, listado en httpslayer5.iolandscape#smi).

### Conformidad histórica

Los informes rastrearán el historial de las versiones del service mesh, las versiones de las especificaciones SMI y su compatibilidad.
