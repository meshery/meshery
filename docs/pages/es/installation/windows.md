---
layout: page
title: Windows
permalink: es/installation/platforms/windows
language: es
type: installation
---

# **Inicio rápido con Windows**

Para configurar y ejecutar Meshery en Windows:

1. <a href="#step1">Configurar Windows y habilitar Docker </a>
2. <a href="#step4">Instalar un cluster Kubernetes</a>
3. <a href="#step5">Instalar Meshery</a>

### **Compatibilidad**

Las siguientes versiones mínimas de Compilación de Windows son requeridas:

<table id="compatibility-table">
  <tr>
    <th id="model">Nombre</th>
    <th id="model">Versión</th>
  </tr>
  <tr>
    <td><a href="#wsl1">WSL1</a></td>
    <td><b>x64</b> - Windows 7 </td>
  </tr>
  <tr>
    <td><a href="#wsl2">WSL2</a></td>
    <td><b>x64</b> - Version 1903, Compilación 18362; <b>ARM 64</b> - Version 2004, Compilación 19041</td>
  </tr>
  <tr>
    <td><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18945">Kernel Personalizado</a></td>
    <td>Compilación 18945</td>
  </tr>
  <tr>
    <td><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-19013">Kernel con módulos requeridos K8s</a></td>
    <td>Compilación 19013</td>
  </tr>
</table>

**Nota**
<br />Ejecute el siguiente comando en Powershell para comprobar la compilación y la versión de Windows:

```powershell
[System.Environment]::OSVersion.Version
```

### **Pasos**

Realice los siguientes pasos en orden:

#### 1. <a name="step1" href="https://docs.microsoft.com/en-us/windows/wsl/install-win10"><b>Instalar el subsistema de Windows para Linux (WSL)</b></a>

Abra la Powershell en modo administrador y ejecute:

```powershell
dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
Restart-Computer
```

##### **Elegir tu versión WSL:**

<h6><b><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18917" name="wsl2">WSL2</a></b> (Recomendado)</h6>
Elija la versión por defecto a `WSL2`, la cual será heredada por cualquier distro que desee usar.

**Habilite la característica VM (Virtual Machine)**:

```powershell
dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart
```

**Colocar WSL2 como la versión default**:

```bash
wsl --set-default-version 2
```

<h6><b><a href="https://docs.microsoft.com/en-us/windows/wsl/install-win10" name="wsl1"> WSL1 </a></b></h6>

<b>Advertencia:</b>
Se recomienda actualizar a <a href="#wsl2">WSL2</a> ya que WSL1 no soporta la aplicación Docker Desktop para Windows. En su lugar, solo soporta la versión obsoleta, [Docker Toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/).

Si aun desea continuar, siga las instrucciones para <button class="toggle-button" onclick="HideToggleFunction()"><b>WSL1</b></button>

<div id="hiddendiv">
<p>
1. La versión por defecto de WSL está establecida en WSL1 por defecto. Puedes ir al siguiente link para <a href="https://docs.microsoft.com/en-us/windows/wsl/install-win10#install-your-linux-distribution-of-choice">instalar la distro</a> de tu elección. <br /><br />

2. <b><a href="https://docs.docker.com/toolbox/toolbox_install_windows/">Docker Toolbox</a></b> <br />

<b>Advertencia</b>: Docker Toolbox es una versión obsoleta. Se recomienda actualizar el sistema e instalar la aplicación Docker Desktop con WSL2. <br/><br />

Docker Toolbox utiliza características específicas del kernel de Linux, y no puede funcionar de forma nativa en Windows. En su lugar, crea y utiliza una pequeña VM de Linux en tu máquina junto con <a href="https://docs.docker.com/machine/overview/"><code>docker-machine</code> </a>, y utiliza VirtualBox para ejecutar Docker. <br />

<ul>
<li> Diríjase a <a href="https://github.com/docker/toolbox/releases">Toolbox Releases</a> y descargue la última versión del archivo <code>.exe</code></li>
<li> Siga estas <a href="https://docs.docker.com/toolbox/toolbox_install_windows/#step-2-install-docker-toolbox">instrucciones</a> para configurar con éxito la aplicación Docker Toolbox. </li>
</ul>

</p>
</div>

#### 2. <b>[Instalar la nueva distro](https://docs.microsoft.com/en-us/windows/wsl/install-win10#install-your-linux-distribution-of-choice)</b>

En este tutorial, [Ubuntu 18.04](https://www.microsoft.com/en-us/p/ubuntu-1804-lts/9n9tngvndl3q?activetab=pivot:overviewtab) será la distro utilizada. Siéntase libre de usar cualquier distro a su elección.

#### 3. <b>habilitar Docker</b>

La aplicación Docker Desktop para Windows incluye un completo conjunto de herramientas, incluyendo Docker Engine, el cliente Docker CLI, Docker Compose, Notary, Kubernetes, y un Credential Helper.

<table id="compatibility-table">
  <tr>
    <th id="model">Versión Windows 10</th>
    <th id="model">Docker Desktop</th>
  </tr>
  <tr>
    <td>Pro/Education/Enterprise</td>
    <td><a href="https://docs.docker.com/docker-for-windows/install/">Docker Desktop para Windows Pro</a></td>
  </tr>
  <tr>
    <td>Home</td>
    <td><a href="https://docs.docker.com/docker-for-windows/install-windows-home/">Docker Desktop para Windows Home</a></td>
  </tr>
</table>

#### 4. <a name="step4"> <b>Instalar un cluster Kubernetes</b></a>

Una vez que Docker esté instalado, el siguiente paso será instalar un cluster Kubernetes. En esta guía, se usará [K3d](https://github.com/rancher/k3d) ya que sólo depende de Docker.

```bash
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
k3d cluster create
export KUBECONFIG="$(k3d kubeconfig get 'k3s-default')"
```

#### 5. <a name="step5"><b>Instalar Meshery</b></a>

Siga los [pasos de instalación]({{ site.baseurl }}/es/installation/windows.md) para instalar el CLI mesheryctl. Luego, ejecute:

```bash
./mesheryctl system start
```
