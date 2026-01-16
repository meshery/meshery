---
layout: page
title: Contributing to Meshery Policies
permalink: project/contributing/contributing-policies
abstract: How to contribute to Meshery Policies
language: en
type: project
category: contributing
list: include
---

## Background
Meshery has a built-in policy engine, based on [Open Policy Agent (OPA)](https://www.openpolicyagent.org/docs/latest/). Meshery uses the [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) query language to create these [policies](https://docs.meshery.io/concepts/logical/policies). 

## Prerequisites
To start contributing to Meshery Policy Engine, make sure you have [OPA CLI](https://www.openpolicyagent.org/docs/latest/#running-opa) installed on your system. You will also need to clone the [Meshery Server](https://github.com/meshery/meshery/) project from Github.

## Components of Meshery Policy Engine

Rego policies are the declarative logic behind the policy engine. These policies define how the underlying engine understands relationships within Meshery Designs. Rego policies can be found [here](https://github.com/meshery/meshery/tree/master/server/meshmodel/meshery-core/0.7.2/v1.0.0/policies) in the Meshery Github repository.

Currently the Rego policies are invoked from the Go code in Meshery Server. This requires configuring the OPA context as seen [here](https://github.com/meshery/meshkit/blob/master/models/meshmodel/core/policies/rego_policy_relationship.go).

The Meshery Policy Engine is invoked by calling the [evaluation endpoint](https://docs.meshery.io/reference/rest-apis#api-meshmodels-relationships-evaluate). This endpoint is passed a Meshery Design in JSON which it passes to the Meshery Policy Engine for validation.

## Working with Meshery Policy Engine

Working with the Meshery Policy Engine can feel complex because of the large number of different relationships a Meshery Design might contain. Contributors must then be careful to test any changes thoroughly to avoid unexpected results. When working with the Rego policies themselves it makes testing and development easier to be able to test the policies directly without running the entire Meshery Server. When working on the policies keep the following in mind:

- Since we will be running these policies using the OPA CLI directly without the Go bootstrap code we need to provide the necessary context for the policy engine. This means providing a data structure containing the valid Meshery Relationships the policy engine will evaluate against. 
- The input to the policy engine is a Meshery Design as a JSON data structure. When testing the Meshery Policy Engine you will need to provide this design as input. If you are working on a bug you will need to get the design from the issue or directly from the individual who reported the bug.

## Executing OPA 

You can evaluate the Rego policies against test data stored in the `policies/test` folder:

{% capture code_content %}make rego-eval{% endcapture %}
{% include code.html code=code_content %}

