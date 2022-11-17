---
layout: default
title: Actualizar `mesheryctl` y Meshery
description: Cómo Meshery y todos sus componentes
permalink: es/guides/upgrade
language: es
display-title: "false"
type: Guides
---

# Actualizando Meshery

## Actualizar el Servidor, Adaptadores, e Interfaz Gráfica de Meshery

Varios componentes de Meshery deberán actualizarse a medida que estén disponibles nuevas versiones. Meshery consta de una serie de componentes que incluyen un servidor, adaptadores, UI y CLI. Como aplicación, Meshery es una composición de diferentes componentes funcionales.

<a href="{{site.baseurl}}/assets/img/architecture/upgrading-meshery.svg">
    <img src="{{site.baseurl}}/assets/img/architecture/upgrading-meshery.svg" width="20%" />
</a>

Algunos de los componentes deben actualizarse simultáneamente, mientras que otros pueden actualizarse de forma independiente. La siguiente tabla muestra los componentes, sus versiones y las unidades de implementación (grupos de implementación).

### Versionado de componentes Meshery

<table class="mesherycomponents">
    <tr>
        <th>Components</th>
        <th>Sub-component</th>
        <th>Considering or Updating</th>
    </tr>
    <tr>
        <td class="childcomponent">Meshery Adapters</td>
        <td>Any and All Adapters</td>
        <td>Docker Deployment: Watchtower updates this component in accordance with the user’s release channel subscription.</td>
    </tr>
    <tr>
        <td rowspan="3" class="childcomponent">Meshery Server</td>
        <td>Meshery UI</td>
        <td rowspan="3">Manages lifecycle of Meshery Operator; Adapters, UI, Load Generators, Database.<br /><br />
Docker Deployment: Watchtower updates this component in accordance with the user’s release channel subscription.</td>
    </tr>
    <tr>
        <td>Load Generators</td>
    </tr>
    <tr>
        <td>Database</td>
    </tr>
    <tr>
        <td rowspan="2" class="childcomponent">Meshery Operator</td>
        <td>MeshSync</td>
        <td>Meshery Operator manages the lifecycle of this component and its sub-components.</td>
    </tr>
    <tr>
        <td>Meshery Broker</td>
        <td>Meshery Operator manages the lifecycle of this event bus component.</td>
    </tr>
    <tr>
        <td class="childcomponent">`mesheryctl`</td>
        <td></td>
        <td><code>mesheryctl</code> manages the lifecycle of Meshery Server. <br /><br />
        <ul> 
            <li><code>system start</code> calls system update by default, which updates server and existing adapters, but doesn’t update meshery.yaml.</li>
            <li><code>system reset</code> retrieving docker-compose.yaml from GitHub (use git tag to reset to the right Meshery version).</li>
            <li><code>system context</code> manages config.yaml, which manages meshery.yaml. </li>
            <li><code>mesheryctl</code> should generally be checking for latest release and informing user.</li>
        </ul>
        </td>
    </tr>
    <tr>
        <td rowspan="2" class="childcomponent"><a style="color:white;" ref="/extensibility/providers">Remote Providers</a></td>
        <td>Meshery Cloud</td>
        <td>Process Extension: Integrators manage the lifecycle of their Remote Providers. Process is unique per provider.</td>
    </tr>
    <tr>
        <td>Meshery Cloud</td>
        <td> Static Extension: Integrators manage the lifecycle of their Meshery Extensions. Process is unique per provider.</td>
    </tr>
</table>

Sub-components deploy as a unit, however, they do not share the same version number.

### Despliegues Docker de Meshery

Para actualizar el servidor Meshery, los adaptadores y la interfaz de usuario, ejecute el siguiente comando:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system upgrade
 </div></div>
 </pre>

### Despliegues Kubernetes de Meshery

Utilice `kubectl apply` o` helm` para actualizar los manifiestos de la aplicación Meshery en su clúster de Kubernetes.

## Actualizando `mesheryctl`

El cliente de línea de comandos de Meshery está disponible en diferentes administradores de paquetes. Utilice las instrucciones relevantes para su entorno.

### Actualizando `mesheryctl` usando Homebrew

<p>Para actualizar `mesheryctl`, ejecute el siguiente comando:</p>

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew upgrade mesheryctl
 </div></div>
 </pre>

### Actualizando `mesheryctl` usando Bash

Actualice `mesheryctl` y ejecute Meshery en Mac o Linux con este script:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 curl -L https://meshery.io/install | bash -
 </div></div>
 </pre>

### Actualice `mesheryctl` usando Scoop

Para actualizar `mesheryctl`, ejecute el siguiente comando:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 scoop update mesheryctl
 </div></div>
 </pre>
