# mesheryctl-system-logout

Source: /pr-preview/pr-21670/reference/references/mesheryctl/system/logout/

# mesheryctl system logout

Remove authentication for Meshery Server

## Synopsis


Remove authentication for Meshery Server

This command removes the authentication token from the user's filesystem
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system logout [flags]

</div>
</div>
</pre> 

## Examples

Logout current session with your Meshery Provider.
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl system logout

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help   help for logout

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

Go back to [command reference index](/pr-preview/pr-21670/reference/references/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-21670/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.
