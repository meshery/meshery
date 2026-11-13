# KubeCon Workshop: Publish a Design and Build a Custom Model

> A 60-minute Meshery Playground workshop. Choose Path A (Kanvas) or Path B (Meshery UI and mesheryctl). Every attendee publishes a simple Kubernetes design to the Catalog, then adds a custom KubeCon logo model to it.

---

LLMS index: [llms.txt](/pr-preview/pr-23457/llms.txt)

---

This is a 60-minute hands-on session for Meshery ContribFests. Every attendee first creates a very simple Kubernetes design in the Meshery Playground, publishes it to the Catalog, then builds a custom annotation-only model with the KubeCon logo and adds that component to the same design.

- **Where attendees work:** https://playground.meshery.io in a browser. Path A uses [Kanvas Designer](https://playground.meshery.io/extension/meshmap) and the [Registry](https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models). Path B uses **Configuration** > **Designs** and, for Lab 2, `mesheryctl`.
- **Outcome:** A saved design (a Kubernetes Deployment and Service) submitted to the Catalog, and a custom KubeCon logo model registered in the Playground and added to that design
- **Learning goals:** build and validate a Design; publish it to the Meshery Catalog; create or import an annotation-only Model; and incorporate one of that model's components into the design.
- **Audience:** Kubernetes users; no Meshery experience needed
- **Resources:** Meshery Models explainers ([architecture slides](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit?slide=id.g28d3a0f18b4_0_292#slide=id.g28d3a0f18b4_0_292), [Models and Components like Legos](https://meshery.io/blog/2024/03/2024-03-29-models-designs-legos/)).

## Choose your path

This workshop is a choose-your-own-adventure. Pick **Path A** or **Path B** at the start, and stay on that path for both labs. The outcome is the same: a published design, and a KubeCon logo model on that design.

<div class="tab-container tab-level-1"><input id="choose-workshop-path-10-1" type="radio" name="choose-workshop-path-10" checked>
    <label for="choose-workshop-path-10-1">Path A: Kanvas
    </label>
    <section class="tabbed" id="choose-workshop-path-10-content-1"><p><strong>Path A</strong> uses the <a href="https://docs.layer5.io/kanvas/">Kanvas</a> extension on <a href="https://playground.meshery.io">https://playground.meshery.io</a> for both labs.</p>
<ul>
<li>Lab 1: you build the Deployment and Service visually in <a href="https://playground.meshery.io/extension/meshmap">Kanvas Designer</a>, then Validate and Publish to Catalog.</li>
<li>Lab 2: you open the Registry&rsquo;s <strong>Create Model</strong> wizard in the Playground, mark the model as visual annotation only, and upload the KubeCon logos. You then pin the component in Kanvas Designer and drop it onto your Lab 1 design.</li>
</ul>
<p>You work in the browser. You do not need <code>mesheryctl</code> on Path A.</p>
</section><input id="choose-workshop-path-10-2" type="radio" name="choose-workshop-path-10">
    <label for="choose-workshop-path-10-2">Path B: Meshery UI and mesheryctl
    </label>
    <section class="tabbed" id="choose-workshop-path-10-content-2"><p><strong>Path B</strong> uses only Meshery UI and/or <code>mesheryctl</code>. It does not use the Kanvas extension.</p>
<ul>
<li>Lab 1: you import a ready-made Kubernetes manifest (a Deployment and a Service) from the Designs page or with <code>mesheryctl design import</code>, then Validate and publish from Meshery UI.</li>
<li>Lab 2: you scaffold, package, and import the KubeCon logo model with <code>mesheryctl</code> (or import the ready-made archive in the Registry), then add that component to your Lab 1 design by appending it to the exported design file and updating the design.</li>
</ul>
<p>You will point <code>mesheryctl</code> at the Playground before the Path B labs.</p>
</section></div>

## Prerequisites

### Attendees

- A laptop with a current Chrome or Safari.
- A GitHub or Google account, or an email address, to sign in to the Playground. _Sign-in is required: the Playground runs live clusters and every user needs a named account._
  - Optional but faster: sign in at https://playground.meshery.io before the session.
- Basic Kubernetes vocabulary: Deployment, Service, label, selector.
- **Path A:**
    1. a laptop or desktop with screen large enough for drag and drop in Kanvas. Tablets and phones are not enough.
- **Path B:**
    1. a laptop or desktop.
    1. [install `mesheryctl`](https://docs.meshery.io/installation/mesheryctl/) on your computer if you will use the CLI steps. You can complete Lab 1 with Meshery UI alone. Lab 2 is smoother with `mesheryctl`, and has a UI import fallback.
    1. [`jq`](https://jqlang.org/download/) for Lab 2, which uses it to turn each logo SVG into a JSON string. The Lab 2 CLI steps need a POSIX shell (macOS Terminal, Linux, or WSL / Git Bash on Windows). macOS 15 and later include `jq`; on Linux run `sudo apt install jq` or your distribution's equivalent.

<details>
<summary>Instructor prerequisites</summary>

- A Playground account that has completed both labs, on both paths, within 24 hours of the session, including a Catalog publish request. On Path A, walk through Create Model and Import Model in the Registry, then pin and drop the logo in Kanvas. On Path B, import the ready-made manifest and run `mesheryctl model import` into the Playground.
- Projector at 1080p or better, browser zoom at 125% so the back row can read the Kanvas form and the Registry wizard.
- This page open on a second screen, and the attendee steps (Labs 1 and 2, both paths) on a slide or QR code.
- A note to the Meshery Playground maintainers that a room of attendees will submit Catalog publish requests and import a small annotation model.

</details>

<details>
<summary>Instructor prep</summary>

1. **Dry run both labs on both paths.** Use a fresh design name such as `kubecon-instructor-design` and, on Path B, a fresh model name such as `kubecon-instructor-logo` so you hit the same screens attendees will. Note any UI label that differs from this page and tell the room at the start.
2. **Confirm the demo image still exists.** The Deployment uses `meshery/meshery-milestone:latest`. Docker Hub listed only the `latest` tag when this page was last checked on Nov 12, 2026.
3. **Confirm the logo files and the Path B manifest on this page download.** [Color SVG](kubecon-logo-color.svg), [white SVG](kubecon-logo-white.svg), the [ready-made model archive](kubecon-logo.tar.gz), and [kubecon-design.yaml](kubecon-design.yaml) should all open or download from the published docs page.
4. **Pre-build a finished design.** Keep a completed Lab 2 design open in a second tab. If the network or Playground stalls, you can keep teaching from it.
5. **Prepare the handle list.** Ask attendees to use their GitHub handle (or first name plus last initial) in every name, so dozens of designs and models never collide.
6. **Brief helpers.** Helpers own sign-in problems, path choice, Kanvas configuration, the Registry wizard, and `mesheryctl` login so the instructor keeps the room moving. Agree on a raised-hand or sticky-note signal.

</details>

## Welcome and sign-in

<details>
<summary>Instructor guide: welcome and sign-in</summary>

1. **Set the goal in one sentence.** "In 60 minutes you will publish a simple Kubernetes design to the Catalog, then add your own KubeCon logo model to it."
2. **Logistics.** Wi-Fi, the Playground URL, the helper signal, and the naming rule: put your handle in every name.
3. **Path choice.** Point at Path A (Kanvas, first tab) and Path B (Meshery UI and mesheryctl). Ask the room to pick one and stay on it for both labs.
4. **Live sign-in.** Open https://playground.meshery.io on the projector and sign in so the room sees the provider choice.
5. **Hands-on sign-in.** Helpers walk the room. Anyone not signed in by the end of this block pairs with a neighbor for Lab 1 and keeps trying in parallel.

</details>

### Attendees

1. **Open the Playground.** Go to https://playground.meshery.io. You are sent to the Meshery Cloud sign-in page.
2. **Sign in.** Continue with GitHub or Google, or create an account with email. You land back in Meshery.
3. **Check you are in.** You see the Meshery dashboard with your avatar in the top right. Raise your hand if you do not.
4. **Pick your path.** Stay on Path A or Path B for both labs. You can switch later if you need to, but you will save time by picking one now.

## Concepts

A design is something you build and can publish. A model is a package of components you can drop into a design. Today's custom model is annotation-only: it is a logo, not a workload.

| Term | Meaning |
|---|---|
| Design | A declarative document of the components you want, their configuration and their relationships. You can validate, dry run, deploy, share and publish it. |
| Catalog | The public library of published designs. Publishing submits your design for review so others can find and clone it. |
| Model | A versioned package that teaches Meshery about a set of components and their styling. Kubernetes is a model; today's KubeCon logo will be a model too. |
| Component | One kind of thing inside a model. In Lab 1 that is a Deployment or a Service. In Lab 2 it is a KubeCon logo (or, on Path A, the component registered with your model). |
| Annotation | A visual-only component. It appears on a design for diagramming and is not deployed to a cluster. |
| Kanvas | The visual designer in the Playground (Path A) where you drag components onto a canvas to build a design. |
| Registry | Where Meshery keeps every model it knows. You will create or import yours into it in Lab 2. |

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Key message</div>


Designs are saved to your account and survive. Deployments in the Playground clusters are reset daily. Today's real artifacts are the design you publish and the model you import, not a running workload.
</div>


## Lab 1: Create and publish your design

Attendees build an extremely simple design: one Kubernetes Deployment that runs the Meshery milestone Nginx image, and one Kubernetes Service that targets it. They save, validate, and publish the design to the Catalog.

<div class="tab-container tab-level-1"><input id="lab-1-paths-11-1" type="radio" name="lab-1-paths-11" checked>
    <label for="lab-1-paths-11-1">Path A: Kanvas
    </label>
    <section class="tabbed" id="lab-1-paths-11-content-1"><ol>
<li><strong>Open Kanvas Designer.</strong> Go to <a href="https://playground.meshery.io/extension/meshmap">https://playground.meshery.io/extension/meshmap</a>. Kanvas opens in Designer mode, the first of the two mode icons at the top center of the canvas, on a new design named <em>Untitled Design</em>.</li>
<li><strong>Create and name the design.</strong> Click the name <em>Untitled Design</em> at the top, select all of its text (clicking only places the cursor, so typing appends to the old name), and type <code>kubecon-&lt;your-handle&gt;-design</code>, for example <code>kubecon-octocat-design</code>. The save indicator beside the name is a cloud with a tick when the design is saved. Kanvas saves automatically about 3 seconds after you stop editing.</li>
<li><strong>Add a Deployment.</strong> From the Kubernetes model, drag a Deployment onto the canvas. Open its configuration form, give it a name, and set the container image to <code>meshery/meshery-milestone:latest</code>. That image is a small Nginx app; Docker Hub currently publishes it under the <code>latest</code> tag only. Add a label such as <code>app: kubecon-&lt;your-handle&gt;</code>.</li>
<li><strong>Add a Service.</strong> Drag a Service onto the canvas. In its configuration form, set the selector to the same <code>app: kubecon-&lt;your-handle&gt;</code> value so the Service targets the Deployment, and set port and targetPort to <code>80</code>.</li>
<li><strong>Connect them.</strong> Tap the Service, pick a relationship type from the edge picker, and drag to the Deployment. Kanvas highlights compatible targets.</li>
<li><strong>Validate.</strong> Run Validate and fix anything it reports. Validate checks every component against its schema and does not contact a cluster. The count in the result comes from the design, so it includes components whose definitions could not be fetched; a clean result does not guarantee every component was checked.</li>
<li><strong>Confirm it is saved.</strong> Check the save indicator shows a cloud with a tick.</li>
<li><strong>Publish to the Catalog.</strong> Open Details for the design: click the info (&ldquo;i&rdquo;) button, or choose Details from the File menu, or go to <strong>Configuration</strong> &gt; <strong>Designs</strong> and open the design&rsquo;s <strong>Design Information</strong>. Fill in Type, Technology, Description and Caveats (for example: Type &ldquo;Workload&rdquo;, Technology &ldquo;Kubernetes&rdquo;, a one-sentence description of your Deployment and Service). Click <strong>Publish to Catalog</strong>.</li>
</ol>
</section><input id="lab-1-paths-11-2" type="radio" name="lab-1-paths-11">
    <label for="lab-1-paths-11-2">Path B: Meshery UI and mesheryctl
    </label>
    <section class="tabbed" id="lab-1-paths-11-content-2"><p>Download the ready-made manifest first: <a href="/pr-preview/pr-23457/guides/tutorials/playground/kubecon-build-your-own-model-and-design/kubecon-design.yaml">kubecon-design.yaml</a>. It defines a Deployment that uses <code>meshery/meshery-milestone:latest</code> and a Service that selects it on port 80. You can edit the resource names and the <code>app: kubecon-demo</code> label to include your handle before you import, or leave them as they are and just give the design a unique name.</p>
<p>If you will use <code>mesheryctl</code> in this lab or in Lab 2, point it at the Playground once, before you import anything. Signing in to the Playground in your browser does not configure <code>mesheryctl</code>.</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl system context create playground --url https://playground.meshery.io --set
</span></span><span class="line"><span class="cl">mesheryctl system token create playground-token -f ~/.meshery/playground-auth.json --set
</span></span><span class="line"><span class="cl">mesheryctl system login
</span></span><span class="line"><span class="cl">chmod <span class="m">600</span> ~/.meshery/playground-auth.json
</span></span><span class="line"><span class="cl">mesheryctl system context view
</span></span></code></pre></div><p>The token file is separate so signing in does not overwrite the token your other contexts use. <code>mesheryctl system login</code> opens your browser to sign in; finish there and return to the terminal. The token is a credential, so <code>chmod 600</code> keeps it readable only by you (<code>mesheryctl</code> v1.0.76 and earlier write it world-readable). Confirm <code>mesheryctl system context view</code> shows the Playground endpoint before you import or delete anything.</p>
<h3 id="import-the-manifest-as-a-design">
	Import the manifest as a design
</h3>
<p>Use the UI, or <code>mesheryctl</code>, or both.</p>
<p><strong>Meshery UI</strong></p>
<ol>
<li>Go to <strong>Configuration</strong> &gt; <strong>Designs</strong>, or open <a href="https://playground.meshery.io/configuration/designs?view=table">https://playground.meshery.io/configuration/designs?view=table</a>. The page opens in grid view; the steps in this lab use the table view, so click <strong>Switch View</strong> (top right) if you see cards.</li>
<li>Click <strong>Import Design</strong>.</li>
<li>Enter the <strong>Design Name</strong> <code>kubecon-&lt;your-handle&gt;-design</code>, for example <code>kubecon-octocat-design</code>.</li>
<li>Choose <strong>File Upload</strong>, then pick <code>kubecon-design.yaml</code> from your computer.</li>
<li>Click <strong>Import</strong>. Importing does not deploy anything.</li>
</ol>
<p><strong>mesheryctl</strong></p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl design import -f kubecon-design.yaml -n kubecon-&lt;your-handle&gt;-design
</span></span></code></pre></div><p><code>-f</code> is a local path or a direct URL. <code>-n</code> names the design; without it, <code>mesheryctl</code> uses the file name. <code>-s</code> is optional because Meshery Server detects the format. If you pass it, use <code>-s &quot;Kubernetes Manifest&quot;</code> (the value is case-insensitive). Confirm the import:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl design list
</span></span><span class="line"><span class="cl">mesheryctl design view kubecon-&lt;your-handle&gt;-design
</span></span></code></pre></div><p><code>mesheryctl design view</code> prints YAML by default. Use <code>-o json</code> if you want JSON.</p>
<h3 id="optional-build-the-design-in-the-design-configurator">
	Optional: build the design in the Design Configurator
</h3>
<p>Instead of importing the manifest, you can create the design from forms in Meshery UI (still without Kanvas):</p>
<ol>
<li>Go to <strong>Configuration</strong> &gt; <strong>Designs</strong> and click <strong>Create Design</strong>. The Design Configurator opens at <code>/configuration/designs/configurator</code>.</li>
<li>Enter the <strong>Design Name</strong> <code>kubecon-&lt;your-handle&gt;-design</code>.</li>
<li>In <strong>Select Category</strong>, choose <em>Orchestration &amp; Management</em>. In <strong>Select Model</strong>, choose <em>Kubernetes</em>.</li>
<li>Expand the Deployment form, set the container image to <code>meshery/meshery-milestone:latest</code>, and add a label such as <code>app: kubecon-&lt;your-handle&gt;</code>.</li>
<li>Expand the Service form, set the selector to that same label, and set port and targetPort to <code>80</code>.</li>
<li>Click the <strong>Save Design as New File</strong> icon. The design appears on the Designs page.</li>
</ol>
<h3 id="validate-and-dry-run">
	Validate and dry run
</h3>
<p>On the Designs page in table view, open the action menu (⋮) for your design. In grid view the same actions are split across the card: <strong>Action</strong> holds Validate, Dry Run, Deploy and Evaluate; the pencil is Edit, the download icon is Download, and the info icon is Design Information.</p>
<ol>
<li>Choose <strong>Validate Design</strong>. This checks the design against component schemas and does not contact a cluster. Fix anything it reports, then validate again.</li>
<li>Optionally choose <strong>Dry Run</strong>. Dry Run talks to a cluster and is not required to finish the lab.</li>
</ol>
<p>There is no <code>mesheryctl design validate</code> or <code>mesheryctl design check</code> command. <code>mesheryctl design evaluate</code> runs relationship-evaluation policies and writes an evaluated design to a file; it is not a substitute for <strong>Validate Design</strong> and is not required in this lab.</p>
<h3 id="publish-to-the-catalog">
	Publish to the Catalog
</h3>
<p>The Designs table does not have a separate Publish action. Publish from Design Information:</p>
<ol>
<li>On <strong>Configuration</strong> &gt; <strong>Designs</strong>, open the action menu for your design and choose <strong>Design Information</strong>.</li>
<li>Fill in Type, Technology, Description and Caveats (for example: Type &ldquo;Workload&rdquo;, Technology &ldquo;Kubernetes&rdquo;, a one-sentence description of your Deployment and Service). Technology maps to the form&rsquo;s compatibility field; pick at least one.</li>
<li>Click <strong>Publish to Catalog</strong>. You must be signed in, you must own the design, and you need the Publish Design permission.</li>
</ol>
<p>There is no <code>mesheryctl</code> command that publishes a design to the Catalog.</p>
</section></div>

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">What 'published' means in this room</div>


If your account is an Organization Admin or Provider Admin, the design is published immediately, and the UI reports that the design was published. For a regular user, the request is queued for review (catalog request status `pending`). The UI reports that the design is queued for publishing and that maintainers have been notified. The design stays in My Designs, with its existing visibility (private or public), until an Organization Admin or Provider Admin approves or rejects the request. You will get an email when that happens. While the request is pending, you can keep editing the same design; the submitted request updates automatically, and you do not need to publish again. After a design is approved it is locked, its visibility becomes published, and it is listed in the Catalog. You would clone it to make a new version. For this lab, success is submitting the publish request, not waiting for approval.
</div>


<details>
<summary>Instructor guide: Lab 1</summary>

**Path A.** Open Kanvas Designer at https://playground.meshery.io/extension/meshmap and start a design named `kubecon-instructor-design`. Point at the save indicator. Drag a Deployment, set the image to `meshery/meshery-milestone:latest`, and apply a label such as `app: kubecon-instructor`. Drag a Service, set the selector to that same label, and set port and targetPort to 80. Draw an edge from the Service to the Deployment. Run Validate and show a saved state. Open Details and click Publish to Catalog.

**Path B.** Open https://playground.meshery.io/configuration/designs, import [kubecon-design.yaml](kubecon-design.yaml) with **Import Design**, or run `mesheryctl design import -f kubecon-design.yaml -n kubecon-instructor-design`. Show **Validate Design** in the action menu, then **Design Information** and **Publish to Catalog**.

Tell the room about a queued request versus an admin publish, and that submitting the request is the lab goal.

</details>

## Lab 2: Create a custom model and add it to your design

Attendees create an annotation-only model whose visual identity is the KubeCon logo, import or register that model in the Playground, and add it to the Lab 1 design.

Download these files from this page first:

- [KubeCon logo, color](kubecon-logo-color.svg) (`svgColor`, Path A: Logo (Light Theme))
- [KubeCon logo, white](kubecon-logo-white.svg) (`svgWhite`, Path A: Logo (Dark Theme))
- Fallback (and the working Path A model file): [ready-made model archive](kubecon-logo.tar.gz)

The logos are the official KubeCon + CloudNativeCon North America 2026 artwork from the [CNCF artwork repository](https://github.com/cncf/artwork).

<div class="tab-container tab-level-1"><input id="lab-2-paths-12-1" type="radio" name="lab-2-paths-12" checked>
    <label for="lab-2-paths-12-1">Path A: Kanvas
    </label>
    <section class="tabbed" id="lab-2-paths-12-content-1"><p>Path A builds the model in the Playground Registry, then adds the component in Kanvas Designer. You do not need <code>mesheryctl</code>.</p>
<h3 id="open-the-create-model-wizard">
	Open the Create Model wizard
</h3>
<ol>
<li>Open the Registry: click your avatar in the top right and choose <strong>Registry</strong>, or go to <strong>Settings</strong> &gt; <strong>Registry</strong>, or open <a href="https://playground.meshery.io/settings?settingsCategory=Registry&amp;tab=Models">https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models</a>.</li>
<li>On the <strong>Models</strong> tab, click <strong>Create Model</strong>. The modal title is <strong>Create Model</strong>.</li>
</ol>
<h3 id="walk-through-the-wizard">
	Walk through the wizard
</h3>
<p>The stepper labels are <strong>Model Details</strong>, <strong>Model Categorization</strong>, <strong>Styling</strong>, <strong>Source</strong>, <strong>Additional Details</strong>, and <strong>Finalize Generation</strong>. Use <strong>Next</strong> to advance. <strong>Generate</strong> runs on the summary step and registers the model immediately. The wizard has no separate Immediate Registration toggle.</p>
<ol>
<li><strong>Model Details.</strong> Enter a unique <strong>Model Name</strong> such as <code>kubecon-&lt;your-handle&gt;-logo</code> (lowercase letters, digits, and hyphens only, no spaces) and a <strong>Model Display Name</strong> such as <code>KubeCon Logo (&lt;your-handle&gt;)</code>.</li>
<li><strong>Model Categorization.</strong> Choose a <strong>Category</strong> and <strong>Subcategory</strong>. If nothing fits, choose <em>Uncategorized</em> for both, so the model is easy to find later.</li>
<li><strong>Styling.</strong> Upload <a href="/pr-preview/pr-23457/guides/tutorials/playground/kubecon-build-your-own-model-and-design/kubecon-logo-color.svg">kubecon-logo-color.svg</a> as <strong>Logo (Light Theme)</strong> and <a href="/pr-preview/pr-23457/guides/tutorials/playground/kubecon-build-your-own-model-and-design/kubecon-logo-white.svg">kubecon-logo-white.svg</a> as <strong>Logo (Dark Theme)</strong>. The inputs accept <code>.svg</code> files only (the browser reports <code>image/svg+xml</code>). Optionally set <strong>Primary Color</strong>, <strong>Secondary Color</strong>, and <strong>Shape</strong>. If you leave colors and shape empty, Meshery uses its defaults. Light Theme maps to <code>svgColor</code>; Dark Theme maps to <code>svgWhite</code>.</li>
<li><strong>Source.</strong> Choose <strong>GitHub</strong> or <strong>Artifact Hub</strong>, then enter a <strong>Model URL</strong>. This step is required: <strong>Next</strong> stays disabled until the source and a valid URL are set.</li>
<li><strong>Additional Details.</strong> Check <strong>The components in this model are visual annotations only.</strong> The wizard&rsquo;s help text and the generation summary call this option <strong>Visual Annotation Only</strong>.</li>
<li><strong>Finalize Generation.</strong> Review the summary, including <strong>Visual Annotation Only: Yes</strong>, then click <strong>Generate</strong>. When it finishes, click <strong>Finish</strong>.</li>
</ol>
<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">The wizard cannot create a component from an SVG alone</div>
<p>Create Model always sources components from a GitHub or Artifact Hub URL. On Generate, Meshery processes that URL, extracts component definitions (for example, CRDs), and registers the model. Your uploaded logos and colors style those sourced components. Checking Visual Annotation Only marks those sourced components as annotations; it does not invent a new component from the SVGs. There is no wizard field for a single custom component kind, and <strong>Source</strong> cannot be skipped.</p>
<p>Do not invent a dummy URL just to click Generate. That would either fail, or register unrelated components styled with the KubeCon logo, which is not the lab goal.</p>
<p>The working Path A step for a single KubeCon Logo component is <strong>Import Model</strong> with the ready-made archive, next.</p>
</div>
<h3 id="import-the-ready-made-kubecon-logo-model">
	Import the ready-made KubeCon logo model
</h3>
<ol>
<li>Stay on the Registry <strong>Models</strong> tab and click <strong>Import Model</strong>.</li>
<li>Choose <strong>File Import</strong>, pick <a href="/pr-preview/pr-23457/guides/tutorials/playground/kubecon-build-your-own-model-and-design/kubecon-logo.tar.gz">kubecon-logo.tar.gz</a> from your computer, and import it.</li>
<li>Confirm the model appears in the Registry. This fallback model is named <code>kubecon-logo</code> and is shared, so prefer a handle-specific model when you have one. The archive already has <code>isAnnotation</code> set to <code>true</code> and the SVGs inlined, with a single component whose kind and display name are <code>kubecon-logo</code> and <code>KubeCon Logo</code>.</li>
</ol>
<h3 id="add-the-logo-to-your-lab-1-design-in-kanvas">
	Add the logo to your Lab 1 design in Kanvas
</h3>
<ol>
<li><strong>Reopen your design</strong> in Kanvas Designer at <a href="https://playground.meshery.io/extension/meshmap">https://playground.meshery.io/extension/meshmap</a>.</li>
<li><strong>Find the component.</strong> Search the component list for <code>KubeCon Logo</code> or your model display name. Hover it and click the pin icon (Pin To Dock) so it sits in the dock at the bottom. If it does not appear, refresh the Kanvas tab and search again.</li>
<li><strong>Drop it on the canvas.</strong> Drag the KubeCon Logo component onto the design near your Deployment and Service. It is annotation-only, so it is a visual label, not a workload.</li>
<li><strong>Save.</strong> Wait for the cloud-with-a-tick save indicator.</li>
<li><strong>Optional: update the Catalog request.</strong> If your Lab 1 publish request is still pending review, this save is enough; the submitted design updates automatically. If the design was already approved and is locked in the Catalog, clone it, add the logo to the clone, and publish the new design.</li>
</ol>
</section><input id="lab-2-paths-12-2" type="radio" name="lab-2-paths-12">
    <label for="lab-2-paths-12-2">Path B: Meshery UI and mesheryctl
    </label>
    <section class="tabbed" id="lab-2-paths-12-content-2"><p>Path B builds the model with <code>mesheryctl</code> (or imports the ready-made archive in the Registry) and adds the component without Kanvas. If you did not already point <code>mesheryctl</code> at the Playground in Lab 1, do that first. Signing in to the Playground in your browser does not configure <code>mesheryctl</code>.</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl system context create playground --url https://playground.meshery.io --set
</span></span><span class="line"><span class="cl">mesheryctl system token create playground-token -f ~/.meshery/playground-auth.json --set
</span></span><span class="line"><span class="cl">mesheryctl system login
</span></span><span class="line"><span class="cl">chmod <span class="m">600</span> ~/.meshery/playground-auth.json
</span></span><span class="line"><span class="cl">mesheryctl system context view
</span></span></code></pre></div><p>Confirm <code>mesheryctl system context view</code> shows the Playground endpoint before you import or delete anything.</p>
<h3 id="scaffold-edit-package-and-import">
	Scaffold, edit, package and import
</h3>
<ol>
<li>
<p><strong>Scaffold the model.</strong> Pick a unique model name that includes your handle, using lowercase letters, digits and hyphens only:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl model init kubecon-&lt;your-handle&gt;-logo --version v0.1.0
</span></span></code></pre></div><p>That creates <code>kubecon-&lt;your-handle&gt;-logo/v0.1.0/</code>, with <code>model.json</code>, a sample file under <code>components/</code>, a sample <code>relationship.json</code> under <code>relationships/</code> (you remove it in step 3), and empty <code>connections/</code> and <code>credentials/</code> directories.</p>
</li>
<li>
<p><strong>Know what you are changing.</strong> You turn the scaffold into an annotation-only model with one component:</p>
<ul>
<li>In <code>model.json</code>: set <code>metadata.isAnnotation</code> to <code>true</code>, <code>registrant.kind</code> to <code>meshery</code>, a short <code>description</code>, <code>version</code> and <code>model.version</code> to <code>v0.1.0</code>, and put the color SVG in <code>metadata.svgColor</code> and the white SVG in <code>metadata.svgWhite</code>. Leave <code>schemaVersion</code> as <code>models.meshery.io/v1beta2</code> and the category as <em>Uncategorized</em>.</li>
<li>In the component: set <code>displayName</code> to <code>KubeCon Logo</code>, <code>component.kind</code> to <code>kubecon-logo</code>, <code>metadata.isAnnotation</code> to <code>true</code>, the same SVGs in <code>styles.svgColor</code> and <code>styles.svgWhite</code>, the nested <code>model.name</code> to your model name, and <code>schemaVersion</code> to <code>components.meshery.io/v1beta1</code>. The scaffold writes <code>components.meshery.io/v1beta3</code>, which the Playground&rsquo;s Meshery Server rejects, so this change is required (<a href="https://github.com/meshery/meshery/issues/20901">meshery#20901</a>).</li>
<li>Delete the sample relationship. Its <code>relationships.meshery.io/v1beta3</code> schema version is also rejected, and the model only needs the component.</li>
</ul>
<p>An SVG is full of double quotes and line breaks, so pasting it into a JSON string by hand breaks the file. Let <code>jq</code> do the escaping instead.</p>
</li>
<li>
<p><strong>Apply the edits with <code>jq</code>.</strong> Run this from the directory where you ran <code>model init</code>. Set <code>HANDLE</code> to your handle and <code>DL</code> to the folder holding the two downloaded SVGs:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl"><span class="o">(</span>
</span></span><span class="line"><span class="cl"><span class="nv">HANDLE</span><span class="o">=</span>&lt;your-handle&gt;
</span></span><span class="line"><span class="cl"><span class="nv">DL</span><span class="o">=</span>~/Downloads
</span></span><span class="line"><span class="cl"><span class="nv">MODEL</span><span class="o">=</span>kubecon-<span class="nv">$HANDLE</span>-logo
</span></span><span class="line"><span class="cl"><span class="nb">cd</span> <span class="s2">&#34;</span><span class="nv">$MODEL</span><span class="s2">/v0.1.0&#34;</span> <span class="o">||</span> <span class="o">{</span> <span class="nb">echo</span> <span class="s2">&#34;model directory not found: </span><span class="nv">$MODEL</span><span class="s2">/v0.1.0&#34;</span><span class="p">;</span> <span class="nb">exit</span> 1<span class="p">;</span> <span class="o">}</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="nb">test</span> -r <span class="s2">&#34;</span><span class="nv">$DL</span><span class="s2">/kubecon-logo-color.svg&#34;</span> <span class="o">&amp;&amp;</span> <span class="nb">test</span> -r <span class="s2">&#34;</span><span class="nv">$DL</span><span class="s2">/kubecon-logo-white.svg&#34;</span> <span class="o">||</span> <span class="o">{</span> <span class="nb">echo</span> <span class="s2">&#34;SVG files not found in </span><span class="nv">$DL</span><span class="s2">&#34;</span><span class="p">;</span> <span class="nb">exit</span> 1<span class="p">;</span> <span class="o">}</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl"><span class="c1"># Remove only the XML declaration; keep any markup on the same line.</span>
</span></span><span class="line"><span class="cl"><span class="nv">COLOR</span><span class="o">=</span><span class="k">$(</span>sed <span class="s1">&#39;s/&lt;?xml[^&gt;]*?&gt;//&#39;</span> <span class="s2">&#34;</span><span class="nv">$DL</span><span class="s2">/kubecon-logo-color.svg&#34;</span><span class="k">)</span>
</span></span><span class="line"><span class="cl"><span class="nv">WHITE</span><span class="o">=</span><span class="k">$(</span>sed <span class="s1">&#39;s/&lt;?xml[^&gt;]*?&gt;//&#39;</span> <span class="s2">&#34;</span><span class="nv">$DL</span><span class="s2">/kubecon-logo-white.svg&#34;</span><span class="k">)</span>
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">jq --arg c <span class="s2">&#34;</span><span class="nv">$COLOR</span><span class="s2">&#34;</span> --arg w <span class="s2">&#34;</span><span class="nv">$WHITE</span><span class="s2">&#34;</span> <span class="s1">&#39;
</span></span></span><span class="line"><span class="cl"><span class="s1">  .version = &#34;v0.1.0&#34; | .model.version = &#34;v0.1.0&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .description = &#34;KubeCon logo annotation model&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .registrant.kind = &#34;meshery&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .metadata.isAnnotation = true
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .metadata.svgColor = $c | .metadata.svgWhite = $w
</span></span></span><span class="line"><span class="cl"><span class="s1">&#39;</span> model.json &gt; model.json.tmp <span class="o">&amp;&amp;</span> mv model.json.tmp model.json
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">jq --arg c <span class="s2">&#34;</span><span class="nv">$COLOR</span><span class="s2">&#34;</span> --arg w <span class="s2">&#34;</span><span class="nv">$WHITE</span><span class="s2">&#34;</span> --arg name <span class="s2">&#34;</span><span class="nv">$MODEL</span><span class="s2">&#34;</span> <span class="s1">&#39;
</span></span></span><span class="line"><span class="cl"><span class="s1">  .schemaVersion = &#34;components.meshery.io/v1beta1&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .version = &#34;v0.1.0&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .displayName = &#34;KubeCon Logo&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .component.kind = &#34;kubecon-logo&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .component.version = &#34;v0.1.0&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .metadata.isAnnotation = true
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .styles.svgColor = $c | .styles.svgWhite = $w
</span></span></span><span class="line"><span class="cl"><span class="s1">  | .model.name = $name | .model.version = &#34;v0.1.0&#34; | .model.model.version = &#34;v0.1.0&#34;
</span></span></span><span class="line"><span class="cl"><span class="s1">&#39;</span> components/component.json &gt; components/kubecon-logo.json
</span></span><span class="line"><span class="cl">
</span></span><span class="line"><span class="cl">rm components/component.json relationships/relationship.json
</span></span><span class="line"><span class="cl"><span class="nb">cd</span> ../..
</span></span><span class="line"><span class="cl"><span class="o">)</span>
</span></span></code></pre></div><p>You now have <code>model.json</code> and a single <code>components/kubecon-logo.json</code>. If you prefer to edit by hand, make the same changes in your editor. Generate each SVG value with <code>sed 's/&lt;?xml[^&gt;]*?&gt;//' kubecon-logo-color.svg | jq -Rs .</code> (and the same for the white SVG), then paste the whole output, including its surrounding quotes.</p>
</li>
<li>
<p><strong>Package the model as an OCI archive:</strong></p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl model build kubecon-&lt;your-handle&gt;-logo/v0.1.0
</span></span></code></pre></div><p>This writes <code>kubecon-&lt;your-handle&gt;-logo-v0-1-0.tar</code> in your current directory.</p>
</li>
<li>
<p><strong>Import it into the Playground:</strong></p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl model import -f kubecon-&lt;your-handle&gt;-logo-v0-1-0.tar
</span></span><span class="line"><span class="cl">mesheryctl model view kubecon-&lt;your-handle&gt;-logo
</span></span></code></pre></div><p>You can also import from the UI: avatar, then <strong>Registry</strong>, then <strong>Import Model</strong>, then <strong>File Import</strong>, and choose the <code>.tar</code> file. Either route registers the model in the shared Playground instance.</p>
</li>
</ol>
<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Fallback if init or build fails</div>
<p>Download the <a href="/pr-preview/pr-23457/guides/tutorials/playground/kubecon-build-your-own-model-and-design/kubecon-logo.tar.gz">ready-made model archive</a> from this page and import it with <code>mesheryctl model import -f kubecon-logo.tar.gz</code>, or use Import Model in the Registry. This fallback model is named <code>kubecon-logo</code> and is shared, so prefer the handle-specific path above when you can. The archive is a model directory with <code>model.json</code> and one annotation component; <code>isAnnotation</code> is already set to <code>true</code> and the SVGs are already inlined.</p>
</div>
<h3 id="add-the-logo-to-your-lab-1-design-without-kanvas">
	Add the logo to your Lab 1 design (without Kanvas)
</h3>
<p>Without Kanvas you will not see the logo on a canvas. You add the annotation component to the design document, update the design, and confirm it with <code>mesheryctl design view</code>.</p>
<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">The Design Configurator cannot add an annotation component</div>
<p>The Design Configurator lists your model under <em>Uncategorized</em>, but shows the component as <strong>KubeCon Logo (contains invalid schema)</strong> and cannot expand it. An annotation component has no configuration schema (<code>component.schema</code> is empty), so there is no form to fill. Add it with the export-and-update steps below instead.</p>
</div>
<p><strong>Export, add the component, and update</strong></p>
<ol>
<li>
<p>Export the design. In Meshery UI, on the Designs page, choose <strong>Download</strong> and pick <strong>Meshery Design (yaml)</strong>. Or run:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl design <span class="nb">export</span> kubecon-&lt;your-handle&gt;-design
</span></span></code></pre></div><p>The default <code>--type</code> is <code>current</code> (the design as Meshery stores it now). The file is written as <code>kubecon-&lt;your-handle&gt;-design_&lt;id&gt;.yaml</code> in your current directory.</p>
</li>
<li>
<p>Open the exported YAML. It should have a top-level <code>components:</code> list and a <code>schemaVersion</code> such as <code>designs.meshery.io/v1beta3</code>. Append this component to the end of that top-level list (not to a nested <code>components:</code> key inside a model), indented like the existing entries. Replace the model <code>name</code> with your handle-specific model name, or use <code>kubecon-logo</code> if you imported the fallback archive. Generate a fresh UUID for <code>id</code> (for example with <code>uuidgen</code>):</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-yaml" data-lang="yaml"><span class="line"><span class="cl">- <span class="nt">id</span><span class="p">:</span><span class="w"> </span><span class="m">00000000-0000-0000-0000-000000000000</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">schemaVersion</span><span class="p">:</span><span class="w"> </span><span class="l">components.meshery.io/v1beta1</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">version</span><span class="p">:</span><span class="w"> </span><span class="l">v0.1.0</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">displayName</span><span class="p">:</span><span class="w"> </span><span class="l">KubeCon Logo</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">format</span><span class="p">:</span><span class="w"> </span><span class="l">JSON</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">status</span><span class="p">:</span><span class="w"> </span><span class="l">enabled</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">model</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">schemaVersion</span><span class="p">:</span><span class="w"> </span><span class="l">models.meshery.io/v1beta2</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">version</span><span class="p">:</span><span class="w"> </span><span class="l">v0.1.0</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">name</span><span class="p">:</span><span class="w"> </span><span class="l">kubecon-&lt;your-handle&gt;-logo</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">displayName</span><span class="p">:</span><span class="w"> </span><span class="l">KubeCon Logo</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">status</span><span class="p">:</span><span class="w"> </span><span class="l">enabled</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">metadata</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">isAnnotation</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">model</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">      </span><span class="nt">version</span><span class="p">:</span><span class="w"> </span><span class="l">v0.1.0</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">metadata</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">isAnnotation</span><span class="p">:</span><span class="w"> </span><span class="kc">true</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">component</span><span class="p">:</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">kind</span><span class="p">:</span><span class="w"> </span><span class="l">kubecon-logo</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">version</span><span class="p">:</span><span class="w"> </span><span class="l">v0.1.0</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">    </span><span class="nt">schema</span><span class="p">:</span><span class="w"> </span><span class="s2">&#34;&#34;</span><span class="w">
</span></span></span><span class="line"><span class="cl"><span class="w">  </span><span class="nt">configuration</span><span class="p">:</span><span class="w"> </span>{}<span class="w">
</span></span></span></code></pre></div><p>This is the v1beta1 component shape the ready-made <code>kubecon-logo</code> model uses. Meshery keeps the component as written: the import does not fill in the model&rsquo;s <code>id</code> or the SVG <code>styles</code> from the Registry.</p>
</li>
<li>
<p>Update the existing design: on the Designs page click its Edit (pencil) action to open the Design Configurator, replace the YAML in the editor with your edited file, and click the <strong>Update Design</strong> icon (not <strong>Save Design as New File</strong>, unless you want a copy).</p>
<p><code>mesheryctl design import</code> always creates a <strong>new</strong> design (it posts to <code>/api/pattern/import</code> and does not update an existing ID). If you re-import the edited file, give it a new name, for example:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl design import -f kubecon-&lt;your-handle&gt;-design-with-logo.yaml -s design -n kubecon-&lt;your-handle&gt;-design-with-logo
</span></span></code></pre></div></li>
<li>
<p>Confirm the logo component is in the design:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-sh" data-lang="sh"><span class="line"><span class="cl">mesheryctl design view kubecon-&lt;your-handle&gt;-design
</span></span></code></pre></div></li>
<li>
<p><strong>Optional: update the Catalog request.</strong> If your Lab 1 publish request is still pending review, updating the same design is enough; the submitted request updates automatically. If you created a new design by re-importing, publish that new design. If the original design was already approved and is locked, clone it, add the logo to the clone, and publish the new design.</p>
</li>
</ol>
</section></div>

<details>
<summary>Instructor guide: Lab 2</summary>

**Path A.** Show the logo downloads on this page. Open the Registry at https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models, walk through **Create Model** far enough to show Model Details, Styling (the two SVG uploads), Visual Annotation Only, and the required Source step. Tell the room, honestly, that Generate sources components from a GitHub or Artifact Hub URL and cannot create a single logo component from the SVGs alone. Switch to **Import Model** and import `kubecon-logo.tar.gz`. In Kanvas Designer, pin the component, drop it on the Lab 1 design, and save.

**Path B.** Run `mesheryctl model init`, walk through the field changes in step 2, run the `jq` edits in step 3, then `model build` and `model import -f` against the Playground context. If the CLI path is slow for the room, switch early to the ready-made `kubecon-logo.tar.gz` and Import Model in the Registry. Show the export, the appended YAML snippet, and **Update Design** in the Design Configurator. Confirm with `mesheryctl design view`.

Remind people that a pending Catalog request on the same design picks up the edit automatically.

</details>

## Wrap-up and next steps

1. **Recap.** A Deployment and Service design, submitted to the Catalog; a custom annotation-only KubeCon logo model registered in the Playground and added to that design. Deployments in the Playground reset daily; designs and models you own do not.
2. **One quick share if time.** If the labs finished early, one volunteer shows their design. On Path A, show the Service-to-Deployment edge and the logo on the canvas. On Path B, show the design on the Designs page and the logo component in **Design Information** or `mesheryctl design view`.
3. **Next steps.** Contribute a model upstream, try the stretch goals at home, join the community Slack at https://slack.layer5.io, and keep using your Playground account.

## What success looks like

Each attendee can show all of the following before leaving.

- Signed in to the Playground with their own account
- A design named with their handle contains a Deployment using `meshery/meshery-milestone:latest` and a Service that selects it on port 80
- Validate has been run (Kanvas Validate on Path A, or **Validate Design** on the Designs page on Path B), and the design is saved
- A Catalog publish request has been submitted (the design is published, or the request is queued / pending review)
- A KubeCon logo model is listed in the Registry (their handle-specific model, or the shared fallback)
- The logo component has been added to the Lab 1 design, and the design has been saved or updated again (visible on the Kanvas canvas on Path A, or in the design file / `mesheryctl design view` on Path B)

<details>
<summary>Instructor and helper notes</summary>

Helpers can use this list as a walk-around check during the last few minutes.

Room-level target: at least 80% of attendees complete both labs. Track it with a show of hands after Lab 1 and after Lab 2. If useful, also ask how many are on Path A versus Path B.

</details>

## Troubleshooting and fallbacks

| Symptom | Fix |
|---|---|
| Stuck on the sign-in page | Try another provider (GitHub, Google or email). If the venue blocks it, pair with a neighbor and keep trying on a phone hotspot. |
| Model name rejected | Use lowercase letters, digits, and hyphens only, no spaces or underscores. |
| Create Model will not leave the Source step | The wizard requires a GitHub or Artifact Hub URL and generates components from that URL. It cannot create a component from an SVG alone. Use **Import Model** with [kubecon-logo.tar.gz](kubecon-logo.tar.gz) instead. |
| `mesheryctl` talks to the wrong server | Run `mesheryctl system context view` and confirm the endpoint is `https://playground.meshery.io` before import or delete. Recreate the context and token if it is not. |
| `model import` says *No valid component or relationship found* | The component still has the scaffold's `schemaVersion: components.meshery.io/v1beta3`. Change it to `components.meshery.io/v1beta1`, then run `model build` and `model import` again. |
| `model import` reports *encountered error for 1 entity (1 relationship)* | The scaffold's sample `relationships/relationship.json` is still there. Delete it, rebuild and re-import. The logo component itself imported. |
| `model build` or `model import` fails | Open the Notification Center or the CLI error. Check that `metadata.isAnnotation` is `true`, that each SVG is a JSON-escaped string (see step 3), and that the component `schemaVersion` is `components.meshery.io/v1beta1`. If it stays stuck, import the ready-made [kubecon-logo.tar.gz](kubecon-logo.tar.gz). |
| `jq: command not found` | Install `jq` (see Prerequisites), or use the ready-made [kubecon-logo.tar.gz](kubecon-logo.tar.gz). |
| Model shows version `v0.0.1` | `model init` from `mesheryctl` v1.0.76 and earlier does not write `--version` into the files. Setting `version` and `model.version` in step 3 fixes it. |
| Model not in the Kanvas component list | Refresh the Kanvas tab, then search by display name (`KubeCon Logo`). If you used Import Model, confirm the model is listed in the Registry first. |
| KubeCon Logo shows "(contains invalid schema)" in the Design Configurator | Expected: an annotation component has no configuration form. Use the export-and-update path in Lab 2, Path B. |
| `design import` says "Invalid design source type" | Use `-s design` (or `kubernetes manifest`, `helm chart`, `docker compose`), or leave `-s` out and let Meshery Server detect the format. |
| Validate Design or Design Information is not in the menu | You are in grid view. Click **Switch View**, or use the card's **Action** menu and info icon. |
| Service does not target the Deployment | The Service selector must match a label on the Deployment, for example `app: kubecon-<your-handle>` or `app: kubecon-demo` if you imported the ready-made manifest unchanged. Port and targetPort should both be `80`. |
| Publish to Catalog is disabled | Sign in again and reopen **Design Information**. Publish requires a signed-in account that owns the design and has the Publish Design permission. |
| Catalog request still says it is queued or pending | That is expected for a regular user. The lab is complete once the request is submitted. Keep editing if you are still in Lab 2; a pending request on the same design updates automatically. |
| Name already taken | Someone else used it. Add a suffix such as `-two`. |
| Save indicator shows Unable to save | Stop editing until it clears. Usually a dropped connection or an expired session; sign in again. |
| Playground slow for the whole room | Attendees who already have a design move to publish or to Lab 2. If model import is slow, use the ready-made archive or pin any already-imported `kubecon-logo` model. |

<details>
<summary>Instructor notes: fallbacks</summary>

- If the Playground is slow for the whole room, continue teaching from the pre-built design.
- The UI steps come from current Meshery, Kanvas, and Registry source. The instructor's morning dry run is the check; update the session if a label differs.
- Path A attendees who cannot use Kanvas can switch to Path B mid-session. Path B attendees who want a canvas can switch to Path A after the model is imported and pick up at "Add the logo to your Lab 1 design in Kanvas."

</details>

## Stretch goals

These are optional. Use them if you finish a lab early, or after the session.

- Finish the Create Model wizard with a real CRD source, such as `git://github.com/cert-manager/cert-manager/master/deploy/crds`, to see how Generate builds components from a URL. This is separate from the KubeCon logo component.
- Dry run or deploy only the Deployment and Service, then look at them in Operator mode (Path A) or from the Designs page **Dry Run** / **Deploy** actions (Path B). The Playground clusters reset daily.
- Export the design to keep a file copy. On Path A, open Details from the File menu to share it with a neighbor. On Path B, use **Download** on the Designs page or `mesheryctl design export`.
- Path A: wrap the Deployment and Service in a Namespace and let Kanvas create the parent-child relationship.

## Related docs

- [Creating a Meshery Design](/pr-preview/pr-23457/guides/configuration-management/creating-a-meshery-design/)
- [Importing and Exporting Designs](/pr-preview/pr-23457/guides/configuration-management/import-export-designs/)
- [Catalog](/pr-preview/pr-23457/concepts/architecture/catalog/)
- [Publishing Designs](https://docs.layer5.io/kanvas/designer/publishing-designs/)
- [Create Models](/pr-preview/pr-23457/guides/configuration-management/creating-models/)
- [Importing Models](/pr-preview/pr-23457/guides/configuration-management/importing-models/)
- [mesheryctl design import](/pr-preview/pr-23457/reference/mesheryctl/design/import/)
- [mesheryctl design list](/pr-preview/pr-23457/reference/mesheryctl/design/list/)
- [mesheryctl design view](/pr-preview/pr-23457/reference/mesheryctl/design/view/)
- [mesheryctl design export](/pr-preview/pr-23457/reference/mesheryctl/design/export/)
- [mesheryctl model init](/pr-preview/pr-23457/reference/mesheryctl/model/init/)
- [mesheryctl model build](/pr-preview/pr-23457/reference/mesheryctl/model/build/)
- [mesheryctl model import](/pr-preview/pr-23457/reference/mesheryctl/model/import/)
- [Using Meshery Playground](/pr-preview/pr-23457/installation/playground/)
- Kanvas: [Getting started](https://docs.layer5.io/kanvas/getting-started/), [Creating relationships](https://docs.layer5.io/kanvas/getting-started/creating-relationships/), [Pin a model to the dock](https://docs.layer5.io/kanvas/designer/pin-model-to-dock/), [Saving designs](https://docs.layer5.io/kanvas/designer/saving-designs/)
- [Contributor Training](https://layer5.io/learn/workshops)
