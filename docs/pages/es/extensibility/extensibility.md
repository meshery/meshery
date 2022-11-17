---
layout: default
title: Extensibilidad
permalink: es/extensibility
type: extensibility
abstract: 'La arquitectura Meshery es extensible. Meshery proporciona varios puntos de extensión para trabajar con diferentes mallas de servicio a través de <a href="extensibility#adapters">adaptadores</a>, <a href="extensibility#load-generators">generadores de carga</a> and <a href="extensibility#providers">providers</a>.'
redirect_from: es/reference/extensibility
language: es
---

Meshery tiene una arquitectura extensible con varios puntos de extensión. Meshery proporciona varios puntos de extensión para trabajar con diferentes mallas de servicio a través de [adaptadores]({{site.baseurl}}/es/extensibility/adapters), diferentes [generadores de carga]({{site.baseurl}}/es/extensibility/load-generators) y diferentes [proveedores]({{site.baseurl}}/es/extensibility/providers). Meshery también ofrece un REST API.

**Principios rectores de la extensibilidad**

Los siguientes principios se mantienen en el diseño de la extensibilidad de Meshery.

1. Reconozca que los diferentes entornos de implementación tienen diferentes sistemas para integrarse.
1. Ofrezca una experiencia predeterminada que proporcione la experiencia de usuario óptima.

## Puntos de Extensión

Meshery no es solo una aplicación. Es un conjunto de microservicios donde el componente central se llama Meshery. Los integradores pueden extender Meshery aprovechando los puntos de extensión designados. Los puntos de extensión vienen en varias formas y están disponibles a través de la arquitectura de Meshery.

![Puntos de Extensión de Meshery ]({{site.baseurl}}/assets/img/architecture/meshery_extension_points.svg)

_Figura: Puntos de extensión disponibles en Meshery_

Los siguientes puntos de extensión están incorporados actualmente a Meshery.

## Tipos de Puntos de Extensión

1. [Proveedores]({{site.baseurl}}/es/extensibility/providers)
1. [Generadores de Carga]({{site.baseurl}}/es/extensibility/load-generators)
1. [Adaptadores]({{site.baseurl}}/es/extensibility/adapters)
1. [API REST]({{site.baseurl}}/es/extensibility/api#rest)
1. [API GraphQL]({{site.baseurl}}/es/extensibility/api#graphql)
