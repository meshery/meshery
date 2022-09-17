---
layout: default
title: "Extensibilidad: adaptadores de malla de servicio"
permalink: es/extensibility/adapters
type: extensibility
abstract: 'La arquitectura Meshery es extensible. Meshery proporciona varios puntos de extensión para trabajar con diferentes mallas de servicio a través de <a href="extensibility#adapters">adaptadores</a>, <a href="extensibility#load-generators"> generadores de carga</a> and <a href="extensibility#providers">providers</a>.'
#redirect_from: extensibility
language: es
---

## Principios rectores del diseño de adaptadores

Los adaptadores permiten a Meshery interactuar con las diferentes mallas de servicio. Revise la lista de todos los disponibles [service mesh adapters](service-meshes/adapters). Ver el [Meshery Architecture](architecture) diagramas para imágenes sobre cómo los adaptadores se relacionan con otros componentes de Meshery.

Meshery mantiene los siguientes principios rectores para el diseño de adaptadores:

1. **Los adaptadores permiten a Meshery interactuar con las diferentes mallas de servicio, exponiendo su valor diferenciado a los usuarios.**

- Se debe alentar a los proyectos de mallas de servicio a mantener sus propios adaptadores. Permitirles exponer sus capacidades diferenciadas fomenta esto.

1. **Los adaptadores deben evitar la reinvención de las ruedas, pero deben buscar aprovechar la funcionalidad proporcionada por las mallas de servicio bajo administración.**

- Esto reduce los costos de mantenimiento y mejora la confiabilidad.

### Capacidades del adaptador

Meshery se comunica con adaptadores a través de grpc. Los adaptadores establecen comunicación con la malla de servicios. Los adaptadores tienen un conjunto predefinido de operaciones que se agrupan en función de tipos de operaciones predefinidos.

Los tipos de operación predefinidos son:

- Instalar en pc
- Aplicación de muestra
- Config
- Validar
- Custom

## Meshery Adapter Codebase Overview

[Bibliotecas comunes](https://docs.google.com/presentation/d/1uQU7e_evJ8IMIzlLoBi3jQSRvpKsl_-K1COVGjJVs30/edit#) se utilizan para evitar la duplicación de códigos y aplicar DRY.

### [MeshKit](https://github.com/layer5io/meshkit)

La jerarquía de código es conectable e independiente entre sí. Puede haber N cantidad de paquetes dependiendo del caso de uso.

- `errors/` - contiene las implementaciones y los controladores de errores y los códigos de error que se utilizan en los proyectos.
- `logger/` - contiene las implementaciones del controlador de registro y los atributos personalizados para agregar, si corresponde.
- `utils/` - contiene todas las funciones de utilidad que son específicas de los proyectos de malla y se deben usar de forma genérica en todos ellos.
- `tracing/` - contiene las implementaciones de los controladores de rastreo con diferentes proveedores de rastreo como jaeger, newrelic, etc.

Cada paquete dentro de un meshkit es una implementación de interfaz de controlador, la implementación podría ser de cualquier paquete de terceros o del go-kit.

### [Biblioteca de adaptadores de Meshery](https://github.com/meshery/meshery-adapter-library)

Esta sección contiene una descripción general de alto nivel de la biblioteca de adaptadores de malla, su propósito y arquitectura. Para obtener más detalles, se remite al lector a la documentación y al código del repositorio.

El propósito principal de la biblioteca de adaptadores de malla es:

- proporcionar un conjunto de interfaces, algunas con implementaciones predeterminadas, para ser utilizadas y extendidas por los adaptadores.
- Implementar preocupaciones transversales comunes como registro, errores y rastreo.
- proporcionar un mini framework que implemente el servidor gRPC que permite conectar la configuración específica de la malla y - las operaciones implementadas en los adaptadores.
- proporcionar puntos de extensión de middleware

#### Descripción general y uso

La biblioteca consta de interfaces e implementaciones predeterminadas para la funcionalidad principal y común de un adaptador. También proporciona un mini-framework que ejecuta el servicio del adaptador gRPC, llamando a las funciones de los controladores inyectados por el código del adaptador. Esto se representa en un estilo UML-ish en la figura siguiente. La biblioteca se utiliza en todos los adaptadores de Meshery.

## Contribución a los adaptadores Meshery

Con el[CONTRIBUTING.md](https://github.com/layer5io/meshery/blob/master/CONTRIBUTING.md#adapter) Tenga en cuenta que el desarrollo sigue el flujo de trabajo de solicitud de bifurcación y extracción habitual que se describe aquí, consulte también Proceso de GitHub. Al bifurcar, GitHub desactiva todos los flujos de trabajo. Es una buena práctica y seguro activarlos de manera que el código se valide en cada inserción. Esto requiere que el filtro de ramas para "al enviar" esté configurado en "\ _ \ _" para que se active también en las ramas que contengan "/" en su nombre. Las acciones se parametrizan mediante secretos (consulte Estrategia de creación y lanzamiento). La imagen de Docker solo se crea y se envía a Docker Hub si se inserta una etiqueta y se configura la información de autenticación correspondiente. El único secreto que debe establecerse en cada bifurcación es GO_VERSION, especificado en la estrategia de compilación y lanzamiento; de lo contrario, se usa la versión predeterminada de la acción correspondiente.

Cada compromiso debe estar firmado, consulte {{site.baseurl}}project/contributingContributing Overview.

### Ejecutando un adaptador como contenedor

La prueba de los cambios locales que se ejecutan como un contenedor se puede lograr de dos maneras:

1. Defina la dirección del adaptador en la interfaz de usuario: a menos que el contenedor en ejecución tenga el nombre especificado en el destino de ejecución de la ventana acoplable en el Makefile, el contenedor debe eliminarse manualmente primero. Luego, ejecute `make docker` seguido de` make docker-run`. Luego, conéctese al adaptador en la interfaz de usuario en "Configuración> Mallas de servicio" usando `localhost: <puerto>` si el servidor de malla se está ejecutando como un binario, o <docker IP address>:<port> si se está ejecutando como un contenedor docker.
1. Usando mesheryctl: En `~ / .meshery / meshery.yaml`, cambie la etiqueta que especifica la imagen del adaptador a“ más reciente ”. Ejecute make docker, seguido de `mesheryctl system start --skip-update`. Esto supone que el inicio del sistema mesheryctl se ha ejecutado al menos una vez antes.

### Ejecutando un adaptador como proceso

Otra forma de probar los cambios locales es ejecutar el adaptador como un proceso. Para hacer esto, clone el repositorio de meshery e inicie meshery usando `make server-local`. Inicie el adaptador desde su IDE o ejecutando make run. Luego, en la interfaz de malla, agregue el adaptador usando “localhost: <PORT>”.

### Creación de un nuevo adaptador Meshery

Meshery usa adaptadores para administrar e interactuar con diferentes mallas de servicio. Los adaptadores Meshery están escritos en Go. Ya sea que esté creando un nuevo adaptador o modificando uno existente, asegúrese de leer la especificación de diseño de [Adaptadores Meshery](https://docs.google.com/document/d/1b8JAMzr3Rntu7CudRaYv6r6ccACJONAB5t7ISCaPNuA/edit#). Para nuevos adaptadores, comience con la plantilla de repositorio (https://github.com/layer5io/layer5-repo-template).

1. Obtenga el archivo de especificaciones proto buf del repositorio de Meshery:
   `wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto`
1. Generar codigo
   1. Usando Go como ejemplo, haga lo siguiente:
      - agregando GOPATH to PATH: `export PATH=$PATH:$GOPATH/bin`
      - instalar grpc: `go get -u google.golang.org/grpc`
      - instalar el complemento de protocolo para ir: `go get -u github.com/golang/protobuf/protoc-gen-go`
      - Generar código Go: `protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/`
   1. Para otros idiomas, consulte gRPC.io para obtener guías específicas de idiomas.
1. Implemente los métodos de servicio y exponga el servidor gRPC en un puerto de su elección (por ejemplo, 10000).

Consejo: El [Adaptador Meshery para Istio](https://github.com/layer5io/meshery-istio) es un buen adaptador de referencia para usar como ejemplo de Adaptador Meshery.
