# Choosing a Content Type

> How to decide whether a page is a tutorial, a how-to guide, an explanation, or reference, with tests and Meshery examples.

---

LLMS index: [llms.txt](/pr-preview/pr-22678/llms.txt)

---

Meshery documentation follows the [Diátaxis](https://diataxis.fr) framework, which sorts documentation into four types: tutorials, how-to guides, explanation (Meshery calls these Concepts), and reference. Use this page to decide which type a page is before you write or restructure it.

## One page, one type

Diátaxis does not require each of its four types to have its own top-level section in the table of contents. Meshery keeps its top level small and bundled, so a tutorial track can live inside an existing section, such as getting started, as a clearly labeled subsection.

What Diátaxis does require is:

- Each page is exactly one type.
- Each page is clearly labeled as that type, through its section, title, or introduction, so readers know what they are getting.

## Tutorials and how-to guides are different documents

The easiest types to confuse are tutorials and how-to guides, because both are written as steps. The difference is not how much hand-holding a page does, and the two do not come in pairs. A tutorial is not a more detailed how-to guide, and a how-to guide is not a shorter tutorial.

They differ in who owns the outcome:

- In a **tutorial**, the author picks the goal, the environment, and the result. The reader is a learner whose real goal is to become able to do things.
- In a **how-to guide**, the reader brings a real goal from their own work, and the guide only has to get them there.

Diátaxis describes this as learning versus doing: study versus work. That difference changes what goes on the page, not just how much detail it has.

| | Tutorial | How-to guide |
| --- | --- | --- |
| **Environment** | Controls it, for example "use the Playground and this sample design", so every step is guaranteed to work. | Assumes the reader's real environment, with its own versions and choices, and has to handle variation. |
| **Choices** | Makes every choice for the reader and offers no alternatives. | Branches where real readers differ, for example "if you use EKS, do X". |
| **Explanation** | Keeps it to a minimum and links out to Concepts. | Has none, just the steps. |
| **Order and completeness** | One guided path, where each step builds on the last. | Covers only the task, and the reader can start partway through. |
| **Success** | The learner ends up with the promised result and the confidence to try more. | The reader's own task is done. |

### Examples from Meshery

Because the two types serve different needs, a subject can have a tutorial, a how-to guide, both, or neither.

- **Tutorial only.** "Build your first design in Kanvas" uses a fixed sample application in the Meshery Playground. Nobody arrives with "build the sample design" as a real task, so a how-to version would serve no one.
- **How-to only.** "Configure Meshery behind a corporate proxy" or "Upgrade Meshery Server in a production cluster" are real operational jobs. Nobody learns Meshery by doing them, so a tutorial version would be artificial.
- **Both, as different documents.** For relationships, a tutorial walks a learner through defining a simple hierarchical relationship between two sample components and watching it take effect, with every value given. The how-to guide, "Define a relationship for your model", assumes you have your own model, tells you how to choose `kind`, `type`, and `subType`, and covers edge-binding versus hierarchical relationships. The two share a subject, but neither is a more detailed version of the other.

## Tests for tutorial versus how-to guide

To decide which type a piece of content is, ask these questions in order:

1. **Whose goal is it?** If the author chose the outcome as an exercise, it is a tutorial. If the reader brings it from their own work, it is a how-to guide.
2. **Would someone search for this outcome by name while doing their job?** "How do I connect an AKS cluster?" means yes, so it is a how-to guide. "Build the sample design" means no, so it is a tutorial.
3. **Does it work only in a controlled setup?** If it relies on the Playground, sample data, or fixed versions, it is a tutorial. If it has to work in whatever setup the reader has, it is a how-to guide.
4. **Does it make every choice for the reader?** If it removes decisions, it is a tutorial. If it offers branches or options, it is a how-to guide.
5. **What does "done" look like?** If the reader can now do things they could not do before, it is a tutorial. If their task is finished, it is a how-to guide.
6. **Would a competent user skip straight to step 5?** If skipping breaks the lesson, it is a tutorial. If skipping is normal, it is a how-to guide.

## When the answers are mixed

If the tests give mixed answers, the page is usually two documents fused together and should be split into a tutorial and a how-to guide that link to each other. When you split, move, or merge a published page, add redirects as described in [Contributing to Meshery Docs](/pr-preview/pr-22678/project/contributing/contributing-docs/docs/).

Two more checks catch content that is neither a tutorial nor a how-to guide:

- If the content is mostly about **why** (background, design, how the parts fit together), it belongs in **Concepts**.
- If the content is mostly about **what the values are** (fields, flags, options, defaults, API objects), it belongs in **Reference**.

See [Meshery Documentation Structure and Organization](/pr-preview/pr-22678/project/contributing/contributing-docs-structure/) for how these content types map onto the sections of Meshery Docs.
