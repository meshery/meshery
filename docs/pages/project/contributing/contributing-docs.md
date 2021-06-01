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

{% include alert.html type="info" title="Jekyll" content="The Meshery Docs site is built using Jekyll - a simple static site generator. You can learn more about Jekyll and setting up your development environment in the [Jekyll Docs](https://jekyllrb.com/docs/)." %}

- First [install Ruby](https://jekyllrb.com/docs/installation/), then install Jekyll and Bundler.

**Note:** Windows users can run Jekyll by following the [Windows Installation Guide](https://jekyllrb.com/docs/installation/windows/) and also installing Ruby Version Manager [RVM](https://rvm.io). RVM is a command-line tool which allows you to work with multiple Ruby environments on your local machine. Alternatively, if you're running Windows 10 version 1903 Build 18362 or higher, you can upgrade to Windows Subsystem for Linux [WSL](https://docs.microsoft.com/en-us/windows/wsl/install-win10) and run Jekyll in Linux instead.

Alternatively, if you are running Windows 10, you may install the Windows Subsystem for Linux:

- [WSL1](https://docs.microsoft.com/en-us/windows/wsl/install-win10) for Windows build version 1607 or higher.

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
  $ bundle install
  ```

### Serve the site

- Serve the code locally
  ```bash
  $ make site
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
