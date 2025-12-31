---
layout: default
title: mesheryctl-connection
permalink: reference/mesheryctl/connection
redirect_from: reference/mesheryctl/connection/
type: reference
display-title: "false"
language: en
command: connection
subcommand: nil
---

# mesheryctl connection

Manage Meshery connections

## Synopsis

View and manage your Meshery connection.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/connection
<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection [flags]

</div>
</pre> 

## Examples

### Count all available connections

Use this command to display the total number of connections configured in Meshery.
This is useful for quick verification or scripting.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection --count
</div>
</pre>

---

### List all connections

Lists all available connections along with their ID, name, type, and status.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection list
</div>
</pre>

Example output:
<pre class='codeblock-pre'>
<div class='codeblock'>
ID          NAME               TYPE          STATUS
4f3a2c9e    local-k8s          kubernetes    connected
9b7d1a22    prod-cluster       kubernetes    disconnected
</div>
</pre>

---

### Delete a connection using connection ID

Deletes a specific connection using its unique connection ID.
Use `mesheryctl connection list` to find the connection ID.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl connection delete [connection_id]
</div>
</pre>

Successful deletion confirmation:
<pre class='codeblock-pre'>
<div class='codeblock'>
Connection deleted successfully
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

* [mesheryctl connection create](/reference/mesheryctl/connection/create)
* [mesheryctl connection delete](/reference/mesheryctl/connection/delete)
* [mesheryctl connection list](/reference/mesheryctl/connection/list)
* [mesheryctl connection view](/reference/mesheryctl/connection/view)

Go back to [command reference index](/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/project/contributing/contributing-cli#preserving-manually-added-documentation) for guidance.
