---
layout: default
title: "Extensibilidad: APIs"
permalink: es/extensibility/api
type: extensibility
abstract: "La arquitectura Meshery es extensible y ofrece una variedad de puntos de extensión y API REST y GraphQL."
#redirect_from: extensibility
language: es
---

## APIs de Meshery

Cada una de las API de Meshery está sujeta al siguiente sistema de autenticación y autorización.

### Autenticación

Las solicitudes a cualquiera de los puntos finales de la API deben estar autenticadas e incluir un token de acceso JWT válido en los encabezados HTTP. El tipo de autenticación está determinado por el [Proveedor](#providers). Use el Provider Local, "None", pone Meshery en modo de usuario-único y no requiere autenticación.

### Autorización

Actualmente, Meshery solo requiere un token válido para permitir que los clientes invoquen sus API.

### Endpoints (puntos finales del API)

Cada uno de los puntos finales de la API se expone a través de [server.go](https://github.com/layer5io/meshery/blob/master/router/server.go). Los puntos finales se agrupan por función (p.ej. /api/mesh or /api/perf).

Alternativamente, [Proveedores Remotos](./providers) pueden extender los puntos finales de Meshery detrás del endpoint `/api/extensions/`.

## GraphQL

Meshery proporciona una API GraphQL disponible a través del puerto predeterminado de `/tcp`.

## REST

Meshery proporciona una API REST disponible a través del puerto predeterminado de `9081/tcp`.
