---
layout: page
title: Contribuir Flujo
permalink: es/project/contributing-flow
language: es
type: project
---

# Trabajando por Fork

Simplemente dirígete a la página de GitHub y haz clic en el botón "Fork". Es así de simple. Una vez que hayas hecho eso, puedes usar tu cliente git favorito para clonar tu repositorio o simplemente dirigirse directamente a la línea de comando:

## Clona tu fork en tu máquina local

```
git clone git@github.com:USERNAME/FORKED-PROJECT.git
```

## Mantener Tu Fork Actualizado

Si bien este no es un paso necesario, si planeas hacer algo más que una pequeña solución rápida, querrás asegurarte de mantener tu fork actualizado rastreando el repositorio "upstream" original al que hiciste "fork". Para hacer esto, deberás agregar un remoto:

## Agregar repositorio 'upstream' a la lista de remotos

```
git remote add upstream https://github.com/meshery/meshery.git
```

("meshery" se usa como repositorio de ejemplo. Asegúrate de hacer referencia al repositorio _actual_ al que está contribuyendo, por ejemplo, "meshery-linkerd").

## Verifica el nuevo repositorio remoto llamado 'upstream'

```
git remote -v
```

Siempre que quieras actualizar tu fork con los últimos cambios de 'upstream', primero deberás buscar las ramas del repositorio 'upstream' y los últimos "commits" para llevarlos a tu repositorio:

## Obtener desde el remoto 'upstream'

```
git fetch upstream
```

## Ver todas las ramas, incluidas las de 'upstream'

```
git branch -va
```

Ahora, revisa tu rama "master" y combina la rama "master" del repositorio 'upstream':

## Haz "checkout" a tu rama "master" y haz "merge" de 'upstream'

```
git checkout master
git merge upstream/master
```

Si no hay "commits" únicos en la rama "master" local, git simplemente realizará un "fast-forward". Sin embargo, si has realizado cambios en "master" (en la gran mayoría de los casos, probablemente no deberías hacerlo), consulta la siguiente sección, es posible que tengas que lidiar con conflictos. Al hacerlo, ten cuidado de respetar los cambios realizados anteriormente.

Ahora, tu rama "master" local está actualizada con todo lo modificado en 'upstream'.

**Crear una "Branch"** (haciendo tu trabajo)
Siempre que comiences a trabajar en una nueva función o corrección de "bugs", es importante que crees una nueva rama. No solo es un flujo de trabajo de git adecuado, sino que también mantiene tus cambios organizados y separados de la rama "master" para que puedas enviar y administrar fácilmente múltiples "pull requests" para cada tarea que completes.

Para crear una nueva rama y comenzar a trabajar en ella, realiza el siguiente flujo.

## Consulta la rama "master": quieres que tu nueva rama provenga de la rama "master".

```
git checkout master
```

## Crea una nueva rama (dale a tu rama su propio nombre informativo simple)

Para mejoras, utiliza `feature/tu_nombre/issue#` o `feature/tu_nombre/nombre_del_feature`

Para "bugs" usa `bug/tu_nombre/issue#` o `bug/tu_nombre/nombre_del_bug`

```
git branch feature/jdoe/567
```

## Cambia a tu nueva rama

```
git checkout feature/jdoe/567
```

Ahora, ve a la ciudad hackeando y haciendo los cambios que quieras.

## Enviar tus cambios (un "Pull Request")

Antes de enviar tu "pull request", es posible que quieras hacer algunas cosas para limpiar tu rama y hacer que sea lo más simple posible para que el "maintainer" del repositorio original pruebe, acepte y haga "merge" de tu trabajo.

En el tiempo que has estado trabajando en tus cambios, si se han hecho "commits" en la rama "master" 'upstream', deberás hacer "rebase" a tu rama de desarrollo para que al hacerle "merge" sea un "fast-forward" simple que no requiera ninguna trabajo de resolución de conflictos.

## Obten los cambios de "master" 'upstream' y haz "merge" con la rama "master" de tu repositorio

```
git fetch upstream
git checkout master
git merge upstream/master
```

## Si hubo nuevos "commits", haz "rebase" a tu rama de desarrollo

```
git checkout feature/jdoe/567
git rebase master
```

Ahora, puede ser deseable reducir algunos de tus "commits" más pequeños al juntarlos en una pequeña cantidad de "commits" más grandes y cohesivos. Puedes hacer esto con un "rebase" interactivo:

## Haz "rebase" a todos tus "commits" en tu rama de desarrollo

```
git checkout
git rebase -i master
```

Esto abrirá un editor de texto donde puedes especificar que "commits" aplastar.

## Entrega

Una vez que hayas hecho tus "commits" y enviado todos tus cambios a GitHub, ve a la página de tu "fork" en GitHub, selecciona tu rama de desarrollo y haz clic en el botón de "pull request". Si necesitas realizar algún ajuste en tu "pull request", simplemente envía las actualizaciones a GitHub. Tu "pull request" rastreará automáticamente los cambios en tu rama de desarrollo y se actualizará.
