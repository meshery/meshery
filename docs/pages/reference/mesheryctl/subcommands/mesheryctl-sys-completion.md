---
layout: default
title: mesheryctl system completion
permalink: reference/mesheryctl/system/completion
type: reference
display-title: "false"
language: en
lang: en
# image: /assets/img/platforms/brew.png
---

<!-- Copy this template to create individual doc pages for each mesheryctl commands -->

<!-- Name of the command -->
#  mesheryctl system completion

## Description

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.completion.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
{% endfor %}

<!-- Basic usage of the command -->
<pre class="codeblock-pre">
  <div class="codeblock">
    mesheryctl system completion [bash|zsh|fish]
  </div>
</pre>

## Examples

{% for subcommand_hash in site.data.mesheryctlcommands.lifecycle.system.completion.command %}{% assign subcommand = subcommand_hash[1] %}
{{ subcommand.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ subcommand.usage }}
  </div>
</pre>
{% endfor %}
{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.completion.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
  {{ flag.usage }}
  </div>
</pre>
{% endfor %}

<br/>

bash <= 3.2
<pre class="codeblock-pre">
  <div class="codeblock">
  source /dev/stdin <<< "$(mesheryctl system completion bash)"
  </div>
</pre>

<br/>

bash >= 4.0
<pre class="codeblock-pre">
  <div class="codeblock">
  source <(mesheryctl system completion bash)
  </div>
</pre>

<br/>

bash <= 3.2 on OSX
<pre class="codeblock-pre">
  <div class="codeblock">
  brew install bash-completion # ensure you have bash-completion 1.3+
  mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
  </div>
</pre>

<br/>

bash >= 4.0 on OSX
<pre class="codeblock-pre">
  <div class="codeblock">
  brew install bash-completion@2
  mesheryctl system completion bash > $(brew --prefix)/etc/bash_completion.d/mesheryctl
  </div>
</pre>

<br/>

**zsh**
<br/>
If shell completion is not already enabled in your environment you will need to enable it.  You can execute the following once:
Might need to start a new shell for this setup to take effect.
<pre class="codeblock-pre">
  <div class="codeblock">
  echo "autoload -U compinit; compinit" >> ~/.zshrc
  source <(mesheryctl system completion zsh)
  </div>
</pre>

<br/>

zsh on OSX / oh-my-zsh
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl system completion zsh > "${fpath[1]}/_mesheryctl"
  </div>
</pre>

<br/>  

fish
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl system completion fish | source
  </div>
</pre>

<br/>
To load fish shell completions for each session, execute once:
<pre class="codeblock-pre">
  <div class="codeblock">
  mesheryctl system completion fish > ~/.config/fish/completions/mesheryctl.fish
  </div>
</pre> 
  
<br/>


<!-- Options/Flags available in this command -->
## Options

{% for flag_hash in site.data.mesheryctlcommands.lifecycle.system.completion.flag %}{% assign flag = flag_hash[1] %}
{{ flag.description }}
<pre class="codeblock-pre">
  <div class="codeblock">
    {{ flag.flag }}
  </div>
</pre>
{% endfor %}
<br/>
