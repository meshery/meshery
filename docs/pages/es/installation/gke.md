---
layout: page
title: GKE
permalink: es/installation/platforms/gke
language: es
type: installation
---

# Inicio rápido con Google Kubernetes Engine (GKE)

Para otorgar a Meshery el acceso necesario a su instancia de Kubernetes administrada,
Meshery deberá tener asignada una `ServiceAccount` ( Cuenta de servicio ). Se puede utilizar una `ServiceAccount` existente o crear una nueva. Asegúrese de que la `ServiceAccount` que utilize tenga asignada el rol de `cluster-admin`.

Meshery usará esta `ServiceAccount` para interactuar con su instancia administrada de Kubernetes. Utilice cualquiera de los dos métodos siguientes para preparar un archivo kubeconfig compatible:

- [Configuración automática](#configuración-automática-recomendada)
- [Configuración manual](#configuración-manual-opcional)

#### **Configuración Automática** (Recomendada)

1. En su navegador, navegue hasta Meshery (por ejemplo, `http://localhost:9081`) e inicie sesión.
1. Descargue su token de autenticación de Meshery haciendo clic en **Get Token** debajo de su perfil de usuario.
1. Utilice este token de autenticación para ejecutar el siguiente comando:
   ```
   $ mesheryctl system config gke --token <RUTA AL ARCHIVO DE TOKEN>
   ```

Este comando actualiza su kubeconfig para proporcionar a Meshery acceso a su instancia administrada de Kubernetes.
Una vez configurado, proceda con Meshery (`mesheryctl system start`).

#### **Configuración Manual** (Opcional)

Si el procedimiento [Configuración automática](#configuración-automática-recomendada) falla ó si desea preparar manualmente su archivo kubeconfig para proporcionar a Meshery el acceso necesario a su instancia administrada de Kubernetes, realice las siguientes acciones:

1. Descargue el script de shell [generate_kubeconfig_gke.sh](./generate_kubeconfig_gke.sh).
1. Ejecute este script de shell que identifica el nombre de la cuenta de servicio y los argumentos del espacio de nombres, de esta manera:

   ```sh
   ./generate_kubeconfig_gke.sh cluster-admin-sa-gke default
   ```

1. Una vez que el script haya terminado por completo, puede proceder a iniciar Meshery con la configuración compatible con GKE ejecutando:

   ```sh
   $ mesheryctl system start
   ```

1. En su navegador, navegue hasta Meshery (por ejemplo,, `http://localhost:9081`) e inicie su sesión.
1. En la sección Settings-->Environment, introduzca el archivo generado (`config-cluster-admin-sa-gke-default.yaml`) como el archivo kubeconfig.

Meshery ahora debería estar conectado con su instancia administrada de Kubernetes. Entre a visualizar las [guías Meshery]({{ site.baseurl }}/guides) para que obtenga consejos de uso avanzado.
