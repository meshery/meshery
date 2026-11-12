# KubeCon Workshop: Publish a Design and Build a Custom Model

> A 60-minute Meshery Playground workshop in which every attendee publishes a simple Kubernetes design to the Catalog, then adds a custom KubeCon logo model to it.

---

LLMS index: [llms.txt](/docs-preview/pr-preview/pr-23341/llms.txt)

---

This is a 60-minute hands-on session for Meshery ContribFests. Every attendee first creates a very simple Kubernetes design in the Meshery Playground, publishes it to the Catalog, then builds a custom annotation-only model with the KubeCon logo and adds that component to the same design.

Learning goals: build and validate a Design in Kanvas; publish it to the Meshery Catalog; scaffold, package and import an annotation-only Model with `mesheryctl`; and incorporate one of that model's components into the design.

## At a glance

- **Duration:** 60 minutes
- **Where attendees work:** https://playground.meshery.io in a browser, plus `mesheryctl` for Lab 2
- **Outcome:** A saved design (a Kubernetes Deployment and Service) submitted to the Catalog, and a custom KubeCon logo model registered in the Playground and dropped onto that design
- **Audience:** Kubernetes users; no Meshery experience needed
- **Format:** Short talk, two guided labs, wrap-up
- **Resources:** Meshery Models explainers ([architecture slides](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit?slide=id.g28d3a0f18b4_0_292#slide=id.g28d3a0f18b4_0_292))

## Timed agenda

| Min | Block | Attendees leave with |
|---|---|---|
| 8 | Welcome, goals, logistics, and sign-in | A working Playground session |
| 5 | Concepts: designs, the Catalog, models and annotations | The mental model for the labs |
| 22 | Lab 1: Create and publish your design | A saved, validated Deployment and Service, submitted to the Catalog |
| 20 | Lab 2: Create a custom model and add it to your design | The KubeCon logo on that design |
| 5 | Wrap-up and next steps | Links to keep going |

The five blocks are 8, 5, 22, 20 and 5 minutes, which add up to 60. If Lab 1 runs long, take the time from wrap-up, never from Lab 2.

## Prerequisites

### Attendees

- A laptop with a current Chrome, Edge, Firefox or Safari. Tablets and phones are not enough for drag and drop.
- [Install `mesheryctl`](https://docs.meshery.io/installation/mesheryctl/) on your computer. You will use it in Lab 2.
- A GitHub or Google account, or an email address, to sign in to the Playground. Sign-in is required: the Playground runs live clusters and every user needs a named account.
- Optional but faster: sign in at https://playground.meshery.io before the session.
- Basic Kubernetes vocabulary: Deployment, Service, label, selector.

<details>
<summary>Instructor prerequisites</summary>

- A Playground account that has completed both labs end to end within 24 hours of the session, including a Catalog publish request and a `mesheryctl model import` into the Playground.
- Projector at 1080p or better, browser zoom at 125% so the back row can read the Kanvas form.
- This page open on a second screen, and the attendee steps (Labs 1 and 2) on a slide or QR code.
- A note to the Meshery Playground maintainers that a room of attendees will submit Catalog publish requests and import a small annotation model.

</details>

<details>
<summary>Instructor prep (morning of)</summary>

1. **Dry run both labs on the Playground.** Use a fresh design name such as `kubecon-instructor-design` and a fresh model name such as `kubecon-instructor-logo` so you hit the same screens attendees will. Note any UI label that differs from this page and tell the room at the start.
2. **Confirm the demo image still exists.** The Deployment uses `meshery/meshery-milestone:latest`. Docker Hub listed only the `latest` tag when this page was last checked on Nov 12, 2026.
3. **Confirm the logo files on this page download.** [Color SVG](kubecon-logo-color.svg), [white SVG](kubecon-logo-white.svg), and the [ready-made model archive](kubecon-logo.tar.gz) should all open or download from the published docs page.
4. **Pre-build a finished design.** Keep a completed Lab 2 design open in a second tab. If the network or Playground stalls, you can keep teaching from it.
5. **Prepare the handle list.** Ask attendees to use their GitHub handle (or first name plus last initial) in every name, so dozens of designs and models never collide.
6. **Brief helpers.** Helpers own sign-in problems, Kanvas configuration, and `mesheryctl` login so the instructor keeps the room moving. Agree on a raised-hand or sticky-note signal.

</details>

## Welcome and sign-in

<details>
<summary>Instructor guide: welcome and sign-in</summary>

1. **Set the goal in one sentence.** "In 60 minutes you will publish a simple Kubernetes design to the Catalog, then add your own KubeCon logo model to it."
2. **Logistics.** Wi-Fi, the Playground URL, the helper signal, and the naming rule: put your handle in every name.
3. **Live sign-in.** Open https://playground.meshery.io on the projector and sign in so the room sees the provider choice.
4. **Hands-on sign-in.** Helpers walk the room. Anyone not signed in by the end of this block pairs with a neighbor for Lab 1 and keeps trying in parallel.

</details>

### Attendees

1. **Open the Playground.** Go to https://playground.meshery.io. You are sent to the Meshery Cloud sign-in page.
2. **Sign in.** Continue with GitHub or Google, or create an account with email. You land back in Meshery.
3. **Check you are in.** You see the Meshery dashboard with your avatar in the top right. Raise your hand if you do not.

## Concepts

A design is something you build and can publish. A model is a package of components you can drop into a design. Today's custom model is annotation-only: it is a logo, not a workload.

| Term | Meaning |
|---|---|
| Design | A declarative document of the components you want, their configuration and their relationships. You can validate, dry run, deploy, share and publish it. |
| Catalog | The public library of published designs. Publishing submits your design for review so others can find and clone it. |
| Model | A versioned package that teaches Meshery about a set of components and their styling. Kubernetes is a model; today's KubeCon logo will be a model too. |
| Component | One kind of thing inside a model. In Lab 1 that is a Deployment or a Service. In Lab 2 it is a single KubeCon logo. |
| Annotation | A visual-only component. It appears on the canvas for diagramming and is not deployed to a cluster. |
| Kanvas | The visual designer in the Playground where you drag components onto a canvas to build a design. |
| Registry | Where Meshery keeps every model it knows. You will import yours into it in Lab 2. |

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Key message</div>


Designs are saved to your account and survive. Deployments in the Playground clusters are reset daily. Today's real artifacts are the design you publish and the model you import, not a running workload.
</div>


## Lab 1: Create and publish your design

Attendees build an extremely simple design in Kanvas: one Kubernetes Deployment that runs the Meshery milestone Nginx image, and one Kubernetes Service that targets it. They save, validate, and publish the design to the Catalog.

### Attendee steps

1. **Open Kanvas Designer.** Go to https://playground.meshery.io/extension/meshmap and choose the Designer tab.
2. **Create and name the design.** Start a new design and name it `kubecon-<your-handle>-design`, for example `kubecon-octocat-design`. The save indicator beside the name is a cloud with a tick when the design is saved. Kanvas saves automatically about 3 seconds after you stop editing.
3. **Add a Deployment.** From the Kubernetes model, drag a Deployment onto the canvas. Open its configuration form, give it a name, and set the container image to `meshery/meshery-milestone:latest`. That image is a small Nginx app; Docker Hub currently publishes it under the `latest` tag only. Add a label such as `app: kubecon-<your-handle>`.
4. **Add a Service.** Drag a Service onto the canvas. In its configuration form, set the selector to the same `app: kubecon-<your-handle>` value so the Service targets the Deployment, and set port and targetPort to `80`.
5. **Connect them.** Tap the Service, pick a relationship type from the edge picker, and drag to the Deployment. Kanvas highlights compatible targets.
6. **Validate.** Run Validate and fix anything it reports. Validate checks every component against its schema and does not contact a cluster. The count in the result comes from the design, so it includes components whose definitions could not be fetched; a clean result does not guarantee every component was checked.
7. **Confirm it is saved.** Check the save indicator shows a cloud with a tick.
8. **Publish to the Catalog.** Open Details for the design: click the info ("i") button, or choose Details from the File menu, or go to Configuration, then Designs, and open the design's details there. Fill in Type, Technology, Description and Caveats (for example: Type "Workload", Technology "Kubernetes", a one-sentence description of your Deployment and Service). Click **Publish To Catalog**.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">What 'published' means in this room</div>


If your account is an Organization Admin or Provider Admin, the design is published immediately. For a regular user, the design enters **Pending Review** and stays in My Designs until an Organization Admin or Provider Admin approves or rejects it. You will get an email when that happens. While the request is pending, you can keep editing the design and the submitted request updates automatically; you do not need to publish again. After a design is approved it is locked and listed in the Catalog, and you would clone it to make a new version. For this lab, success is submitting the publish request, not waiting for approval.
</div>


<details>
<summary>Instructor guide: Lab 1</summary>

1. Open Kanvas Designer at https://playground.meshery.io/extension/meshmap and start a design named `kubecon-instructor-design`. Point at the save indicator.
2. Drag a Deployment, set the image to `meshery/meshery-milestone:latest`, and apply a label such as `app: kubecon-instructor`.
3. Drag a Service, set the selector to that same label, and set port and targetPort to 80. Draw an edge from the Service to the Deployment.
4. Run Validate and show a saved state.
5. Open Details and click Publish To Catalog. Tell the room about Pending Review versus an admin publish, and that submitting the request is the lab goal.

</details>

## Lab 2: Create a custom model and add it to your design

Attendees create an annotation-only model that contains a single component: the KubeCon logo. They import that model into the Playground, open the Lab 1 design, and drop the logo onto the canvas.

Download these files from this page first:

- [KubeCon logo, color](kubecon-logo-color.svg) (`svgColor`)
- [KubeCon logo, white](kubecon-logo-white.svg) (`svgWhite`)
- Optional fallback: [ready-made model archive](kubecon-logo.tar.gz)

The logos are the official KubeCon + CloudNativeCon North America 2026 artwork from the [CNCF artwork repository](https://github.com/cncf/artwork).

### Point mesheryctl at the Playground

Signing in to the Playground in your browser does not configure `mesheryctl`. Do this once:

```sh
mesheryctl system context create playground --url https://playground.meshery.io --set
mesheryctl system token create playground-token -f ~/.meshery/playground-auth.json --set
mesheryctl system login
chmod 600 ~/.meshery/playground-auth.json
mesheryctl system context view
```

The token file is separate so signing in does not overwrite the token your other contexts use. Confirm `mesheryctl system context view` shows the Playground endpoint before you import or delete anything, so the command does not target another server.

### Scaffold, edit, package and import

1. **Scaffold the model.** Pick a unique model name that includes your handle, using lowercase letters, digits and hyphens only:

   ```sh
   mesheryctl model init kubecon-<your-handle>-logo --version v0.1.0
   ```

   That creates `kubecon-<your-handle>-logo/v0.1.0/`, with `model.json`, a sample file under `components/`, and empty `relationships/`, `connections/` and `credentials/` directories.

2. **Edit `model.json`.** Set `metadata.isAnnotation` to `true`. Set `registrant.kind` to `meshery`, give the model a short description, and paste the color SVG into `metadata.svgColor` and the white SVG into `metadata.svgWhite`. Strip any XML header (`<?xml ... ?>`) and collapse each SVG to a single line before pasting. Leave `schemaVersion` as `models.meshery.io/v1beta2`.

3. **Replace the sample component.** In `components/`, keep a single JSON file (you can rename `component.json` to `kubecon-logo.json`). Set:

   - `displayName` to `KubeCon Logo`
   - `component.kind` to `kubecon-logo`
   - `metadata.isAnnotation` to `true`
   - `styles.svgColor` and `styles.svgWhite` to the same inlined SVGs
   - `schemaVersion` to `components.meshery.io/v1beta1` (this wire format is the most widely accepted by current Meshery Servers)
   - the nested `model.name` to the same name you passed to `model init`

   You can delete the sample file under `relationships/` if you do not need it. A model imports if it contains at least one valid component.

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

   You can also import from the UI: avatar, then Registry, then Import Model, then File Import, and choose the `.tar` file. Either route registers the model in the shared Playground instance.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Fallback if init or build fails</div>


Download the [ready-made model archive](kubecon-logo.tar.gz) from this page and import it with `mesheryctl model import -f kubecon-logo.tar.gz`, or use Import Model in the Registry. This fallback model is named `kubecon-logo` and is shared, so prefer the handle-specific path above when you can. The archive is a model directory with `model.json` and one annotation component; `isAnnotation` is already set to `true` and the SVGs are already inlined.
</div>


### Add the logo to your Lab 1 design

1. **Reopen your design** in Kanvas Designer.
2. **Find the component.** Search the component list for `KubeCon Logo` or your model display name. Hover it and click the pin icon (Pin To Dock) so it sits in the dock at the bottom. If it does not appear, refresh the Kanvas tab and search again.
3. **Drop it on the canvas.** Drag the KubeCon Logo component onto the design near your Deployment and Service. It is annotation-only, so it is a visual label, not a workload.
4. **Save.** Wait for the cloud-with-a-tick save indicator.
5. **Optional: update the Catalog request.** If your Lab 1 publish request is still pending review, this save is enough; the submitted design updates automatically. If the design was already approved and is locked in the Catalog, clone it, add the logo to the clone, and publish the new design.

<details>
<summary>Instructor guide: Lab 2</summary>

1. Show the logo downloads on this page, then run `mesheryctl model init`, edit `isAnnotation` and the SVG fields, run `model build`, and `model import -f` against the Playground context.
2. If the CLI path is slow for the room, switch early to the ready-made `kubecon-logo.tar.gz` and Import Model in the Registry.
3. Pin the component, drop it on the Lab 1 design, and save. Remind people that a pending Catalog request picks up the edit automatically.

</details>

## Wrap-up and next steps

1. **Recap.** A Deployment and Service design, submitted to the Catalog; a custom annotation-only KubeCon logo model imported into the Playground and placed on that design. Deployments in the Playground reset daily; designs and models you own do not.
2. **One quick share if time.** If the labs finished early, one volunteer shows their design: the Service-to-Deployment edge, and the logo on the canvas.
3. **Next steps.** Contribute a model upstream, try the stretch goals at home, join the community Slack at https://slack.layer5.io, and keep using your Playground account.

## What success looks like

Each attendee can show all of the following before leaving.

- Signed in to the Playground with their own account
- A design named with their handle contains a Deployment using `meshery/meshery-milestone:latest` and a Service that selects it on port 80
- Validate has been run, and the save indicator shows saved
- A Catalog publish request has been submitted (the design is published, or it is in Pending Review)
- A KubeCon logo model is listed in the Registry (their handle-specific model, or the shared fallback)
- The logo component is on the Lab 1 design and the design has been saved again

<details>
<summary>Instructor and helper notes</summary>

Helpers can use this list as a walk-around check during the last few minutes.

Room-level target: at least 80% of attendees complete both labs. Track it with a show of hands after Lab 1 and after Lab 2.

</details>

## Troubleshooting and fallbacks

| Symptom | Fix |
|---|---|
| Stuck on the sign-in page | Try another provider (GitHub, Google or email). If the venue blocks it, pair with a neighbor and keep trying on a phone hotspot. |
| Model name rejected | Use lowercase letters, digits, and hyphens only, no spaces or underscores. |
| `mesheryctl` talks to the wrong server | Run `mesheryctl system context view` and confirm the endpoint is `https://playground.meshery.io` before import or delete. Recreate the context and token if it is not. |
| `model build` or `model import` fails | Open the Notification Center or the CLI error. Check that `metadata.isAnnotation` is `true`, that each SVG is a single-line string, and that the component `schemaVersion` is `components.meshery.io/v1beta1`. If it stays stuck, import the ready-made [kubecon-logo.tar.gz](kubecon-logo.tar.gz). |
| Model not in the Kanvas component list | Refresh the Kanvas tab, then search by display name (`KubeCon Logo`). |
| Service does not target the Deployment | The Service selector must match a label on the Deployment, for example `app: kubecon-<your-handle>`. Port and targetPort should both be `80`. |
| Publish To Catalog is disabled | Sign in again and reopen Details. Publish requires a signed-in account with the Publish Design permission. |
| Catalog request still says Pending Review | That is expected for a regular user. The lab is complete once the request is submitted. Keep editing if you are still in Lab 2; the request updates automatically. |
| Name already taken | Someone else used it. Add a suffix such as `-two`. |
| Save indicator shows Unable to save | Stop editing until it clears. Usually a dropped connection or an expired session; sign in again. |
| Playground slow for the whole room | Attendees who already have a design move to publish or to Lab 2. If model import is slow, use the ready-made archive or pin any already-imported `kubecon-logo` model. |

<details>
<summary>Instructor notes: fallbacks</summary>

- If the Playground is slow for the whole room, continue teaching from the pre-built design.
- The UI steps come from current Meshery and Kanvas documentation. The instructor's morning dry run is the check; update the session if a label differs.

</details>

## Stretch goals

These are optional. Use them if you finish a lab early, or after the session.

- Generate a model from cert-manager CRDs with the Registry's Create Model wizard, using `git://github.com/cert-manager/cert-manager/master/deploy/crds`.
- Dry run or deploy only the Deployment and Service, then look at them in Operator mode. The Playground clusters reset daily.
- Export the design to keep a file copy, or open Details from the File menu to share it with a neighbor.
- Wrap the Deployment and Service in a Namespace and let Kanvas create the parent-child relationship.

## Related docs

- [Creating a Meshery Design](/docs-preview/pr-preview/pr-23341/guides/configuration-management/creating-a-meshery-design/)
- [Catalog](/docs-preview/pr-preview/pr-23341/concepts/architecture/catalog/)
- [Publishing Designs](https://docs.layer5.io/kanvas/designer/publishing-designs/)
- [Create Models](/docs-preview/pr-preview/pr-23341/guides/configuration-management/creating-models/)
- [Importing Models](/docs-preview/pr-preview/pr-23341/guides/configuration-management/importing-models/)
- [mesheryctl model init](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/model/init/)
- [mesheryctl model build](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/model/build/)
- [mesheryctl model import](/docs-preview/pr-preview/pr-23341/reference/mesheryctl/model/import/)
- [Using Meshery Playground](/docs-preview/pr-preview/pr-23341/installation/playground/)
- Kanvas: [Getting started](https://docs.layer5.io/kanvas/getting-started/), [Creating relationships](https://docs.layer5.io/kanvas/getting-started/creating-relationships/), [Pin a model to the dock](https://docs.layer5.io/kanvas/designer/pin-model-to-dock/), [Saving designs](https://docs.layer5.io/kanvas/designer/saving-designs/)
- [Contributor Training](https://layer5.io/learn/workshops)
