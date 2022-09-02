---
layout: page
title: Descripción General
permalink: es/overview
language: es
---

# Introducción a Meshery

El plano de gestión de service mesh adopta, opera y desarrolla sobre diferentes service meshes.
Meshery facilita el aprendizaje sobre la funcionalidad y el desempeño de los service mesh e incorpora la recopilación y visualización de métricas de las aplicaciones que se ejecutan dentro de un service mesh o multiples service meshes.
Meshery proporciona estas funcionalidades de alto nivel:

1. Gestión del desempeño de service meshes.
2. Gestión de configuración de service meshes.
   - Mejores prácticas de configuración.
3. Gestión del ciclo de vida de service meshes.
4. Interoperabilidad y asociación de service meshes.

<iframe class="container" width="560" height="315" src="https://www.youtube.com/embed/CFj1O_uyhhs" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

<div style="text-align:center;width:100%"><emphasis>Presentado en Service Mesh Day 2019</emphasis></div>

<h2>¿Qué desafíos resuelve Meshery?</h2>
<b>Gestión de service mesh - uno o múltiples service meshes.</b>

Siempre que haya que responder a estas preguntas sobre el desempeño, son subjetivas a la carga de trabajo y la infraestructura específicas utilizadas para la medición. Ante este desafío, el proyecto Envoy, por ejemplo, se niega a publicar datos de desempeño porque tales pruebas pueden ser:

- Confusas
- Malinterpretadas.

Más allá de la necesidad del desempeño y el alto consumo de datos bajo una permutación de diferentes cargas de trabajo (aplicaciones), los tipos y tamaños de los recursos de infraestructura, la necesidad de un proyecto multifuncional, y las comparaciones entre semejantes, son anhelados para facilitar una comparación de las diferencias de comportamiento entre service meshes y la selección de su uso. Los proyectos individuales son reservados en publicar resultados de pruebas de otras competencias de service mesh. Es necesario un análisis independiente, imparcial y creíble.

La intención de Meshery es ser un proveedor y una utilidad neutra de proyectos para evaluar de manera uniforme el desempeño de los service mesh. Entre proyectos de service mesh y servicios de proxy (y sorprendentemente, dentro de un mismo proyecto), existen diferentes números de herramientas y resultados. Meshery permite elegir un conjunto eficiente de herramientas para su ecosistema al proporcionar evaluaciones y métricas de desempeño.

1. Al aprovechar Meshery, podrás lograr una comparación de desempeño de los service meshes semejantes.
2. Realiza un seguimiento del desempeño de un service mesh de un release a otro.
3. Comprende las diferencias de comportamiento entre los service meshes.
4. Realiza un seguimiento del desempeño de su aplicación de una versión a otra.

## Meshery es para Entusiastas y Operadores

Sea la opción para implementar desde el Día 0 o el mantenimiento del Día 2 de despliegue, Meshery tiene capacidades útiles para cualquier circunstancia. La audicencia dirigida para un proyecto con Meshery sería cualquier tipo de operador que aproveche el service mesh en su ecosistema; incluyendo desarrolladores, ingenieros de devops, quienes toman decisiones, arquitectos, y organizaciones que dependen de una plataforma de microservicios..

## Meshery es para la gestión de desempeño: pruebas y comparaciones de mercado

Meshery ayuda a los usuarios a sopesar el valor de la implementación de un service mesh con la alta carga involucrada a la ejecución del mismo. Meshery proporciona un análisis estadístico de la latencia de las solicitudes y el rendimiento que se observan a través de las permutaciones de su carga de trabajo, infraestructura y la configuración del service mesh.
Además de solicitar la latencia y el rendimiento, Meshery también rastrea la sobrecarga de memoria y el CPU en los nodos del clúster. Mide el plano de datos y el plano de control con diferentes conjuntos de cargas de trabajo e infraestructuras.

<a href="https://raw.githubusercontent.com/meshery/meshery/master/docs/assets/img/readme/meshery_lifecycle_management.png"><img alt="Layer5 Service Mesh Community" src="{{ site.baseurl }}{% link assets/img/readme/meshery_lifecycle_management.png %}"  width="100%" align="center"/></a>
Establece un desempeño de marca y rastrea el desempeño como estándar a medida que cambia en el tiempo.

## Meshery es para cualquier service mesh
La diversidad de infraestructuras es una realidad para cualquier organización. Ya sea que se esté ejecutando un solo service mesh o varios tipos de service meshes, en Meshery encontrarás que soporta una diversidad de infraestructura (o bien, la falta de una).


- **Adaptadores disponibles de service mesh** - Adaptadores de service mesh soportados por Mashery.

| Plataforma    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Estado        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "stable" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }})                     |       {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
<br>
- **Adaptadores en progreso de service mesh** - Adaptadores de service mesh sometidos por la comunidad y en desarrollo

| Plataforma    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Estado        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "beta" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
<br>
- **Adaptadores de service mesh que buscamos apoyo** - Adaptadores de service mesh que se encuentran en búsqueda de ayuda de la comunidad.

| Plataforma        &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;     | Estado        |
| :------------               | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "alpha" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

## Comunidad
Este proyecto es hecho por la comunidad y toda colaboración es bienvenida! [Fork here on Github](https://github.com/layer5io/meshery)

- Únete a [la reunión semanal de la comunidad](https://docs.google.com/document/d/1c07UO9dS7_tFD-ClCWHIrEzRnzUJoFQ10EzfJTpS7FY/edit?usp=sharing) [Viernes de 10am a 11am Horario Central](/assets/projects/meshery/Meshery-Community-Meeting.ics).

  - Observa las [grabaciones de la comunidad](https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0) y suscríbete al [calendario de la comunidad](https://bit.ly/2SbrRhe).

- Ingresa al [drive de la comunidad](https://layer5.io/newcomers) (solicita acceso).

# FAQ

## ¿Por qué usar Meshery?

- Porque es un proyecto de código abierto, proveedor neutral de proyectos que facilita las pruebas en meshes.
- Porque fortio no está empaquetado en una utilidad de prueba de mesh, sino que es solo un generador de carga en sí mismo.
- Debido a que regpatrol es de código cerrado, el código binario no es liberado, embebido para un mesh y es producido por un proveedor de ese mesh.

## ¿Por qué crear Meshery y no usar otra herramienta de marca?

Meshery está diseñado específicamente para facilitar la evaluación comparativa de los service mesh y sus cargas de trabajo, al contrario de otras herramientas. Sin embargo también existen otras herramientas que se utilizan para la evaluación comparativa de los service meshes, como regpatrol. Regpatrol es utilizado por IBM, el cual no es código abierto ni está disponible en su forma binaria para su uso, posee las siguientes diferencias con Meshery:

- Telemetría: regpatrol obtiene la telemetría del adaptador Mixer Prometheus y utiliza el agente de nodo propietario de IBM.
- Meshery se obtiene del adaptador Mixer Prometheus y utiliza el exportador de nodos Prometheus.
- Tipo de tráfico: regpatrol usa JMeter, que puede analizar las respuestas y realizar pruebas funcionales.
- Meshery está usando fortio, que es solo para pruebas de generación y rendimiento.

# Recursos

## Presentaciones de Meshery

- [O'Reilly OSCON 2020](https://conferences.oreilly.com/oscon/oscon-or)
- [O'Reilly Infrastructure & Ops 2020](https://conferences.oreilly.com/infrastructure-ops/io-ca/public/schedule/speaker/226795)
- [InnoTech Dallas 2020](https://innotechdallas2020.sched.com/event/aN7E/a-management-plane-for-service-meshes)
- [KubeCon EU 2020](https://kccnceu20.sched.com/event/Zetg/discreetly-studying-the-effects-of-individual-traffic-control-functions-lee-calcote-layer5?iframe=no&w=100%&sidebar=yes&bg=no)
- DockerCon 2020 ([deck](https://calcotestudios.com/talks/decks/slides-dockercon-2020-service-meshing-with-docker-desktop-and-webassembly.html), [video](https://www.youtube.com/watch?v=5BrbbKZOctw&list=PL3A-A6hPO2IN_HSU0pSfijBboiHggs5mC&index=4&t=0s))
- [Open Source 101 at Home](https://calcotestudios.com/talks/decks/slides-open-source-101-at-home-solving-the-service-mesh-adopters-dilemma.html)
- [Docker Captains Roundtable 2020](https://calcotestudios.com/talks/decks/slides-docker-captains-2020-meshery-the-multi-service-mesh-manager.html)
- [Cloud Native Austin 2020](https://www.meetup.com/Cloud-Native-Austin/events/267784090/)
- NSMCon 2019 talk ([video](https://www.youtube.com/watch?v=4xKixsDTtdM), [deck](https://calcotestudios.com/talks/decks/slides-nsmcon-kubecon-na-2019-adopting-network-service-mesh-with-meshery.html))
- [Service Mesh Day 2019](https://youtu.be/CFj1O_uyhhs)
- [DockerCon 2019 Open Source Summit](https://www.docker.com/dockercon/2019-videos?watch=open-source-summit-service-mesh)
- KubeCon EU 2019 ([video](https://www.youtube.com/watch?v=LxP-yHrKL4M&list=PLYjO73_1efChX9NuRaU7WocTbgrfvCoPE), [deck](https://calcotestudios.com/talks/decks/slides-kubecon-eu-2019-service-meshes-at-what-cost.html))
- [KubeCon EU 2019 Istio Founders Meetup](https://calcotestudios.com/talks/decks/slides-istio-meetup-kubecon-eu-2019-istio-at-scale-large-and-small.html)
- [Cloud Native Rejekts EU 2019](https://calcotestudios.com/talks/decks/slides-cloud-native-rejekts-2019-evaluating-service-meshes.html)
- [Container World 2019](https://calcotestudios.com/talks/decks/slides-container-world-2019-service-meshes-but-at-what-cost.html)
- Solving the Service Mesh Adopter’s Dilemma ([deck](https://calcotestudios.com/talks/decks/slides-open-source-101-at-home-solving-the-service-mesh-adopters-dilemma.html), [event](https://opensource101.com/sessions/solving-the-service-mesh-adopters-dilemma/),[video](https://www.youtube.com/watch?v=Q1zSWbO0RmI&list=PL3A-A6hPO2IN_HSU0pSfijBboiHggs5mC&index=2&t=0s))

## Otros Recursos

- [Comparación de Service Meshes](https://layer5.io/landscape)
- [Herramientas de Service meshes](https://layer5.io/landscape#tools)
- [Libros sobre Service Mesh](https://layer5.io/books)
- [Workshops sobre Service Mesh](https://layer5.io/workshops)

