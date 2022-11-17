---
layout: default
title: Interpretación de los resultados de la prueba de rendimiento
description: Esta guía tiene como objetivo ayudar a los usuarios a obtener información sobre cómo deberían verse los resultados de las pruebas de rendimiento.
permalink: es/guides/interpreting-performance-test-results
type: guides
language: es
---

Una vez que se prueba el rendimiento y se ejecuta, se realiza el análisis disco y los resultados de la prueba persisten. Resultados o descargables y persistentes en el [Desempeño de la Malla de Servicio](https://smp-spec.io/) (SMP) format.

## Vista Gráfica

<img src="{{ site.baseurl }}/assets/img/performance-management/dashboard.png" />

<img src="{{ site.baseurl }}/assets/img/performance-management/chart.png" />

## Pruebas

Para obtener la mejor experiencia, debe comparar entre dos o más pruebas de configuración similar. Las pruebas que están configuradas con un alto grado de variación (por ejemplo, una prueba se ejecutó durante 5 minutos, mientras que otra prueba se ejecutó durante 1 hora) producirán comparaciones de las que es más difícil extrapolar información.

La latencia y el rendimiento son las dos señales examinadas más significativamente.

Meshery utilizará diferentes algoritmos para calcular los resultados según el generador de carga que se utilizó para ejecutar la prueba.
