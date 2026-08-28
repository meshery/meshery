---
title: mesheryctl-connection
display_title: false
command: connection
subcommand: nil
---

# mesheryctl connection

Manage Meshery connections

## Synopsis

View and manage your Meshery connection.

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection [flags]</code>
</div>
</pre> 

## Examples

Display total count of all available connections
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection --count</code>
</div>
</pre> 

Create a new Kubernetes connection using a specific type
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection create --type aks</code>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection create --type eks</code>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection create --type gke</code>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection create --type minikube</code>
</div>
</pre> 

List all the connection
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection list</code>
</div>
</pre> 

<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection list --count</code>
</div>
</pre> 

Delete a connection
<pre class='codeblock-pre'>
<div class='codeblock'>
<code class='clipboardjs'>mesheryctl connection delete [connection_id]</code>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -c, --count   Display the count of total available connections
  -h, --help    help for connection

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
