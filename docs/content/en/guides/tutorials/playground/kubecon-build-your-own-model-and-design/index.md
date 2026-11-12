---
title: "KubeCon Workshop: Build Your Own Meshery Model and Design"
description: A 60-minute Meshery Playground workshop in which every attendee creates their own model and a design that uses it.
model: playground
params:
    kind: workshop
categories: [tutorials]
weight: 10
aliases:
- /guides/tutorials/kubecon-build-your-own-model-and-design
---

This is a 60-minute hands-on session for Meshery ContribFests. Every attendee creates their own Meshery model from real CRDs and then their own Meshery design that uses it, entirely in the Meshery Playground.

Learning goals: explain what a Meshery model, component, relationship and design are; generate a model from a CRD source with the Registry's Create Model wizard; build, configure and validate a design in Kanvas that uses that model; and keep the work after the session.

## At a glance

- **Duration:** 60 minutes
- **Where attendees work:** https://playground.meshery.io (browser only, optionally install `mesheryctl`)
- **Outcome:** One model per attendee, registered in the shared Playground instance, and one design saved to the attendee's account
- **Audience:** Kubernetes users; no Meshery experience needed
- **Format:** Short talk, two guided labs, wrap-up
- **Resources:** Meshery Models explainers ([architecture slides](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit?slide=id.g28d3a0f18b4_0_292#slide=id.g28d3a0f18b4_0_292))

## Timed agenda

| Clock | Min | Block | Attendees leave with |
|---|---|---|---|
| 0:00 | 3 | Welcome, goals, logistics | Wi-Fi, Playground link, helper flag |
| 0:03 | 7 | Sign in to the Playground | A working Playground session |
| 0:10 | 5 | Concepts: models, components, relationships, designs | The mental model for the labs |
| 0:15 | 18 | Lab 1: Create your model | Their own model in the Registry |
| 0:33 | 20 | Lab 2: Create your design | A saved, validated design using their model |
| 0:53 | 7 | Wrap-up and next steps | Links to keep going |

The six blocks are 3, 7, 5, 18, 20 and 7 minutes, which add up to 60. If Lab 1 runs long, take the time from wrap-up, never from Lab 2.

## Prerequisites

### Attendees

- A laptop with a current Chrome, Edge, Firefox or Safari. Tablets and phones are not enough for drag and drop.
- [Install `mesheryctl`](https://docs.meshery.io/installation/mesheryctl/) on your computer.
- A GitHub or Google account, or an email address, to sign in to the Playground. Sign-in is required: the Playground runs live clusters and every user needs a named account.
- Optional but faster: sign in at https://playground.meshery.io before the session.
- Basic Kubernetes vocabulary: Namespace, Deployment, Service, CRD.

<details>
<summary>Instructor prerequisites</summary>

- A Playground account that has completed both labs end to end within 24 hours of the session.
- Projector at 1080p or better, browser zoom at 125% so the back row can read the Registry form.
- This page open on a second screen, and the attendee steps (Labs 1 and 2) on a slide or QR code.
- A note to the Meshery Playground maintainers that a room of attendees will generate models at roughly 0:15 to 0:33.

</details>

<details>
<summary>Instructor prep (morning of)</summary>

1. **Dry run both labs on the Playground.** Use a fresh model name such as `kubecon-instructor-certs` so you hit the same screens attendees will. Note any UI label that differs from this page and tell the room at the start.
2. **Confirm the cert-manager CRD source still resolves.** Open `https://github.com/cert-manager/cert-manager/tree/master/deploy/crds` and confirm the folder exists on the default branch. It was verified on Nov 11, 2026. Also confirm the stretch-goal sources if you plan to point fast finishers at them.
3. **Pre-build a finished design.** Keep a completed Lab 2 design open in a second tab. If the network or Playground stalls, you can keep teaching from it.
4. **Prepare the handle list.** Ask attendees to use their GitHub handle (or first name plus last initial) in every name, so dozens of models and designs never collide.
5. **Brief helpers.** Helpers own sign-in problems and Lab 1 errors so the instructor keeps the room moving. Agree on a raised-hand or sticky-note signal.

</details>

## Module 1: Welcome and sign-in (0:00 to 0:10)

<details>
<summary>Instructor guide: welcome and sign-in</summary>

1. **Set the goal in one sentence (1 min).** "In 60 minutes you will publish your own Meshery model and use it in your own design."
2. **Logistics (2 min).** Wi-Fi, the Playground URL, the helper signal, and the naming rule: put your handle in every name.
3. **Live sign-in (1 min).** Open https://playground.meshery.io on the projector and sign in so the room sees the provider choice.
4. **Hands-on sign-in (6 min).** Helpers walk the room. Anyone not signed in by 0:10 pairs with a neighbor for Lab 1 and keeps trying in parallel.

</details>

### Attendees

1. **Open the Playground.** Go to https://playground.meshery.io. You are sent to the Meshery Cloud sign-in page.
2. **Sign in.** Continue with GitHub or Google, or create an account with email. You land back in Meshery.
3. **Check you are in.** You see the Meshery dashboard with your avatar in the top right. Raise your hand if you do not.

## Module 2: Concepts in five minutes (0:10 to 0:15)

A model is a package; a design is something you build from many packages.

<details>
<summary>Instructor guide: concepts</summary>

Teach only what the labs use. Draw it as one picture: a model is a package; a design is something you build from many packages. Use the right-hand column of the table below as your script.

</details>

| Term | Meaning |
|---|---|
| Model | A versioned package that teaches Meshery about one technology: its components, relationships and styling. Kubernetes is a model; cert-manager is a model. |
| Component | One kind of thing inside a model, usually one Kubernetes kind or CRD: Deployment, Certificate, Issuer. |
| Relationship | How components connect: parent and child (a Pod inside a Namespace), edges (a Service to a Deployment), and shared labels. |
| Registry | Where Meshery keeps every model it knows. You will add yours to it in Lab 1. |
| Design | A declarative document of the components you want, their configuration and their relationships. You can validate, dry run, deploy, share and export it. |
| Kanvas | The visual designer in the Playground where you drag components onto a canvas to build a design. |

{{% alert color="info" title="Key message" %}}
Designs are saved to your account and survive. Deployments in the Playground clusters are reset daily. Today's real artifact is the model and the design, not a running workload.
{{% /alert %}}

## Lab 1: Create your model (0:15 to 0:33)

Attendees generate a cert-manager model from a public folder of CRDs with the Registry's Create Model wizard. Everyone uses the same source so the room stays on one demo path.

Source URL for the wizard: `git://github.com/cert-manager/cert-manager/master/deploy/crds`

The format is `git://github.com/[organization]/[repository]/[branch]/path/to/crds`, as documented in [Create Models]({{< ref "guides/configuration-management/creating-models/index.md" >}}). This is the documented example and the one shown in the live demo. Other CRD sources are listed under Stretch goals.

<details>
<summary>Instructor demo (5 min)</summary>

1. **Open the Registry.** Avatar in the top right, then Registry (or Settings, then Registry). Stay on the Models tab. Direct link: https://playground.meshery.io/settings?settingsCategory=Registry&tab=Models
2. **Click Create Model.** Point out that this is the wizard; Import Model is for files you already have.
3. **Fill in names.** Model Name `kubecon-instructor-certs` (lowercase letters, digits and hyphens only), Model Display Name `Cert Manager (instructor)`.
4. **Category and subcategory.** Pick the closest match. If nothing fits, choose Uncategorized. Styling is optional: skip it or pick a distinctive color so the model is easy to spot later.
5. **Source.** Paste the cert-manager URL above.
6. **Additional settings.** Leave Visual Annotation Only off, because these are real infrastructure components. Clicking Generate registers the model immediately; there is no separate Immediate Registration toggle.
7. **Generate and verify.** Click Generate, show the confirmation and the Notification Center entry, then open the model in the Registry and show its components (Certificate, Issuer, ClusterIssuer and others).

</details>

### Attendee steps (13 min)

1. **Open the Registry.** Click your avatar (top right), choose Registry, and stay on the Models tab.
2. **Start the wizard.** Click Create Model.
3. **Name it with your handle.** Model Name: `kubecon-<your-handle>-certs`, for example `kubecon-octocat-certs`. Display Name: something readable, for example `Cert Manager (octocat)`.
4. **Choose a category.** Pick the closest Category and Subcategory, or Uncategorized. Style is optional; choose a color you will recognize in the designer if you have time.
5. **Paste the cert-manager source URL.** Use the URL from the start of this lab.
6. **Additional settings.** Leave Visual Annotation Only off. Clicking Generate registers the model immediately.
7. **Generate.** Click Generate and wait for the confirmation. Check the Notification Center if you miss it.
8. **Verify.** Find your model in the Models list (search your handle). Open it and count its components. Raise your hand when you can see them.

{{% alert color="warning" title="Model names are permanent" %}}
A generated model cannot be deleted from the UI. That is why everyone uses their handle in the name. To stop using a model later, mark it as ignored, or delete it with `mesheryctl model delete <model-name>`. Signing in to the Playground in your browser does not configure `mesheryctl`: point it at the Playground first with `mesheryctl system context create playground --url https://playground.meshery.io --set`, give the Playground its own token file with `mesheryctl system token create playground-token -f ~/.meshery/playground-auth.json --set` so signing in does not overwrite the token your other contexts use, sign in with `mesheryctl system login`, run `chmod 600 ~/.meshery/playground-auth.json` so only you can read the token on a shared machine, and confirm `mesheryctl system context view` shows the Playground endpoint before you delete anything, so the command does not target another server.
{{% /alert %}}

## Lab 2: Create your design (0:33 to 0:53)

Attendees build a design in Kanvas that combines their new cert-manager model with core Kubernetes components, configures them, adds relationships and validates the result.

<details>
<summary>Instructor demo (5 min)</summary>

1. **Open Kanvas Designer.** Go to https://playground.meshery.io/extension/meshmap and choose the Designer tab.
2. **Create and name the design.** Start a new design and name it `kubecon-instructor-design`. Point at the save indicator beside the name: a cloud with a tick means saved.
3. **Find your model.** Search the component list for your display name. Hover it and click the pin icon (Pin To Dock) so it sits in the dock at the bottom.
4. **Add a Namespace.** From the Kubernetes model, drag a Namespace onto the canvas and name it `kubecon-demo`.
5. **Add your model's components inside it.** From your pinned model, drag an Issuer and a Certificate and drop them inside the Namespace. Dropping a child inside a parent creates a hierarchical relationship.
6. **Configure.** Select each component and fill its configuration form. For cert-manager: give the Issuer a name and a self-signed spec; give the Certificate a `secretName` and an `issuerRef` naming that Issuer. Required fields are marked.
7. **Add a workload and an edge.** Drag a Deployment and a Service into the Namespace. Tap the Service, pick a relationship type from the edge picker, and drag to the Deployment. Kanvas highlights compatible targets.
8. **Validate.** Run Validate. It checks every component against its schema and does not contact a cluster. The count in the result comes from the design, so it includes components whose definitions could not be fetched; a clean result does not guarantee every component was checked. Fix anything it flags.

</details>

### Attendee steps (15 min)

1. **Open Kanvas Designer.** Go to https://playground.meshery.io/extension/meshmap and choose Designer.
2. **New design.** Create a design named `kubecon-<your-handle>-design`.
3. **Pin your model.** Search the component list for your model's display name, hover it and click Pin To Dock.
4. **Add a Namespace.** Drag a Namespace from the Kubernetes model onto the canvas and name it.
5. **Use your model.** Drag an Issuer and a Certificate from your pinned model into the Namespace.
6. **Configure them.** Select each component and fill in its required fields in the configuration form. Give everything a name. For cert-manager: give the Issuer a self-signed spec; give the Certificate a `secretName` and an `issuerRef` naming that Issuer.
7. **Add a workload and connect it.** Add a Deployment and a Service. Create an edge from the Service to the Deployment with the edge picker.
8. **Validate.** Run Validate and fix anything it reports. The validated component count includes components whose definitions cannot be fetched, so a matching count does not confirm that every component was checked.
9. **Confirm it is saved.** Check the save indicator shows a cloud with a tick. Kanvas saves automatically about 3 seconds after you stop editing.

## Wrap-up and next steps (0:53 to 1:00)

1. **Recap (1 min).** Model in the Registry on the shared Playground instance, design in Kanvas saved to your account. Deployments reset; designs do not.
2. **One quick share if time (2 min).** If the labs finished early, one volunteer shows their design: which relationships did Kanvas create, and what did Validate catch?
3. **Next steps (4 min).** Contribute a model upstream, try the stretch goals and the CLI scaffold route at home, join the community Slack at https://slack.layer5.io, and keep using your Playground account.

## Stretch goals (take-home)

These are optional. Use them if you finish a lab early, or after the session.

- **Generate a model from a different CRD source** and compare component counts with your cert-manager model:

  | Pick | Technology | Source URL for the wizard |
  |---|---|---|
  | B | Argo Rollouts | `git://github.com/argoproj/argo-rollouts/master/manifests/crds` |
  | C | KEDA | `git://github.com/kedacore/keda/main/config/crd/bases` |
  | D | Gateway API | `git://github.com/kubernetes-sigs/gateway-api/main/config/crd/standard` |
  | E | External Secrets | `git://github.com/external-secrets/external-secrets/main/config/crds/bases` |

  Suggested design pairs if you use one of these in Kanvas: Argo Rollouts Rollout plus AnalysisTemplate; KEDA ScaledObject plus TriggerAuthentication; Gateway API Gateway plus HTTPRoute; External Secrets SecretStore plus ExternalSecret.
- Read the CLI scaffold route for hand-authored models: `mesheryctl model init`, `mesheryctl model build`, then `mesheryctl model import`. This needs a local Meshery Server, so it is homework, not part of today's lab.
- Dry run the design against the Playground cluster. A dry run sends the design to the cluster without applying it. Custom resources fail if the cluster lacks that CRD, which is a good discussion point.
- Deploy only the Namespace, Deployment and Service, then look at them in Operator mode. Remind yourself that the Playground clusters reset daily.
- Export the design to keep a file copy, or open Details from the File menu to share it with a neighbor.
- Alternative route without Kanvas: in Meshery UI go to Configuration, then Designs, then Create Design, pick your model in Select Model and fill the generated forms while the YAML updates beside them.

## What success looks like

Each attendee can show all of the following before leaving.

- Signed in to the Playground with their own account (by 0:10)
- A model named with their handle appears in Registry, Models, with its components listed (by 0:33)
- That model is pinned to the Kanvas dock (Lab 2)
- A design named with their handle contains a Namespace, at least two components from their model, a Deployment and a Service (Lab 2)
- At least one parent and child relationship and one edge relationship are visible (Lab 2)
- Validate passes, or every remaining finding is understood (by 0:53)
- The save indicator shows saved, and the design reopens from their design list (by 0:53)

<details>
<summary>Instructor and helper notes</summary>

Helpers can use this list as a walk-around check during the last seven minutes.

Room-level target: at least 80% of attendees complete both labs. Track it with a show of hands at 0:33 and 0:53.

</details>

## Troubleshooting and fallbacks

| Symptom | Fix |
|---|---|
| Stuck on the sign-in page | Try another provider (GitHub, Google or email). If the venue blocks it, pair with a neighbor and keep trying on a phone hotspot. |
| Model name rejected | Use lowercase letters, digits, and hyphens only, no spaces or underscores. |
| Generate fails or never confirms | Open the Notification Center for the error. Re-check the cert-manager source URL format and branch name. If it stays stuck, try a stretch-goal source. |
| Model not in the Kanvas component list | Refresh the Kanvas tab, then search by display name. Generate already registers the model; there is no Immediate Registration toggle. |
| Name already taken | Someone else used it. Add a suffix such as `-two`. |
| Save indicator shows Unable to save | Stop editing until it clears. Usually a dropped connection or an expired session; sign in again. |
| Playground slow for the whole room | Attendees who already have a model move to Lab 2; others use an existing Registry model such as Kubernetes plus any model already listed, and still build their own design. |

<details>
<summary>Instructor notes: fallbacks</summary>

- If the Playground is slow for the whole room, continue teaching from the pre-built design.
- The UI steps come from current Meshery and Kanvas documentation. The instructor's morning dry run is the check; update the session if a label differs.

</details>

## Related docs

- [Create Models]({{< ref "guides/configuration-management/creating-models/index.md" >}})
- [Importing Models]({{< ref "guides/configuration-management/importing-models/index.md" >}})
- [Creating a Meshery Design]({{< ref "guides/configuration-management/creating-a-meshery-design.md" >}})
- [Using Meshery Playground]({{< ref "installation/playground.md" >}})
- Kanvas: [Getting started](https://docs.layer5.io/kanvas/getting-started/), [Creating relationships](https://docs.layer5.io/kanvas/getting-started/creating-relationships/), [Pin a model to the dock](https://docs.layer5.io/kanvas/designer/pin-model-to-dock/), [Saving designs](https://docs.layer5.io/kanvas/designer/saving-designs/)
- [Contributor Training](https://layer5.io/learn/workshops)
