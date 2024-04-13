---
layout: page
title: Contributing to Meshery Docs
permalink: project/contributing/contributing-docs
abstract: How to contribute to Meshery Docs.
language: en
type: project
category: contributing
list: include
---

{% include alert.html type="info" title="Important" content="Before contributing, please review the <a href='https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md#documentation-contribution-flow'>Documentation Contribution Flow</a>." %}

Use the following steps to set up your development environment depending on your Operating System of choice, fork the repository, install dependencies and run the site locally.  You can then make changes, test locally, sign-off and commit, and push the changes made for review.

{% include alert.html type="info" title="Meshery Documentation Design Specification" content="See the <a href='https://docs.google.com/document/d/17guuaxb0xsfutBCzyj2CT6OZiFnMu9w4PzoILXhRXSo/edit#'>Meshery Documentation Design Specification</a> which serves to provide an overview of the tooling and approach used to create Meshery’s documentation and it information architecture." %}

## Documentation Framework

Meshery documentation is built using the following components:

- Framework - [Jekyll](https://jekyllrb.com)
- Theme - [https://github.com/vsoch/docsy-jekyll](https://github.com/vsoch/docsy-jekyll)
- Repo - [https://github.com/meshery/meshery/tree/master/docs](https://github.com/meshery/meshery/tree/master/docs)
- Site - [https://docs.meshery.io](https://docs.meshery.io)
- AWS API GW - an instance is configured to redirect from docs.meshery.io to meshery.layer5.io, because of the repo location of where the docs currently reside.

## Set up your development environment

{% include alert.html type="info" title="Jekyll" content="The Meshery Docs site is built using Jekyll - a simple static site generator. Jekyll can be installed on different platforms like Windows, Linux, and MacOS by the following steps " %}

### For Windows

{% include alert.html type="light" title="Note" content="Windows users can run Jekyll by following the <a href='https://jekyllrb.com/docs/installation/windows/'>Windows Installation Guide</a> and also installing Ruby Version Manager <a href='https://rvm.io'>RVM</a>. RVM is a command-line tool which allows you to work with multiple Ruby environments on your local machine. Alternatively, if you're running Windows 10 version 1903 Build 18362 or higher, you can upgrade to Windows Subsystem for Linux <a href='https://docs.microsoft.com/en-us/windows/wsl/install-win10'>WSL</a> and run Jekyll in Linux instead." %}

- Fire up your WSL VM and install the ruby version manager (RVM):

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">sudo apt update
  sudo apt install curl g++ gnupg gcc autoconf automake bison build-essential libc6-dev \
   	libffi-dev libgdbm-dev libncurses5-dev libsqlite3-dev libtool \
   	libyaml-dev make pkg-config sqlite3 zlib1g-dev libgmp-dev \
   	libreadline-dev libssl-dev
  sudo gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
  curl -sSL https://get.rvm.io | sudo bash -s stable
  sudo usermod -a -G rvm `whoami`
  </code>
  </div></pre>

  
  If `gpg --keyserver` gives an error, you can use:

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">sudo gpg --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
  </code>
  </div></pre>
  
  OR

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">sudo gpg2 --keyserver hkp://pool.sks-keyservers.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3 7D2BAF1CF37B13E2069D6956105BD0E739499BDB %}
  </code>
  </div></pre>

  Restart your WSL VM before moving forward to install Ruby.

- To install Ruby, run:

    <pre class="codeblock-pre"><div class="codeblock">
    <code class="clipboardjs">bash
    rvm install ruby
    rvm --default use ruby 2.7.5
    gem update
    gem install jekyll bundler
    </code>
    </div></pre>
  
- Update the Gemfile located in meshery/docs directory

    <pre class="codeblock-pre"><div class="codeblock">
    <code class="clipboardjs">ruby '3.0.2'
    </code>
    </div></pre>

  <strong>Note:</strong> In place of `3.0.2` add your installed version
  
- Also add this to the next line in the _Gemfile_.
    <pre class="codeblock-pre"><div class="codeblock">
    <code class="clipboardjs">gem 'wdm','>=0.1.0' if Gem.win_platform?
    </code>
    </div></pre>
  
{% include alert.html type="warning" title="Don't Commit Gemfile" content="This is just a workaround for your local machine. So, do not commit or push the modified Gemfile or Gemfile.lock during Pull Requests." %}
  
- Next, navigate to the following folder 
  ```C:\Ruby24-x64\lib\ruby\gems\2.4.0\gems\eventmachine-1.2.5-x64-mingw32\lib```
  
- Add ```require 'em/pure_ruby' ``` in the first line of the <strong>eventmachine.rb</strong> file

### For Linux

- Prerequisites
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">sudo apt-get update
  sudo apt-get install autoconf bison build-essential libssl-dev libyaml-dev libreadline6-dev zlib1g-dev libncurses5-dev libffi-dev libgdbm3 libgdbm-dev
  </code>
  </div></pre>

#### Installing rbenv

- Cloning the rbenv repository

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">git clone https://github.com/rbenv/rbenv.git ~/.rbenv
  </code>
  </div></pre>

{% include alert.html type="info" title="Note" content="Change bashrc with your shell specific rc file, for eg: if you are using zsh then the filename is zshrc." %}

- Setting the path

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">echo 'export PATH="$HOME/.rbenv/bin:$PATH"' >> ~/.bashrc
  </code>
  </div></pre>

- rbenv init

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">echo 'eval "$(rbenv init -)"' >> ~/.bashrc
  </code>
  </div></pre>

- Reload your bashrc

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">source ~/.bashrc
  </code>
  </div></pre>

- Check installation

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">type rbenv
  </code>
  </div></pre>


#### Install Ruby

- rbenv install version

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">rbenv install 3.2.2
  </code>
  </div></pre>

- To list all the versions that can be installed

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">rbenv install --list
  </code>
  </div></pre>

- Set which Ruby version you want to use

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">rbenv global version
  </code>
  </div></pre>

- Check Ruby installation

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">ruby -v
  </code>
  </div></pre>

### For MacOS

- Use docs here [Jekyll installation](https://jekyllrb.com/docs/installation/macos/)

{% include alert.html type="light" title="Note" content="In case of any installation issues, use the <a href='http://discuss.meshery.io)'>discussion forum</a>." %}

### Get the code

- Fork and then clone the [Meshery repository](https://github.com/meshery/meshery)

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">git clone https://github.com/YOUR-USERNAME/meshery
  </code>
  </div></pre>

- Change to the docs directory
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">cd docs
  </code>
  </div></pre>

- Install any Ruby dependencies
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">gem install bundler
  bundle install
  </code>
  </div></pre>

{% include alert.html type="info" title="Note" content="If you are a Mac user you do not need to install the Ruby dependencies, after moving on to the docs directory, you can serve the site." %}

### Serve the site

- Serve the code locally.
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make docs
  </code>
  </div></pre>

- If that gives an error run:
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">bundle exec jekyll serve --drafts --config _config_dev.yml
  </code>
  </div></pre>

  
  _From the `Makefile`, this command runs `$ bundle exec jekyll serve --drafts --livereload --config _config_dev.yml`. If this command results in an error, try running the server without Livereload using the command: `$ bundle exec jekyll serve --drafts --config _config_dev.yml`. Note that you will have to manually restart the server to reflect any changes made without Livereload. There are two Jekyll configurations, `jekyll serve` for local development and `jekyll build` to generate site artefacts for production deployment._

### Using Docker

If you've Docker and `make` installed in your system, then you can serve the site locally. This doesn't require the need to install _Jekyll_ and _Ruby_ in your system. 

{% include alert.html type="info" content="This may not work in Windows." %}

Run the following command from the _docs_ folder, else it will fail.

{% capture code_content %}make docker{% endcapture %}
{% include code.html code=code_content %}


#### Troubleshooting Note

While performing any of the above steps, if you receive an error about mismatching ruby versions similar to the example below, follow one of the steps mentioned below.

`Your ruby version is x.x.x but your Gemfile specified 2.7.x`

The reason for this error is because Jekyll always considers the exact version of Ruby unlike JavaScript.

- Install the required Ruby version by using `rvm` or by any means given above.
- If you're unable to install the required Ruby version, then manually configure the `Gemfile` (This is not recommended and should be done only if the above steps fail). Modify the ruby version inside the Gemfile similar to the example below:

  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">source "https://rubygems.org"
  ruby '3.1.1' //Change to the version you have installed
  </code>
  </div></pre>

The `Gemfile.lock` will be updated automatically once the `make docs` command is run. 
(For Windows, run `bundle exec jekyll serve` if WSL2 isn't present)

{% include alert.html type="warning" title="Don't Commit Gemfile" content="If you have followed the third step then please don't commit the changes made on `Gemfile` and `Gemfile.lock` in your branch to preserve integrity, else the CI action will fail to generate the site preview during PR." %}

### Using Gitpod
- Ensure you have an account on <a href="https://www.gitpod.io/" class="meshery-light">Gitpod</a> and add the <a href="https://chromewebstore.google.com/detail/gitpod/dodmmooeoklaejobgleioelladacbeki" class="meshery-light">Gitpod extension</a> to your browser.
- Open your forked Meshery repository on GitHub.
- Click on the "**Gitpod**" button in the top right corner of the repository page (it is only visible with the <a href="https://chromewebstore.google.com/detail/gitpod/dodmmooeoklaejobgleioelladacbeki" class="meshery-light">Gitpod browser extension</a> installed).

{% include alert.html type="dark" title="About Gitpod" content="Gitpod will automatically clone and open the repository for you in a browser-based version of Visual Studio Code. It will also automatically build the project for you on launch, comes with Docker and other tools pre-installed, making it one of the fastest ways to spin up a development environment for <a href='https://github.com/meshery/meshery'>Meshery.</a>" %}

- After opening the project on Gitpod, change to the docs directory.
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">cd docs
  </code>
  </div></pre>

- Serve the code locally.
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">make docs
  </code>
  </div></pre>

You should be able to access the site on port `4000`. If you want to access it in your localhost read the [docs for port-forwarding using ssh](https://www.gitpod.io/docs/configure/workspaces/ports#local-port-forwarding-via-ssh).

### Make Necessary Changes
- Make changes as required based on the issue you are solving.
- Ensure to verify that your changes reflect correctly by serving the site locally using `make docs`.

{% include alert.html type="light" title="Note" content="If the issue requires making new doc page that replaces the old page, please don't forget to add a redirect link on the old page. This redirect link field should have the link of the new page created." %}


### Create a Pull Request

- After making changes, don't forget to commit with the sign-off flag (-s)!
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">git commit -s -m “my commit message w/signoff”
  </code>
  </div></pre>

- Once all changes have been committed, push the changes.
  
  <pre class="codeblock-pre"><div class="codeblock">
  <code class="clipboardjs">git push origin [branch-name]
  </code>
  </div></pre>

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

If the code snippet does not contain any special characters then, it can be used directly as below:
```{% raw %}{% include code.html code="code_snippet_here" %}{% endraw %}````

If the code snippet has special characters then it must be captured first in variable `code_content`,  and then pass it to the component. 
```
{% raw %}{% capture code_content%}code_snippet_here{% endcapture %}{% endraw %}
{% raw %}{% include code.html code=code_content %}{% endraw %}
```

Don't use `code` component when the snippet is in between an ordered list, this breaks the order and next item in the list will start numbering from 1. Instead, use `<pre class="codeblock_pre">...</pre>` method described above."

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

{% include alert.html type="light" title="Note" content="For contributing `mesheryctl` reference section, refer <a href='/contributing-cli'>Contributing CLI</a>" %}


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
7. Open a pull request (in your web browser) against the repo: https://github.com/meshery/meshery.

### Navigation Table of Contents in Sidebar (toc)

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

In this example, **Group 1** is a ``parent`` section within the Table of Contents, and **Thing 1**, **Thing 2**, and **Thing 3** are ``children`` of **Group 1**. This is a direct hierarchy with a single level of children under the parent.

``Grandchildren`` are added in the sidebars in the following manner:

```
    toc:
  - title: Group 1
    subfolderitems:
      - page: Thing 1
        url: /thing1.html
        subfolderitems:
          - page: Subthing 1.1
            url: /subthing1_1.html
          - page: Subthing 1.2
            url: /subthing1_2.html
      - page: Thing 2
        url: /thing2.html
      - page: Thing 3
        url: /thing3.html
```

The output of the code snippet would be:

```
    Group 1
      Thing 1
       Subthing 1.1
       Subthing 1.2
      Thing 2
      Thing 3
```

Here, **Subthing 1.1** and **Subthing 1.2** are the ``grandchildren`` of **Thing 1**.

In `docs/_includes/sidebar.html` contains three levels of navigation hierarchy.

- **Parent**: It serves as a top level category for related content.
- **Children**: They are immediate subsections or topics that fall under the parent section. 
- **Grandchildren**: They are nested under Thing 1 and provide a more detailed breakdown of information within the child section. Grandchildren are used to organize content further, offering a more detailed structure for a specific topic.

These sections create a hierarchical and organized navigation experience for readers.

### Alerts

{% include alert.html type="info" title="What is an alert?" content="An alert is a box that can stand out to indicate important information. You can choose from levels success, warning, danger, info, and primary. This example is an info box, and the code for another might look like this:" %}

```
{%raw%}{% include alert.html type="info" title="Here is another!" %}{%endraw%}
```

Just for fun, here are all the types:

{% include alert.html type="warning" content="This is a warning" %}
{% include alert.html type="danger" content="This alerts danger!" %}
{% include alert.html type="success" content="This alerts success" %}
{% include alert.html type="info" content="This is useful information." %}
{% include alert.html type="primary" content="This is a primary alert" %}
{% include alert.html type="secondary" content="This is a secondary alert" %}

{% include alert.html type="light" content="This is a light alert" %}
{% include alert.html type="dark" content="This is a dark alert" %}

#### alert.html

Meshery Docs has a common include file [alert.html](https://github.com/meshery/meshery/blob/master/docs/_includes/alert.html), to provide consistent formatting for notes, warnings, and various informative callouts intended for the readers.

To use the `alert.html` feature in our documentation include the following code:
``` 
  {%raw%}{% include alert.html type="info" title="Here is another!" %}{%endraw%}
```


Other supported alert types include `warning`, `danger`,`success`,`primary`, `secondary`, `light`, `dark` .

### Quotes

You can include block quotes to emphasize text. 

> Here is an example. Isn't this much more prominent to the user?

## Development


### Suggested Reading

Disable suggested reading by setting the `suggested-reading` frontmatter variable to `false`.

### Editable Intra-page Table of Contents Toolbar

Control the display of this intra-page navigator with either page level or layout level frontmatter variables:

`display-toolbar`

Set to `true` (make "editable" toolbar visible) or `false` (hide "editable" toolbar)

### `if` conditional

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

Two helpful resources:

1. Liquid Docs - [https://shopify.github.io/liquid/](https://shopify.github.io/liquid/)
2. Jekyll Docs - [https://jekyllrb.com/docs/](https://jekyllrb.com/docs/)
