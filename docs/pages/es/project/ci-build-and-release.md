---
layout: default
title: Construcción y Lanzamiento (CI)
permalink: es/project/build-and-release
language: es
---

El sistema de creación y lanzamiento de Meshery incorpora muchas herramientas, organizadas en diferentes flujos de trabajo, cada uno de los cuales se desencadena por diferentes eventos. El sistema de compilación y lanzamiento de Meshery no se ejecuta en un horario, sino que se basa en eventos. Las acciones de GitHub se utilizan para definir los flujos de trabajo de CI de Meshery. Las nuevas compilaciones de Meshery y sus diversos componentes se generan automáticamente tras la inserción, el lanzamiento y otros eventos similares, generalmente en relación con sus respectivas ramas maestras.


## Artefactos

Hoy, los adaptadores Meshery y Meshery se lanzan como imágenes de contenedor Docker, disponibles en Docker Hub. Los adaptadores de Meshery son adaptadores fuera de proceso (lo que significa que no se compilan en el binario principal de Meshery) y, como tales, son artefactos de compilación independientes. automáticamente usando acciones de GitHub.

### Repositorios de artefactos

Los artefactos producidos en los procesos de construcción se publican y conservan en diferentes repositorios públicos y en diferentes formatos.

| Ubicación     | Proyecto       | Repositorio    |
| ------------- | ------------- | ------------- |
| Docker Hub    | Meshery       | [https://hub.docker.com/r/layer5/meshery](https://hub.docker.com/r/layer5/meshery) |
| GitHub        | mesheryctl    | [https://github.com/layer5io/meshery/releases](https://github.com/layer5io/meshery/releases) |
| Docker Hub    | Meshery Adapter for \<service-mesh\> | https://hub.docker.com/r/layer5/meshery-\<service-mesh\> |
| Docs          | Meshery Documentation | [https://docs.meshery.io](https://docs.meshery.io) |
| GitHub        | [Service Mesh Performance](https://smp-spec.io) | [https://github.com/layer5io/service-mesh-performance](https://github.com/layer5io/service-mesh-performance) |

## Secretos

Algunas partes del flujo de trabajo requieren secretos para realizar sus tareas. Estos secretos se definen dentro de los repositorios respectivos y son accesibles a los flujos de trabajo durante el tiempo de ejecución. Los secretos definidos actualmente incluyen:

- `DOCKER_USERNAME`: nombre de usuario del usuario de Docker Hub con los privilegios adecuados para enviar imágenes
- `DOCKER_PASSWORD`: Contraseña para el usuario de Docker Hub
- `GO_VERSION`: A partir del 9 de diciembre de 2020 es 1,15
- `IMAGE_NAME`: nombre de imagen apropiado para cada una de las imágenes del contenedor de Docker. Todos están bajo la organización `layer5io`.
- `SLACK_BOT_TOKEN`: Usado para la notificación de nuevas estrellas de GitHub entregadas al repositorio de Meshery.
- CYPRESS_RECORD_KEY`: Se utiliza para la integración con la cuenta Layer5 en Cypress.
- `GLOBAL_TOKEN`: se utiliza para transmitir de forma segura los resultados de las pruebas de rendimiento para el proveedor Ninguno.

El usuario de Docker Hub, `mesheryci`, pertenece al equipo" ciusers "en Docker Hub y actúa como la cuenta de servicio bajo la cual se envían estas compilaciones automatizadas. Cada vez que se crea un nuevo repositorio de Docker Hub, tenemos que otorgar permisos de "Administrador" (para actualizar el archivo README en el repositorio de Docker Hub) al equipo de usuarios.

## Comprobaciones y pruebas
El flujo de trabajo de CI de Meshery incorpora varias comprobaciones (lista parcial a continuación) durante las fusiones y / o confirmaciones de cualquier rama y solicitudes de extracción a la rama maestra para evitar que el código roto se fusione con el maestro.

En conjunto, los repositorios de Meshery generalmente tendrán un flujo de trabajo de CI para confirmaciones y solicitudes de extracción que constan de las siguientes acciones:

- Control de pelusa (golint)
- Comprobación de análisis estático (comprobación estática)
- Veterinario (govet)
- Controles de seguridad (gosec)
- Pruebas unitarias (ir a pruebas)
- Construir (ir a construir)
- Lanzamiento de binarios a través de GoReleaser (solo para mesheryctl en el repositorio de Meshery)
- Compilación, etiquetado y envío de Docker

## Construcciones automatizadas

Todos los repositorios de Meshery GitHub están configurados con acciones de GitHub. Cada vez que se envía una solicitud de extracción contra la rama maestra de cualquier repositorio, se invocarán las acciones de GitHub de ese repositorio (ya sea que el RP esté combinado o no). Los flujos de trabajo definidos en el repositorio de Meshery generalmente (pero no siempre) realizarán las siguientes acciones:


1. active una compilación de Docker para crear una imagen de contenedor de Docker
1. Genere dos etiquetas de Docker:
   1. una etiqueta que contiene el SHA de git merge
   1. una etiqueta que contiene la etiqueta git de esa versión en particular (si hay una)
1. Asigne cada una de estas dos etiquetas a la nueva imagen del contenedor, así como a la última etiqueta.
1. Inserte las nuevas etiquetas y la imagen de Docker en Docker Hub.

### Construyendo `mesheryctl`

Como caso especial, el repositorio de mallas contiene un artefacto adicional producido durante cada compilación. Este artefacto es mesheryctl, que se construye como un binario ejecutable. Para facilitar el trabajo de construir mesheryctl para una combinación de diferentes arquitecturas de plataforma y sistemas operativos, estamos usando [GoReleaser] (https://goreleaser.com). Independientemente de la rama, por cada confirmación de git y envío de git al repositorio de meshery, GoReleaser ejecutará y generará el sistema operativo y los binarios específicos de arch (pero NO los publicará en GitHub). Aunque los binarios de mesheryctl se crean cada vez que se fusiona una solicitud de extracción con el maestro, solo se publican (conservan) los artefactos de canal estables.

### Liberando `mesheryctl` en GitHub

Solo cuando esté presente una etiqueta git que contenga un número de versión semántica (es una confirmación en la rama maestra), GoReleaser se ejecutará, generará los archivos y también publicará los archivos en [versiones de GitHub de Meshery] (https://github.com/layer5io / meshery / releases) automáticamente. GoReleaser está configurado para generar artefactos para el siguiente sistema operativo, combinación ARCH:

- Darwin - i386, x86_64
- Linux - i386, x86_64
- Windows - i386, x86_64
- FreeBSD - i386, x86_64

Los artefactos estarán disponibles como un archivo tar.gz para todos los sistemas operativos. mesheryctl se incluye en paquetes para administradores de paquetes de uso común: homebrew y scoop.

#### Homebrew

GoReleaser facilita la creación de una fórmula de preparación para mesheryctl. El repositorio [homebrew-tap] (https://github.com/layer5io/homebrew-tap) es la ubicación de las fórmulas de preparación de Layer5.

#### Scoop

GoReleaser facilita la creación de una aplicación Scoop para mesheryctl. El repositorio [scoop-bucket](https://github.com/layer5io/scoop-bucket) es la ubicación para el Scoop bucket de Layer5.

## Versionado de Lanzamientos

Seguimos las versiones semánticas de uso común para las versiones de Meshery, Meshery Adapter y Performance Benchmark Specification. Dado un número de versión MAYOR.MENOR.PARCHE.CONSTRUIDO, incremente:

- versión  MAYOR - cambios importantes con un potencial poco común de cambios de API incompatibles.
- versión  MENOR - agregue funcionalidad de una manera compatible con versiones anteriores.
- versión PARCHE - principalmente para correcciones de errores y seguridad.
- AlPHA/BETA/RC - se utiliza para facilitar las pruebas tempranas de una próxima versión.

### Versionado de Componentes

Meshery comprende una serie de componentes que incluyen un servidor, adaptadores, interfaz de usuario y CLI. Como aplicación, Meshery es una composición de estos diferentes componentes funcionales. Si bien todos los componentes de Meshery generalmente se implementan como una unidad colectiva (juntos), cada componente se versiona de forma independiente, para permitir que estén acoplados libremente e iterar en la funcionalidad de forma independiente. Algunos de los componentes deben actualizarse simultáneamente, mientras que otros pueden actualizarse de forma independiente. Consulte la guía [Actualizando Versión Meshery](es/guide/upgrade) para más  información.

Las etiquetas de lanzamiento de GitHub contendrán un número de versión semántica. Los números de versión semántica deberán administrarse manualmente etiquetando una confirmación relevante en la rama maestra con un número de versión semántica (ejemplo: v1.2.3).

## Proceso de lanzamiento

La documentación de las versiones de Meshery contiene una tabla de versiones y notas de la versión y debe actualizarse con cada versión.

### Lanzamientos Automatizados

Las versiones las activa manualmente un miembro del equipo de versiones que publica una versión. El miembro del equipo de publicación debe asignar los nombres y las etiquetas de las versiones. Los flujos de trabajo de acciones de GitHub activarán y se encargarán de ejecutar los pasos necesarios y de publicar todos los artefactos (por ejemplo, imágenes binarias y acoplables).

### Activadores de flujo de trabajo

Los siguientes eventos activarán uno o más flujos de trabajo:

1. Lanzamiento etiquetado
1. Confirmación enviada a la rama principal
1. PR abierto o compromiso enviado a la rama PR
1. PR se fusionó con la rama principal

### Notas de Lanzamiento

Si bien el uso de Acciones de GitHub facilita las compilaciones automatizadas, ReleaseDrafter está ayudando a facilitar las notas de la versión automatizadas y el control de versiones.
### Generando Notas de la Versión

ReleaseDrafter genera una etiqueta de GitHub y un borrador de lanzamiento. La acción ReleaseDrafter se activará y redactará automáticamente notas de la versión de acuerdo con la configuración establecida. ReleaseDrafter borra las versiones tan pronto como una confirmación se convierte en maestra después de la versión anterior. La acción de GitHub, ReleaseDrafter, es compatible con las versiones semánticas y se utiliza para incrementar automáticamente el número de versión semántica mirando la versión de la versión anterior.

#### Publicación automatizada de notas de la versión

La publicación de notas de la versión en Meshery Docs está automatizada. Activado por un evento de lanzamiento, un flujo de trabajo comprobará el repositorio de Meshery, copiará las notas de la versión redactadas automáticamente en una colección de Jekyll en Meshery Docs y generará una solicitud de extracción.

#### Etiquetadora automática de solicitud de extracción

Un bot etiquetador de problemas de GitHub está configurado para asignar etiquetas automáticamente a los problemas según los archivos que hayan cambiado en qué directorios. Por ejemplo, una solicitud de extracción con cambios en los archivos de la carpeta "/ docs / **" recibirá la etiqueta "area / docs". La presencia de la etiqueta "area / docs" se utiliza para activar las compilaciones de documentación y las compilaciones de Netlify de Meshery Docs. Se asignan y utilizan etiquetas similares para desencadenar flujos de trabajo o se utilizan como indicadores condicionales en los flujos de trabajo para determinar qué flujos de trabajo o qué pasos de un flujo de trabajo ejecutar.

## Canales de lanzamiento

Los artefactos de las compilaciones para Meshery y sus componentes se publican en dos canales de lanzamiento diferentes, por lo que se pueden proporcionar controles mejorados tanto a los usuarios de Meshery como a los desarrolladores de Meshery. Los dos canales de liberación son canales de liberación * edge * y * estable *.

En relación con las versiones estables, las versiones de borde ocurren con mucha más frecuencia. Las versiones de borde se realizan con cada fusión para master, a menos que la fusión para master sea para una versión estable. Las versiones estables se realizan con cada combinación para dominar cuando una etiqueta de versión de GitHub también está presente en el flujo de trabajo.

### Canal estable

El siguiente es un ejemplo de los canales de lanzamiento y las etiquetas de la ventana acoplable que se utilizan para diferenciarlos. La última etiqueta se aplicará solo a las imágenes del canal de versión estable. Aquí hay dos lanzamientos con dos imágenes diferentes.

**Última imagen estable**

- layer5/meshery:stable-latest
- layer5/meshery:stable-v0.4.1
- layer5/meshery:stable-324vdgb (sha)

**Imagen Estable Anterior**

- layer5/meshery:stable-v0.4.0
- layer5/meshery:stable-289d02 (sha)

Cada imagen de la ventana acoplable construida recibe las etiquetas de borde o las pestañas estables. El conjunto de etiquetas de imagen asignadas se determina en función de si hay una etiqueta de liberación o no. En otras palabras, las imágenes de la ventana acoplable de canal estable obtienen las etiquetas "estables" solo en presencia de una etiqueta de lanzamiento (por ejemplo, v0.4.1).

### Canal Edge

El canal de liberación edge generalmente contiene código menos probado, menos "horneado". La razón principal de la "ventaja" es permitir que los contribuyentes y usuarios avanzados accedan a las funciones más temprano que tarde. Algunas funciones necesitan pruebas que se facilitan mejor al permitir que los usuarios con tolerancia y paciencia las prueben.

Las versiones estables y perimetrales se publican en el mismo repositorio de Docker Hub. Los repositorios de Docker Hub diferencian los canales de lanzamiento por etiqueta de imagen. Se sigue la siguiente convención de etiquetado de imágenes de Docker:

**Última Imagen Edge**

- layer5/meshery:edge-latest
- layer5/meshery:edge-289d02 (sha)

**Image Edge Anterior**

- layer5/meshery:edge-324vdgb (sha)


### Cambio entre canales de liberación de Meshery

Los usuarios están autorizados a cambiar entre los canales de lanzamiento cuando quieran.

#### Cambio de canales de liberación mediante mesheryctl

Los usuarios pueden usar mesheryctl para cambiar entre canales de lanzamiento, p. Ej. `canal del sistema mesheryctl [stable|edge]`. Alternativamente, los usuarios pueden cambiar manualmente entre canales actualizando las etiquetas de imagen de la ventana acoplable en sus archivos de manifiesto meshery.yaml/Kubernetes. Este comando genera un meshery.yml (un archivo docker-compose) con etiquetas de publicación apropiadas para el canal para las diferentes imágenes del contenedor Docker.

#### Visualización del canal de Lanzamiento y la Información de la Versión en la Interfaz de Usuario de Meshery

A los usuarios se les muestra la suscripción al canal de lanzamiento de su implementación de Meshery en una nueva configuración en el área de Preferencias de la interfaz de usuario de Meshery, de modo que las personas pueden usar la interfaz de usuario de manera alternativa para cambiar entre canales si lo desean. Los números de versión de los adaptadores Meshery también se muestran en la interfaz de usuario.

## Cadencia de Liberación

Los lanzamientos menores del proyecto Meshery se lanzan con frecuencia (una media mensual) y los lanzamientos de parches se realizan bajo demanda entre esos momentos. El proyecto aún no tiene lanzamientos a largo plazo que se mantengan con correcciones de errores. Las correcciones de errores y los parches se publicarán según sea necesario en la última versión de lanzamiento.

### Soporte de lanzamiento

El soporte general de la comunidad y el soporte comercial de Layer5 están disponibles. Por separado, terceros y socios pueden ofrecer soluciones de soporte a más largo plazo.

#### Pre v1.0

El proyecto se centra en la funcionalidad, la calidad y la adopción, al tiempo que conserva la flexibilidad para los cambios en la arquitectura.

#### Post v1.0

Una vez que se ha realizado una versión 1.0, aproximadamente una vez al mes, los responsables del proyecto tomarán una de estas compilaciones diarias y la ejecutarán a través de una serie de pruebas de calificación adicionales y etiquetarán la compilación como una versión estable. Aproximadamente una vez al trimestre, los encargados del mantenimiento del proyecto toman una de estas versiones estables, realizan un montón de pruebas más y etiquetan la compilación como una versión de soporte a largo plazo (LTS). Finalmente, si encontramos algún problema con una versión de LTS, publicamos parches.

Los diferentes tipos (Daily, Stable, LTS) representan diferentes niveles de calidad de producto y diferentes niveles de soporte del equipo de Meshery. En este contexto, soporte significa que produciremos versiones de parches para problemas críticos y ofreceremos asistencia técnica.
