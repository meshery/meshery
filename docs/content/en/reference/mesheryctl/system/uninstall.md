---
title: mesheryctl-system-uninstall
display_title: false
command: system
subcommand: uninstall
categories: [mesheryctl-system]
aliases:
- /reference/references/mesheryctl/system/uninstall/
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
when a normal uninstall fails. A normal uninstall stops at the first failure, so
a failed run leaves the cluster untouched from that point on. --force instead
attempts every cleanup stage and reports all the failures it hit.

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
      --config string    path to config file (default "~/.meshery/config.yaml")
  -c, --context string   (optional) temporarily change the current context.
  -v, --verbose          verbose output
  -y, --yes              (optional) assume yes for user interactive prompts.

</div>
</pre>

## What uninstall removes

For a Kubernetes installation, uninstall deletes the Broker and MeshSync custom
resource instances, uninstalls the Meshery Helm release, deletes the `brokers.meshery.io`
and `meshsyncs.meshery.io` custom resource definitions, and finally deletes the
`meshery` namespace and waits for it to disappear. The custom resources go first
because they cannot be found once the Helm release is gone.

## Docker cleanup limitations

A Docker installation is removed with Docker Compose, so uninstall can only
remove containers that the compose file in `~/.meshery/meshery.yaml` declares:

- Containers whose service is no longer in that compose file, or that were started
  with a different project name, are left running.
- Named volumes are kept, so the Meshery database and other persistent state
  survive the uninstall.
- Downloaded images and the contents of `~/.meshery` (the compose file, the
  config file and the seeded content) are kept as well.

Remove what you no longer need with `docker volume rm <volume>`,
`docker image rm <image>` and, for the local state, `rm -rf ~/.meshery`. If no
compose file is present, uninstall reports that Meshery is not installed and
leaves Docker untouched: it does not create a compose file to look for, and no
container, volume or image is removed.

{{% mesheryctl/system-uninstall-failure %}}

## See Also

Go back to [command reference index]({{< ref "reference/mesheryctl/_index.md" >}}), if you want to add content manually to the CLI documentation, please refer to the [instruction]({{< ref "project/contributing/cli/cli.md#documenting" >}}) for guidance.
