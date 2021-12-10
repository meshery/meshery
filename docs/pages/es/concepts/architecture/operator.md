---
layout: default
title: Operador
permalink: es/concepts/architecture/operator
type: concepts
redirect_from: architecture/operator
abstract: "Meshery Operator controla y administra el ciclo de vida de componentes desplegados dentro de un clúster de Kubernetes"
language: es
display-title: "false"
list: include
---

# Meshery Operator <img style="width:1em; inline; margin-bottom:10px;" src="{{ site.baseurl }}/assets/img/architecture/B203EFA85E89491B.png"/>

Meshery Operator es el operador multi-service mesh (un controlador personalizado de Kubernetes) que administra MeshSync y su agente de mensajes.

Meshery Operator es un administrador de controladores de Kubernetes, conocido como un Operador de Kubernetes. Meshery Operator administra el ciclo de vida de cualquier componente de Meshery que es desplegado o ejecutado dentro de un clúster de Kubernetes.

## Despliegues

Es recomendable desplegar un Meshery Operator por clúster.

[![Meshery Operator y MeshSync]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg
)]({{ site.baseurl }}/assets/img/architecture/meshery-operator-and-meshsync.svg)

### Secuencia de inicialización

[![Meshery Operator y MeshSync]({{ site.baseurl }}/assets/img/architecture/meshery-operator-deployment-sequence.svg
)]({{ site.baseurl }}/assets/img/architecture/meshery-operator-deployment-sequence.svg)

## Controladores administrados por Meshery Operator

### Controlador de Agente

El Meshery broker es uno de los componentes centrales de la arquitectura de Meshery. Este controlador administra el ciclo de vida del agente que Meshery usa para la transmisión de datos a través del clúster y el mundo exterior.

Ver [Meshery Broker]({{site.baseurl}}/architecture/broker) para más información.

### Controlador de MeshSync

El controlador de MeshSync administra el ciclo de vida de MeshSync que es desplegado para la sincronización de recursos para el clúster.

Ver [MeshSync]({{site.baseurl}}/architecture/meshsync) para más información.

## Preguntas más Frecuentes del Operador

### ¿Cuándo es desplegado Meshery Operator y cuándo es eliminado?  
Como un controlador personalizado de Kubernetes, Meshery Operator es provisionado y desprovisionado cuando el Meshery Server está conectado o desconectado del clúster de Kubernetes. Las conexiones del Meshery Operator al clúster de Kubeernetes se controlan utilizando clientes de Meshery Server: `mesheryctl` o Meshery UI. Este comportamiento descrito a continuación es consistente si su despliegue de Meshery está usando Docker o Kubernetes como la plataforma para hospedar el despliegue de Meshery.

**CLI de Meshery**
`mesheryctl` inicia la conexión del clúster de Kubernetes cuando se ejecuta `mesheryctl system start` y desconectado cuando se ejecuta `mesheryctl system stop`. Este comportamiento es consistente si su despliegue de Meshery está usando Docker o Kubernetes como la plataforma para hospedar el despliegue de Meshery.

**UI de Meshery**
Meshery UI ofrece mayor control granular sobre el despliegue de Meshery Operator en el cual puede eliminar Meshery Operator de un clúster de Kubernetes sin desconectar Meshery Server del clúster de Kubernetes. Puede controlar el despliegue de Meshery Operator usando el interruptor de encendido/apagado que se encuentra en la sección de Configuración de Meshery Operator.

### ¿El Meshery Operator usa un SDK o framework?
Sí, Meshery Operator usa el Operator SDK.
