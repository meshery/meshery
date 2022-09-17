---
layout: page
title: Vulnerabilidades de seguridad
permalink: es/project/security-vulnerabilities
description: Cómo el equipo de Meshery maneja las vulnerabilidades de seguridad.
language: es
type: project
---

# Informar una vulnerabilidad

Estamos muy agradecidos con quienes investigan temas de seguridad y quienes informan
sobre las vulnerabilidades de seguridad de Meshery. Investigamos cada informe a fondo.

Para realizar un informe, envía un correo electrónico
a la lista de correo privada [meshery-security-vulns-reports@layer5.io](mailto:meshery-security-vulns-reports@layer5.io) con los detalles de la vulnerabilidad.
Para los bugs normales del producto que no están relacionados con vulnerabilidades de seguridad latentes, diríjete al repositorio correspondiente
y envía un [nuevo issue](https://github.com/meshery/meshery/issues/new/choose).

### ¿Cuándo informar una vulnerabilidad de seguridad?

Envíanos un informe siempre que:

- Creas que Meshery tiene una vulnerabilidad de seguridad potencial.
- No estás seguro de cómo una vulnerabilidad afecta a Meshery.
- Creas que una vulnerabilidad está presente en otro proyecto que Meshery
  depende de (Docker, por ejemplo).

### ¿Cuándo no informar una vulnerabilidad de seguridad?

No envíes un informe de vulnerabilidad si:

- Necesitas ayuda para ajustar los componentes de Meshery para la seguridad.
- Necesitas ayuda para aplicar actualizaciones relacionadas con la seguridad.
- Tu problema no está relacionado con la seguridad.

## Evaluación

El equipo de Meshery reconoce y analiza cada informe de vulnerabilidad dentro de los 10 días hábiles.

Cualquier información de vulnerabilidad que compartas con el equipo de Meshery permanece
dentro del proyecto Meshery. No divulgamos la información a otros
proyectos. Solo compartimos la información necesaria para solucionar el problema.

Mantenemos al informante actualizado a medida que se aborda el estado del problema de seguridad.

## Arreglando el problema

Una vez que se ha distinguido por completo una vulnerabilidad de seguridad, el equipo de Meshery desarrolla una solución.
El desarrollo y la prueba de la solución se realizan en un repositorio privado de GitHub para evitar
divulgación prematura de la vulnerabilidad.

## Divulgación temprana

El proyecto Meshery mantiene una lista de correo para la divulgación temprana privada de vulnerabilidades de seguridad.
La lista se utiliza para proporcionar información procesable para socios cercanos de Meshery. La lista no está destinada
para que las personas se enteren de los problemas de seguridad.

## Divulgación pública

El día elegido para la divulgación pública, se lleva a cabo una secuencia de actividades lo más rápido posible:

- Los cambios se combinan, desde el repositorio privado de GitHub que contiene la corrección, en el conjunto apropiado de ramas públicas.

- El equipo de Meshery se asegura de que todos los binarios necesarios se creen y se publiquen rápidamente.

- Una vez que los binarios están disponibles, se envía un anuncio en los siguientes canales:

  - El [Blog de Meshery](https://meshery.io/blog/)
  - El [Feed del Twitter de Meshery](https://twitter.com/mesheryio)
  - El canal de #announcements en Slack

En la medida de lo posible, este anuncio será procesable e incluirá cualquier medida de mitigación que los clientes puedan tomar antes de
actualizar a una versión fija.
