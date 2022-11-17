---
layout: page
title: Minikube
permalink: es/installation/platforms/minikube
language: es
type: installation
---

# Inicio Rápido con Minikube

Para configurar y ejecutar Meshery en Minikube:

1. [Iniciar Minikube](#1-iniciar-minikube)
1. [Configurar Meshery para usar minkube](#2-configurar-meshery-para-usar-minikube)
1. [Ejecutar Meshery](#3-configurar-meshery)

### **Compatibilidad**

Las siguientes versiones mínimas de componentes son requieridas:

<table id="compatibility-table">
  <tr>
    <th id="model">Nombre</th>
    <th id="model">Versión</th> 
  </tr>
  <tr>
    <td><a href="https://kubernetes.io/docs/tasks/tools/install-minikube/">Minikube</a></td>
    <td>1.0.0 </td>
  </tr>
  <tr>
    <td><a href="https://istio.io/docs/setup/kubernetes/prepare/platform-setup/minikube/">Kubernetes</a></td>
    <td>1.14.1</td>
  </tr>
  <tr>
    <td><a href="https://kubernetes.io/docs/tasks/tools/install-kubectl/">kubectl</a></td>
    <td>1.14.1</td>
  </tr>
</table>

### **Pasos**

Ejecute los siguientes pasos en este orden:

#### 1. **iniciar minikube**

```bash
minikube start --cpus 4 --memory 8192 --kubernetes-version=v1.14.1
```

_Nota: la memoria mínima requerida es --memory=4096 (para despliegues de Istio solamente)_

_Note: Si está utilizando manejador docker (docker driver), después de completar los pasos de instalación de Meshery, ejecute el siguiente comando para establecer la conectividad entre Meshery y el servidor Kubernetes._

```
docker network connect bridge meshery_meshery_1
```

#### 2. **Configurar Meshery para usar minikube**

1. Inicie sesión en Meshery. En el menú de su perfil de usuario, haga clic en `Get Token`.
2. Use `mesheryctl` y así configurar Meshery para usar minikube. Ejecute:

```sh
mesheryctl system config minikube -t ~/Downloads/auth.json
```

**También puede generar y cargar manualmente el archivo kubeconfig para que Meshery lo use:**

Este archivo de configuración será utilizado por Meshery.

```
kubectl config view --minify --flatten > config_minikube.yaml
```

```
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: < cert shortcutted >
    server: https://192.168.99.100:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate-data: <cert shortcutted >
    client-key-data: < key shortcutted >
```

Nota: Asegúrese que _current-context_ (contexto actual) esté establecido a `minikube`.

#### 3. **Configurar Meshery**

Siga los [pasos de instalación]({{ site.baseurl }}/installation) para instalar el CLI (Intérprete de Línea de Comandos) mesheryctl.

Meshery ahora debería estar conectado con su instancia administrada de Kubernetes. Eche un vistazo a [guías Meshery]({{ site.baseurl }}/guides) para consejos de uso advanzado.
