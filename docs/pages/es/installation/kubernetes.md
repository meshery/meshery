---
layout: page
title: Kubernetes
permalink: es/installation/platforms/kubernetes
language: es
type: installation
---

# Inicio Rápido con Kubernetes

Para preparar y ejecutar Meshery en Kubernetes:

- [ Use Helm y configure un cluster Kubernetes ](#usando-helm)

- [Ejecutar Meshery en un cluster Kubernetes existente](#usando-el-manifiesto-de-kubernetes)

#### **Usando Helm**

##### 1. **Helm v3**

Ejecute lo siguiente:

```sh
$ git clone https://github.com/meshery/meshery.git; cd meshery
$ kubectl create namespace meshery
$ helm install meshery --namespace meshery install/kubernetes/helm/meshery
```

##### 2. **Helm v2**

Ejecute lo siguiente:

```sh
$ git clone https://github.com/meshery/meshery.git; cd meshery
$ kubectl create namespace meshery
$ helm template meshery --namespace meshery install/kubernetes/helm/meshery | kubectl apply -f -
```

#### **Usando el Manifiesto de Kubernetes**

Meshery también puede ser desplegado(deployed) en un clúster de Kubernetes existente.Vea [tabla de compatibilidades](#matriz-de-compatibilades) para la compatibilidad de versiones. Para instalar Meshery en su cluster, clone el repositorio de Meshery:

```sh
$ git clone https://github.com/meshery/meshery.git; cd meshery
```

Cree un nombre de espacio (namespace) como un nuevo espacio lógico para hospedar Meshery y sus componentes:

```sh
$ kubectl create ns meshery
```

Todos los yamls de despliegue necesarios para desplegar(deploy) Meshery se incluyen en la carpeta ʻinstall / deployment_yamls / k8s` dentro de la carpeta Meshery clonada. Para implementar los yamls en el clúster, ejecute el siguiente comando:

```sh
$ kubectl -n meshery apply -f install/deployment_yamls/k8s
```

Una vez que se desplieguen los archivos yaml, necesitamos exponer el servicio `meshery` para poder acceder al servicio desde fuera del clúster. Hay varias formas de exponer un servicio en Kubernetes. Aquí describiremos 3 formas comunes en las que podemos exponer un servicio:

- **Ingress** - Si su cluster Kubernetes tiene un Controlador Ingress (Ingress Controller) funcional, entonces usted puede configurar un ingreso para exponer Meshery:

```sh
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
name: meshery-ingress
annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
rules:
- host: *
    http:
    paths:
    - path: /
        backend:
        serviceName: meshery-service
        servicePort: 9081
```

- **Balanceador de Cargas** -
  Si su clúster de Kubernetes tiene un balanceador de carga externo, esta podría ser la ruta lógica.

- **NodePort** -
  Si su clúster no tiene un Controlador Ingress (Ingress Controller) o un balanceador de carga, use NodePort para exponer Meshery:

```sh
apiVersion: v1
kind: Service
spec:
    type: NodePort
```

Meshery debe ahora ser conectado con su instancia administrada de Kubernetes. Eche un vistazo a las guías [guías Meshery]({{ site.baseurl }}/guides) para tips de uso advanzado.
