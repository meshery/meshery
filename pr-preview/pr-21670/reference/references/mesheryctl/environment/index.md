# mesheryctl-environment

Source: /pr-preview/pr-21670/reference/references/mesheryctl/environment/

# mesheryctl environment

Manage environments

## Synopsis

Create, delete, list of view details of environment(s) of a specific organization

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl environment [flags]

</div>
</div>
</pre> 

## Examples

Create an environment in an organization
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl environment create --orgId [orgId] --name [name] --description [description]

</div>
</div>
</pre> 

Delete an environment in an organization
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl environment delete environment-id

</div>
</div>
</pre> 

List of registered environments in an organization
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl environment list --orgId [orgId]

</div>
</div>
</pre> 

View a particular environment
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl environment view --orgId [orgId]

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for environment

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

Go back to [command reference index](/pr-preview/pr-21670/reference/references/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-21670/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.
