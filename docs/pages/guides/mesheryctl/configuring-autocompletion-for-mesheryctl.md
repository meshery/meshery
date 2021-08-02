---
layout: default
title: Configuring Autocompletion for `mesheryctl`
permalink: guides/mesheryctl/configuring-autocompletion-for-mesheryctl
language: en
type: Guides
category: mesheryctl
---

If you would like to have `mesheryctl` commands automatically completed for use as you use `mesheryctl`, then use the following instructions to configure automatic completion within your environment.

## Autocompletion for Bash

### bash <= 3.2

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source /dev/stdin <<< "$(mesheryctl system completion bash)"
 </div></div>
 </pre>

### bash >= 4.0

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source <(mesheryctl system completion bash)
 </div></div>
 </pre>

### bash <= 3.2 on MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew install bash-completion # ensure you have bash-completion 1.3+
 mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
 </div></div>
 </pre>

### bash >= 4.0 on MacOS

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 brew install bash-completion@2
 mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
 </div></div>
 </pre>

## Autocompletion for zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 source <(mesheryctl system completion zsh)
 </div></div>
 </pre><br>

If shell completion is not already enabled in your environment you will need to enable it. You can execute the following once:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 ~/.zshrc > echo "autoload -U compinit; compinit" 
 </div></div>
 </pre>

_Note_ : You might need to restart your shell for this setup to take effect.

#### zsh on MacOS and Oh My zsh

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion zsh > "${fpath[1]}/_mesheryctl"
 </div></div>
 </pre>

### Autocompletion for fish

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion fish | source
 </div></div>
 </pre><br>

To load fish shell completions for each session, execute once:

 <pre class="codeblock-pre"><div class="codeblock">
 <div class="clipboardjs">
 mesheryctl system completion fish > ~/.config/fish/completions/mesheryctl.fish
 </div></div>
 </pre>

# Suggested Reading

For an exhaustive list of `mesheryctl` commands and syntax:

- See [`mesheryctl` Command Reference]({{ site.baseurl }}/reference/mesheryctl).

Guides to using Meshery's various features and components.

{% assign sorted_guides = site.pages | sort: "type" | reverse %}

<ul>
  <li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading mesheryctl</a></lI>
  {% for item in sorted_guides %}
  {% if item.type=="Guides" and item.category=="mesheryctl" and item.list!="exclude" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>
