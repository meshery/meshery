# mesheryctl-design

LLMS index: [llms.txt](/pr-preview/pr-22666/llms.txt)

---

# mesheryctl design

Manage cloud native designs

## Synopsis

Manage cloud and cloud native infrastructure using predefined designs.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design [flags]

</div>
</div>
</pre> 

## Examples

Apply design file:
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design apply --file [path to design file | URL of the file]

</div>
</div>
</pre> 

Delete design file:
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design delete --file [path to design file]

</div>
</div>
</pre> 

View design file:
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design view [design name | ID]

</div>
</div>
</pre> 

List all designs:
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl design list

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
  -h, --help           help for design
  -t, --token string   Path to token file default from current context

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

Go back to [command reference index](/pr-preview/pr-22666/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-22666/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.

---

Section pages:

- [mesheryctl-design-apply](/pr-preview/pr-22666/reference/mesheryctl/design/apply/)
- [mesheryctl-design-delete](/pr-preview/pr-22666/reference/mesheryctl/design/delete/)
- [mesheryctl-design-deploy](/pr-preview/pr-22666/reference/mesheryctl/design/deploy/)
- [mesheryctl-design-evaluate](/pr-preview/pr-22666/reference/mesheryctl/design/evaluate/)
- [mesheryctl-design-export](/pr-preview/pr-22666/reference/mesheryctl/design/export/)
- [mesheryctl-design-import](/pr-preview/pr-22666/reference/mesheryctl/design/import/)
- [mesheryctl-design-list](/pr-preview/pr-22666/reference/mesheryctl/design/list/)
- [mesheryctl-design-undeploy](/pr-preview/pr-22666/reference/mesheryctl/design/undeploy/)
- [mesheryctl-design-view](/pr-preview/pr-22666/reference/mesheryctl/design/view/)
