---
layout: page
title: KinD
permalink: es/installation/platforms/kind
language: es
type: installation
---

# Inicio Rápido con KinD

Para configurar y ejecutar Meshery en KinD:

- [Instalar kinD](#instalación)
- [Crear un cluster Kubernetes con kinD](#crear-cluster-usando-kind)
  - [Acceder al cluster kinD](#Accediendo-el-cluster-kind)
- [Alternativamente, Ejecute Helm](#usando-helm)

### **Instalación**

- En Mac / Linux mediante Homebrew (Recomendado):

```powershell
brew install kind
```

- En macOS / Linux mediante curl:

```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.8.1/kind-$(uname)-amd64
chmod +x ./kind
mv ./kind /some-dir-in-your-PATH/kind
```

Si está ejecutando Ubuntu en WSL2, use la distro `Docker Ubuntu` para instalar `Docker`.

#### **Crear cluster usando KinD**

Con el fin de construir exitosamente el servidor Meshery en su servidor local, siga las instrucciones específicas a su Sistema Operativo para completar la creación de un cluster KinD.

###### 1. **KinD en WSL2**

Primero que nada, obtendremos la dirección ip de su distro WSL2 mediante:

```bash
ip addr | grep eth0
```

Usted verá una salida como:

```bash
4: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    inet 172.1.1.1/20 brd 172.1.1.255 scope global eth0
```

Copie la dirección IP, la usaremos en el siguiente paso.

Lo siguiente es, crear un archivo que tenga el nombre `kind_cluster.yaml` y capture la dirección ip en el campo `apiServerAddress`:

```bash
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "172.1.1.1"
```

Ahora cree el cluster KinD con el archivo de configuración `kind_cluster.yaml`:

```
kind create cluster --config kind_cluster.yaml --name kind --wait 300s
```

Usted verá:

```bash
Creating cluster "kind" ...
 • Ensuring node image (kindest/node:v1.17.0) 🖼  ...
 ✓ Ensuring node image (kindest/node:v1.17.0) 🖼
 • Preparing nodes 📦   ...
 ✓ Preparing nodes 📦
 • Writing configuration 📜  ...
 ✓ Writing configuration 📜
 • Starting control-plane 🕹️  ...
 ✓ Starting control-plane 🕹️
 • Installing CNI 🔌  ...
 ✓ Installing CNI 🔌
 • Installing StorageClass 💾  ...
 ✓ Installing StorageClass 💾
 • Waiting ≤ 5m0s for control-plane = Ready ⏳  ...
 ✓ Waiting ≤ 5m0s for control-plane = Ready ⏳
 • Ready after 59s 💚
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Not sure what to do next? 😅 Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

###### 2. **KinD en otros sistemas**

Crear un cluster Kubernetes es tan simple como ejecutar `kind create cluster`.

Para obtener más configuración de la instalación, consulte la documentación oficial de KinD.

#### **Accediendo el cluster KinD**

De forma predeterminada, la configuración de acceso al clúster se almacena en ${HOME}/.kube/config si la variable de entorno $KUBECONFIG no está establecida. Usted puede establecer esta variable de ambiente `KUBECONFIG` con el siguiente comando:

```bash
export KUBECONFIG=${HOME}/.kube/config
```

Utilice el siguiente comando para verificar la conexión del clúster y asegurarse de que el clúster al que se que conectó, es el clúster que fue creado por KinD:

```bash
kubectl cluster-info --context kind-kind
```

Para borrar su cluster use:

```bash
kind delete cluster --name kind
```

#### **Usando Helm**

##### **Helm v3**

Recomendamos encarecidamente utilizar Helm v3, debido a que la versión actual ya no incluye el componente Tiller(https://helm.sh/blog/helm-3-preview-pt2/#helm). Es más ligero y seguro de esta manera.

Ejecute lo siguienteg:

```bash
$ git clone https://github.com/meshery/meshery.git; cd meshery
$ kubectl create namespace meshery
$ helm install meshery --namespace meshery install/kubernetes/helm/meshery
```

- **NodePort** - Si su cluster no cuenta con un Controlador Ingress (Ingress Controller) o un balanceador de cargas, entonces utilice NodePort para exponer Meshery y que pueda ser modificado en `values.yaml`:

```bash
service:
  type: NodePort
  port: 9081
  annotations: {}
```
