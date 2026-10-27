# Importing and Exporting Designs

> Manage and transfer Meshery designs seamlessly through import and export using CLI and UI.

---

LLMS index: [llms.txt](/pr-preview/pr-22622/llms.txt)

---

# Importing and Exporting Designs

Import a design to bring infrastructure configuration you already have into Meshery. Export a design to back it up, move it to another Meshery Server, or use it outside Meshery.

A [design](/pr-preview/pr-22622/concepts/logical/designs/) holds **components**, such as Deployments and Services, and the **relationships** between them.

## What happens when you import

You can import any of these formats:

- Meshery Design (YAML, or an OCI image exported from Meshery)
- Kubernetes manifest (YAML)
- Helm chart (`.tgz` archive)
- Docker Compose file
- Kustomize package (archive)

Meshery Server identifies the format from the file itself. It converts Helm charts, Docker Compose files, and Kustomize packages to a Kubernetes manifest, matches each resource's `apiVersion` and `kind` to a component in the registry, builds a design from the matched components, and saves the design. Importing does not deploy anything. Deploy the design afterward from Meshery UI or with `mesheryctl design deploy`.

## Import a design

<div class="tab-container tab-level-1"><input id="import-design-tabs-7-1" type="radio" name="import-design-tabs-7" checked>
    <label for="import-design-tabs-7-1"><i class="fa fa-desktop" aria-hidden="true"></i>Meshery UI
    </label>
    <section class="tabbed" id="import-design-tabs-7-content-1"><ol>
<li>Go to <strong>Configuration</strong> &gt; <strong>Designs</strong> and click <strong>Import Design</strong>.</li>
<li>Enter a <strong>Design Name</strong>.</li>
<li>Choose an <strong>Upload method</strong>:
<ul>
<li><strong>File Upload</strong>, then pick the file from your computer.</li>
<li><strong>URL Import</strong>, then enter a direct URL to the file, for example <code>https://raw.githubusercontent.com/&lt;org&gt;/&lt;repo&gt;/main/app.yaml</code>. The URL must return the file itself, not a web page.</li>
</ul>
</li>
<li>Click <strong>Import</strong>.</li>
</ol>
</section><input id="import-design-tabs-7-2" type="radio" name="import-design-tabs-7">
    <label for="import-design-tabs-7-2"><i class="fa fa-terminal" aria-hidden="true"></i>mesheryctl
    </label>
    <section class="tabbed" id="import-design-tabs-7-content-2"><div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl design import -f &lt;path-or-URL&gt; <span class="o">[</span>-n &lt;design-name&gt;<span class="o">]</span> <span class="o">[</span>-s &lt;source-type&gt;<span class="o">]</span>
</span></span></code></pre></div><ul>
<li><code>-f</code> is a local path or a direct URL to the file.</li>
<li><code>-n</code> names the design. Without it, <code>mesheryctl</code> uses the file name.</li>
<li><code>-s</code> is optional because Meshery Server detects the format. If you pass it, use one of <code>&quot;Kubernetes Manifest&quot;</code>, <code>&quot;Helm Chart&quot;</code>, <code>&quot;Docker Compose&quot;</code>, or <code>&quot;Design&quot;</code> (case-insensitive).</li>
</ul>
<p>For example, to import a Kubernetes manifest:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl design import -f ./app-deployment.yaml -n my-app
</span></span></code></pre></div><p><code>mesheryctl</code> prints the new design&rsquo;s name and ID.</p>
</section></div>

### Import a Kustomize package

To import a Kustomize package, upload a `.zip`, `.tar`, `.tar.gz`, or `.tgz` archive with a file named `kustomization.yaml` at the root of the archive. Include the files it references in the same archive. For example:

```text
kustomization.yaml
deployment.yaml
service.yaml
```

You can also import the archive from a direct URL. Meshery downloads the file and takes its name from the response's `Content-Disposition` header or the last segment of the URL, so make sure that name ends in `.zip`, `.tar`, `.tar.gz`, or `.tgz`.

## Export a design

<div class="tab-container tab-level-1"><input id="export-design-tabs-8-1" type="radio" name="export-design-tabs-8" checked>
    <label for="export-design-tabs-8-1"><i class="fa fa-desktop" aria-hidden="true"></i>Meshery UI
    </label>
    <section class="tabbed" id="export-design-tabs-8-content-1"><ol>
<li>Go to <strong>Configuration</strong> &gt; <strong>Designs</strong>.</li>
<li>Click <strong>Download</strong> on the design you want to export.</li>
<li>Click the download icon next to the format you want:</li>
</ol>
<table>
  <thead>
      <tr>
          <th>Format</th>
          <th>What you get</th>
      </tr>
  </thead>
  <tbody>
      <tr>
          <td><strong>Meshery Design (yaml)</strong></td>
          <td>A complete Meshery design, including embedded images. Use it for backups or to move a design to another Meshery Server without losing anything.</td>
      </tr>
      <tr>
          <td><strong>Meshery Design (OCI image)</strong></td>
          <td>The design as an OCI image that you can push to and pull from a container registry.</td>
      </tr>
      <tr>
          <td><strong>Kubernetes Manifest (yaml)</strong></td>
          <td>The design&rsquo;s Kubernetes resources as a plain manifest, for use with <code>kubectl</code> or other tools.</td>
      </tr>
      <tr>
          <td><strong>Helm Chart (tar.gz)</strong></td>
          <td>The design packaged as a Helm chart.</td>
      </tr>
  </tbody>
</table>
</section><input id="export-design-tabs-8-2" type="radio" name="export-design-tabs-8">
    <label for="export-design-tabs-8-2"><i class="fa fa-terminal" aria-hidden="true"></i>mesheryctl
    </label>
    <section class="tabbed" id="export-design-tabs-8-content-2"><div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl design <span class="nb">export</span> &lt;design-name-or-ID&gt; <span class="o">[</span>--type &lt;type&gt;<span class="o">]</span> <span class="o">[</span>-o &lt;directory&gt;<span class="o">]</span>
</span></span></code></pre></div><ul>
<li><code>--type</code> is one of:
<ul>
<li><code>current</code> (default): the design as Meshery stores it now.</li>
<li><code>original</code>: the file you originally imported, such as the Helm chart or manifest.</li>
<li><code>oci</code>: the design as an OCI image.</li>
</ul>
</li>
<li><code>-o</code> is the directory to save the file in. The default is the current directory.</li>
</ul>
<p>For example, to export a design as an OCI image into <code>./exports</code>:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl design <span class="nb">export</span> my-app --type oci -o ./exports
</span></span></code></pre></div><p>To export a design as a Helm chart or a Kubernetes manifest, use Meshery UI.</p>
</section></div>

## See also

- [Creating a Meshery Design](/pr-preview/pr-22622/guides/configuration-management/creating-a-meshery-design/)
- [`mesheryctl design` reference](/pr-preview/pr-22622/reference/mesheryctl/design/)
