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

## Policy Evaluation

The relationships are a powerful way to design your infrastructure and each of them are backed by one or more policies. Policies evaluate the designs for potential relationships and the decide whether to create/delete/update the relationships.


[![Meshery Models Policy Evaluation]({{ site.baseurl }}/assets/img/concepts/meshery-models-policy-evaluation.svg
)]({{ site.baseurl }}/assets/img/concepts/meshery-models-policy-evaluation.svg)

Meshery Server has a built-in policy engine, based on Open Policy Agent (OPA). Currently, Meshery Server is the only place where the policy evals occur. Policy evaluation is invoked each time a design is updated, and each time a design is imported. By default, policies evaluate for all registered relationships.

In any given Meshery deployment, you can reference and search the full set of registered policies (in Meshery's internal registry) in using either of Meshery's client interfaces.

{% include alert.html type="info" title="Viewing All Registered Relationships" content='<p>You can view all registered relationships using either Meshery UI or Meshery CLI.</p><dl><dt>Using Meshery UI...</dt><dd>Navigate to <i>Settings</i>, then to <i>Registry</i></dd>.<dt>Using Meshery CLI...</dt><dd><code>mesheryctl policy list</code></dd></dl>' %}

<!-- There are different points in time in which policy evaluations are invoked

1. Each time the design is updated.
2. A Design/HelmChart/K8s Manifest/Docker Compose app is imported/uploaded.
3. Ad-hoc invocation from the Actions Center (coming soon). 
-->

### How are conflicts resolved?

In the event of a conflict or tie, Meshery relies on Open Policy Agent's [reconciliation behavior](https://www.openpolicyagent.org/docs/latest/faq/) for conflict resolution.

{% include alert.html type="warning" title="Conflict Resolution" content="It may happen that certain eval decisions contain results such that two different components create a conflicting relationship with same component. While this is semantically correct, the visual representation of the relationship in such cases may be undesirable, and you may see relationships and components being redrawn depending upon how the client / Meshery UI visualizes the relationships." %}

{% include alert.html type="dark" title="Future Feature" content="Policy evaluation in WASM runtime is on roadmap for Meshery v0.8.3." %}
