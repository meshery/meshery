---
layout: default
title: Policies
permalink: concepts/logical/policies
type: concepts
abstract: "Meshery Policies enable you with a broad set of controls and governance of the behavior of systems under Meshery's management."
language: en
list: include
---

Policies offer an evaluation algorithm to ensure desired behavior enforcement. Policies can be applied to components and relationships, defining rules and actions based on predefined conditions.

[![Meshery Models Policy Evaluation]({{ site.baseurl }}/assets/img/concepts/meshery-models-policy-evaluation.svg
)]({{ site.baseurl }}/assets/img/concepts/meshery-models-policy-evaluation.svg)

## Policy Evaluation

The relationships are a powerful way to design your infrastructure and each of them are backed by one or more policies. Policies evaluate the designs for potential relationships and the decide whether to create/delete/update the relationships.

By default, policies evaluate for all registered relationships.
_For navigating all registered relationships, brose the Registry UI under Meshery UI Settings._
You can configure the behaviour and restrict the policy evaluation on a subset of relationships from the Actions Center.
The supported set of relationships can be found in

Where and when does the policy evaluation occur?

The Meshery Server has a Policy Engine built on top of OPA (Open Policy Agent). Currently, Meshery Server is the only place where the policy evals occur. The policy evaluation is invoked each time the design is updated, and each time a Design/HelmChart/K8s Manifest/Docker Compose app is imported/uploaded.


There are different points in time in which policy evaluations are invoked

1. Each time the design is updated.
2. A Design/HelmChart/K8s Manifest/Docker Compose app is imported/uploaded.
3. Ad-hoc invocation from the Actions Center (coming soon).

### How are conflicts resolved?

Refer the OPA FAQs to know how conflicts are resolved. _https://www.openpolicyagent.org/docs/latest/faq/_

**Limitations:**
It may happen that certain eval decisions contain results such that two different components create a conflicting relationship with same component. While this is semantically correct, the visual representation of the relationship in such cases may be undesirable, and you may see relationships and components being redrawn depending upon how the client / Meshery UI visualizes the relationships.

# Itemizing your Policy Definitions in your Meshery deployment

In any given Meshery deployment, you can reference and search the full set of registered policies (in Meshery's internal registry) in using either of Meshery's client interfaces.

**Meshery UI**

- Visit *Setttings* --> *Registry*

**Meshery CLI**

- Run `mesheryctl policy list`

{% include alert.html type="warning" title="Future Feature" content="Policy evaluation in WASM runtime is on roadmap for Meshery v0.8.3." %}

