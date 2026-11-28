# mesheryctl-registry-init

LLMS index: [llms.txt](/pr-preview/pr-24543/llms.txt)

---

# mesheryctl registry init

Scaffold a model-source.yaml for a new model

## Synopsis

Scaffold models/<model-name>/model-source.yaml, the canonical generation inputs for a model,
along with placeholder SVGs and an overrides/ directory, so a new model can be added by pull request.

Run it from the root of a meshery/meshery checkout so the model lands in models/.

<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry init &lt;model-name&gt; [flags]

</div>
</div>
</pre> 

## Examples

Scaffold a model whose source is an Artifact Hub package
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry init cert-manager --source https://artifacthub.io/packages/helm/cert-manager/cert-manager --category "Security &amp; Compliance" --sub-category "Key Management" --description "Cloud native certificate management for Kubernetes."

</div>
</div>
</pre> 

Scaffold a model from a GitHub repository, with an explicit registrant
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry init my-operator --source git://github.com/example/my-operator/main/deploy/crds --registrant meshery --category Provisioning --sub-category "Automation &amp; Configuration" --description "Operator for my workloads."

</div>
</div>
</pre> 

Overwrite an existing model-source.yaml (SVGs and overrides are kept)
<pre class='codeblock-pre'>
<div class='codeblock'>
<div class='clipboardjs'>
mesheryctl registry init cert-manager --source https://artifacthub.io/packages/helm/cert-manager/cert-manager --category "Security &amp; Compliance" --sub-category "Key Management" --description "Cloud native certificate management for Kubernetes." --force

</div>
</div>
</pre> 

## Options

<pre class='codeblock-pre'>
<div class='codeblock'>
      --category string       (required) model category
      --description string    (required) summary description of the model
      --display-name string   human-readable model name (default: derived from the model name)
      --force                 overwrite an existing model-source.yaml (SVGs, overrides and version folders are never touched)
  -h, --help                  help for init
  -p, --path string           directory the model folder is created under (default "models")
      --registrant string     registrant of the model (default: derived from the source kind)
  -s, --source string         (required) upstream URL of the model source: a GitHub repository, an Artifact Hub package, or an HTTP manifest or chart
      --source-kind string    kind of source [github|artifacthub|http] (default: detected from --source)
      --sub-category string   (required) model subcategory

</div>
</pre>

## Options inherited from parent commands

<pre class='codeblock-pre'>
<div class='codeblock'>
      --config string   path to config file (default "~/.meshery/config.yaml")
  -v, --verbose         verbose output

</div>
</pre>

## See Also

Go back to [command reference index](/pr-preview/pr-24543/reference/mesheryctl/), if you want to add content manually to the CLI documentation, please refer to the [instruction](/pr-preview/pr-24543/project/contributing/cli/cli/#documenting) for guidance.
