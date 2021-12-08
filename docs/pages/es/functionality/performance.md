---
layout: page
title: Gestión del rendimiento
permalink: es/functionality/performance-management
type: functionality
language: es
---

La clave del funcionamiento eficiente de cualquier service mesh es la medición y gestión de su rendimiento.

## Generadores de carga

Meshery proporciona a los usuarios la opción de elegir qué generador de carga prefieren usar para una prueba de rendimiento determinada. Los usuarios pueden establecer su configuración basada en su propia preferencia de generador de carga diferente al generador de carga por defecto.

Meshery soporta los siguientes generadores de carga y es [extensible](extensibility) para soportar otros:

- Fortio
- wrk2
- NightHawk

### Fortio

Fortio es un rápida, pequeña (imagen de Docker de 3Mb, dependencias mínimas), reusable, biblioteca de go integrable, así como también una herramienta de línea de comandos y un proceso de servidor. El servidor incluye una Interfaz de Usuario web simple y una representación gráfica de los resultados (un gráfico de latencia único y gráficos comparativos de min, max, avg, qps y percentiles gráficos).

### wrk2

Es una herramienta moderna de evaluación comparativa HTTP capaz de generar una carga significativa cuando se ejecuta en un solo CPU multi-núcleo. Combina un diseño multiproceso con sistemas de notificación de eventos escalables como lo son epoll y kqueue.

### NightHawk

NightHawk es una herramienta de caracterización L7 (HTTP/HTTPS/HTTP2). Actualmente ofrece:

- Un cliente de prueba de carga que soporta HTTP/1.1 y HTTP/2 sobre HTTP y HTTPS (los certificados HTTPS aún no se validan).
- Un servidor de prueba simple capaz de generar tamaños de respuestas dinámicas, así como también inyectar retrasos.
- Un binario para tranformar la salida de NightHawk a formatos conocidos, permitiendo la integración con otros sistemas y paneles.

## Node y métricas del service mesh

Meshery proporciona los resultados de pruebas de rendimiento junto a las métricas del entorno, incluyendo el control del service mesh y las métricas del plano de datos. También incluye las métricas de los recursos del nodo del clúster, para que los operadores puedan comprender fácilmente la sobrecarga del plano de control y el plano de datos de su service mesh en el contexto de la sobrecarga incurrida en los nodos dentro del clúster.

## Grafana y Meshery

Conecte Meshery a su instancia de Grafana existente y Meshery importará los tableros que elija.

<a href="/assets/img/performance-management/meshery-and-grafana.png">
    <img src="/assets/img/performance-management/meshery-and-grafana.png" style="width: 100%" />
</a>

### Conexión a Grafana

Si tiene una clave de API configurada para restringir el acceso a sus tableros de Grafana, deberá ingresar la clave de API cuando establezca la conexión de Meshery con Grafana.

- Importación de tableros de Grafana
  - Importación de tablero existente de Grafana a través de API
  - Importación de tablero personalizado de Grafana a través de yaml
- Configuración de las preferencias del panel de gráficos

## Prometheus y Meshery

Meshery permite a los usuarios conectarse a una o más instancias de Prometheus para recopilar datos telemétricos (en forma de métricas). Estas métricas pueden pertenecer al service mesh, Kubernetes, aplicaciones en la malla o realmente... cualquier métrica que Prometheus haya recolectado.

Una vez que haya conectado Meshery a su (s) despliegue (es) de Prometheus, puede realizar pruebas de conectividad ad-hoc para verificar la comunicación entre Meshery y Prometheus.

## Lectura sugerida

- Guía: [Interpretación de los resultados de la prueba de rendimiento]({{site.baseurl}}/guides/interpreting-performance-test-results)
