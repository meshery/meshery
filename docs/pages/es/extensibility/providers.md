---
layout: default
title: "Extensibilidad: Proveedores"
permalink: es/extensibility/providers
type: extensibility
#redirect_from: architecture/adapters
abstract: "Meshery ofrece soporte para más adaptadores que cualquier otro proyecto o producto en el mundo. Meshery utiliza adaptadores para gestionar las distintas mallas de servicio."
language: es
list: include
---

Meshery ofrece a los proveedores como un punto de extensibilidad. Con un proveedor local integrado (llamado "None"), los proveedores remotos de Meshery están diseñados para ser conectables. Los proveedores remotos ofrecen puntos de extensión a los usuarios / integradores para ofrecer una funcionalidad mejorada, utilizando Meshery como plataforma.

1. **Los puntos de extensibilidad ofrecen una separación clara de las capacidades de código abierto y cerrado.**
   - Meshmap es un ejemplo de una función que se entregará a través de un proveedor remoto.
1. **Los Proveedores Remotos deben poder ofrecer RBAC personalizado, componentes de Interfaz de Usuario (UI) personalizados y componentes de backend personalizados**
   - Es necesario identificar o crear marcos cargables dinámicamente para cumplir con cada uno de estos propósitos.

### Principios de diseño: marco de proveedor remoto de Meshery

El marco de extensibilidad del proveedor remoto de Meshery está diseñado para permitir:

1. **Funcionalidad de interfaz de usuario conectable:**

   - Componentes de interfaz de usuario personalizados fuera del árbol con una experiencia de usuario perfecta.
   - Un sistema de recuperación remota de paquetes de extensión (componentes ReactJS y binarios Golang).

1. **Funcionalidad backend conectable:**

   - Los proveedores remotos tienen una cantidad de capacidades desconocidas para Meshery.

1. **AuthZ Conectable**
   - Diseñe un sistema de control de acceso extensible basado en roles de modo que los proveedores remotos puedan determinar su propio conjunto de controles. Proveedores remotos para devolver JWT con roles personalizados, claves de permisos y llaveros de permisos.

![Proveedoress](/assets/img/providers/provider_screenshot.png)

### ¿Qué funcionalidad realizan los proveedores?

Lo que ofrece un proveedor remoto determinado puede variar ampliamente entre proveedores. Meshery ofrece puntos de extensión que los proveedores remotos pueden usar para inyectar diferentes funcionalidades, una funcionalidad específica para ese proveedor.

- **Autenticación and Autorización**
  - Ejemplos: gestión de sesiones, autenticación de dos factores, integración LDAP.
- **Persistencia a Largo Plazo**
  - Ejemplos: almacenamiento y recuperación de resultados de pruebas de rendimiento.
  - Ejemplos: almacenamiento y recuperación de las preferencias del usuario.
- **Visualización Mejorada**
  - Examples: Creation of a visual service mesh topology.
  - Examples: Different charts (metrics), debug (log viewer), distributed trace explorers.
- **Reporteo**
  - Ejemplos: usar el servidor GraphQL de Meshery para crear nuevos paneles.

## Typos of proveedores

En Meshery se definen dos tipos de proveedores: `local` y `remoto`. El proveedor local está integrado en Meshery. Los proveedores remotos pueden ser implementados por cualquier persona u organización que desee integrarse con Meshery. Cualquier número de proveedores remotos puede estar disponible en su implementación de Meshery.

### Proveedores Remotos

El uso de un proveedor remoto pone a Meshery en modo multiusuario y requiere autenticación de usuario. Utilice un proveedor remoto cuando su uso de Meshery sea continuo o se utilice en un entorno de equipo (utilizado por varias personas).

Nombre: **“Meshery”** (default)

- Refuerza la autenticación del usuario.
- Persistencia a largo plazo de los resultados de las pruebas.
- Guarde la configuración del entorno.
- Recupere los resultados de las pruebas de rendimiento.
- Recuperar los resultados de las pruebas de conformidad.
- De uso gratuito.

### Proveedor Local

El uso del proveedor local, "None", pone a Meshery en modo de usuario único y no requiere autenticación. Utilice el proveedor local cuando su uso de Meshery esté destinado a ser de corta duración.

Nombre: **“None”**

- Sin autenticación de usuario.
- Almacenamiento en contenedor de los resultados de las pruebas. Efímero.
- Configuración del entorno no guardada.
- Sin historial de resultados de pruebas de rendimiento.
- Sin historial de resultados de pruebas de conformidad.
- De uso gratuito.

## Construyendo un Proveedor

Meshery interactúa con los proveedores a través de una interfaz Go. Las implementaciones del proveedor deben colocarse en el código y compilarse juntas hoy. Deberá inyectarse una instancia de proveedor en Meshery cuando se inicie el programa.

Meshery mantiene la implementación de los proveedores remotos por separado para que se introduzcan a través de un proceso separado y se inyecten en Meshery en tiempo de ejecución (OR) y cambien la forma en que funciona el código para que los proveedores invoquen a Meshery.

### Puntos de extensión de proveedor remoto

Entretejidos en la interfaz de usuario basada en web de Meshery hay una variedad de puntos de extensión. Cada punto de extensión está cuidadosamente diseñado para ofrecer una experiencia de usuario perfecta. Cada punto de extensión se identifica con un nombre y un tipo. Los siguientes puntos de extensión de la interfaz de usuario de Meshery están disponibles:

- **Nombre:** navigator
  **Tipo:** Elementos de Menú
  **Descripción:** Se supone que es una extensión de página completa que obtendrá un punto final dedicado en la interfaz de usuario de meshery. Y aparecerá en la barra lateral / navegador de la interfaz de usuario de Meshery. Los elementos del menú pueden referirse a extensiones de página completa.

**Name:** user_prefs
**Type:** Componente Único
**Description:** Se supone que son componentes de reacción remota que se colocarán en una página preexistente y no tendrán un punto final dedicado. A partir de ahora, el único lugar donde se puede cargar esta extensión es la sección "Preferencias del usuario" en la configuración de malla.

**Name:** /extension/<su nombre aquí>
**Type:** Página Completa
Description:

El paquete Provider se descomprime en el sistema de archivos del servidor Meshery en `/app/provider-pkg/<package-name>`.

Los proveedores remotos deben cumplir con los siguientes puntos finales:

1. `/login` - regresa token válido
1. `/logout` - invalida el token
1. `/capabilities` - regresa capabilities.json

## Puntos de extensión de la interfaz de usuario

Todas las extensiones de la interfaz de usuario se alojarán en el punto final. <servidormeshery:puerto/proveedor>

### UserPrefs

El punto de extensión UserPrefs espera y carga un componente para que se muestre en la página /userpreferences.

### Navigator

El punto de extensión Navigator carga un conjunto de elementos de menú que se mostrarán en la barra de menú en el lado izquierdo de la interfaz de usuario de Meshery.

## Ejemplo de Punto Final de Capacidades

Meshery Seerver enviará todas las solicitudes a los puntos finales del proveedor remoto. Los endpoints (puntos finales) se determinan e identifican dinámicamente en la sección "capacidades" del punto final `/capabilities`. Los proveedores como objeto tienen los siguientes atributos (esto debe devolverse como una respuesta al endpoint `/capabilities`):

```json
{
  "provider_type": "remote",
  "package_version": "v0.1.0",
  "package_url": "https://layer5labs.github.io/meshery-extensions-packages/provider.tar.gz",
  "provider_name": "Meshery",
  "provider_description": [
    "Persistent sessions",
    "Save environment setup",
    "Retrieve performance test results",
    "Free use"
  ],
  "extensions": {
    "navigator": [
      {
        "title": "MeshMap",
        "href": {
          "uri": "/meshmap",
          "external": false
        },
        "component": "provider/navigator/meshmap/index.js",
        "icon": "provider/navigator/img/meshmap-icon.svg",
        "link:": true,
        "show": true,
        "children": [
          {
            "title": "View: Single Mesh",
            "href": {
              "uri": "/meshmap/mesh/all",
              "external": false
            },
            "component": "navigator/meshmap/index.js",
            "icon": "navigator/img/singlemesh-icon.svg",
            "link": false,
            "show": true
          }
        ]
      }
    ],
    "user_prefs": [
      {
        "component": "userprefs/meshmap-preferences.js"
      }
    ]
  },
  "capabilities": [
    { "feature": "sync-prefs", "endpoint": "/user/preferences" },
    { "feature": "persist-results", "endpoint": "/results" },
    { "feature": "persist-result", "endpoint": "/result" },
    { "feature": "persist-smi-results", "endpoint": "/smi/results" },
    { "feature": "persist-metrics", "endpoint": "/result/metrics" },
    { "feature": "persist-smp-test-profile", "endpoint": "/user/test-config" }
  ]
}
```

Meshery le permite, como propietario de la malla de servicios, personalizar la implementación de su malla de servicios.

## Administrar su Código de Extensión de Proveedor Remoto

Las extensiones de proveedor remoto se mantienen fuera del árbol de Meshery (servidor e interfaz de usuario). Es posible que deba crear sus extensiones en el mismo entorno y conjunto de dependencias que Meshery. El marco de extensibilidad de Meshery se ha diseñado de manera que las extensiones en el árbol se pueden evitar de forma segura y, al mismo tiempo, proporcionar una plataforma sólida desde la que extender la funcionalidad de Meshery. A menudo, aquí se encuentra la delimitación de la funcionalidad abierta frente a la cerrada dentro de Meshery. Los proveedores remotos pueden traer (complemento) qué funcionalidad desean detrás de esta interfaz extensible (más sobre la extensibilidad de Meshery), al menos hasta el punto en que Meshery ha proporcionado una forma de conectar esa característica.

Ofrecer soporte fuera del árbol para extensiones de Meshery significa que:

1. No es necesario que el código fuente de las extensiones de Meshery sea de código abierto,
1. La responsabilidad por la estabilidad de Meshery se reduce significativamente, evitando errores potenciales en componentes extendidos.

A través de puntos de extensión claramente definidos, las extensiones de Meshery pueden ofrecerse como capacidades de código cerrado que se conectan al código de Meshery de código abierto. Para facilitar la integración de sus extensiones de Meshery, puede automatizar la construcción y liberación de sus repositorios de código separados pero interdependientes. Usted será responsable de mantener sus extensiones basadas en ReactJS y Golang.
