---
layout: page
title: Contributing to Meshery Docs
permalink: es/project/contributing/contributing-docs
description: How to contribute to Meshery Docs.
language: es
type: project
category: contributing
---

Antes de contribuir, revise el [Documentation Contribution Flow](https://github.com/layer5io/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow). En los siguientes pasos, configurará su entorno de desarrollo, bifurcará y clonará el repositorio, ejecutará el sitio localmente y, finalmente, confirmará, aprobará y enviará los cambios realizados para su revisión.

{% include alert.html type="info" title="Meshery Documentation Design Specification" content="Ver el <a href='https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit#'>Especificación de diseño de documentación de Meshery </a>, que sirve para proporcionar una descripción general de las herramientas y el enfoque utilizado para crear la documentación de Meshery y la arquitectura de información de ti." %}

## Marco de documentación

La documentación de Meshery se compone de estos componentes:

- Marco de referencia - Jekyll
- MarcoTema - https://github.com/vsoch/docsy-jekyll
- Repo - https://github.com/layer5io/meshery/tree/master/docs
- DNS - https://meshery.layer5.io/docs
- AWS API GW - una instancia está configurada para redirigir de docs.meshery.io a meshery.layer5.io, debido a la ubicación del repositorio donde residen actualmente los documentos.

## Configura tu entorno de desarrollo

{% include alert.html type="info" title="Jekyll" content="El sitio de Meshery Docs se construye con Jekyll, un generador de sitios estático simple. Jekyll se puede instalar en diferentes plataformas como Windows, Linux y MacOS siguiendo los siguientes pasos " %}

### Para Windows

**Note:** Los usuarios de Windows pueden ejecutar Jekyll siguiendo el [Windows Installation Guide](https://jekyllrb.com/docs/installation/windows/) y también instalando Ruby Version Manager [RVM](https://rvm.io). RVM es una herramienta de línea de comandos que le permite trabajar con múltiples entornos Ruby en su máquina local. Alternativamente, si está ejecutando Windows 10 versión 1903 Build 18362 o superior, puede actualizar al Subsistema de Windows para Linux [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) y ejecute Jekyll en Linux en su lugar.

- Inicie su VM WSL e instale el administrador de versiones de ruby (RVM):

```bash
  sudo apt update
  sudo apt install curl g++ gnupg gcc autoconf automake bison build-essential libc6-dev \
    	libffi-dev libgdbm-dev libncurses5-dev libsqlite3-dev libtool \
    	libyaml-dev make pkg-config sqlite3 zlib1g-dev libgmp-dev \
    	libreadline-dev libssl-dev
  sudo gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
  curl -sSL https://get.rvm.io | sudo bash -s stable
  sudo usermod -a -G rvm `whoami`
```

Si `gpg --keyserver` da un error, puede usar:

```bash
  sudo gpg --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
```

o

```bash
  sudo gpg2 --keyserver hkp://pool.sks-keyservers.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB
```

Reinicie su VM WSL antes de seguir adelante.

- Para instalar Ruby, ejecute:
  ```bash
    rvm install ruby
    rvm --default use ruby 2.7.1
    gem update
    gem install jekyll bundler
  ```

### Para Linux

- Prerrequisitos
  ```bash
    sudo apt-get update
    sudo apt-get install autoconf bison build-essential libssl-dev libyaml-dev libreadline6-dev zlib1g-dev libncurses5-dev libffi-dev libgdbm3 libgdbm-dev
  ```

#### Instalación de rbenv

- Clonación del repositorio rbenv
  ```bash
    clon de git https://github.com/rbenv/rbenv.git ~/.rbenv
  ```
- Marcando el camino
  ```bash
    echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc
  ```
- rbenv init
  ```bash
    echo 'eval "$(rbenv init -)"' >> ~/.bashrc
  ```
- Recarga tu bashrc
  ```bash
    source ~/.bashrc
  ```
  <strong>Nota:</strong> Cambie bashrc con el archivo rc específico de su shell, por ejemplo: si está usando zsh, el nombre del archivo es zshrc.
- Verificar instalación
  ```bash
    type rbenv
  ```

#### Instalar Ruby

- versión de instalación de rbenv

```bash
  rbenv install 2.5.1
```

- Para enumerar todas las versiones que se pueden instalar

```bash
  rbenv install --list-all
```

- Establezca qué versión de Ruby desea usar

```bash
  rbenv global version
```

- Compruebe la instalación de Ruby

```bash
  ruby -v
```

### Para MacOS

- Utilice los documentos aquí [Instalación de Jekyll](https://jekyllrb.com/docs/installation/macos/)

### Obtener el codigo

- Bifurcar y luego clonar el [repositorio de Meshery](https://github.com/layer5io/meshery)
  ```bash
  $ git clone https://github.com/YOUR-USERNAME/meshery
  ```
- Cambiar al directorio de documentos
  ```bash
  $ cd docs
  ```
- Instale las dependencias de Ruby

  ```bash
  $ gem install bundler
  $ bundle install
  ```

<strong> Nota: </strong> si es un usuario de Mac, no necesita instalar las dependencias de Ruby; después de pasar al directorio de documentos, puede servir el sitio.

### Sirva el sitio

- Sirve el código localmente
  ```bash
  $ make docs
  ```
- Si eso da un error, ejecute:

  ```bash
    $ bundle exec jekyll serve
  ```

\_Nota: Desde Makefile, este comando en realidad está ejecutando `$ bundle exec jekyll serve --drafts --livereload`. Hay dos configuraciones de Jekyll, `jekyll serve` para desarrollo local y` jekyll build` cuando necesita generar los artefactos del sitio para la producción.

### Uso de Docker

Si tiene Docker y `make` instalados en su sistema, entonces puede servir el sitio localmente

```
$ make docker-docs
```

Esto no requiere la necesidad de instalar Jekyll y Ruby en su sistema.

** Pero, debe asegurarse de que GNU make funcione en su sistema (puede que no funcione en Windows) **

#### Nota

Mientras realiza el paso anterior, si enfrenta errores con un mensaje como el siguiente ...

`Your ruby version is x.x.x but your Gemfile specified 2.7.x`

Esto se debe a que Jekyll siempre considera la versión exacta de Ruby a diferencia de JavaScript.

Por lo tanto, debe seguir cualquiera de los tres pasos para resolver este problema;

- Instale la versión de Ruby requerida usando `rvm` o por cualquier medio dado arriba
- Alternativamente, si tiene Docker instalado, escriba `make docker-docs` para ver los cambios
- Si no puede instalar la versión de Ruby requerida, configure manualmente el `Gemfile` como se muestra a continuación (¡no se recomienda! Hágalo solo si fallan los dos pasos anteriores):

```
source "https://rubygems.org"
ruby '2.7.1' //to any version you have installed
```

Automáticamente el `Gemfile.lock` se actualizará una vez que se proporcione el` make docs` (para Windows, ejecute `bundle exec jekyll serve` si WSL2 no está presente)

** ADVERTENCIA: Si ha seguido el tercer paso, no confirme los cambios realizados en `Gemfile` y` Gemfile.lock` en su rama para preservar la integridad; de lo contrario, la acción de CI no generará la vista previa del sitio durante las relaciones públicas. **.

### Crear una solicitud de extracción

- Después de realizar cambios, no olvide comprometerse con el sign-off bandera (-s)!
  ```bash
  $ commit -s -m “my commit message w/signoff”
  ```
- Una vez que se hayan confirmado todos los cambios, insértelos.
  ```bash
  $ git push origin <branch-name>
  ```
- Luego, en Github, navega hasta el [repositorio de Meshery](https://github.com/layer5io/meshery) y crea una solicitud de extracción a partir de los cambios introducidos recientemente!

---

- _Ver el[Meshery Documentation Google Doc](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit) para referencia adicional._

## Uso de las funciones de Meshery Docs

### Función de portapapeles

Los complementos de portapapeles más populares, como Clipboard JS, requieren la creación manual de una nueva ID para cada fragmento de código. Aquí se utiliza un enfoque diferente. Para los fragmentos de código, usamos etiquetas html o rebajas de la siguiente manera:

```
   <pre class="codeblock-pre"><div class="codeblock">
   <code class="clipboardjs">
     code snippet
   </code></div></pre>
```

**<pre></pre>** _tags son opcionales a menos que el fragmento de código esté en un formato de párrafo y también le dé un efecto de terminal al código_

**Un bloque completo:**

````
```code snippet```
````

Formateo en línea:

\`code snippet\`: `code snippet`

**Específico del idioma:**

````
```(language name)
  code snippet
```
````

Whenever the code tags are detected, the clipboard javascript file is automatically loaded. Each code element is given a custom id and a clipboard-copy icon to copy the content.

## Resumen del flujo de contribución a la documentación

El siguiente es un resumen conciso de los pasos para contribuir a la documentación de Meshery.

1. Cree una bifurcación, si aún no lo ha hecho, siguiendo los pasos descritos [aquí] (CONTRIBUTING-gitflow.md)
1. En la copia local de su bifurcación, navegue hasta la carpeta de documentos.
   `cd docs`
1. Cree y verifique una nueva rama para realizar cambios dentro de
   `git checkout -b <mis-cambios>`
1. Edite / agregue documentación.
   `vi <página específica> .md`
1. Ejecute el sitio localmente para obtener una vista previa de los cambios.
   `hacer sitio`
1. Confirme, [sign-off] (# commit-signing) y envíe los cambios a su sucursal remota.
   `git push origin <mis-cambios>`
1. Abra una solicitud de extracción (en su navegador web) contra el repositorio: https://github.com/layer5io/meshery.

### Tabla de contenido en la barra lateral (toc)

Las barras laterales usan toc para crear una tabla de contenido. Está escrito de la siguiente manera:

```
    toc:
  - title: Group 1
    subfolderitems:
      - page: Thing 1
        url: /thing1.html
      - page: Thing 2
        url: /thing2.html
      - page: Thing 3
        url: /thing3.html
```

La salida del fragmento de código sería:

```
    Group 1
      Thing 1
      Thing 2
      Thing 3
```

### `if` condicional

Esto ejecuta el bloque de código solo si la condición dada es verdadera. Se ejecuta de la siguiente manera:

```
    {{ "{% if product.title == 'Awesome Shoes' " }}%}
    These shoes are awesome!
    {{ "{% endif " }}%}
```

Si la condición es verdadera, la salida sería:

```
    How are you?
```

### `for` loop

La instrucción for ejecuta un bloque de código repetidamente. Está escrito de la siguiente manera:

```
    {{ "{% for names in collection.names " }}%}
    {{ "{{ name.title "}}}}
    {{ "{% endfor " }}%}
```

El resultado producido por el fragmento de código anterior:

```
    Sam Ham Ethan
```

### Comentario

Los comentarios permiten dejar un bloque de código desatendido, no se ejecutará ninguna declaración entre el comentario de apertura y el de cierre.

### Incluir

La etiqueta anterior se utiliza para insertar un archivo ya renderizado dentro de la plantilla actual. Está escrito de la siguiente manera:

```
    {{ "{% include file.html " }}%}
```

### Asignar

La etiqueta de asignación se utiliza para crear una nueva variable. Está escrito de la siguiente manera:

```
    {{ "{% assign variable1 = true " }}%}
```

{% include suggested-reading.html diffName ="false" language = "es" %}
