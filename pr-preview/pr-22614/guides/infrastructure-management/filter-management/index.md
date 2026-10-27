# Envoy WASM Filter Management

> Import, inspect, publish, and remove Envoy WebAssembly (WASM) filters with Meshery UI and mesheryctl.

---

LLMS index: [llms.txt](/pr-preview/pr-22614/llms.txt)

---

Meshery stores Envoy WebAssembly (WASM) filters alongside your designs, so you can keep the filter binaries your data plane uses in one place, share them, and remove them when you are done. Use this guide to import a filter, inspect it, publish it to the Catalog, and delete it.

## Before you begin

- A running Meshery Server, and `mesheryctl` signed in to it if you want to use the CLI. See [Authenticate with Meshery via CLI](/pr-preview/pr-22614/guides/mesheryctl/authenticate-with-meshery-via-cli/).
- A compiled filter: a `.wasm` file on your machine or a URL that serves one. The community maintains sample filters in [layer5io/wasm-filters](https://github.com/layer5io/wasm-filters).

## Import a filter




<ul class="nav nav-tabs" id="tabs-1" role="tablist">
  <li class="nav-item">
      <button class="nav-link active"
          id="tabs-01-00-tab" data-bs-toggle="tab" data-bs-target="#tabs-01-00" role="tab"
          data-td-tp-persist="meshery ui" aria-controls="tabs-01-00" aria-selected="true">
        Meshery UI
      </button>
    </li><li class="nav-item">
      <button class="nav-link"
          id="tabs-01-01-tab" data-bs-toggle="tab" data-bs-target="#tabs-01-01" role="tab"
          data-td-tp-persist="mesheryctl" aria-controls="tabs-01-01" aria-selected="false">
        mesheryctl
      </button>
    </li>
</ul>

<div class="tab-content" id="tabs-1-content">
    <div class="tab-body tab-pane fade show active"
        id="tabs-01-00" role="tabpanel" aria-labelled-by="tabs-01-00-tab" tabindex="1">
        <ol>
<li>Open <code>/configuration/filters</code> on your Meshery Server (for example <code>http://localhost:9081/configuration/filters</code>).</li>
<li>Click <strong>Import Filters</strong>.</li>
<li>Enter a <strong>Name</strong> for the filter.</li>
<li>Choose an <strong>Upload method</strong>:
<ul>
<li><strong>File Upload</strong>: select the <code>.wasm</code> file on your machine.</li>
<li><strong>URL Upload</strong>: paste the URL that serves the file.</li>
</ul>
</li>
<li>Click <strong>Import</strong>.</li>
</ol>

    </div>
    <div class="tab-body tab-pane fade"
        id="tabs-01-01" role="tabpanel" aria-labelled-by="tabs-01-01-tab" tabindex="1">
        <p>Import from a local file or a URL:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl filter import /path/to/filter.wasm
</span></span><span class="line"><span class="cl">mesheryctl filter import https://example.com/myfilter.wasm
</span></span></code></pre></div><p>Set a name, and pass the filter&rsquo;s configuration either as a file path or as an inline string:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl filter import /path/to/filter.wasm --name <span class="s2">&#34;rate-limit&#34;</span> --wasm-config ./config.json
</span></span></code></pre></div>
    </div>
</div>


## Find and inspect filters




<ul class="nav nav-tabs" id="tabs-2" role="tablist">
  <li class="nav-item">
      <button class="nav-link active"
          id="tabs-02-00-tab" data-bs-toggle="tab" data-bs-target="#tabs-02-00" role="tab"
          data-td-tp-persist="meshery ui" aria-controls="tabs-02-00" aria-selected="true">
        Meshery UI
      </button>
    </li><li class="nav-item">
      <button class="nav-link"
          id="tabs-02-01-tab" data-bs-toggle="tab" data-bs-target="#tabs-02-01" role="tab"
          data-td-tp-persist="mesheryctl" aria-controls="tabs-02-01" aria-selected="false">
        mesheryctl
      </button>
    </li>
</ul>

<div class="tab-content" id="tabs-2-content">
    <div class="tab-body tab-pane fade show active"
        id="tabs-02-00" role="tabpanel" aria-labelled-by="tabs-02-00-tab" tabindex="2">
        <p>The Filters page lists every filter you own, in grid or table view. Filter the list by <strong>Visibility</strong> (Public, Private, and, when the Catalog is enabled, Published). Each filter offers these actions:</p>
<ul>
<li><strong>Config</strong>: edit the filter&rsquo;s configuration.</li>
<li><strong>Download</strong>: download the filter binary.</li>
<li><strong>Filter Information</strong>: view the filter&rsquo;s details.</li>
<li><strong>Clone</strong>: copy a published filter into your own filters.</li>
</ul>

    </div>
    <div class="tab-body tab-pane fade"
        id="tabs-02-01" role="tabpanel" aria-labelled-by="tabs-02-01-tab" tabindex="2">
        <p>List filters, 25 per page, optionally searching by name:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl filter list
</span></span><span class="line"><span class="cl">mesheryctl filter list <span class="s2">&#34;rate&#34;</span>
</span></span><span class="line"><span class="cl">mesheryctl filter list --page <span class="m">2</span>
</span></span><span class="line"><span class="cl">mesheryctl filter list --count
</span></span></code></pre></div><p>View one filter by name, ID, or a unique prefix of either, or view them all:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl filter view <span class="s2">&#34;rate-limit&#34;</span>
</span></span><span class="line"><span class="cl">mesheryctl filter view --all --output-format json --save
</span></span></code></pre></div><p><code>--output-format</code> accepts <code>yaml</code> (the default) or <code>json</code>. <code>--save</code> writes the output to a file.</p>

    </div>
</div>


## Publish a filter to the Catalog

If your provider enables the [Meshery Catalog](https://meshery.io/catalog) and you have permission to publish WASM filters, a **Publish** action appears on each of your filters. Click it, fill in the catalog details, and submit. Published filters go through the provider's approval flow before others can see and clone them. Use **Unpublish** to withdraw one.

## Delete a filter




<ul class="nav nav-tabs" id="tabs-3" role="tablist">
  <li class="nav-item">
      <button class="nav-link active"
          id="tabs-03-00-tab" data-bs-toggle="tab" data-bs-target="#tabs-03-00" role="tab"
          data-td-tp-persist="meshery ui" aria-controls="tabs-03-00" aria-selected="true">
        Meshery UI
      </button>
    </li><li class="nav-item">
      <button class="nav-link"
          id="tabs-03-01-tab" data-bs-toggle="tab" data-bs-target="#tabs-03-01" role="tab"
          data-td-tp-persist="mesheryctl" aria-controls="tabs-03-01" aria-selected="false">
        mesheryctl
      </button>
    </li>
</ul>

<div class="tab-content" id="tabs-3-content">
    <div class="tab-body tab-pane fade show active"
        id="tabs-03-00" role="tabpanel" aria-labelled-by="tabs-03-00-tab" tabindex="3">
        <p>Select one or more filters in table view and click the delete icon, then confirm <strong>Delete</strong>.</p>

    </div>
    <div class="tab-body tab-pane fade"
        id="tabs-03-01" role="tabpanel" aria-labelled-by="tabs-03-01-tab" tabindex="3">
        <p>Delete by name, ID, or a unique prefix of either. If the prefix matches more than one filter, Meshery deletes the first match, so use the full name or ID when in doubt.</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl filter delete <span class="s2">&#34;rate-limit&#34;</span>
</span></span></code></pre></div>
    </div>
</div>


## Frequently Asked Questions

### Q: Why does Meshery provide and manage WASM filters?

Meshery manages WASM filters for data plane proxies for:

- <strong>Agility</strong>: in partnership with service mesh control planes, Meshery can dynamically load, configure, and unload filter(s) in running Envoy proxies.
- <strong>Maintainability</strong>: To extend the functionality you don't have to change the entire codebase.
- <strong>Diversity</strong>: WASM filters can be easily compiled with many popular programming languages like C/C++ and Rust.
- <strong>Reliability and isolation</strong>: WASM filters are deployed into VM and hence they are isolated from the Envoy process. For instance if the WASM filter crashes it will not impact the envoy process.
- <strong>Security</strong>: The WASM filter uses a well-defined API to connect with the host i.e. Envoy proxy which can be modified only a limited number of connections.
- <strong>Speed</strong>: WASM runs at near-native speed.

### Q: What is Envoy proxy?

An <strong>Envoy</strong> is a high-performance, programmable L3/L4 and L7 proxy that a number of service mesh use as their data plane. See the service mesh landscape for a complete list. Connection and request handling are at the core of Envoy's focus. Envoy has a pluggable filter module that allows the incorporation of WebAssembly filters that provide additional traffic filtering intelligence. Often this includes higher-order functions such as access control, transformation, data enrichment, auditing and more.

### Q: What is WASM?

WebAssembly is a low-level assembly-like language with a compact binary format that runs with near-native performance and provides languages such as C/C++, C# and Rust with a compilation target so that they can run in many different environments, including behind Envoy's Application Binary Interface to extend Envoy with new capabilities. WASM programs are distributed in an efficient portable binary instruction format and embed into isolated execution environments.

### Q: What are the disadvantages of WASM filters to be taken into consideration?

Execution of WASM filters requires system resources. Just as executing the same functionality that your filter provides outside of the WASM virtual machine, resources can include overhead of latency, CPU, and memory. See the [Cloud Native Performance](https://smp-spec.io) project for more detail.

### Q: What is the Envoy Proxy WASM SDK?

Envoy proxy runs the WASM filter in a stack-based virtual machine so that the filter memory is isolated from the host environment. The interaction between the embedded host i.e. Envoy proxy, and the WASM filter are done via call backs and functions by the Envoy Proxy WASM SDK. The Envoy Proxy WASM SDK is implemented in various programming languages like C++, RUST, AssemblyScript and Go (still experimental)

### Q: What is Image Hub?

Image Hub is a sample application which was first officially demonstrated at DockerCon 2020. This sample application is written to run on consul to explore the WebAssembly (WASM) modules in the Envoy filter. This application is written in the RUST programming language. These modules are used to implement the multi-tenancy or to implement the per user rate limiting in the applications endpoint.
