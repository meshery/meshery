---
title: mesheryctl-system-uninstall
display_title: false
command: system
subcommand: uninstall
categories: [mesheryctl-system]
---

# mesheryctl system uninstall

Uninstall Meshery

## Synopsis

Uninstall Meshery and completely remove the Meshery installation.

Stop and remove all Meshery containers (Docker) or delete all Meshery resources
from the cluster, including custom resource instances, custom resource
definitions, the Helm release and the Meshery namespace (Kubernetes).

For Docker installations, containers are always removed forcefully, so --force
has no effect. For Kubernetes installations, --force skips the Helm release
uninstallation and forcefully deletes the remaining Meshery resources; use it
when a normal uninstall fails.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system uninstall [flags]

</div>
</div>
</pre>

## Examples

Uninstall Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system uninstall

</div>
</div>
</pre>

(optional) keep the Meshery namespace during uninstallation
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system uninstall --keep-namespace

</div>
</div>
</pre>

Uninstall Meshery forcefully on Kubernetes (use it when system uninstall doesn't work)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system uninstall --force

</div>
</div>
</pre>

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --force            (optional) forcefully delete Kubernetes resources; has no effect for Docker installations
  -h, --help             help for uninstall
      --keep-namespace   (optional) keep the Meshery namespace during uninstallation

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string    path to config file (default "/home/runner/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## See Also

Go back to [command reference index]({{< ref "reference/references/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#preserving-manually-added-documentation" >}}) for guidance.
