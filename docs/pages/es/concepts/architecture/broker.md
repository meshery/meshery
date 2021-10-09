---
layout: default
title: Broker
permalink: es/concepts/architecture/broker
type: concepts
redirect_from: architecture/broker
abstract: "El componente Meshery broker facilita la transmisión de datos entre los componentes del clúster de Kubernetes y el mundo exterior"
language: es
list: include
---

El agente es un controlador personalizado de Kubernetes que provee transmisión de datos a través de componentes independientes de Meshery ya sea que estos componentes estén ejecutándose dentro o fuera del clúster de Kubernetes.

### Preguntas más Frecuentes del Agente

#### ¿Cuántos agentes pueden ejecutarse?
Se recomienda ejecutar una instancia de agente para cada clúster de Kubernetes. Sin embargo, la instancia en sí puede ser escalado de acuerdo al volumen de datos entrantes en cada clúster. El escalamiento es independiente de la cantidad de instancias que se ejecutan.

#### ¿Cómo se ve una configuración de Alta Disponibilidad?
Aprovechamos la funcionalidad de Kubernetes en términos del comportamiento de Alta Disponibilidad. Lo que significa que la instancia del agente se instancía/reinicia por sí misma cuando se produce un problema. En parte, Meshery Operator es también responsable de mantener al agente funcional.

#### ¿Qué características de estado tiene el Agente?
Todos los mensajes que se publican en el agente se persisten en memoria dentro de la instancia del agente hasta que se consuma. El Volumen de persistencia/Espacio del disco no está siendo utilizado actualmente por el Agente.

#### ¿Cómo puedo saber si el Agente está funcionando?, ¿cómo puedo solucionar problemas del Agente?
La instancia del Agente es desplegada dentro del clúster de Kubernetes como un `Statefulset`. En el caso de que el agente parezca que no funciona, aquí están algunos pasos para solucionar problemas de la instancia:

- Asegúrese de que los pods correspondientes al `Statefulset` se encuentren en funcionamiento.
- Asegúrese de que el clúster de Kubernetes tiene soporte para Kubernetes `Service` tipo `LoadBalancer` o `NodePort`.
- Asegurar la conectividad entre el Meshery Server y el endpoint del servicio del Agente.
