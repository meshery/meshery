---
layout: page
title: Contribuir
permalink: es/project/contributing
language: es
type: project
---

# <a name="contributing">Descripción general de contribución</a>

Por favor, ¡hazlo! ¡Gracias por tu ayuda! 🎈

Este proyecto está construido por la comunidad y la colaboración es bienvenida. Se espera que cada persona colaboradora se adhiera al [Código de conducta de CNCF](https://github.com/meshery/meshery/blob/master/CODE_OF_CONDUCT.md).

¿No estás seguro/a por dónde empezar?

Sigue estos pasos y te sentirás como en casa.

1. Consulte la [_Guía de bienvenida de la comunidad_](https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit) para saber cómo, dónde y por qué contribuir.
2. Regístrate para un [_MeshMate_](https://layer5.io/community#meshmate) y encontrar el Mentor (o Mentora) perfecto que te ayude a explorar los proyectos de Layer5 y encontrar tu lugar en la comunidad:

- **Familiarízate** con todos los proyectos de Layer5 (echale un vistazo al [Drive de la Comunidad](https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA) y la [lista de repositorios de Layer5](https://docs.google.com/document/d/1brtiJhdzal_O6NBZU_JQXiBff2InNtmgL_G1JgAiZtk/edit#header=h.uwtb5xf7b5hw): Dedicale tiempo a comprender cada una de las iniciativas de Layer5 a través de descripciones de alto nivel disponibles en el drive de comunidad y a través de discusiones con tu Mesh Mate.
- **Identifica** tu área de interés: Usa el tiempo con tu MeshMate para familiarizarte con la arquitectura y las tecnologías utilizadas en los proyectos. Informa a tu MeshMate de tus habilidades actuales y las habilidades que pretendes desarrollar.
- **Ejecuta** Meshery: Ponte en los zapatos del usuario y recorre todas las características y funciones de Meshery como usuario.
- **Construye** Meshery: Asegurate que tienes un entorno de desarrollo usable.
- **Comunicate** con la comunidad de Layer5 uniéndote a la [cuenta de Slack](http://slack.layer5.io).
- **Contribuye** tomando cualquier issue abierto con la etiqueta [help wanted](https://github.com/meshery/meshery/issues/) y entrale. Si es necesario, crea un [nuevo issue](https://github.com/meshery/meshery/issues/new/choose). Todas las [pull requests](https://github.com/meshery/meshery/pulls) deben hacer referencia a un issue abierto. Incluye palabras clave en las descripciones de tus pull requests, así como mensajes de commits, para [cerrar automáticamente los issues en GitHub](https://docs.github.com/en/free-pro-team@latest/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue).

**Secciones**

- Flujo general de contribuciones
- <a href="#commit-signing">Certificado de origen del desarrollador</a>
- Flujo de contribución de Meshery
  - <a href="#contributing-docs">Documentación de Meshery</a>
  - <a href="#contributing-meshery">Meshery Backend</a>
    - <a href="#adapter">Escribir un Meshery Adapter</a>
  - <a href="#contributing-ui">Meshery UI</a>
    Las pautas de estilo de codificación relevantes son los comentarios de revisión de Go Code Review y la sección de formato y estilo de Go Code: las mejores prácticas para entornos de producción de Peter Bourgon.

## <a name="contributing">Flujo de contribución general</a>

Para contribuir a Meshery, siga el flujo de trabajo de fork-and-pull descrito [aquí](CONTRIBUTING-gitflow.md).

### <a name="commit-signing">Certificado de origen del desarrollador</a>

Para contribuir a este proyecto, debes aceptar el Certificado de
origen de desarrollador (DCO en inglés Developer Certificate of Origin) para cada commit que haces. El DCO es una simple declaración de que tu, como contribuyente, tienes el derecho legal de hacer la contribución.

Consulta el archivo [DCO](https://developercertificate.org) para obtener el texto completo de lo que debes aceptar y cómo funciona [aquí](https://github.com/probot/dco#how-it-works).
Para indicar que estás de acuerdo con el DCO para las contribuciones, simplemente agregues una línea a cada uno de tus commits git:

```
Signed-off-by: Jane Smith <jane.smith@example.com>
```

En la mayoría de los casos, puedes agregar esta firma a tu confirmación automáticamente con la marca `-s` o` --signoff` en `git commit`. Debes usar tu nombre real y un correo electrónico accesible (lo sentimos, no se permiten seudónimos ni contribuciones anónimas). Un ejemplo de firma de una confirmación:

```
$ commit -s -m “my commit message w/signoff”
```

Para asegurarte de que todas tus confirmaciones están firmadas, puedes optar por agregar este alias a su global `.gitconfig`:
_~/.gitconfig_

```
[alias]
  amend = commit -s --amend
  cm = commit -s -m
  commit = commit -s
```

O puedes configurar tu IDE, por ejemplo , Visual Studio Code para cerrar automáticamente las confirmaciones por ti:

<a href="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png"><img src="https://user-images.githubusercontent.com/7570704/64490167-98906400-d25a-11e9-8b8a-5f465b854d49.png" width="50%"><a>

### <a name="contributing-docs">Flujo de contribución a la documentación</a>

¡Por favor contribuye! La documentación de Meshery utiliza páginas de GitHub para alojar el sitio de documentos. Obten más información sobre el [framework de documentación de Meshery](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit?usp=sharing). El flujo de proceso de contribución es el siguente:

1. Crea un fork, si aún no lo has hecho, sigue los pasos descritos [aquí](CONTRIBUTING-gitflow.md)
1. En la copia local del fork, navega a la carpeta docs:
   `cd docs`
1. Crea una nueva rama para realizar cambios:
   `git checkout -b <my-changes>`
1. Edita/agrega documentación:
   `vi <specific page>.md`
1. Ejecuta el sitio localmente para obtener una vista previa de los cambios:
   `make site`
1. Has commit, [sign-off](#commit-signing), y envía los cambios a tu rama remota
   `git push origin <my-changes>`
1. Abre pull request (en su navegador web) en nuestro repositorio principal: https://github.com/meshery/meshery.

### <a name="contributing-meshery">Flujo de contribución de Meshery</a>

Meshery está escrito en `Go` (Golang) y usa los módulos Go. La interfaz de usuario se basa en React y Next.js. Para facilitar la creación y el packaging, se incluye un archivo Makefile en la carpeta principal del repositorio.

Las guías de estilo de programacción relevantes son los [Comentarios de revisión de código de Go](https://code.google.com/p/go-wiki/wiki/CodeReviewComments) y la sección _Formato y estilo_ de Peter Bourgon's [Go:BestPractices for Production](https://peter.bourgon.org/go-in-production/#formatting-and-style).

**Por favor toma en cuenta**: Todos los comandos `make` deben ejecutarse en una terminal desde la carpeta principal de Meshery.

#### Requisitos previos para construir Meshery en tu entorno de desarrollo:

1. Tener la versión 1.11+ de`Go` instalada si quieres compilar y/o hacer cambios en el código existente.
1. La variable de entorno `GOPATH` debe configurarse de manera apropiada.
1. `npm` y `node` deben estar instalados en tu máquina, preferiblemente las últimas versiones.
1. Haz fork de este repoisitorio (`https://github.com/meshery/meshery.git`), clona tu version forked de Meshery a tu maquina local, preferiblemente fuera de tu `GOPATH`. Si llegase a suceder que copiaras Meshery dentro de tu `GOPATH` y tuvieras una version de go `Go` anterior a la versión 1.13, por favor establece la variable de entorno `GO111MODULE=on` para habilitar los modulos de Go.

#### Construir y ejecutar el servidor Meshery

Para construir y ejecutar el código del servidor Meshery, ejecuta el siguiente comando:

```sh
make run-local
```

Cada vez que se realicen cambios en el código GO, tendrás que detener el servidor y ejecutar el comando anterior nuevamente.
Una vez que el servidor Meshery esta en funcionamiento, deberías poder acceder a Meshery en tu `localhost` en el puerto` 9081` en `http://localhost:9081`. Una cosa para tomar en cuenta, es que posiblemente NO veas la [interfaz de usuario de Meshery](# contrib-ui) hasta que también hayas creado el código de la interfaz de usuario.
Después de ejecutar el servidor Meshery, deberás seleccionar tu **Proveedor de nube** navegando a `localhost:9081`. Solo entonces podrás utilizar la interfaz de usuario de Meshery en el puerto `3000`.

#### Creación de una imagen de Docker

Para crear una imagen de Docker de Meshery, asegúrete de tener instalado `Docker` para poder crear la imagen. Ahora, ejecuta el siguiente comando para construir la imagen de Docker:

```sh
make docker
```

#### <a name="adapter">Escribiendo un Adaptador de Meshery</a>

Meshery usa adaptadores para hacer provisiones e interactuar con diferentes meshes de servicio. Sigue estas instrucciones para crear un nuevo adaptador o modificar un adaptador existente.

1. Obtén el archivo proto buf spec del repositorio de Meshery:
   `wget https://raw.githubusercontent.com/meshery/meshery/master/meshes/meshops.proto`
1. Genera el código
   1. Usando Go como ejemplo, haz lo siguiente::
      - agrega GOPATH a tu PATH: `export PATH=$PATH:$GOPATH/bin`
      - instala grpc: `go get -u google.golang.org/grpc`
      - instala el plugin protoc para go: `go get -u github.com/golang/protobuf/protoc-gen-go`
      - Genenra el código de Go: `protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/`
   1. Para otros lenguajes consulta gRPC.io para las guías especificas del lenguaje.
1. Implementa los métodos de servicio (service method) y expón el puerto de gRPC en el puerto de tu elección (por ejemplo 10000).

_Tip:_ El [adaptador de Meshery para Istio](https://github.com/meshery/meshery-istio) es un buen adaptador de referencia para usar como ejemplo de un adaptador Meshery escrito en Go.

### <a name="contributing-ui">Flujo de Contribución al UI</a>

Meshery está escrito en `Go` (Golang) y aprovecha los módulos Go. La interfaz de usuario se basa en React y Next.js. Para facilitar la creación y el empaquetado, se incluye un archivo Makefile en la carpeta principal del repositorio.

#### Instalar las dependencias de UI

Para instalar/actualizar las dependencias de UI:

```
make setup-ui-libs
```

#### Construir y exportar UI

Para construir y exportar el código de UI:

```
make build-ui
```

Ahora que el código de la interfaz de usuario está creado, la interfaz de usuario de Meshery estará disponible en` http: // localhost: 9081`.
Cada vez que se realizan cambios en el código de la interfaz de usuario, el código anterior deberá ejecutarse para reconstruir la interfaz de usuario.

#### Servidor de desarrollo de UI

Si deseas trabajar en la UI, será una buena idea utilizar el servidor de desarrollo de UI incluido. Puedes ejecutar el servidor de desarrollo de UI ejecutando el siguiente comando:

```
make run-ui-dev
```

Asegúrete de tener el servidor Meshery configurado y en funcionamiento en el puerto predeterminado `http://localhost:9081` antes de proceder a acceder y trabajar en el servidor de UI en `http://localhost:3000`.
Cualquier cambio en la interfaz de usuario realizado ahora se volverá a compilar y se servirá automáticamente en el navegador.

#### Ejecutando Meshery desde IDE

Si deseas ejecutar Meshery desde un IDE como Goland, VSCode. establece la variable de entorno descrito aqui:

```
PROVIDER_BASE_URLS="https://meshery.layer5.io"
PORT=9081
DEBUG=true
ADAPTER_URLS=mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001 mesherylocal.layer5.io:10002 mesherylocal.layer5.io:10003 mesherylocal.layer5.io:10004 mesherylocal.layer5.io:10005 mesherylocal.layer5.io:10006 mesherylocal.layer5.io:10007 mesherylocal.layer5.io:10008 mesherylocal.layer5.io:10009
```

argumento de go tool

```shell
-tags draft
```

actualiza /etc/hosts

```shell
127.0.0.1 mesherylocal.layer5.io
```
