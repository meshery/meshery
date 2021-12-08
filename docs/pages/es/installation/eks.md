---
layout: page
title: EKS
permalink: es/installation/platforms/eks
language: es
type: installation
---

# Inicio rápido con el Servicio Elastic Kubernetes de Amazon (EKS)

Para brindar a Meshery el acceso necesario a su instancia de Kubernetes administrada,
Meshery deberá tener asignada una ServiceAccount (Cuenta de servicio). Se puede utilizar una "ServiceAccount" existente o crear una nueva.

1. Asegúrese de que la `ServiceAccount` que utiliza tenga asignada la función `cluster-admin`.
1. Configure Meshery para que se ejecute en EKS:
   - [Configuración automática](#configuración-automática-recomendada)
   - [Configuración manual](#configuración-manualopcional)

_Nota: asegúrese de poder acceder a EKS con `kubectl` siguiendo la [Guía de EKS.](https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html){:target="\_blank"}_

#### **Configuración automática (Recomendada)**

1. En su navegador, navegue hasta Meshery (por ejemplo `http://localhost:9081`) e inicie sesión.
1. Descargue su token de autenticación de Meshery haciendo clic en **Get Token** en la sección de su perfil de usuario.
1. Utilice este token de autenticación para ejecutar el siguiente comando:
   ```
   $ mesheryctl system config eks --token <RUTA A ARCHIVO DE TOKEN>
   ```

Este comando actualiza su kubeconfig para proporcionar a Meshery acceso a su instancia administrada de Kubernetes.
Una vez configurado, prosiga a utilizar Meshery (`mesheryctl system start`).

#### **Configuración Manual(Opcional)**

Si el procedimiento de [Configuración automática](#configuración-automática-recomendada) falla o si desea preparar manualmente su archivo kubeconfig para proveer a Meshery con el acceso necesario a su instancia administrada de Kubernetes, realice las siguientes acciones:

1. Cree una `ServiceAccount` con el rol de `cluster-admin`

   ```sh
   $ kubectl create serviceaccount meshery
   ```

1. Agregando/Vinculando el rol `cluster-admin` a la nueva cuenta de servicio `meshery`

   ```sh
   $ kubectl create clusterrolebinding meshery-binding --clusterrole=cluster-admin \
   --serviceaccount=default:meshery
   ```

1. Obtenga el nombre secreto de `ServiceAccount`.

   ```sh
   $ kubectl get secrets

   NAME                           TYPE                                  DATA   AGE
   default-token-fnfjp            kubernetes.io/service-account-token   3      95d
   meshery-token-5z9xj               kubernetes.io/service-account-token   3      66m
   ```

   _Nota: Aquí el nombre secreto es **meshery-token-5z9xj**_

1. Obtenga el secreto/token:

   ```sh
   $ kubectl describe secret  sa-1-token-5z9xj
   Name:         meshery-token-5z9xj
   Namespace:    default
   Labels:       <none>
   Annotations:  kubernetes.io/service-account.name: meshery
                 kubernetes.io/service-account.uid: 397XXX-XXX-XXXX-XXXXX-XXXXX

   Type:  kubernetes.io/service-account-token

   Data
   ====
   ca.crt:     1025 bytes
   namespace:  7 bytes
   token:      XXXhbGciOiJSUXXXX
   ```

1. Genere un nuevo archivo yaml kubeconfig a usar como entrada para Meshery.
1. Establezca las Credential de configuración usando el `token` generado anteriormente.

   ```sh
   $ kubectl config set-credentials meshery --token=XXXXX

   o/p:User "meshery" set.
   ```

1. Establezca el contexto actual para la nueva cuenta de servicio `meshery`

   ```sh
   $ kubectl config set-context --current --user=meshery

   o/p:
   Context "aws" modified.
   ```

1. Genere el archivo yaml kubeconfig para usar como entrada a Meshery.

   ```sh
   $ kubectl config view --minify --flatten >  config_aws_eks.yaml
   ```

Meshery ahora debería estar conectado con su instancia administrada de Kubernetes. Eche un vistazo a las [Guías Meshery]({{ site.baseurl }}/guides) para ver más consejos de uso avanzados.

