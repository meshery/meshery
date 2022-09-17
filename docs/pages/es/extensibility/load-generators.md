---
layout: default
title: "Extensibilidad: Generadores de Carga"
permalink: es/extensibility/load-generators
type: extensibility
#redirect_from: architecture/adapters
abstract: "Meshery ofrece soporte para más adaptadores que cualquier otro proyecto o producto en el mundo. Meshery utiliza adaptadores para gestionar las distintas mallas de servicio."
language: es
list: include
---

Los usuarios pueden preferir usar un generador de carga sobre el siguiente dada la diferencia de capacidades entre los generadores de carga, por lo que Meshery proporciona una "interfaz de generador de carga" (una interfaz gRPC) detrás de la cual se puede implementar un generador de carga. Meshery ofrece a los usuarios la opción de elegir qué generador de carga prefieren utilizar para una prueba de rendimiento determinada. Los usuarios pueden configurar su propia preferencia de generador de carga diferente al generador de carga predeterminado.

### ¿Qué función cumplen los generadores de carga en Meshery?

Los generadores de carga proporcionarán la capacidad de ejecutar pruebas de carga desde Meshery. A partir de hoy, los generadores de carga están integrados como bibliotecas en Meshery y Meshery invoca las API de los generadores de carga con las opciones de prueba de carga adecuadas para ejecutar la prueba de carga. Por el momento, Meshery tiene soporte para generadores de carga HTTP. El soporte para pruebas de carga de gRPC y TCP está en la hoja de ruta. Meshery tiene integración funcional con fortio, wrk2 y nighthawk.

### ¿Por qué admitir varios generadores de carga?

Diferentes casos de uso y diferentes opiniones requieren diferentes enfoques para el análisis estadístico de los resultados de rendimiento. Por ejemplo, wrk2 da cuenta de un concepto llamado Omisión Coordinada.

### ¿Qué generadores de carga admite Meshery?

1. [fortio](https://github.com/fortio/fortio) - Biblioteca de pruebas de carga Fortio, herramienta de línea de comandos, servidor de eco avanzado y UI web en go (golang). Permite especificar una carga establecida de consulta por segundo y registrar histogramas de latencia y otras estadísticas útiles.
1. [wrk2](https://github.com/giltene/wrk2) - Un rendimiento constante, una variante de wrk de grabación de latencia correcta.
1. [nighthawk](https://github.com/envoyproxy/nighthawk) - Permite a los usuarios ejecutar pruebas de rendimiento distribuidas para imitar mejor los escenarios de sistemas distribuidos del mundo real.

- Vea el proyecto GetNighthawk.
