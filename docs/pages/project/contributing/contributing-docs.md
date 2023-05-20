---
layout: page
title: Contributing to Meshery Docs
permalink: project/contributing/contributing-docs
description: How to contribute to Meshery Docs.
language: en
type: project
category: contributing
---

Before contributing, please review the [Documentation Contribution Flow](https://github.com/layer5io/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow). In the following steps you will set up your development environment, fork and clone the repository, run the site locally, and finally commit, sign-off, and push any changes made for review.

{% include alert.html type="info" title="Meshery Documentation Design Specification" content="See the <a href='https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit#'>Meshery Documentation Design Specification</a> which serves to provide an overview of the tooling and approach used to create Meshery’s documentation and it information architecture." %}

## Documentation Framework

Meshery documentation is made of these components:

- Framework - [Jekyll](https://jekyllrb.com)
- Theme - [https://github.com/vsoch/docsy-jekyll](https://github.com/vsoch/docsy-jekyll)
- Repo - [https://github.com/layer5io/meshery/tree/master/docs](https://github.com/layer5io/meshery/tree/master/docs)
- Site - [https://docs.meshery.io](https://docs.meshery.io)
- AWS API GW - an instance is configured to redirect from docs.meshery.io to meshery.layer5.io, because of the repo location of where the docs currently reside.

## Set up your development environment

{% include alert.html type="info" title="Jekyll" content="The Meshery Docs site is built using Jekyll - a simple static site generator. Jekyll can be installed on different platforms like Windows, Linux, and MacOS by the following steps " %}

### For Windows

**Note:** Windows users can run Jekyll by following the [Windows Installation Guide](https://jekyllrb.com/docs/installation/windows/) and also installing Ruby Version Manager [RVM](https://rvm.io). RVM is a command-line tool which allows you to work with multiple Ruby environments on your local machine. Alternatively, if you're running Windows 10 version 1903 Build 18362 or higher, you can upgrade to Windows Subsystem for Linux [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) and run Jekyll in Linux instead.

- Fire up your WSL VM and install the ruby version manager (RVM):

{% capture code_content %}sudo apt update
sudo apt install curl g++ gnupg gcc autoconf automake bison build-essential libc6-dev \
 	libffi-dev libgdbm-dev libncurses5-dev libsqlite3-dev libtool \
 	libyaml-dev make pkg-config sqlite3 zlib1g-dev libgmp-dev \
 	libreadline-dev libssl-dev
sudo gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
curl -sSL https://get.rvm.io | sudo bash -s stable
sudo usermod -a -G rvm `whoami`{% endcapture %}
{% include code.html code=code_content %}

If `gpg --keyserver` gives an error, you can use:

{% include code.html code="sudo gpg --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3" %}

or

{% include code.html code="sudo gpg2 --keyserver hkp://pool.sks-keyservers.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB" %}

Restart your WSL VM before moving forward.

- For installing Ruby, run:

  ```bash
    rvm install ruby
    rvm --default use ruby 2.7.5
    gem update
    gem install jekyll bundler
  ```
  
- Update the Gemfile located in meshery/docs directory
  ```
    ruby '3.0.2'
  ```
  <strong>Note:</strong> In place of `3.0.2` add your installed version
  
- Also add this to the next line in the Gemfile
  ```
    gem 'wdm','>=0.1.0' if Gem.win_platform?
  ```
  <strong>Note:</strong> This is just a workaround for your local machine so do not commit or push the modified Gemfile or Gemfile.lock during Pull Requests
  
- Next, go to this folder on your device
  ```C:\Ruby24-x64\lib\ruby\gems\2.4.0\gems\eventmachine-1.2.5-x64-mingw32\lib```
  
- Add ```require 'em/pure_ruby' ``` in the first line of the <strong>eventmachine.rb</strong> file

### For Linux

- Prerequisites
{% capture code_content %}sudo apt-get update
sudo apt-get install autoconf bison build-essential libssl-dev libyaml-dev libreadline6-dev zlib1g-dev libncurses5-dev libffi-dev libgdbm3 libgdbm-dev{% endcapture %}
{% include code.html code=code_content %}
#### Installing rbenv

- Cloning the rbenv repository
{% capture code_content %}git clone https://github.com/rbenv/rbenv.git ~/.rbenv{% endcapture %}
{% include code.html code=code_content %}
  <strong>Note:</strong> Change bashrc with your shell specific rc file, for eg: if you are using zsh then the filename is zshrc.
- Setting the path
{% capture code_content %}echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc{% endcapture %}
{% include code.html code=code_content %}
- rbenv init
{% capture code_content %}echo 'eval "$(rbenv init -)"' >> ~/.bashrc{% endcapture %}
{% include code.html code=code_content %}
- Reload your bashrc
{% capture code_content %}source ~/.bashrc{% endcapture %}
{% include code.html code=code_content %}
- Check installation
{% capture code_content %}type rbenv{% endcapture %}
{% include code.html code=code_content %}

#### Install Ruby

- rbenv install version

{% capture code_content %}rbenv install 2.7.5{% endcapture %}
{% include code.html code=code_content %}

- To list all the versions that can be installed

{% capture code_content %}rbenv install --list{% endcapture %}
{% include code.html code=code_content %}

- Set which Ruby version you want to use

{% capture code_content %}rbenv global version{% endcapture %}
{% include code.html code=code_content %}

- Check Ruby installation

{% capture code_content %}ruby -v{% endcapture %}
{% include code.html code=code_content %}

### For MacOS

- Use docs here [Jekyll installation](https://jekyllrb.com/docs/installation/macos/)

### Get the code

- Fork and then clone the [Meshery repository](https://github.com/meshery/meshery)
{% capture code_content %}git clone https://github.com/YOUR-USERNAME/meshery{% endcapture %}
{% include code.html code=code_content %}
- Change to the docs directory
{% capture code_content %}cd docs{% endcapture %}
{% include code.html code=code_content %}
- Install any Ruby dependencies
{% capture code_content %}gem install bundler
bundle install{% endcapture %}
{% include code.html code=code_content %}

  <strong>Note:</strong> If you are a Mac user you do not need to install the Ruby dependencies, after moving on to the docs directory, you can serve the site.

### Serve the site

- Serve the code locally
{% capture code_content %}make docs{% endcapture %}
{% include code.html code=code_content %}
- If that gives an error run:
{% capture code_content %}bundle exec jekyll serve --drafts --config _config_dev.yml{% endcapture %}
{% include code.html code=code_content %}

  _From the Makefile, this command is actually running `$ bundle exec jekyll serve --drafts --livereload --config _config_dev.yml`. If this command causes errors try running the server without Livereload with this command: `$ bundle exec jekyll serve --drafts --config _config_dev.yml`. Just keep in mind you will have to manually restart the server to reflect any changes made without Livereload. There are two Jekyll configuration, `jekyll serve` for developing locally and `jekyll build` when you need to generate the site artefacts for production._

### Using Docker

If you've Docker and `make` installed in your system, then you can serve the site locally

{% capture code_content %}make docker{% endcapture %}
{% include code.html code=code_content %}

This doesn't require the need for installing Jekyll and Ruby in your system

**But, you need to make sure that GNU make is working in your system (might not work in Windows)**

#### Note

While performing the above step, if you're facing errors with a message like below...

`Your ruby version is x.x.x but your Gemfile specified 2.7.x`

This is because Jekyll always considers the exact version of Ruby unlike JavaScript.

So, you need to follow either of the three steps to resolve this problem;

- Install the required Ruby version by using `rvm` or by any means given above
- Alternatively, if you have Docker installed, then type `make docker-docs` to view the changes
- If you're unable to install the required Ruby version, then manually configure the `Gemfile` as below (not recommended! Do only if above two steps fail):

```
source "https://rubygems.org"
ruby '2.7.5' //to any version you have installed
```

Automatically the `Gemfile.lock` will update once the `make docs` is given (for Windows, run `bundle exec jekyll serve` if WSL2 isn't present)

**WARNING: If you have followed the third step then please don't commit the changes made on `Gemfile` and `Gemfile.lock` in your branch to preserve integrity, else the CI action will fail to generate the site preview during PR**.


### Make Necessary Changes
- Make changes as required by the issue you are solving
- Be sure check that your changes appear correctly locally by serving the site using `make docs`

#### Note
- If the issue requires making new doc page that replaces the old page, please don't forget to add a redirect link on the old page. This redirect link field should have the link of the new page created.


### Create a Pull Request

- After making changes, don't forget to commit with the sign-off flag (-s)!
{% capture code_content %}git commit -s -m “my commit message w/signoff”{% endcapture %}
{% include code.html code=code_content %}
- Once all changes have been committed, push the changes.
{% capture code_content %}git push origin <branch-name>{% endcapture %}
{% include code.html code=code_content %}
- Then on Github, navigate to the [Meshery repository](https://github.com/layer5io/meshery) and create a pull request from your recently pushed changes!

---

- _See the [Meshery Documentation Google Doc](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit) for additional reference._

## Using the features of Meshery Docs

### Clipboard Feature

Most popular clipboard plugins like Clipboard JS require the manual creation of a new ID for each code snippet. A different approach is used here. For code snippets, we either use html tags or markdown in the following manner:

{% capture code_content %}<pre class="codeblock-pre"><div class="codeblock">
<code class="clipboardjs">code_snippet_here</code>
</div></pre>{% endcapture %}
{% include code.html code=code_content %}

You can also use the `code` component created for this feature to make it easy to use. It can be used by including code.html and then passing the code snippet to it.

If the code snippet does not cotain any special characters then, it can be used directly like,<br><br>
`{% raw %}{% include code.html code="code_snippet_here" %}{% endraw %}`.

If the code snippet has special characters then we need to capture it first in variable `code_content` like this,<br><br> `{% raw %}{% capture code_content%}code_snippet_here{% endcapture %}{% endraw %}`<br><br> and then pass it to the component like this,<br><br> `{% raw %}{% include code.html code=code_content %}{% endraw %}`.<br><br> That's it.

**NOTE:** Don't use `code` component when the snippet is in between an ordered list, it breaks the order and next item in the list will have numbering from 1. So better use `<pre class="codeblock_pre">...</pre>` method above instead when the snippet is in between an ordered list.

**A full block:**
{% capture code_content %}```code snippet```{% endcapture %}
{% include code.html code=code_content %}

**Inline formatting:**

{% capture code_content %}\`code snippet\`: `code snippet`{% endcapture %}
{% include code.html code=code_content %}

**Language specific:**

{% capture code_content %}```(language name)
code snippet
```{% endcapture %}
{% include code.html code=code_content %}

Whenever the code tags are detected, the clipboard javascript file is automatically loaded. Each code element is given a custom id and a clipboard-copy icon to copy the content.

## Documentation Contribution Flow Summary

**NOTE: For contributing `mesheryctl` reference section, refer [Contributing CLI](/contributing-cli)**


The following is a concise summary of the steps to contribute to Meshery documentation.

1. Create a fork, if you have not already, by following the steps described [here](CONTRIBUTING-gitflow.md)
2. In the local copy of your fork, navigate to the docs folder.
   `cd docs`
3. Create and checkout a new branch to make changes within
   `git checkout -b <my-changes>`
4. Edit/add documentation.
   `vi <specific page>.md`
5. Add redirect link on the old page (only when a new page is created that replaces the old page)
5. Run site locally to preview changes.
   `make docs`
6. Commit, [sign-off](#commit-signing), and push changes to your remote branch.
   `git push origin <my-changes>`
7. Open a pull request (in your web browser) against the repo: https://github.com/layer5io/meshery.

### Table of Contents in Sidebar (toc)

Sidebars use toc to create a table of contents. It is written in the following manner:

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

The output of the code snippet would be:

```
    Group 1
      Thing 1
      Thing 2
      Thing 3
```

### `if` condititional

This executes the block of code only if the given condition is true. It is executed in the following manner:

{% capture code_content %}{{ "{% if product.title == 'Awesome Shoes' " }}%}
  These shoes are awesome!
{{ "{% endif " }}%}
{% endcapture %}
{% include code.html code=code_content %}

If the condition is true, the output would be:

```
    These shoes are awesome!
```

### `for` loop

The for statement executes a block of code repeatedly. It is wriiten in the following manner:

{% capture code_content %}{{ "{% for names in collection.names " }}%}
  {{ "{{ name.title "}}}}
{{ "{% endfor " }}%}{% endcapture %}
{% include code.html code=code_content %}

The output produced by the above code snippet:

```
    Sam Ham Ethan
```

### Comment

Comments allow to leave a block of code unattended, any statements between opening and closing comment would not be executed.

### Include

The above tag is used to insert a already rendered file within the current template. It is written in the following manner:

{% capture code_content %}{{ "{% include file.html " }}%}{% endcapture %}
{% include code.html code=code_content %}

### Assign

The assign tag is used to create a new variable. It is written in the following manner:

{% capture code_content %}{{ "{% assign variable1 = true " }}%}{% endcapture %}
{% include code.html code=code_content %}

<!-- {% include suggested-reading.html %} -->

