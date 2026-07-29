---
title: mesheryctl-registry-purge
display_title: false
command: registry
subcommand: purge
---

# mesheryctl registry purge

Prune old model version directories

## Synopsis

Prune older versions of Meshery Models under ./models, retaining only the most recent --retain versions of each model.

Prerequisite: Execute this command from the root of a meshery/meshery repo fork; it operates on the "./models" directory relative to the current working directory.

The "kubernetes" and "meshery-core" models are always excluded from purging, regardless of flags. Use --exclude to skip additional models.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry purge [flags]

</div>
</pre> 

## Examples

Retain only the latest version of every model (default).
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry purge

</div>
</pre> 

Retain the 3 most recent versions of every model.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry purge --retain 3

</div>
</pre> 

Preview what would be removed without deleting anything.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry purge --dry-run

</div>
</pre> 

Additionally skip specific models.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry purge --exclude aws-ec2-controller,cilium

</div>
</pre> 

Skip the confirmation prompt.
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl registry purge --retain 2 -y

</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --dry-run          print what would be removed without deleting anything
      --exclude string   comma-delimited list of additional model names to exclude from purging
  -h, --help             help for purge
      --retain int       number of most-recent versions to retain per model (default 1)
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "/home/runner/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
