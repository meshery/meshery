---
layout: default
title: Usando mesheryctl
description: Como usar mesheryctl
permalink: es/guides/mesheryctl 
language: es
lang: es
categories: es
type: Guides
---

`mesheryctl` es la interfaz de línea de comandos para administrar Meshery y la interfaz con su funcionalidad mediante un terminal. Los comandos `mesheryctl` se clasifican en tres áreas principales:

- Gestión del ciclo de vida de Meshery (controle el ciclo de vida de Meshery con comandos como `system start`, `stop`, `status`, `reset`).
- Gestión del ciclo de vida de las mallas de servicio
- Gestión del rendimiento de cargas de trabajo y mallas de servicio

<!-- Running `reset` will remove all active container instances, prune pulled images and remove any local volumes created by starting Meshery. -->

## Guías Relacionadas

- Para obtener una lista exhaustiva de comandos y sintaxis, consulte la **[Referencia del Comando `mesheryctl`]({{ site.baseurl }}/guides/mesheryctl-commands)**.
- Para actualizar `mesheryctl`, consulte el **[Guía de Actualización]({{ site.baseurl }}/guides/upgrade)**.

## Instalando `mesheryctl`

### Mac ó Linux

Use su elección de homebrew o bash para instalar `mesheryctl`. Solo necesitas usar uno.
### Homebrew

Instale `mesheryctl` y ejecute Meshery en Mac con Homebrew.

#### Instalando con Homebrew

Para instalar `mesheryctl`, ejecute los siguientes comandos:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew tap layer5io/tap
 brew install mesheryctl
 mesheryctl system start
 </div></div>
 </pre>

**Actualizando con Homebrew**

Para actualizar `mesheryctl`, ejecute el siguiente comando:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew upgrade mesheryctl
 </div></div>
 </pre>

#### Bash

**Instalando con Bash**

Instalando `mesheryctl` y ejecute Meshery en Mac ó Linux con este script:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

**Actualizando con Bash**

Actualice `mesheryctl` y ejecute Meshery en Mac ó Linux con este script:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 curl -L https://git.io/meshery | bash -
 </div></div>
 </pre>

## Windows

### Instalando el binario `mesheryctl`

Descargue y descomprima `mesheryctl` desde la página de [liberaciones de Meshery](https://github.com/layer5io/meshery/releases/latest). Agregue `mesheryctl` a su variable de entorno PATH para facilitar su uso. Después, ejecute:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 ./mesheryctl system start
 </div></div>
 </pre>

### Scoop

Utilice [Scoop](https://scoop.sh) para instalar Meshery en su máquina Windows.

**Instalando con Scoop**

Agregue el Meshery Scoop Bucket e instale:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 scoop bucket add mesheryctl https://github.com/layer5io/scoop-bucket.git
 scoop install mesheryctl
 </div></div>
 </pre>

**Actualizando con Scoop**

Para actualizar `mesheryctl`, ejecute el siguiente comando:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 scoop update mesheryctl
 </div></div>
 </pre>

## Instalación Avanzada

Los usuarios pueden controlar la imagen del contenedor específico y la etiqueta (versión) de Meshery que les gustaría ejecutar editando su *~/.meshery/meshery.yaml* (un archivo de docker compose).

Alineado con la imagen del contenedor de Meshery, en lugar de dejar la etiqueta implícita: stable-latest detrás de image: layer5 / meshery, los usuarios identificarán una etiqueta de imagen específica de la siguiente manera:

```
bash
version: '3'
services:
  meshery:
    image: layer5/meshery:v0.5.0
    labels:
      - "com.centurylinklabs.watchtower.enable=true"
```

# Configurando Autocompletado para `mesheryctl`

Si desea que los comandos `mesheryctl` se completen automáticamente para usarlos mientras usa `mesheryctl`, utilice las siguientes instrucciones para configurar la finalización automática dentro de su entorno.

## Autocompletado para Bash

### bash <= 3.2

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source /dev/stdin <<< "$(mesheryctl system completion bash)"
 </div></div>
 </pre>

### bash >= 4.0

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source <(mesheryctl system completion bash)
 </div></div>
 </pre>

### bash <= 3.2 en MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew install bash-completion # ensure you have bash-completion 1.3+
 mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
 </div></div>
 </pre>

### bash >= 4.0 en MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew install bash-completion@2
 mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
 </div></div>
 </pre>

## Autocompletado para zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source <(mesheryctl system completion zsh)
 </div></div>
 </pre><br>


Si la finalización de shell aún no está habilitada en su entorno, deberá habilitarla. Puede ejecutar lo siguiente una vez:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 ~/.zshrc > echo "autoload -U compinit; compinit" 
 </div></div>
 </pre>
_Nota_ : Es posible que deba reiniciar su shell para que esta configuración surta efecto.

#### zsh en MacOS y Oh My zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion zsh > "${fpath[1]}/_mesheryctl"
 </div></div>
 </pre>

### Autocompletado para fish

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion fish | source
 </div></div>
 </pre><br>


Para cargar finalizaciones de concha de pescado para cada sesión, ejecute una vez:
 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion fish > ~/.config/fish/completions/mesheryctl.fish
 </div></div>
 </pre>

# Lectura Sugerida

- Para obtener una lista exhaustiva de comandos y sintaxis, consulte la **[Referencia del Comando `mesheryctl`]({{ site.baseurl }}/es/guides/mesheryctl-commands)**.
- Para actualizar `mesheryctl`, referirse a la **[Guía de Actualización]({{ site.baseurl }}/es/guides/upgrade)**.
