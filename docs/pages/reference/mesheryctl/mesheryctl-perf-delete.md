---
layout: default
title: mesheryctl-perf-delete
permalink: reference/mesheryctl/perf/delete
redirect_from: reference/mesheryctl/perf/delete/
type: reference
display-title: "false"
language: en
command: perf
subcommand: delete
---

# mesheryctl perf delete

Delete a Performance profile

## Synopsis

Delete Performance profiles by name or delete all profiles
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete [profile-name] [flags]

</div>
</pre> 

## Examples

Delete a specific performance profile
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete meshery-profile

</div>
</pre> 

Delete multiple profiles matching a pattern
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete meshery-profile-

</div>
</pre> 

Delete all performance profiles (with confirmation)
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete --all

</div>
</pre> 

Force delete without confirmation
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete meshery-profile --force

</div>
</pre> 

Force delete all profiles without confirmation
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete --all --force

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --all     (optional) Delete all performance profiles
      --force   (optional) Force delete without confirmation
  -h, --help    help for delete

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string          path to config file (default "/home/runner/.meshery/config.yaml")
  -o, --output-format string   (optional) format to display in [json|yaml]
  -t, --token string           (required) Path to meshery auth config
  -v, --verbose                verbose output
  -y, --yes                    (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.