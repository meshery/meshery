---
layout: page
title: Contributing to Meshery Docs
permalink: project/contributing-docs
description: How to contribute to Meshery Docs.
language: en
type: project
category: contributing
---

Before contributing, please review the [Documentation Contribution Flow](https://github.com/layer5io/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow). In the following steps you will set up your development environment, fork and clone the repository, run the site locally, and finally commit, sign-off, and push any changes made for review.

{% include alert.html type="info" title="Meshery Documentation Design Specification" content="See the <a href='https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit#'>Meshery Documentation Design Specification</a> which serves to provide an overview of the tooling and approach used to create Meshery’s documentation and it information architecture." %}

## Documentation Framework

Meshery documentation is made of these components:

- Framework - Jekyll
- Theme - https://github.com/vsoch/docsy-jekyll
- Repo - https://github.com/layer5io/meshery/tree/master/docs
- DNS - https://meshery.layer5.io/docs
- AWS API GW - an instance is configured to redirect from docs.meshery.io to meshery.layer5.io, because of the repo location of where the docs currently reside.

## Set up your development environment

{% include alert.html type="info" title="Jekyll" content="The Meshery Docs site is built using Jekyll - a simple static site generator. Jekyll can be installed in the different platforms like Windows, Linux, and MacOs by the following steps:" %}

### For Windows

- Fire up your WSL VM and install the ruby version manager (RVM): 
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

  If `gpg --keyserver` gives an error, you can use:
  ```bash
    sudo gpg --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
  ```
  or
  ```bash
    sudo gpg2 --keyserver hkp://pool.sks-keyservers.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB
  ```
  Restart your WSL VM before moving forward.

- For installing Ruby, run:
  ```bash
    rvm install ruby
    rvm --default use ruby 2.7.1
    gem update
    gem install jekyll bundler
  ```

### For Linux

- Prerequisites
  ```bash
    sudo apt-get update
    sudo apt-get install autoconf bison build-essential libssl-dev libyaml-dev libreadline6-dev zlib1g-dev libncurses5-dev libffi-dev libgdbm3 libgdbm-dev
  ```
#### Installing rbenv

- Cloning the rbenv repository
  ```bash
    git clone https://github.com/rbenv/rbenv.git ~/.rbenv
  ```
- Setting the path
  ```bash
    echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc
  ```
- rbenv init
  ```bash
    echo 'eval "$(rbenv init -)"' >> ~/.bashrc
  ```
- Reload your bashrc
  ```bash
    source ~/.bashrc
  ```
  <strong>Note:</strong> Change bashrc with your shell specific rc file, for eg: if you are using zsh then the filename is zshrc.
  
- Check installation
  ```bash
    type rbenv
  ```
#### Install Ruby

- rbenv install version
```bash
  rbenv install 2.5.1
```
- To list all the versions that can be installed
```bash
  rbenv install --list-all
```
- Set which Ruby version you want to use
```bash
  rbenv global version
```
- Check Ruby installation
```bash
  ruby -v
```

### For MacOs

- Use docs here [Jekyll installation](https://jekyllrb.com/docs/installation/macos/)
- After installing Ruby and Jekyll, run the following:
```bash
  bundle install
  bundle config build.puma --with-cflags="-Wno-error=implicit-function-declaration"
  gem install puma:4.3.5 -- --with-cflags="-Wno-error=implicit-function-declaration"
  bundle install
  make file
```

### Get the code

- Fork and then clone the [Meshery repository](https://github.com/layer5io/meshery)
  ```bash
  $ git clone https://github.com/YOUR-USERNAME/meshery
  ```
- Change to the docs directory
  ```bash
  $ cd docs
  ```
- Install any Ruby dependencies
  ```bash
  $ gem install bundler
  $ bundle install
  ```
  <strong>Note:</strong> If you are a Mac user you do not need to install the Ruby dependencies, after moving on to the docs directory, you can serve the site.

### Serve the site

- Serve the code locally
  ```bash
  $ make site
  ```
  - If that gives an error run:
  ```bash
    $ bundle exec jekyll serve
  ```
  _Note: From the Makefile, this command is actually running `$ bundle exec jekyll serve --drafts --livereload`. There are two Jekyll configuration, `jekyll serve` for developing locally and `jekyll build` when you need to generate the site artifacts for production._

### Create a Pull Request

- After making changes, don't forget to commit with the sign-off flag (-s)!
  ```bash
  $ commit -s -m “my commit message w/signoff”
  ```
- Once all changes have been committed, push the changes.
  ```bash
  $ git push origin <branch-name>
  ```
- Then on Github, navigate to the [Meshery repository](https://github.com/layer5io/meshery) and create a pull request from your recently pushed changes!

---

- _See the [Meshery Documentation Google Doc](https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit) for additional reference._

## Using the features of Meshery Docs

### Clipboard Feature

Most popular clipboard plugins like Clipboard JS require the manual creation of a new ID for each code snippet. A different approach is used here. For code snippets, we either use html tags or markdown in the following manner:

```
   <pre><code>
     code snippet
   </code></pre>
``` 
_\<pre\> tags are optional unless the code snippet is in a paragraph format_


Or

````
code snippet          or         ```code snippet```
````

Whenever the code tags are detected, the clipboard javascript file is automatically loaded. Each code element is given a custom id and a clipboard-copy icon to copy the content.

## Documentation Contribution Flow Summary

The following is a concise summary of the steps to contribute to Meshery documentation.

1. Create a fork, if you have not already, by following the steps described [here](CONTRIBUTING-gitflow.md)
1. In the local copy of your fork, navigate to the docs folder.
   `cd docs`
1. Create and checkout a new branch to make changes within
   `git checkout -b <my-changes>`
1. Edit/add documentation.
   `vi <specific page>.md`
1. Run site locally to preview changes.
   `make site`
1. Commit, [sign-off](#commit-signing), and push changes to your remote branch.
   `git push origin <my-changes>`
1. Open a pull request (in your web browser) against the repo: https://github.com/layer5io/meshery.


# Suggested Reading

{% assign sorted_reading = site.pages | sort: page.title | reverse %}

<ul>
  {% for item in sorted_reading %}
  {% if item.type=="project" and item.category=="contributing" and item.list!="exclude" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>


