---
layout: page
title: Lanzamientos
permalink: es/project/releases
language: es
type: project
---

<!-- Una lista de los [lanzamientos de Meshery](https://github.com/layer5io/meshery/releases). Consulta el documento [Estrategia de compilación y lanzamiento](https://docs.google.com/document/d/11nAxYtz2SUusCYZ0JeNRrOLIxkgmmbUVWz63MBZV2oE/edit?usp=sharing) para obtener más detalles. -->

<!-- <div class="card">
  <div class="card-header">
    Destacados
  </div>
  <div class="card-body">
    <h5 class="card-title">Tratamiento de título especial</h5>
    <p class="card-text">Con el texto de apoyo a continuación como una introducción natural a contenido adicional.</p>
    <a href="#" class="btn btn-primary">Ir a algún lugar</a>
  </div>
</div>
<br>
<div class="card">
  <div class="card-header">
    Destacados
  </div>
  <div class="card-body">
    <h5 class="card-title">Tratamiento de título especial</h5>
    <p class="card-text">Con el texto de apoyo a continuación como una introducción natural a contenido adicional.</p>
    <a href="#" class="btn btn-primary">Ir a algún lugar</a>
  </div>
</div>
<br>
<div class="card">
  <div class="card-header">
    Destacados
  </div>
  <div class="card-body">
    <h5 class="card-title">Tratamiento de título especial</h5>
    <p class="card-text">Con el texto de apoyo a continuación como una introducción natural a contenido adicional.</p>
    <a href="#" class="btn btn-primary">Ir a algún lugar</a>
  </div>
</div>
<div class="card">
  <div class="card-header">
    Destacados
  </div>
  <div class="card-body">
    <h5 class="card-title">Tratamiento de título especial</h5>
    <p class="card-text">Con el texto de apoyo a continuación como una introducción natural a contenido adicional.</p>
    <a href="#" class="btn btn-primary">Ir a algún lugar</a>
  </div>
</div>
<div class="card">
  <div class="card-header">
    Destacados
  </div>
  <div class="card-body">
    <h5 class="card-title">Tratamiento de título especial</h5>
    <p class="card-text">Con el texto de apoyo a continuación como una introducción natural a contenido adicional.</p>
    <a href="#" class="btn btn-primary">Ir a algún lugar</a>
  </div>
</div> -->

## Lanzamientos de Meshery

### [v0.4.4](https://github.com/layer5io/meshery/releases/tag/v0.4.4)

Lanzado el 27 de agosto de 2020

**¿Qué hay de nuevo?**

- Características

  - Timeout agregado en golang-ci lint
  - Movido a Gorilla/mux para el enrutamiento del servidor
  - Migrado a golangci-lint
  - Interfaz del generador de carga nighthawk inicializada

- Mantenimiento

  - Se actualizó la rama meshmap a la última versión en la rama master
  - Título de la página en negrita
  - Se movió meshery adapter para osm

- Documentación

  - Se agregaron notas RVM adicionales para usuarios de Windows
  - Se agregó captura de pantalla y descripción de SMI
  - Añadida la versión actual de osm

- Corrección de bugs
  - Detección automática de la configuración de kube
  - Se corrigió la barra social del pie de página

### [v0.4.3](https://github.com/layer5io/meshery/releases/tag/v0.4.3)

Lanzado el 22 de agosto de 2020

**¿Qué hay de nuevo?**

- Características

  - Utiliza Open Service Mesh en lugar de osm
  - Menú desplegable de uso del campo de duración + texto en la preferencia de rendimiento
  - Se agregaron "tooltips" para los elementos del menú de navegación en el estado minimizado
  - Botón de calendario comunitario agregado

- Documentación

  - Páginas de documentos arregladas
  - README.md actualizado
  - Se agregó el enlace de la guía de bienvenida y el enlace para los principiantes
  - Se agregó "Open Service Mesh" al menú de navegación en los documentos de Meshery

- Corrección de bugs
  - smi-conformance: se agregaron dos puntos al principio
  - Se agregó la fila inferior que faltaba

### [v0.4.2](https://github.com/layer5io/meshery/releases/tag/v0.4.2)

Lanzado el 17 de agosto de 2020

**¿Qué hay de nuevo?**

- Características

  - Actualizar las operaciones CRUD en la configuración de prueba
  - Corrección de la "snackbar" no deseada en la pestaña de preferencias perf
  - Agregar validación para el protocolo en URL
  - Inicializar el adaptador OSM para Meshery
  - Agregar un nuevo comando "restart" a mesheryctl

- Mantenimiento

  - Reparar el "ci" roto para pruebas de interfaz de usuario
  - Optimizar imágenes
  - Quitar el paso redundante de "docker push"

- Documentación

  - Quitar enlaces rotos de la tabla de plataformas compatibles
  - Quitar ejemplo de salida de "brew upgrade mesheryctl"
  - Crear archivo smi-conformance-capability.md
  - Reparar los enlaces rotos del sitio de documentación
  - Reparar el enlace roto a SMP en release.md
  - Actualizar el adaptador Meshery Kuma a beta
  - Actualizar los documentos del adaptador
  - Agregar logotipos a la carpeta "assets"
  - Agregar Open Service Mesh a la lista
  - Corregir error tipográfico en la descripción del Adaptador App Mesh

- Corrección de bugs
  - Crear/Quitar/Leer en perfiles de prueba de usuario
  - Corrección para la página de resultados
  - Arreglar el archivo de configuración de releaseDrafter

### [v0.4.1](https://github.com/layer5io/meshery/releases/tag/v0.4.1)

Lanzado el 8 de agosto de 2020

**¿Qué hay de nuevo?**

- Documentación
  - Correcciones de "typos" y errores

### [v0.4.0-beta.4](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.4)

Lanzado el 9 de julio de 2020

**¿Qué hay de nuevo?**

- Corrección de bugs
  - hotfix: No se puede crear la página de configuración de Meshery

### [v0.4.0-beta.3](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.3)

Lanzado el 8 de julio de 2020

**¿Qué hay de nuevo?**

- Corrección de bugs
  - Hotfix del error introducido en la versión v0.4.0-beta.2

### [v0.4.0-beta.2](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.2)

Lanzado el 8 de julio de 2020

**¿Qué hay de nuevo?**

- Meshery
  - Se agregó el archivo kubeconfig inicial.
- Mesheryctl
  - Se quitó el error fatal debido a que faltaba `meshery.yaml`.
- Documentación
  - Logotipo del adaptador Kuma actualizado en Meshery Docs.
  - Se agregaron puertos de red a la tabla de adaptadores.
  - Captura de pantalla de proveedores agregada.
  - Se agregaron todas las plataformas compatibles a la lista.
  - Se agregó Nighthawk como generador de carga.
  - Se arregló el enlace roto en la página de EKS.

### [v0.3.19](https://github.com/layer5io/meshery/releases/tag/v0.3.19)

Lanzado el 7 de julio de 2020

**¿Qué hay de nuevo?**

- Mesheryctl
  - Se corrigió un error importante por no inicializar la carpeta de la aplicación (`~/meshery`) y el archivo de configuración de la aplicación (`meshery.yaml`).

### [v0.4.0-beta.1](https://github.com/layer5io/meshery/releases/tag/v0.4.0-beta.1)

Lanzado el 5 de julio de 2020

**¿Qué hay de nuevo?**

- Meshery
  - Etiqueta "stable-latest" para el nuevo canal de lanzamiento.

### [v0.3.18](https://github.com/layer5io/meshery/releases/tag/v0.3.18)

Lanzado el 4 de julio de 2020

**¿Qué hay de nuevo?**

- Mesheryctl
  - `mesheryctl version` ahora obtiene la información de la versión del servidor, incluida la etiqueta de lanzamiento y git sha.
  - `meshery perf` ahora soporta SPMS como un paquete a través del indicador `--file`.
- Meshery
  - Mejora de la extensibilidad: ahora se proporciona una interfaz abstracta, Load Generator Interface, para facilitar la extensión de Meshery para admitir otros generadores de carga, como Nighthawk.
- Documentos
  - instrucciones de EKS y GKE mejoradas

### [v0.3.17](https://github.com/layer5io/meshery/releases/tag/v0.3.17)

Lanzado el 1 de julio de 2020

**¿Qué hay de nuevo?**

- Meshery
  - Característica: Se estableció un nuevo canal de lanzamiento "stable".
  - Característica: Experiencia de usuario mejorada para encabezados HTTP de prueba de rendimiento avanzado.

### [v0.3.16](https://github.com/layer5io/meshery/releases/tag/v0.3.16)

Lanzado el 22 de junio de 2020

**¿Qué hay de nuevo?**

- Mesheryctl
  - Característica: Se agregó `mesheryctl system` como el nuevo comando para la gestión del ciclo de vida de Meshery.
    - Todos los comandos de gestión del ciclo de vida de Meshery como `start`, `stop`, `reset`, `logs`, y así, ahora se encuentran en `system`.
  - Característica: Se agregó `mesheryctl system config` como un nuevo subcomando para ayudar a configurar Meshery para comunicarse con las instancias levantadas de Kubernetes en GKE, AKS y EKS.

### [v0.3.15](https://github.com/layer5io/meshery/releases/tag/v0.3.15)

Lanzado el 1 de mayo de 2020

**¿Qué hay de nuevo?**

- Mesheryctl
  - Característica: Se agregó la funcionalidad de `mesheryctl perf --file <smp.yaml>` para perfiles de prueba de rendimiento [con formato SMP](https://github.com/layer5io/service-mesh-performance).
  - Característica: Cambio de comportamiento `mesheryctl perf` para usar tokens y admitir JWT.

### [v0.3.14](https://github.com/layer5io/meshery/releases/tag/v0.3.14)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Cambio de comportamiento en `mesheryctl start` para buscar nuevas imágenes del servidor Meshery por defecto.

### [v0.3.13](https://github.com/layer5io/meshery/releases/tag/v0.3.13)

**¿Qué hay de nuevo?**

- Mesheryctl
  - "Point release" para el soporte ARM inicial de mesheryctl.

### [v0.3.12](https://github.com/layer5io/meshery/releases/tag/v0.3.12)

**¿Qué hay de nuevo?**

- Mesheryctl
  - "Scoop Bucket" y "Scoop package" iniciales para `mesheryctl`.

### [v0.3.11](https://github.com/layer5io/meshery/releases/tag/v0.3.11)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Soluciona problemas menores de la experiencia del usuario en el comando `mesheryctl perf`. Consulta [Comandos CLI de Meshery y Documentación](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#) como referencia.

### [v0.3.10](https://github.com/layer5io/meshery/releases/tag/v0.3.10)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Introduce el comando `mesheryctl perf`. Consulta [Comandos CLI de Meshery y Documentación](https://docs.google.com/document/d/1xRlFpElRmybJ3WacgPKXgCSiQ2poJl3iCCV1dAalf0k/edit#) como referencia.

### [v0.3.9](https://github.com/layer5io/meshery/releases/tag/v0.3.12)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Se agregó `mesheryctl version` para proporcionar el número de versión del lado del servidor.
  - Mejora de `mesheryctl logs | stop | start` para proporcionar la gramática adecuada en situaciones en las que Meshery está detenido o Docker no está presente.
- Servidor Meshery
  - Se quitó la información superflua más allá de la dirección IP y el puerto en los "endpoints" de Grafana y Prometheus ([#612](https://github.com/layer5io/meshery/issues/612))

### [v0.3.8](https://github.com/layer5io/meshery/releases/tag/v0.3.8)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Se quitó la salida de línea de comando extraña. En esta versión se mejora la claridad de la interacción CLI con `mesheryctl`.

### [v0.3.7](https://github.com/layer5io/meshery/releases/tag/v0.3.7)

**¿Qué hay de nuevo?**

- Servidor Meshery
  - Ahora se admiten pruebas de conectividad ad-hoc para Prometheus. Los usuarios pueden hacer clic en el chip Prometheus y hacer que Meshery verifique su capacidad para conectarse a la instancia de Prometheus configurada.

### [v0.3.6](https://github.com/layer5io/meshery/releases/tag/v0.3.6)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Ahora se admiten pruebas de conectividad ad-hoc para Grafana. Los usuarios pueden hacer clic en el chip Grafana y hacer que Meshery verifique su capacidad para conectarse a la instancia de Grafana configurada.

### [v0.3.5](https://github.com/layer5io/meshery/releases/tag/v0.3.5)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Se quitó `init` como comando expuesto a los usuarios. La funcionalidad de este comando se usa internamente para `mesheryctl start`. Un nuevo comando `start --check` proporcionará la funcionalidad de verificación previa en el lugar de `init`.

### [v0.3.4](https://github.com/layer5io/meshery/releases/tag/v0.3.4)

**¿Qué hay de nuevo?**

- Mesheryctl
  - `mesheryctl version` ahora se ha mejorado con la novedad de mostrar el git commit (sha) de la versión mesheryctl.

### [v0.3.3](https://github.com/layer5io/meshery/releases/tag/v0.3.3)

**¿Qué hay de nuevo?**

- Servidor Meshery
  - Providers: una nueva construcción de proyecto que permite a los usuarios seleccionar el proveedor de autenticación, almacenamiento a largo plazo, etc.

### [v0.3.2](https://github.com/layer5io/meshery/releases/tag/v0.3.2)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Se agregó `mesheryctl version` como un nuevo subcomando.

### [v0.3.1](https://github.com/layer5io/meshery/releases/tag/v0.3.1)

**¿Qué hay de nuevo?**

- Servidor Meshery
  - Soporte para wrk2 como generador de carga alternativo.

### [v0.2.4](https://github.com/layer5io/meshery/releases/tag/v0.2.4)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Soporte de Homebrew disponible para mesheryctl.
- Servidor Meshery
  - Adaptador Meshery para Octarine lanzado como estable.
- Documentación
  - Guía de inicio rápido revisada para Mac, Linux y Windows.
  - Soporte WSL2 publicado.
  - El script de generación de kubeconfig GKE cambió a `--decode`.

### [v0.2.3](https://github.com/layer5io/meshery/releases/tag/v0.2.3)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Salida mejorada del comando `status` en Windows.
- Servidor Meshery
  - Posibilidad de hacer "deploy" de Meshery en Istio.
  - "Adapter Chips": Mover el número de puerto del adaptador al tooltip.

### [v0.2.2](https://github.com/layer5io/meshery/releases/tag/v0.2.2)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Comando `update` mejorado, sin sobrescribir en la configuración local mientras se ejecuta el comando `start`.
- Meshery UI
  - Una nueva vista modal para organizar y mostrar los resultados de rendimiento en formato tabular.

### [v0.2.1](https://github.com/layer5io/meshery/releases/tag/v0.2.1)

**¿Qué hay de nuevo?**

- Mesheryctl
  - Revisión de la configuración de Kubernetes dentro del clúster y fuera del clúster.

## [v0.2.0](https://github.com/layer5io/meshery/releases/tag/v0.2.1)

**¿Qué hay de nuevo?**

- Adaptadores Meshery
  - Adaptador introducido para Network Service Mesh.
- Servidor Meshery
  - Capacidad para ejecutar pruebas de rendimiento de forma asincrónica.
  - Recopila y conserva métricas de nodos.

## [v0.1.6](https://github.com/layer5io/meshery/releases/tag/v0.1.6)

**¿Qué hay de nuevo?**

- Registro de cambios

* Nueva interfaz de usuario para administrar la conexión de Meshery al clúster de Kubernetes.
* Nueva compatibilidad de `mesheryctl` para Windows para abrir el navegador predeterminado al `iniciar`.
* Nueva aplicación de muestra agregada al adaptador istio: aplicación Hipster (aplicación de demostración de microservicios de Google).

## [v0.1.5](https://github.com/layer5io/meshery/releases/tag/v0.1.5)

**¿Qué hay de nuevo?**

- Registro de cambios
- Mejoras de UX

* `mesheryctl start` ahora espera a que los contenedores de la aplicación Meshery estén activos antes de iniciar el navegador del usuario.
* `mesheryctl stop` ahora muestra el progreso del comando similar a la experiencia cuando se usa el script bash `meshery`.

## [v0.1.4](https://github.com/layer5io/meshery/releases/tag/v0.1.4)

**¿Qué hay de nuevo?**

- Registro de cambios

## [v0.1.3](https://github.com/layer5io/meshery/releases/tag/v0.1.3)

**¿Qué hay de nuevo?**

- Migrar desde Configurar Meshery a la página Configuración.

## [v0.1.2](https://github.com/layer5io/meshery/releases/tag/v0.1.2)

**¿Qué hay de nuevo?**

- Sincronización del almacenamiento local del navegador

* Almacenamiento de sesiones en memoria de Meshery.

## [v0.1.1](https://github.com/layer5io/meshery/releases/tag/v0.1.1)

**¿Qué hay de nuevo?**

- Parche para el bug Alpine.

## [v0.1.0](https://github.com/layer5io/meshery/releases/tag/v0.1.0)

**What's new**

- Istio inicial completamente funcional.

## [v0.0.9](https://github.com/layer5io/meshery/releases/tag/v0.0.9)

**¿Qué hay de nuevo?**

- Sitio de documentación segregado

* contenido presentado.

## [v0.0.8](https://github.com/layer5io/meshery/releases/tag/v0.0.8)

**¿Qué hay de nuevo?**

- Posibilidad de importar json del board de Grafana

* Integración con Prometheus directamente para métricas.

## [v0.0.7](https://github.com/layer5io/meshery/releases/tag/v0.0.7)

**¿Qué hay de nuevo?**

- Migrado lejos de los gráficos de iframe incrustados de grafana

* Ahora usando Chartjs para gráficos.
* Ahora usando C3 para gráficos.

## [v0.0.6](https://github.com/layer5io/meshery/releases/tag/v0.0.6)

**¿Qué hay de nuevo?**

- Adaptador pre-alfa de Consul.

* Soporte para implementar el adaptador SMI Istio.
* Asegurar la compatibilidad para aplicar los manifiestos SMI a través de Meshery.

## [v0.0.5](https://github.com/layer5io/meshery/releases/tag/v0.0.5)

**¿Qué hay de nuevo?**

- Adaptador Linkerd pre-alpha.

* Capacidad para filtrar resultados.

## [v0.0.4](https://github.com/layer5io/meshery/releases/tag/v0.0.4)

**¿Qué hay de nuevo?**
-Posibilidad de ver resultados persistentes.

- Integración y soporte para gráficos Grafana.
- Embedding panels in iframe.

## [v0.0.3](https://github.com/layer5io/meshery/releases/tag/v0.0.3)

**¿Qué hay de nuevo?**

- Versión inicial con soporte de adaptadores Meshery.

* Lanzamiento de una versión pre-alfa del adaptador Istio.

## [v0.0.2](https://github.com/layer5io/meshery/releases/tag/v0.0.2)

**¿Qué hay de nuevo?**

- Capacidad para admitir la ejecución de yaml personalizado en Kubernetes con Istio.

## [v0.0.1](https://github.com/layer5io/meshery/releases/tag/v0.0.1)

**New release**

- Versión inicial de Meshery
  - Conectar a Kubernetes.
  - Ejecutar comandos preconfigurados en Kubernetes con Istio.

<!-- <table class="responsive-table hover striped">
  <thead>
    <th class="centered">Versión</th>
    <th>Descripción</th>
    <th style="white-space: nowrap;">Fecha de lanzamiento</th>
  </thead>
  <tbody>
      <tr>
        <td class="centered">0.3.15</td>
        <td>
          <em>Mesheryctl</em> -Característica: Se agregó mesheryctl perf --file <smp.yaml> soporte para perfiles de prueba de rendimiento con formato SMP. Cambio de comportamiento de mesheryctl perf para usar tokens y admitir JWT.
        </td>
        <td>May 1, 2020</td>
      </tr>
      <tr>
      <td class="centered">0.3.14</td>
      <td>
        <em>Mesheryctl</em> - Cambio de comportamiento en mesheryctl, comienza a buscar nuevas imágenes del servidor Meshery de forma predeterminada.
      </td>
      <td>24 abr 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.13</td>
      <td>
        <em>Mesheryctl</em> - "Point release" para el soporte ARM inicial de mesheryctl.
      </td>
      <td>16 abr 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.12</td>
      <td>
        <em>Mesheryctl</em> - Paquete inicial Scoop Bucket y Scoop para mesheryctl.
      </td>
      <td>10 abr, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.11</td>
      <td>
        <em>Mesheryctl</em> - Soluciona problemas menores de experiencia del usuario en el comando mesheryctl perf. Consulta Comandos CLI de Meshery y Documentación como referencia.
      </td>
      <td>25 mar, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.10</td>
      <td>
        <em>Mesheryctl</em> - Introduce el comando mesheryctl perf. Consulta Comandos CLI de Meshery y Documentación como referencia.
      </td>
      <td>9 mar, 2020</td>
    </tr>
      <tr>
      <td class="centered">0.3.9</td>
      <td>
        <em>Mesheryctl</em> - Se agregó la versión de mesheryctl para proporcionar el número de versión del lado del servidor, mejora de mesheryctl logs | stop | start para proporcionar la gramática adecuada en situaciones en las que Meshery se detiene o Docker no está presente.
        <em>Meshery</em> - Se quitó la información superflua más allá de la dirección IP y el puerto en los puntos finales de Grafana y Prometheus.
      </td>
      <td>3 feb, 2020</td>
    </tr>
    <tr>
      <td class="centered">0.3.8</td>
      <td>
        <em>Mesheryctl</em> - Se quitó la salida de línea de comando extraña. La claridad de la interacción CLI con mesheryctl se mejora en esta versión.
      </td>
      <td>17 ene, 2020</td>
    </tr>
    <tr>
      <td class="centered">0.3.7</td>
      <td>
        <em>Meshery</em> - Ahora se admiten pruebas de conectividad ad-hoc para Prometheus. Los usuarios pueden hacer clic en el chip Prometheus y hacer que Meshery verifique su capacidad para conectarse a la instancia de Prometheus configurada.
      </td>
      <td>15 ene, 2020</td>
    </tr>
     <tr>
      <td class="centered">0.3.6</td>
      <td>
        <em>Meshery</em> - Ahora se admiten pruebas de conectividad ad-hoc para Grafana. Los usuarios pueden hacer clic en el chip Grafana y hacer que Meshery verifique su capacidad para conectarse a la instancia de Grafana configurada.
      </td>
      <td>12 ene, 2020</td>
    </tr>
     <tr>
      <td class="centered">0.3.5</td>
      <td>
        <em>Mesheryctl</em> - Se quitó init como comando expuesto a los usuarios. La funcionalidad de este comando se usa internamente para mesheryctl start. Un nuevo comando start --check proporcionará la funcionalidad de verificación previa en lugar de init.
      </td>
      <td>12 ene, 2020</td>
    </tr>
     <tr>
      <td class="centered">0.3.4</td>
      <td>
        <em>Mesheryctl</em> - La versión de mesheryctl ahora está mejorada con la novedad de mostrar el git commit (sha) de la versión de mesheryctl.
      </td>
      <td>30 dec, 2019</td>
    </tr>
     <tr>
      <td class="centered">0.3.3</td>
      <td>
        <em>Meshery</em> - Providers (una nueva construcción de proyecto que permite a los usuarios seleccionar el proveedor de autenticación, almacenamiento a largo plazo, etc.).
      </td>
      <td>20 dec, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.3.2</td>
      <td>
        <em>Mesheryctl</em> - agrega la versión mesheryctl como un nuevo subcomando.
      </td>
      <td>29 nov, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.3.1</td>
      <td>
        <em>Meshery</em> - Soporte para wrk2 como generador de carga alternativo.
      </td>
      <td>12 nov, 2019</td>
    </tr>
    <tr><td colspan="3"><strong>v0.3.0</strong></td></tr>
    <tr>
      <td class="centered">0.2.4</td>
      <td>
        <em>Meshery</em>
          - Adaptador Meshery para Octarine lanzado como estable.
        <br />
        <em>mesheryctl</em>
          - ahora disponible a través de homebrew.
        <br />
        <em>Documentación</em>
        - inicio rápido revisado para Mac, Linux _y_ Windows.
        - Soporte WSL2 publicado.
        - El script de generación de kubeconfig GKE cambió a `--decode`.
      </td>
      <td>5 nov, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.3</td>
      <td>
        <em>mesheryctl</em>
          - salida de `status` mejorada en Windows
        <br />
        <em>Meshery</em>
          - Posibilidad de implementar Meshery en Istio.
          - "Adapter Chips": Mover el número de puerto del adaptador al tooltip.
        <br />
        <em>Documentos</em>
          - /search ya no se redirige a github.io.
      </td>
      <td>3 nov, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.2</td>
      <td>Mesheryctl: verbosidad mejorada del comando update; ya no se sobrescribe el archivo .meshery.yml local cuando se ejecuta el inicio o los registros. Documentos: revisión del sitio de documentos con un nuevo tema jekyll (gracias @venilnoronha). Pruebas de rendimiento: una nueva vista modal para organizar y mostrar los resultados de rendimiento en formato tabular. </td>
      <td>26 oct, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.1</td>
      <td>Instalación de Meshery: revisión de la configuración de Kubernetes dentro del clúster contra a fuera del clúster.</td>
      <td>23 oct, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.2.0</td>
      <td>Adaptador Meshery para Network Service Mesh: el adaptador ahora está en versión beta; Aprovisionamiento de NSM. Mejoras en las pruebas de rendimiento: las pruebas de rendimiento se ejecutan de forma asincrónica, notificando al usuario cuando los resultados de la prueba están disponibles; recopilar y conservar métricas de nodo. Service Mesh Sync: soporte para descubrir el tipo de service mesh. Mejoras de rendimiento mediante ajustes de memoria y creación de perfiles de código.</td>
      <td>22 oct, 2019</td>
    </tr>
        <tr><td colspan="3"><strong>v0.2.0</strong></td></tr>
    <tr>
      <td class="centered">0.1.6</td>
      <td>Nueva interfaz de usuario para administrar la conexión de Meshery al clúster de Kubernetes. Nueva compatibilidad con mesheryctl para Windows para abrir el navegador predeterminado al iniciar.</td>
      <td>13 oct, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.5</td>
      <td>Mejoras de UX: mesheryctl start ahora espera a que los contenedores de la aplicación meshery estén activos antes de iniciar el navegador del usuario. Este nuevo comportamiento asegura que los usuarios no experimenten un mensaje 404; mesheryctl stop ahora muestra el progreso del comando similar a la experiencia cuando se usa el script meshery bash.</td>
      <td>20 sep, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.4</td>
      <td>Actualizar README.md para el lanzamiento.</td>
      <td>12 sep, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.3</td>
      <td>Migrar opción de Configurar Meshery a la página Configuración.</td>
      <td>27 jun, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.2</td>
      <td>Sincronización del almacenamiento local del navegador con el almacenamiento de sesiones en memoria de Meshery.</td>
      <td>14 jun, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.1.1</td>
      <td>Parche para el bug Alpine.</td>
      <td>31 may, 2019</td>
    </tr>
    <tr><td colspan="3"><strong>v0.1.0</strong></td></tr>
    <tr>
      <td class="centered">0.0.9</td>
      <td>Sitio de documentación segregado y contenido presentado.</td>
      <td>2 may, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.8</td>
      <td>Capacidad para importar el json del board de Grafana pero integrándose con Prometheus directamente para métricas.</td>
      <td>15 abr, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.7</td>
      <td>Se migró de los gráficos de iframe incrustados de grafana a usar Chartjs y C3 para los gráficos.</td>
      <td>20 mar, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.5</td>
      <td>Adaptador Linkerd pre-alpha. Capacidad para filtrar resultados.</td>
      <td>28 feb, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.4</td>
      <td>Posibilidad de ver resultados persistentes. Integración y soporte para gráficos de Grafana y paneles incrustados en iframe.</td>
      <td>28 feb, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.3</td>
      <td>Versión inicial con soporte de adaptadores Meshery y lanzamiento de una versión pre-alfa del adaptador Istio.</td>
      <td>21 ene, 2019</td>
    </tr>
    <tr>
      <td class="centered">0.0.2</td>
      <td>Capacidad para admitir la ejecución de yaml personalizado en Kubernetes con Istio.</td>
      <td>30 nov, 2018</td>
    </tr>
    <tr>
      <td class="centered">0.0.1</td>
      <td>Versión inicial de Meshery. Conéctese a Kubernetes y ejecute comandos preconfigurados en Kubernetes con Istio.</td>
      <td>16 nov, 2018</td>
    </tr>
  </tbody>
</table>
 -->
