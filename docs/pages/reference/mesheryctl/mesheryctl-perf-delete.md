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

Delete performance profile

## Synopsis

Delete a performance profile by name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete [profile-name] [flags]

</div>
</pre> 

## Examples

Delete a performance profile by name
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete profile-name

</div>
</pre> 

Delete a performance profile with confirmation prompt
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete my-test-profile

</div>
</pre> 

Delete a performance profile and display output in JSON format
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl perf delete my-test-profile -o json

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for delete

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