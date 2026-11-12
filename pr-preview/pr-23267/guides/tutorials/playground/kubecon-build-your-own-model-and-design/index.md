# KubeCon Workshop: Build Your Own Meshery Model and Design

> A 90-minute Meshery Playground workshop in which every attendee creates their own model and a design that uses it.

---

LLMS index: [llms.txt](/pr-preview/pr-23267/llms.txt)

---

A 90-minute hands-on session for Thursday, Nov 12, 2026 (Mountain Time). Every attendee creates their own Meshery model from real CRDs and then their own Meshery design that uses it, entirely in the Meshery Playground.

Learning goals: explain what a Meshery model, component, relationship and design are; generate a model from a CRD source with the Registry's Create Model wizard; build, configure and validate a design in Kanvas that uses that model; and keep the work after the session.

## At a glance

- **When:** Thursday, Nov 12, 2026, 90 minutes (America/Denver)
- **Where attendees work:** https://playground.meshery.io (browser only, nothing to install)
- **Outcome:** One model per attendee, registered in the shared Playground instance, and one design saved to the attendee's account
- **Audience:** Kubernetes users; no Meshery experience needed
- **Format:** Short talk, two guided labs, show and tell
- **Staffing:** 1 instructor plus 1 helper per 20 to 25 attendees

## Timed agenda

| Clock | Min | Block | Attendees leave with |
|---|---|---|---|
| 0:00 | 5 | Welcome, goals, logistics | Wi-Fi, Playground link, helper flag |
| 0:05 | 10 | Sign in to the Playground | A working Playground session |
| 0:15 | 10 | Concepts: models, components, relationships, designs | The mental model for the labs |
| 0:25 | 25 | Lab 1: Create your model | Their own model in the Registry |
| 0:50 | 5 | Buffer and stretch | Stragglers caught up |
| 0:55 | 25 | Lab 2: Create your design | A saved, validated design using their model |
| 1:20 | 7 | Show and tell | Two or three designs shown live |
| 1:27 | 3 | Wrap-up and next steps | Links to keep going |

The 5-minute buffer at 0:50 is deliberate. If Lab 1 runs long, take it from there and from show and tell, never from Lab 2.

## Prerequisites

### Attendees

- A laptop with a current Chrome, Edge, Firefox or Safari. Tablets and phones are not enough for drag and drop.
- A GitHub or Google account, or an email address, to sign in to the Playground. Sign-in is required: the Playground runs live clusters and every user needs a named account.
- Optional but faster: sign in at https://playground.meshery.io before the session.
- Basic Kubernetes vocabulary: Namespace, Deployment, Service, CRD.

### Instructor

- A Playground account that has completed both labs end to end within 24 hours of the session.
- Projector at 1080p or better, browser zoom at 125% so the back row can read the Registry form.
- This page open on a second screen, and the attendee steps (Labs 1 and 2) on a slide or QR code.
- A note to the Meshery Playground maintainers that a room of attendees will generate models at roughly 0:25 to 0:50.

## Instructor prep (morning of)

1. **Dry run both labs on the Playground.** Use a fresh model name such as `kubecon-instructor-certs` so you hit the same screens attendees will. Note any UI label that differs from this page and tell the room at the start.
2. **Confirm the CRD sources still resolve.** Open each source repository in the Lab 1 table and confirm the folder exists on its default branch. All five were verified on Nov 11, 2026.
3. **Pre-build a finished design.** Keep a completed Lab 2 design open in a second tab. If the network or Playground stalls, you can keep teaching from it.
4. **Prepare the handle list.** Ask attendees to use their GitHub handle (or first name plus last initial) in every name, so dozens of models and designs never collide.
5. **Brief helpers.** Helpers own sign-in problems and Lab 1 errors so the instructor keeps the room moving. Agree on a raised-hand or sticky-note signal.

## Module 1: Welcome and sign-in (0:00 to 0:15)

### Instructor

1. **Set the goal in one sentence (1 min).** "In 90 minutes you will publish your own Meshery model and use it in your own design."
2. **Logistics (2 min).** Wi-Fi, the Playground URL, the helper signal, and the naming rule: put your handle in every name.
3. **Live sign-in (2 min).** Open https://playground.meshery.io on the projector and sign in so the room sees the provider choice.
4. **Hands-on sign-in (10 min).** Helpers walk the room. Anyone not signed in by 0:15 pairs with a neighbor for Lab 1 and keeps trying in parallel.

### Attendees

1. **Open the Playground.** Go to https://playground.meshery.io. You are sent to the Meshery Cloud sign-in page.
2. **Sign in.** Continue with GitHub or Google, or create an account with email. You land back in Meshery.
3. **Check you are in.** You see the Meshery dashboard with your avatar in the top right. Raise your hand if you do not.

## Module 2: Concepts in ten minutes (0:15 to 0:25)

Teach only what the labs use. Draw it as one picture: a model is a package; a design is something you build from many packages.

| Term | Say this |
|---|---|
| Model | A versioned package that teaches Meshery about one technology: its components, relationships and styling. Kubernetes is a model; cert-manager is a model. |
| Component | One kind of thing inside a model, usually one Kubernetes kind or CRD: Deployment, Certificate, Issuer. |
| Relationship | How components connect: parent and child (a Pod inside a Namespace), edges (a Service to a Deployment), and shared labels. |
| Registry | Where Meshery keeps every model it knows. You will add yours to it in Lab 1. |
| Design | A declarative document of the components you want, their configuration and their relationships. You can validate, dry run, deploy, share and export it. |
| Kanvas | The visual designer in the Playground where you drag components onto a canvas to build a design. |

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Key message</div>


Designs are saved to your account and survive. Deployments in the Playground clusters are reset daily. Today's real artifact is the model and the design, not a running workload.
</div>


## Lab 1: Create your model (0:25 to 0:50)

Attendees generate a model from a public folder of CRDs with the Registry's Create Model wizard. Each attendee picks one source below; spreading picks across the room keeps the Playground responsive.

| Pick | Technology | Source URL for the wizard |
|---|---|---|
| A | cert-manager | `git://github.com/cert-manager/cert-manager/master/deploy/crds` |
| B | Argo Rollouts | `git://github.com/argoproj/argo-rollouts/master/manifests/crds` |
| C | KEDA | `git://github.com/kedacore/keda/main/config/crd/bases` |
| D | Gateway API | `git://github.com/kubernetes-sigs/gateway-api/main/config/crd/standard` |
| E | External Secrets | `git://github.com/external-secrets/external-secrets/main/config/crds/bases` |

The format is `git://github.com/[organization]/[repository]/[branch]/path/to/crds`, as documented in [Create Models](/pr-preview/pr-23267/guides/configuration-management/creating-models/). Option A is the documented example and the one the instructor demos.

### Instructor demo (8 min)

1. **Open the Registry.** Avatar in the top right, then Registry (or Settings, then Registry). Stay on the Models tab. Direct link: https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models
2. **Click Create Model.** Point out that this is the wizard; Import Model is for files you already have.
3. **Fill in names.** Model Name `kubecon-instructor-certs` (lowercase and hyphens only), Model Display Name `Cert Manager (instructor)`.
4. **Category and subcategory.** Pick the closest match. If nothing fits, choose Uncategorized.
5. **Styling.** Optional: logos, primary and secondary colors, icon shape. Skipping uses Meshery's defaults. Suggest a distinctive color so it is easy to spot later.
6. **Source.** Paste the Option A URL.
7. **Additional settings.** Turn on Immediate Registration so the model is usable right away. Leave Visual Annotation Only off, because these are real infrastructure components.
8. **Generate and verify.** Click Generate, show the confirmation and the Notification Center entry, then open the model in the Registry and show its components (Certificate, Issuer, ClusterIssuer and others).

### Attendee steps (15 min)

1. **Open the Registry.** Click your avatar (top right), choose Registry, and stay on the Models tab.
2. **Start the wizard.** Click Create Model.
3. **Name it with your handle.** Model Name: `kubecon-<your-handle>-<tech>`, for example `kubecon-octocat-certs`. Display Name: something readable, for example `Cert Manager (octocat)`.
4. **Choose a category.** Pick the closest Category and Subcategory, or Uncategorized.
5. **Style it (optional).** Choose a color you will recognize in the designer.
6. **Paste your source URL.** Use the URL for the pick your table group was assigned (A to E).
7. **Register immediately.** Turn on Immediate Registration. Leave Visual Annotation Only off.
8. **Generate.** Click Generate and wait for the confirmation. Check the Notification Center if you miss it.
9. **Verify.** Find your model in the Models list (search your handle). Open it and count its components. Raise your hand when you can see them.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Tell the room</div>


A generated model cannot be deleted from the UI. That is why everyone uses their handle in the name. To stop using a model later, mark it as ignored, or delete it with `mesheryctl model delete <model-name>`. Signing in to the Playground in your browser does not configure `mesheryctl`: point it at the Playground first with `mesheryctl system context create playground --url https://playground.meshery.io --set`, give the Playground its own token file with `mesheryctl system token create playground-token -f ~/.meshery/playground-auth.json --set` so signing in does not overwrite the token your other contexts use, sign in with `mesheryctl system login`, run `chmod 600 ~/.meshery/playground-auth.json` so only you can read the token on a shared machine, and confirm `mesheryctl system context view` shows the Playground endpoint before you delete anything, so the command does not target another server.
</div>


### Stretch for fast finishers

- Generate a second model from a different pick and compare component counts.
- Read the CLI scaffold route for hand-authored models: `mesheryctl model init`, `mesheryctl model build`, then `mesheryctl model import`. This needs a local Meshery Server, so it is homework, not part of today's lab.

## Lab 2: Create your design (0:55 to 1:20)

Attendees build a design in Kanvas that combines their new model with core Kubernetes components, configures them, adds relationships and validates the result.

### Instructor demo (8 min)

1. **Open Kanvas Designer.** Go to https://playground.meshery.io/extension/meshmap and choose the Designer tab.
2. **Create and name the design.** Start a new design and name it `kubecon-instructor-design`. Point at the save indicator beside the name: a cloud with a tick means saved.
3. **Find your model.** Search the component list for your display name. Hover it and click the pin icon (Pin To Dock) so it sits in the dock at the bottom.
4. **Add a Namespace.** From the Kubernetes model, drag a Namespace onto the canvas and name it `kubecon-demo`.
5. **Add your model's components inside it.** From your pinned model, drag an Issuer and a Certificate and drop them inside the Namespace. Dropping a child inside a parent creates a hierarchical relationship.
6. **Configure.** Select each component and fill its configuration form. For cert-manager: give the Issuer a name and a self-signed spec; give the Certificate a `secretName` and an `issuerRef` naming that Issuer. Required fields are marked.
7. **Add a workload and an edge.** Drag a Deployment and a Service into the Namespace. Tap the Service, pick a relationship type from the edge picker, and drag to the Deployment. Kanvas highlights compatible targets.
8. **Validate.** Run Validate. It checks every component against its schema and does not contact a cluster. The count in the result comes from the design, so it includes components whose definitions could not be fetched; a clean result does not guarantee every component was checked. Fix anything it flags.

### Attendee steps (15 min)

1. **Open Kanvas Designer.** Go to https://playground.meshery.io/extension/meshmap and choose Designer.
2. **New design.** Create a design named `kubecon-<your-handle>-design`.
3. **Pin your model.** Search the component list for your model's display name, hover it and click Pin To Dock.
4. **Add a Namespace.** Drag a Namespace from the Kubernetes model onto the canvas and name it.
5. **Use your model.** Drag at least two components from your pinned model into the Namespace. Suggested pairs: cert-manager Issuer plus Certificate; Argo Rollouts Rollout plus AnalysisTemplate; KEDA ScaledObject plus TriggerAuthentication; Gateway API Gateway plus HTTPRoute; External Secrets SecretStore plus ExternalSecret.
6. **Configure them.** Select each component and fill in its required fields in the configuration form. Give everything a name.
7. **Add a workload and connect it.** Add a Deployment and a Service. Create an edge from the Service to the Deployment with the edge picker.
8. **Validate.** Run Validate and fix anything it reports. The validated component count includes components whose definitions cannot be fetched, so a matching count does not confirm that every component was checked.
9. **Confirm it is saved.** Check the save indicator shows a cloud with a tick. Kanvas saves automatically about 3 seconds after you stop editing.

### Stretch for fast finishers

- Dry run the design against the Playground cluster. A dry run sends the design to the cluster without applying it. Custom resources fail if the cluster lacks that CRD, which is a good discussion point.
- Deploy only the Namespace, Deployment and Service, then look at them in Operator mode. Remind people the Playground clusters reset daily.
- Export the design to keep a file copy, or open Details from the File menu to share it with a neighbor.
- Alternative route without Kanvas: in Meshery UI go to Configuration, then Designs, then Create Design, pick your model in Select Model and fill the generated forms while the YAML updates beside them.

## Show and tell, wrap-up (1:20 to 1:30)

1. **Show two or three designs (7 min).** Invite volunteers to share their screen or have them share the design with you. Ask: which model did you build, which relationships did Kanvas create, what did Validate catch?
2. **Recap (1 min).** Model in the Registry on the shared Playground instance, design in Kanvas saved to your account. Deployments reset; designs do not.
3. **Next steps (2 min).** Contribute a model upstream, try the CLI scaffold route at home, join the community Slack at https://slack.layer5.io, and keep using your Playground account.

## What success looks like

Each attendee can show all of the following before leaving. Helpers can use this as a walk-around check during the last ten minutes.

- Signed in to the Playground with their own account (by 0:15)
- A model named with their handle appears in Registry, Models, with its components listed (by 0:50)
- That model is pinned to the Kanvas dock (Lab 2)
- A design named with their handle contains a Namespace, at least two components from their model, a Deployment and a Service (Lab 2)
- At least one parent and child relationship and one edge relationship are visible (Lab 2)
- Validate passes, or every remaining finding is understood (by 1:20)
- The save indicator shows saved, and the design reopens from their design list (by 1:20)

Room-level target: at least 80% of attendees complete both labs. Track it with a show of hands at 0:50 and 1:20.

## Troubleshooting and fallbacks

| Symptom | Fix |
|---|---|
| Stuck on the sign-in page | Try another provider (GitHub, Google or email). If the venue blocks it, pair with a neighbor and keep trying on a phone hotspot. |
| Model name rejected | Use lowercase letters and hyphens only, no spaces or underscores. If your handle contains digits, leave them out of the model name. |
| Generate fails or never confirms | Open the Notification Center for the error. Re-check the source URL format and branch name. Switch to a different pick from the table. |
| Model not in the Kanvas component list | Confirm Immediate Registration was on. Refresh the Kanvas tab, then search by display name. |
| Name already taken | Someone else used it. Add a suffix such as `-two`. |
| Save indicator shows Unable to save | Stop editing until it clears. Usually a dropped connection or an expired session; sign in again. |
| Playground slow for the whole room | Instructor continues from the pre-built design. Attendees who already have a model move to Lab 2; others use an existing Registry model such as Kubernetes plus any model already listed, and still build their own design. |

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Instructor note</div>


The UI steps come from current Meshery and Kanvas documentation. The instructor's morning dry run is the check; update the session if a label differs.
</div>


## Related docs

- [Create Models](/pr-preview/pr-23267/guides/configuration-management/creating-models/)
- [Importing Models](/pr-preview/pr-23267/guides/configuration-management/importing-models/)
- [Creating a Meshery Design](/pr-preview/pr-23267/guides/configuration-management/creating-a-meshery-design/)
- [Using Meshery Playground](/pr-preview/pr-23267/installation/playground/)
- Kanvas: [Getting started](https://docs.layer5.io/kanvas/getting-started/), [Creating relationships](https://docs.layer5.io/kanvas/getting-started/creating-relationships/), [Pin a model to the dock](https://docs.layer5.io/kanvas/designer/pin-model-to-dock/), [Saving designs](https://docs.layer5.io/kanvas/designer/saving-designs/)
- [Contributor Training](https://layer5.io/learn/workshops)
