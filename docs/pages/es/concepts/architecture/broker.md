---
layout: default
title: Broker
permalink: /es/concepts/architecture/broker
type: concepts
redirect_from: architecture/broker
abstract: "El componente broker de Meshery facilita la transmisión de datos entre los componentes del clúster de Kubernetes y el mundo exterior."
language: es
list: include
---

Broker es un controlador personalizado Kubernetes que provee transmisión de datos a través de componentes independientes de Meshery, ya sea que esos componentes se ejecuten dentro o fuera del clúster de Kubernetes.

### Preguntas Frequentes del Broker

#### ¿Cuántos Brokers pueden ejecutarse?
Se recomienda ejecutar una instancia de agente para cada clúster de kubernetes. Sin embargo, la instancia en sí se puede escalar en función del volumen de datos entrantes en cada uno de los clústeres. El escalado es independiente del número de instancias en ejecución.

#### ¿Cómo se ve una configuración HA?
Aprovechamos la funcionalidad de kubernetes en términos del comportamiento de alta disponibilidad(High Avaliavility - HA). Es decir, la instancia del agente se crea una instancia / se reinicia por sí sola cuando ocurre un problema. En parte, Meshery-Operator también es responsable de mantener al corredor en funcionamiento.

#### ¿Qué caracteristicas con estado tiene el Broker?
Todos los mensajes que se publican en el intermediario se conservan en la memoria dentro de la instancia del intermediario hasta que se consumen. El Broker no está utilizando actualmente el volumen persistente / espacio en disco.

#### ¿Cómo saber si el Broker está funcionando? ¿Cómo solucionar problemas con el Broker?
La instancia de Broker se implementa dentro del clúster de kubernetes como un "Statefulset". En el caso de que el agente no parezca funcionar, aquí hay algunos pasos para solucionar el problema de la instancia:

- Asegúrese de que los pods correspondientes al `Statefulset` estén en funcionamiento.
- Asegúrese de que el clúster de kubernetes sea compatible con el tipo de `Servicio` de kubernetes` LoadBalancer` o `NodePort`.
- Asegure la conectividad entre Meshery-Server y el punto final del servicio Broker.

