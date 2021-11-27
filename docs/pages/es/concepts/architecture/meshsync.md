---
layout: default
title: MeshSync
permalink: es/concepts/architecture/meshsync
type: concepts
redirect_from: architecture/meshsync
abstract: "Meshery ofrece soporte para la sincronización de estado de clúster/service mesh de Kubernetes con la ayuda de MeshSync."
language: es
list: include
---

<p style="display:block">
<img src="{{site.baseurl}}/assets/img/meshsync/meshsync.svg" align="left" 
    style="margin-right:1rem;margin-bottom:.5rem;" width="20%" />

MeshSync es un controlador de Kubernetes personalizado que provee descubrimiento escalonado y sincronización continua con Meshery Server en cuanto al estado del clúster de Kubernetes, service meshes y sus cargas de trabajo.

</p>

La instantánea de trabajo de MeshSync del estado de cada clúster y service mesh bajo administración se almacena en la memoria y se actualiza continuamente.

MeshSync es administrado por el <a href="{{site.baseurl}}/concepts/architecture/operator">Meshery Operator</a>.

## ¿Cuáles son las responsabilidades principales de MeshSync?

### Descubrir los recursos existentes dentro del clúster de Kubernetes

Los recursos que están presentes dentro del clúster se descubren eficientemente con la ayuda de pipelines. Los datos se contruyen en un formato particular específico para Meshery y se publican a través de diferentes partes de la arquitectura.

### Escuchar los eventos/cambios en cada componente

MeshSync implementa varios informadores/oyentes en cada recurso para escuchar cambios que ocurren en ellos. Estos se reflejan en tiempo real y se actualizan en sus respectivas áreas.

## Preguntas más Frecuentes de MeshSync
