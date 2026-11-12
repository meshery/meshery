---
title: "KubeCon Workshop: Publish a Design and Build a Custom Model"
description: A 60-minute Meshery Playground workshop. Choose Path A (Kanvas) or Path B (Meshery UI and mesheryctl). Every attendee publishes a simple Kubernetes design to the Catalog, then adds a custom KubeCon logo model to it.
model: playground
params:
    kind: workshop
categories: [tutorials]
weight: 10
aliases:
- /guides/tutorials/kubecon-build-your-own-model-and-design
---

This is a 60-minute hands-on session for Meshery ContribFests. Every attendee first creates a very simple Kubernetes design in the Meshery Playground, publishes it to the Catalog, then builds a custom annotation-only model with the KubeCon logo and adds that component to the same design.

- **Where attendees work:** https://playground.meshery.io in a browser. Path A uses [Kanvas Designer](https://playground.meshery.io/extension/meshmap) and the [Registry](https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models). Path B uses **Configuration** > **Designs** and, for Lab 2, `mesheryctl`.
- **Outcome:** A saved design (a Kubernetes Deployment and Service) submitted to the Catalog, and a custom KubeCon logo model registered in the Playground and added to that design
- **Learning goals:** build and validate a Design; publish it to the Meshery Catalog; create or import an annotation-only Model; and incorporate one of that model's components into the design.
- **Audience:** Kubernetes users; no Meshery experience needed
- **Resources:** Meshery Models explainers ([architecture slides](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit?slide=id.g28d3a0f18b4_0_292#slide=id.g28d3a0f18b4_0_292), [Models and Components like Legos](https://meshery.io/blog/2024/03/2024-03-29-models-designs-legos/)).

## Choose your path

This workshop is a choose-your-own-adventure. Pick **Path A** or **Path B** at the start, and stay on that path for both labs. The outcome is the same: a published design, and a KubeCon logo model on that design.

{{< tabs id="choose-workshop-path" >}}
Path A: Kanvas

**Path A** uses the [Kanvas](https://docs.layer5.io/kanvas/) extension on https://playground.meshery.io for both labs.

- Lab 1: you build the Deployment and Service visually in [Kanvas Designer](https://playground.meshery.io/extension/meshmap), then Validate and Publish To Catalog.
- Lab 2: you open the Registry's **Create Model** wizard in the Playground, mark the model as visual annotation only, and upload the KubeCon logos. You then pin the component in Kanvas Designer and drop it onto your Lab 1 design.

You work in the browser. You do not need `mesheryctl` on Path A.
<!-- tab -->
Path B: Meshery UI and mesheryctl

**Path B** uses only Meshery UI and/or `mesheryctl`. It does not use the Kanvas extension.

- Lab 1: you import a ready-made Kubernetes manifest (a Deployment and a Service) from the Designs page or with `mesheryctl design import`, then Validate and publish from Meshery UI.
- Lab 2: you scaffold, package, and import the KubeCon logo model with `mesheryctl` (or import the ready-made archive in the Registry), then add that component to your Lab 1 design by editing the design in the Design Configurator or by updating the design file.

You will point `mesheryctl` at the Playground before the Path B labs.
{{< /tabs >}}

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

{{% alert color="info" title="Key message" %}}
Designs are saved to your account and survive. Deployments in the Playground clusters are reset daily. Today's real artifacts are the design you publish and the model you import, not a running workload.
{{% /alert %}}

## Lab 1: Create and publish your design

Attendees build an extremely simple design: one Kubernetes Deployment that runs the Meshery milestone Nginx image, and one Kubernetes Service that targets it. They save, validate, and publish the design to the Catalog.

{{< tabs id="lab-1-paths" >}}
Path A: Kanvas

1. **Open Kanvas Designer.** Go to https://playground.meshery.io/extension/meshmap. Kanvas opens in Designer mode, the first of the two mode icons at the top center of the canvas, on a new design named *Untitled Design*.
2. **Create and name the design.** Click the name *Untitled Design* at the top, select all of its text (clicking only places the cursor, so typing appends to the old name), and type `kubecon-<your-handle>-design`, for example `kubecon-octocat-design`. The save indicator beside the name is a cloud with a tick when the design is saved. Kanvas saves automatically about 3 seconds after you stop editing.
3. **Add a Deployment.** From the Kubernetes model, drag a Deployment onto the canvas. Open its configuration form, give it a name, and set the container image to `meshery/meshery-milestone:latest`. That image is a small Nginx app; Docker Hub currently publishes it under the `latest` tag only. Add a label such as `app: kubecon-<your-handle>`.
4. **Add a Service.** Drag a Service onto the canvas. In its configuration form, set the selector to the same `app: kubecon-<your-handle>` value so the Service targets the Deployment, and set port and targetPort to `80`.
5. **Connect them.** Tap the Service, pick a relationship type from the edge picker, and drag to the Deployment. Kanvas highlights compatible targets.
6. **Validate.** Run Validate and fix anything it reports. Validate checks every component against its schema and does not contact a cluster. The count in the result comes from the design, so it includes components whose definitions could not be fetched; a clean result does not guarantee every component was checked.
7. **Confirm it is saved.** Check the save indicator shows a cloud with a tick.
8. **Publish to the Catalog.** Open Details for the design: click the info ("i") button, or choose Details from the File menu, or go to **Configuration** > **Designs** and open the design's **Design Information**. Fill in Type, Technology, Description and Caveats (for example: Type "Workload", Technology "Kubernetes", a one-sentence description of your Deployment and Service). Click **Publish to Catalog**.
<!-- tab -->
Path B: Meshery UI and mesheryctl

Download the ready-made manifest first: [kubecon-design.yaml](kubecon-design.yaml). It defines a Deployment that uses `meshery/meshery-milestone:latest` and a Service that selects it on port 80. You can edit the resource names and the `app: kubecon-demo` label to include your handle before you import, or leave them as they are and just give the design a unique name.

If you will use `mesheryctl` in this lab or in Lab 2, point it at the Playground once, before you import anything. Signing in to the Playground in your browser does not configure `mesheryctl`.

```sh
mesheryctl system context create playground --url https://playground.meshery.io --set
mesheryctl system token create playground-token -f ~/.meshery/playground-auth.json --set
mesheryctl system login
chmod 600 ~/.meshery/playground-auth.json
mesheryctl system context view
```

The token file is separate so signing in does not overwrite the token your other contexts use. Confirm `mesheryctl system context view` shows the Playground endpoint before you import or delete anything.

### Import the manifest as a design

Use the UI, or `mesheryctl`, or both.

**Meshery UI**

1. Go to **Configuration** > **Designs**, or open https://playground.meshery.io/configuration/designs.
2. Click **Import Design**.
3. Enter the **Design Name** `kubecon-<your-handle>-design`, for example `kubecon-octocat-design`.
4. Choose **File Upload**, then pick `kubecon-design.yaml` from your computer.
5. Click **Import**. Importing does not deploy anything.

**mesheryctl**

```sh
mesheryctl design import -f kubecon-design.yaml -n kubecon-<your-handle>-design
```

`-f` is a local path or a direct URL. `-n` names the design; without it, `mesheryctl` uses the file name. `-s` is optional because Meshery Server detects the format. If you pass it, use `-s "Kubernetes Manifest"` (the value is case-insensitive). Confirm the import:

```sh
mesheryctl design list
mesheryctl design view kubecon-<your-handle>-design
```

`mesheryctl design view` prints YAML by default. Use `-o json` if you want JSON.

### Optional: build the design in the Design Configurator

Instead of importing the manifest, you can create the design from forms in Meshery UI (still without Kanvas):

1. Go to **Configuration** > **Designs** and click **Create Design**. The Design Configurator opens at `/configuration/designs/configurator`.
2. Enter the **Design Name** `kubecon-<your-handle>-design`.
3. In **Select Category**, choose *Orchestration & Management*. In **Select Model**, choose *Kubernetes*.
4. Expand the Deployment form, set the container image to `meshery/meshery-milestone:latest`, and add a label such as `app: kubecon-<your-handle>`.
5. Expand the Service form, set the selector to that same label, and set port and targetPort to `80`.
6. Click the **Save Design as New File** icon. The design appears on the Designs page.

### Validate and dry run

On the Designs page, open the action menu for your design:

1. Choose **Validate Design**. This checks the design against component schemas and does not contact a cluster. Fix anything it reports, then validate again.
2. Optionally choose **Dry Run**. Dry Run talks to a cluster and is not required to finish the lab.

There is no `mesheryctl design validate` or `mesheryctl design check` command. `mesheryctl design evaluate` runs relationship-evaluation policies and writes an evaluated design to a file; it is not a substitute for **Validate Design** and is not required in this lab.

### Publish to the Catalog

The Designs table does not have a separate Publish action. Publish from Design Information:

1. On **Configuration** > **Designs**, open the action menu for your design and choose **Design Information**.
2. Fill in Type, Technology, Description and Caveats (for example: Type "Workload", Technology "Kubernetes", a one-sentence description of your Deployment and Service). Technology maps to the form's compatibility field; pick at least one.
3. Click **Publish to Catalog**. You must be signed in, you must own the design, and you need the Publish Design permission.

There is no `mesheryctl` command that publishes a design to the Catalog.
{{< /tabs >}}

{{% alert color="info" title="What 'published' means in this room" %}}
If your account is an Organization Admin or Provider Admin, the design is published immediately, and the UI reports that the design was published. For a regular user, the request is queued for review (catalog request status `pending`). The UI reports that the design is queued for publishing and that maintainers have been notified. The design stays in My Designs, with its existing visibility (private or public), until an Organization Admin or Provider Admin approves or rejects the request. You will get an email when that happens. While the request is pending, you can keep editing the same design; the submitted request updates automatically, and you do not need to publish again. After a design is approved it is locked, its visibility becomes published, and it is listed in the Catalog. You would clone it to make a new version. For this lab, success is submitting the publish request, not waiting for approval.
{{% /alert %}}

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

{{< tabs id="lab-2-paths" >}}
Path A: Kanvas

Path A builds the model in the Playground Registry, then adds the component in Kanvas Designer. You do not need `mesheryctl`.

### Open the Create Model wizard

1. Open the Registry: click your avatar in the top right and choose **Registry**, or go to **Settings** > **Registry**, or open https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models.
2. On the **Models** tab, click **Create Model**. The modal title is **Create Model**.

### Walk through the wizard

The stepper labels are **Model Details**, **Model Categorization**, **Styling**, **Source**, **Additional Details**, and **Finalize Generation**. Use **Next** to advance. **Generate** runs on the summary step and registers the model immediately. The wizard has no separate Immediate Registration toggle.

1. **Model Details.** Enter a unique **Model Name** such as `kubecon-<your-handle>-logo` (lowercase letters, digits, and hyphens only, no spaces) and a **Model Display Name** such as `KubeCon Logo (<your-handle>)`.
2. **Model Categorization.** Choose a **Category** and **Subcategory**. If nothing fits, choose *Uncategorized* for both, so the model is easy to find later.
3. **Styling.** Upload [kubecon-logo-color.svg](kubecon-logo-color.svg) as **Logo (Light Theme)** and [kubecon-logo-white.svg](kubecon-logo-white.svg) as **Logo (Dark Theme)**. The inputs accept `.svg` files only (the browser reports `image/svg+xml`). Optionally set **Primary Color**, **Secondary Color**, and **Shape**. If you leave colors and shape empty, Meshery uses its defaults. Light Theme maps to `svgColor`; Dark Theme maps to `svgWhite`.
4. **Source.** Choose **GitHub** or **Artifact Hub**, then enter a **Model URL**. This step is required: **Next** stays disabled until the source and a valid URL are set.
5. **Additional Details.** Check **The components in this model are visual annotations only.** The wizard's help text and the generation summary call this option **Visual Annotation Only**.
6. **Finalize Generation.** Review the summary, including **Visual Annotation Only: Yes**, then click **Generate**. When it finishes, click **Finish**.

{{% alert color="warning" title="The wizard cannot create a component from an SVG alone" %}}
Create Model always sources components from a GitHub or Artifact Hub URL. On Generate, Meshery processes that URL, extracts component definitions (for example, CRDs), and registers the model. Your uploaded logos and colors style those sourced components. Checking Visual Annotation Only marks those sourced components as annotations; it does not invent a new component from the SVGs. There is no wizard field for a single custom component kind, and **Source** cannot be skipped.

Do not invent a dummy URL just to click Generate. That would either fail, or register unrelated components styled with the KubeCon logo, which is not the lab goal.

The working Path A step for a single KubeCon Logo component is **Import Model** with the ready-made archive, next.
{{% /alert %}}

### Import the ready-made KubeCon logo model

1. Stay on the Registry **Models** tab and click **Import Model**.
2. Choose **File Import**, pick [kubecon-logo.tar.gz](kubecon-logo.tar.gz) from your computer, and import it.
3. Confirm the model appears in the Registry. This fallback model is named `kubecon-logo` and is shared, so prefer a handle-specific model when you have one. The archive already has `isAnnotation` set to `true` and the SVGs inlined, with a single component whose kind and display name are `kubecon-logo` and `KubeCon Logo`.

### Add the logo to your Lab 1 design in Kanvas

1. **Reopen your design** in Kanvas Designer at https://playground.meshery.io/extension/meshmap.
2. **Find the component.** Search the component list for `KubeCon Logo` or your model display name. Hover it and click the pin icon (Pin To Dock) so it sits in the dock at the bottom. If it does not appear, refresh the Kanvas tab and search again.
3. **Drop it on the canvas.** Drag the KubeCon Logo component onto the design near your Deployment and Service. It is annotation-only, so it is a visual label, not a workload.
4. **Save.** Wait for the cloud-with-a-tick save indicator.
5. **Optional: update the Catalog request.** If your Lab 1 publish request is still pending review, this save is enough; the submitted design updates automatically. If the design was already approved and is locked in the Catalog, clone it, add the logo to the clone, and publish the new design.
<!-- tab -->
Path B: Meshery UI and mesheryctl

Path B builds the model with `mesheryctl` (or imports the ready-made archive in the Registry) and adds the component without Kanvas. If you did not already point `mesheryctl` at the Playground in Lab 1, do that first. Signing in to the Playground in your browser does not configure `mesheryctl`.

```sh
mesheryctl system context create playground --url https://playground.meshery.io --set
mesheryctl system token create playground-token -f ~/.meshery/playground-auth.json --set
mesheryctl system login
chmod 600 ~/.meshery/playground-auth.json
mesheryctl system context view
```

Confirm `mesheryctl system context view` shows the Playground endpoint before you import or delete anything.

### Scaffold, edit, package and import

1. **Scaffold the model.** Pick a unique model name that includes your handle, using lowercase letters, digits and hyphens only:

   ```sh
   mesheryctl model init kubecon-<your-handle>-logo --version v0.1.0
   ```

   That creates `kubecon-<your-handle>-logo/v0.1.0/`, with `model.json`, a sample file under `components/`, a sample `relationship.json` under `relationships/` (you remove it in step 3), and empty `connections/` and `credentials/` directories.

2. **Edit `model.json`.** Set `metadata.isAnnotation` to `true`. Set `registrant.kind` to `meshery`, give the model a short description, and set `metadata.svgColor` to the color SVG and `metadata.svgWhite` to the white SVG. Each SVG must go in as a JSON-escaped string: its attributes contain `"` quotes, so pasting the raw markup, even on one line, makes `model.json` invalid. Generate each value with `sed '/^<?xml/d' kubecon-logo-color.svg | jq -Rs .` (and the same for `kubecon-logo-white.svg`), which drops the XML header and escapes quotes and newlines, then paste the whole output, including its surrounding quotes, as the field's value. Set both `version` and `model.version` to `v0.1.0`: `model init` names the folder after `--version` but writes `v0.0.1` into `model.json`. Leave `schemaVersion` as `models.meshery.io/v1beta2`. If you set a category, *Uncategorized* is a safe choice.

3. **Replace the sample component.** In `components/`, keep a single JSON file (you can rename `component.json` to `kubecon-logo.json`). Set:

   - `displayName` to `KubeCon Logo`
   - `component.kind` to `kubecon-logo`
   - `metadata.isAnnotation` to `true`
   - `styles.svgColor` and `styles.svgWhite` to the same JSON-escaped SVG strings
   - `schemaVersion` to `components.meshery.io/v1beta1`. This edit is required: the scaffold writes `components.meshery.io/v1beta3`, which the Playground's Meshery Server rejects, and `model import` then fails with *No valid component or relationship found in the Model provided*.
   - the nested `model.name` to the same name you passed to `model init` (the scaffold writes `untitled-model`)

   Delete the sample `relationships/relationship.json`. It uses `relationships.meshery.io/v1beta3`, which the server also rejects: the model still imports, but the import reports *encountered error for 1 entity (1 relationship)*. A model imports if it contains at least one valid component.

4. **Package the model as an OCI archive:**

   ```sh
   mesheryctl model build kubecon-<your-handle>-logo/v0.1.0
   ```

   This writes `kubecon-<your-handle>-logo-v0-1-0.tar` in your current directory.

5. **Import it into the Playground:**

   ```sh
   mesheryctl model import -f kubecon-<your-handle>-logo-v0-1-0.tar
   mesheryctl model view kubecon-<your-handle>-logo
   ```

   You can also import from the UI: avatar, then **Registry**, then **Import Model**, then **File Import**, and choose the `.tar` file. Either route registers the model in the shared Playground instance.

{{% alert color="info" title="Fallback if init or build fails" %}}
Download the [ready-made model archive](kubecon-logo.tar.gz) from this page and import it with `mesheryctl model import -f kubecon-logo.tar.gz`, or use Import Model in the Registry. This fallback model is named `kubecon-logo` and is shared, so prefer the handle-specific path above when you can. The archive is a model directory with `model.json` and one annotation component; `isAnnotation` is already set to `true` and the SVGs are already inlined.
{{% /alert %}}

### Add the logo to your Lab 1 design (without Kanvas)

Without Kanvas you will not see the logo on a canvas. You can still add the annotation component to the design document, save or update it, and confirm it with `mesheryctl design view`.

**Preferred: add it in the Design Configurator**

1. Go to **Configuration** > **Designs**, open the action menu for your Lab 1 design, and choose **Edit**. The Design Configurator opens at `/configuration/designs/configurator?design_id=...`.
2. In **Select Category**, look for the category you gave the model (the ready-made archive uses *Uncategorized*). In **Select Model**, choose your KubeCon logo model.
3. Expand the **KubeCon Logo** component form, give it a name if prompted, and confirm it appears in the YAML editor next to the forms.
4. Click the **Update Design** icon (not **Save Design as New File**, unless you want a copy).

If the model does not appear in the category or model dropdowns, refresh the page and search again. Whether a just-imported annotation-only model always appears in the Design Configurator has not been verified in the Playground; if it does not, use the YAML path below.

**YAML path: export, add a component, and update**

1. Export the design. In Meshery UI, on the Designs page, choose **Download** and pick **Meshery Design (yaml)**. Or run:

   ```sh
   mesheryctl design export kubecon-<your-handle>-design
   ```

   The default `--type` is `current` (the design as Meshery stores it now).

2. Open the exported YAML. It should have a top-level `components:` list and a `schemaVersion` such as `designs.meshery.io/v1beta1`. Append this component to that list. Replace the model `name` with your handle-specific model name, or use `kubecon-logo` if you imported the fallback archive. Generate a fresh UUID for `id` (for example with `uuidgen`):

   ```yaml
   - id: 00000000-0000-0000-0000-000000000000
     schemaVersion: components.meshery.io/v1beta1
     version: v0.1.0
     displayName: KubeCon Logo
     format: JSON
     status: enabled
     model:
       schemaVersion: models.meshery.io/v1beta2
       version: v0.1.0
       name: kubecon-<your-handle>-logo
       displayName: KubeCon Logo
       status: enabled
       metadata:
         isAnnotation: true
       model:
         version: v0.1.0
     metadata:
       isAnnotation: true
     component:
       kind: kubecon-logo
       version: v0.1.0
       schema: ""
     configuration: {}
   ```

   This snippet matches the v1beta1 component shape used in published catalog designs and in the ready-made `kubecon-logo` model (`schemaVersion`, `component.kind`, nested `model.name` / `model.version`, and `metadata.isAnnotation`). It omits Kanvas-only `styles.position` fields. The Registry should supply the component's SVG styling when the referenced model is registered. This exact append-and-update path has been checked against the schema and against Meshery source, but it has not been dry-run in the Playground.

3. Update the existing design: open it in the Design Configurator, paste the edited YAML into the editor, and click **Update Design**.

   `mesheryctl design import` always creates a **new** design (it posts to `/api/pattern/import` and does not update an existing ID). If you re-import the edited file, give it a new name, for example:

   ```sh
   mesheryctl design import -f kubecon-<your-handle>-design-with-logo.yaml -s "Meshery Design" -n kubecon-<your-handle>-design-with-logo
   ```

4. Confirm the logo component is in the design:

   ```sh
   mesheryctl design view kubecon-<your-handle>-design
   ```

5. **Optional: update the Catalog request.** If your Lab 1 publish request is still pending review, updating the same design is enough; the submitted request updates automatically. If you created a new design by re-importing, publish that new design. If the original design was already approved and is locked, clone it, add the logo to the clone, and publish the new design.
{{< /tabs >}}

<details>
<summary>Instructor guide: Lab 2</summary>

**Path A.** Show the logo downloads on this page. Open the Registry at https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models, walk through **Create Model** far enough to show Model Details, Styling (the two SVG uploads), Visual Annotation Only, and the required Source step. Tell the room, honestly, that Generate sources components from a GitHub or Artifact Hub URL and cannot create a single logo component from the SVGs alone. Switch to **Import Model** and import `kubecon-logo.tar.gz`. In Kanvas Designer, pin the component, drop it on the Lab 1 design, and save.

**Path B.** Run `mesheryctl model init`, edit `isAnnotation` and the SVG fields, run `model build`, and `model import -f` against the Playground context. If the CLI path is slow for the room, switch early to the ready-made `kubecon-logo.tar.gz` and Import Model in the Registry. Show adding the component in the Design Configurator, or the YAML snippet and **Update Design**. Confirm with `mesheryctl design view`.

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
| `model build` or `model import` fails | Open the Notification Center or the CLI error. Check that `metadata.isAnnotation` is `true`, that each SVG is a JSON-escaped string (see step 2), and that the component `schemaVersion` is `components.meshery.io/v1beta1`. If it stays stuck, import the ready-made [kubecon-logo.tar.gz](kubecon-logo.tar.gz). |
| Model not in the Kanvas component list | Refresh the Kanvas tab, then search by display name (`KubeCon Logo`). If you used Import Model, confirm the model is listed in the Registry first. |
| Model not in the Design Configurator dropdowns | Refresh the configurator tab, then look under the category you chose (the ready-made archive uses *Uncategorized*). If it still does not appear, use the YAML export-and-update path. |
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

- [Creating a Meshery Design]({{< ref "guides/configuration-management/creating-a-meshery-design.md" >}})
- [Importing and Exporting Designs]({{< ref "guides/configuration-management/import-export-designs.md" >}})
- [Catalog]({{< ref "concepts/architecture/catalog/index.md" >}})
- [Publishing Designs](https://docs.layer5.io/kanvas/designer/publishing-designs/)
- [Create Models]({{< ref "guides/configuration-management/creating-models/index.md" >}})
- [Importing Models]({{< ref "guides/configuration-management/importing-models/index.md" >}})
- [mesheryctl design import]({{< ref "reference/mesheryctl/design/import.md" >}})
- [mesheryctl design list]({{< ref "reference/mesheryctl/design/list.md" >}})
- [mesheryctl design view]({{< ref "reference/mesheryctl/design/view.md" >}})
- [mesheryctl design export]({{< ref "reference/mesheryctl/design/export.md" >}})
- [mesheryctl model init]({{< ref "reference/mesheryctl/model/init.md" >}})
- [mesheryctl model build]({{< ref "reference/mesheryctl/model/build.md" >}})
- [mesheryctl model import]({{< ref "reference/mesheryctl/model/import.md" >}})
- [Using Meshery Playground]({{< ref "installation/playground.md" >}})
- Kanvas: [Getting started](https://docs.layer5.io/kanvas/getting-started/), [Creating relationships](https://docs.layer5.io/kanvas/getting-started/creating-relationships/), [Pin a model to the dock](https://docs.layer5.io/kanvas/designer/pin-model-to-dock/), [Saving designs](https://docs.layer5.io/kanvas/designer/saving-designs/)
- [Contributor Training](https://layer5.io/learn/workshops)
