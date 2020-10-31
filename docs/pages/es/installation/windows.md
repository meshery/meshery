---
layout: page
title: Windows
permalink: es/installation/platforms/windows
---

## **Inicio Rápido_con Windows**

Para configurar y ejecutar Meshery en Windows:

1. <a href="#step1">Configure Windows y habilite Docker </a>
2. <a href="#step4">Instale un cluster Kubernetes encima </a>
3. <a href="#step5">Ejecute Meshery</a>

### **Compatibilidad**

Las siguientes versiones de Windows son las requeridas:

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
    <td><b>x64</b> - Version 1903, Build 18362; <b>ARM 64</b> - Version 2004, Build 19041</td>
  </tr>
  <tr>
    <td><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18945">Custom Kernel</a></td>
    <td>Build 18945</td>
  </tr>
  <tr>
    <td><a href="https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-19013">Kernel con módules K8s requeridos</a></td>
    <td>Build 19013</td>
  </tr>
</table>

**Note**
<br />Ejecute el siguiente comando para verificar su número de compilación y versión de Windows:

```powershell
[System.Environment]::OSVersion.Version
```

### **Pasos**

Lleve a cabo estos pasos en este orden:

#### 1. <a name="step1" href="https://docs.microsoft.com/es-es/windows/wsl/install-win10"><b>Instalar Subsistema Windows para Linux (WSL)</b></a>

abra su Powershell en modo administrador y ejecute:

```powershell
dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
Restart-Computer
```

##### **Eligiendo su versión WSL:**

<h6><b><a href="https://docs.microsoft.com/es-es/windows/wsl/release-notes#build-18917" name="wsl2">WSL2</a></b> (Recomendado)</h6>
Establezca la versión por defecto a `WSL2`, la cual será heredada por cualquier otra distribución que desee utilizar.

**Habilite la característica VM (Virtual Machine)**:

```powershell
dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart
```

**Establezca WSL2 como la versión por defecto**:

```bash
wsl --set-default-version 2
```

<h6><b><a href="https://docs.microsoft.com/es-es/windows/wsl/install-win10" name="wsl1"> WSL1 </a></b></h6>

<b>Precaución:</b>
Es recomendado que actualice a <a href="#wsl2">WSL2</a> ya que WSL1 no brinda soporte para la aplicación Docker de Escritorio (Docker Desktop) para Windows En vez de esto, solo soporta la versión obsoleta , [Docker Toolbox](https://docs.docker.com/toolbox/toolbox_install_windows/).

si desea contonuar, siga las instrucciones para <button onclick="HideToggleFunction()"><b>WSL1</b></button>

<div id="hiddendiv">
<p>
1. La versión predeterminada de WSL se establece en WSL1 de forma predeterminada. Puede continuar a <a href="https://docs.microsoft.com/es-es/windows/wsl/install-win10#install-your-linux-distribution-of-choice">instale la distribución</a> de su preferencia. <br /><br />

2. <b><a href="https://docs.docker.com/toolbox/toolbox_install_windows/">Caja de Herramientas Docker (Docker Toolbox)</a></b> <br />

<b>Precaución</b>: Docker Toolbox es ya una versión obsoleta. Se recomienda actualizar su sistema e instalar la aplicación Docker Desktop con WSL2. <br/><br />

Docker Toolbox usa funciones de kernel específicas de Linux y no se puede ejecutar de forma nativa en Windows. En cambio, cree y utilice una pequeña máquina virtual Linux en su máquina junto con una<a href="https://docs.docker.com/machine/overview/"><code>máquina
docker</code></a>, y utilice VirtualBox para ejecutar Docker. <br />
<ul>
<li> Ir a <a href="https://github.com/docker/toolbox/releases">Liberaciones Toolbox</a> y descargue el archivo liberado <code>.exe</code> mas actual.</li>
<li> Siga estas <a href="https://docs.docker.com/toolbox/toolbox_install_windows/#step-2-install-docker-toolbox">instrucciones</a> para configurar exitosamente la aplicación Docker Toolbox.</li>
</ul>

</p>
</div>

#### 2. <b>[Instale una nueva distribución](https://docs.microsoft.com/en-us/windows/wsl/install-win10#install-your-linux-distribution-of-choice)</b>

En este tutorial, [Ubuntu 18.04](https://www.microsoft.com/es-es/p/ubuntu-1804-lts/9n9tngvndl3q?activetab=pivot:overviewtab) será la distribución utilizada. Siéntase en libertad de usar la distribución de su preferencia.

#### 3. <b>Habilitar Docker</b>

La aplicación Docker Desktop para Windows incluye un conjunto completo de herramientas, que incluyen Docker Engine, Docker CLI client, Docker Compose, Notary, Kubernetes y Credential Helper.

<table id="compatibility-table">
  <tr>
    <th id="model">Versión de Windows 10</th>
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

Una vez que Docker esté instalado, el siguiente paso será instalar un clúster de Kubernetes.
En este tutorial, [K3d](https://github.com/rancher/k3d) será utilizado, ya que se basa solo en Docker.

```bash
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
k3d cluster create
export KUBECONFIG="$(k3d kubeconfig get 'k3s-default')"
```

#### 5. <a name="step5"><b>Configurar Meshery</b></a>

Siga los [pasos de instalación steps](/docs/installation#windows) para instalar el CLI (Intérprete de Línea de Comandos) mesheryctl. Despues, ejecute:

```bash
./mesheryctl system start
```
