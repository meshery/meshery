---
layout: default
title: Implementación de aplicaciones de muestra
description: Esta guía tiene como objetivo ayudar a los usuarios a comprender mejor las aplicaciones de muestra.
permalink: es/guides/sample-apps
language: es
type: Guides
---

Las aplicaciones de muestra se utilizan para interactuar y ejemplificar las características de su malla de servicios. A menudo son una colección de microservicios que el usuario puede usar como un área de juegos de prueba para experimentar y aprender sobre la malla de servicios y su exhaustivo conjunto de características.
Antes de implementar una aplicación de muestra sobre su malla de servicios, la aplicación debe estar expuesta y permitir el acceso externo a los servicios disponibles en un clúster. Hay una gran variedad de formas de hacer esto, específicas para la malla de servicios que está utilizando.

Una forma popular de exponer su clúster es usar [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/), un objeto API que define reglas que permiten el acceso externo a los servicios en un clúster.

- [Configure Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/#the-ingress-resource)
- [Configure Ingress en Minikube](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/)

## Despliegue una aplicación de muestra en Meshery

1. Vaya a la página de administración de cualquier [malla de servicio]({{ site.baseurl }}/service-meshes) e instale cualquiera de sus versiones estables.
   <a href="{{ site.baseurl }}/assets/img/adapters/consul/consul-install.png"><img alt="Consul-install" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/consul/consul-install.png" /></a>

2. Pulse (+) en **Manage Sample Application Lifecycle**. Ahora podrá ver un menú desplegable con las aplicaciones de muestra disponibles.
   <a href="{{ site.baseurl }}/assets/img/adapters/consul/consul-sample-app.png"><img alt="ImageHub sample app" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/consul/consul-sample-app.png" /></a>

3. Haga clic en la aplicación de muestra que desea implementar. Esto puede tardar hasta un minuto. Se le notificará cuando se haya implementado la aplicación de muestra
   <a href="{{ site.baseurl }}/assets/img/adapters/consul/consul-imagehub-success.png"><img alt="ImageHub deployed" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/consul/consul-imagehub-success.png" /></a>

### [BookInfo](https://github.com/layer5io/istio-service-mesh-workshop/blob/master/lab-2/README.md#what-is-the-bookinfo-application)

Originalmente creado por Istio, BookInfo es una aplicación de muestra que, al implementarse, muestra información sobre un libro, similar a una entrada de catálogo única de una librería en línea. En la página se muestra una descripción del libro, los detalles del libro (ISBN, número de páginas, etc.) y algunas reseñas de libros. La aplicación consta de cuatro microservicios:

- **productpage**: El microservicio productpage llama los detalles y revisa los microservicios para popular la página.
- **details**: El microservicio details contiene información de libros.
- **reviews**: El microservicio reviews contiene revisiones de libros. También llama las el microservicio de calificaciones (ratings).
- **ratings**: El microservicio ratings contiene información del ranqueo de los libros que acompaña la revisión de los libros.

OncUna vez que BookInfo es deplegado, usted puede usar Meshery para aplicar configuraciones personalizadas al control de tráfico, inyectar latencia, llevar a cabo routeo basado en contexto, etcétera.

<a href="{{ site.baseurl }}/assets/img/adapters/sample-apps/bookinfo.png"><img alt="BookInfo sample app" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/sample-apps/bookinfo.png" /></a>

### [Emojivoto](https://github.com/BuoyantIO/emojivoto)

Emojivoto es una aplicación de microservicios, originalmente construida por Linkerd que permite a los usuarios votar por su emoji favorito, y monitorea los votos recividos en el tablero de clasificación. La aplicación se compone de tres microservicios:

- **emojivoto-web**: Web frontend y REST API
- **emojivoto-emoji-svc**: API gRPC para la búsqueda y listado de emoji
- **emojivoto-voting-svc**: API gRPC para el voto y el tablero de clasificación

<a href="{{ site.baseurl }}/assets/img/adapters/sample-apps/emojivoto.png"><img alt="Emojivoto" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/sample-apps/emojivoto.png" /></a>

### [ImageHub](https://layer5.io/projects/image-hub)

Image Hub es una aplicación de muestra para explorar los módulos de WebAssembly que se utilizan como filtros de Envoy. La aplicación se escribió originalmente para ejecutarse en Consul. Sin embargo, no depende de Consul y se puede implementar en cualquier malla de servicios. Estos módulos se pueden utilizar para implementar la tenencia múltiple o para implementar la limitación de la tasa de usuario en los puntos finales de su aplicación, sin interferir con la infraestructura de su aplicación.
Siga este tutorial para configurar [ImageHub con Ingress](https://github.com/layer5io/image-hub#use-image-hub)

<a href="{{ site.baseurl }}/assets/img/adapters/sample-apps/imagehub-on-consul.png"><img alt="Imagehub-on-Consul" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/sample-apps/imagehub-on-consul.png" /></a>

### [HTTPBin](https://httpbin.org)

HttpBin es un servicio de solicitud y respuesta HTTP simple que responde a muchos tipos de solicitudes http / https, incluidos los métodos de solicitud http estándar (o verbos) utilizados por REST.

<a href="{{ site.baseurl }}/assets/img/adapters/sample-apps/httpbin.png"><img alt="httpbin" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/sample-apps/httpbin.png" /></a>

### [Linkerd Books](https://github.com/BuoyantIO/booksapp)

Linkerd Books es una aplicación de muestra basada en Ruby. Está diseñado para demostrar las diversas propuestas de valor, incluida la depuración, la observabilidad y el monitoreo de su malla de servicios. Puede usarse para determinar la eficiencia de su malla y para depurar.

<a href="{{ site.baseurl }}/assets/img/adapters/sample-apps/linkerd-books.png"><img alt="Linkerd Books" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/sample-apps/linkerd-books.png" /></a>

### [Online Boutique](https://github.com/GoogleCloudPlatform/microservices-demo)

Online Boutique es una aplicación nativa de la nube de muestra, originalmente creada por Google. Se compone de 10 microservicios y se puede utilizar para mostrar y trabajar con Kubernetes, Istio, gRPC y OpenCensus. En la implementación, ejecuta una aplicación de demostración de comercio electrónico basada en la web, un ejemplo de la cual se puede ver a continuación:

<a href="{{ site.baseurl }}/assets/img/adapters/sample-apps/online-boutique.png"><img alt="Online Boutique" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/adapters/sample-apps/online-boutique.png" /></a>

<!--Sample apps specific to NSM can be found on {{ site.baseurl }}/service-meshesadapters/nsm/nsm-->
