# mesheryctl-component

LLMS index: [llms.txt](/pr-preview/pr-22677/llms.txt)

---

# mesheryctl component

Manage Meshery components

## Synopsis

List, search and view component(s) and detailed informations

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component [flags]

</div>
</div>
</pre> 

## Examples

Display number of available components in Meshery
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component --count

</div>
</div>
</pre> 

List available component(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component list

</div>
</div>
</pre> 

Search for component(s)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component search [component-name]

</div>
</div>
</pre> 

View a specific component
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl component view [component-name | component-id]

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --count   (optional) Get the number of components in total
  -h, --help    help for component

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

Go back to [command reference index](/pr-preview/pr-22677/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-22677/project/contributing/cli/cli/#preserving-manually-added-documentation) for guidance.

---

Section pages:

- [mesheryctl-component-list](/pr-preview/pr-22677/reference/mesheryctl/component/list/)
- [mesheryctl-component-search](/pr-preview/pr-22677/reference/mesheryctl/component/search/)
- [mesheryctl-component-view](/pr-preview/pr-22677/reference/mesheryctl/component/view/)
