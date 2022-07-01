---
layout: page
title: Installation Guide
permalink: es/installation
language: es
type: installation
language: en
list: exclude
---

<a name="getting-started"></a>

## Inicio Rápido

Poner a correr Meshery en un sistema habilitado con Docker, es fácil. Usa la interfaz de linea de comandos Meshery, `mesheryctl`, para empezar Meshery en cualquiera de sus [plataformas soportadas]({{ site.baseurl }}/installation/platforms).

## Usando `mesheryctl`

`mesheryctl` es una interfaz de línea de comandos para administrar un Desplegado (deployment) Meshery. `mesheryctl` le permite controlar el ciclo de vida de Meshery con comandos como `start`, `stop`, `status`, `reset`. Ejecutando `reset` removerá todas las instancias de contenedor activas, poda las imágenes jaladas y remueve todos los volúmenes locales creados iniciando Meshery.

### Mac ó Linux

Emplea tu elección de homebrew o bash para instalar `mesheryctl`. Tu solo necesitas usarlo una vez.

#### Homebrew

Instala `mesheryctl` y ejecuta Meshery en Mac con Homebrew.

**Instalando con Homebrew**

Para instalar `mesheryctl`, ejecute el siguiente comando:

```bash
brew tap layer5io/tap
brew install mesheryctl
mesheryctl system start
```

**Actualizando con Homebrew**

Para actualizar `mesheryctl`, ejecute el siguiente comando:

```bash
brew upgrade mesheryctl
```

#### Bash

**Instalando con Bash**

Instala `mesheryctl` y ejecuta Meshery en Mac ó Linux con este script:

```bash
curl -L https://meshery.io/install | bash -
```

**Actualizando con Bash**

Actualiza `mesheryctl` y corre Meshery en Mac ó Linux con este script:

```bash
curl -L https://meshery.io/install | bash -
```

### Windows

#### Instalando el binario `mesheryctl`

Descarga y descomprime `mesheryctl` desde la página de [Liberaciones Meshery](https://github.com/meshery/meshery/releases/latest). Agrega `mesheryctl` a tu PATH para facilitar el uso. Después, ejecuta:

```bash
./mesheryctl system start
```

#### Scoop

Usa [Scoop](https://scoop.sh) para instalar Meshery en tu máquina Windows.

**Instalando con Scoop**

Agrega el Meshery Scoop Bucket e instala:

```bash
scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
scoop install mesheryctl
```

**Actualizando con Scoop**

Para actualizar `mesheryctl`, ejecuta el siguiente comando:

```bash
scoop update mesheryctl
```

# Instalación Avanzada

Los Usuarios pueden controlar la imagen de contenedor específico y etiqueta(versión) de Meshery que quieran ejecutar mediante la edición de su archivo local `~/.meshery/meshery.yaml` (un archivo de docker compose).
Alineado con la imagen de contenedor Meshery, en lugar de dejar la etiqueta implícita `:stable-latest` detrás de la imagen: layer5/meshery, los usuarios, en vez de esto, identificarán una etiqueta de imagen específica así:

```bash
version: '3'
services:
  meshery:
    image: layer5/meshery:v0.5.0
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
```

---

Cuando Meshery está levantado y corriendo, las instrucciones para acceder Meshery serán impresas en la pantalla y su navegador web por default deberá ser dirigido a la pantalla de inicio de Meshery.

