---
layout: page
title: Docker
permalink: es/installation/platforms/docker
language: es
type: installation
---

# Inicio rápido con Docker

_Nota: se necesita un mínimo de 4 GB de RAM para los despliegues de [Istio (y la aplicación de muestra BookInfo)](/docs/service-meshes/adapters/istio/istio)._

### **Pasos**

Siga estos pasos de instalación para usar Docker y Docker Compose para ejecutar Meshery. Los usuarios a menudo eligen abordar la instalación de esta manera para ejecutar Meshery en su máquina local. Si necesita instalar `docker`, consulte [Introducción a Docker](https://docs.docker.com/get-started/) y si necesita instalar `docker-compose`, consulte [Instalación de Docker Compose](https://docs.docker.com/compose/install/).

El repositorio de Meshery incluye un archivo `docker-compose.yaml`. Podemos usar `docker-compose` para activar todos los servicios de Meshery ejecutando:

 <pre class="codeblock-pre">
 <div class="codeblock"><div class="clipboardjs">
 $ curl -L https://meshery.io/install | bash -
 </div></div>
 </pre>

Una vez que haya verificado que todos los servicios están en funcionamiento, la interfaz de usuario de Meshery será accesible en su máquina local en el puerto 9081. Abra su navegador y acceda a Meshery en [`http://localhost:9081`](http://localhost:9081).
Se le redireccionará a una página de inicio de sesión social donde puede elegir uno de los métodos de inicio de sesión social disponibles para iniciar sesión en Meshery.

Al iniciar Meshery correctamente, las instrucciones para acceder a Meshery se imprimirán en la pantalla. Eche un vistazo a las [guías Meshery]({{ site.baseurl }}/guides) para obtener consejos de uso avanzados.
