---
layout: enhanced
title: "Extensibility: Authorization"
permalink: extensibility/authorization
type: Extensibility
abstract: "Meshery architecture is extensible. Meshery provides several extension points for working with different cloud native projects via authorization, adapters, load generators and providers."
language: en
list: include
---

Meshery features an extensible authorization system that offers the ability to deliver fine-grained access control across it's web-based user interface, [Meshery UI]({{site.baseurl}}/concepts/architecture).

## Authorization Keys

The extensible authorization system consistes of a large set of keys. Each key uniquely represents a specific capability, for example, the ability to view a [Connection](/concepts/logical/connections), edit or delete a Connection. With the help of these keys, the system evaluates the permissions during runtime and renders UI both helping offer a secure management system and a customizable user experience.

{% include alert.html type="info" title="Note" content="The extensible authorization system is available to both Local and Remote Providers. Depending upon your chosen <a href='/extensibility/providers'>Remote Provider</a>, keys, clustering of them, assigning them to user groups, not just individual users or to user roles may be offered." %}

## Authorization Framework

Meshery utilizes CASL (JS-based permission framework) to evaluate any given user's set of session keys against the built-in keyhooks populated through each invidual Meshery UI page. This allows for granular control over the UI, empowering you to tailor your Meshery experience to your organization's needs by limiting access to specific features and functionalities based on the user's assigned keys.

<a href="/assets/img/permission-in-UI.png">
  <img style="width:min(100%,800px)" src="/assets/img/permission-in-UI.png" />
</a>

### Introduction to CASL.js

[CASL.js](https://casl.js.org) is an isomorphic authorization JavaScript library which restricts what resources a given client is allowed to access. It's designed to be incrementally adoptable and can easily scale between a simple claim based and fully featured subject and attribute based authorization. It makes it easy to manage and share permissions/keys across UI components, API services, and database queries.

An example of how CASL evaluate permission in UI.
{% capture code_content %}<React.Fragment>
	{!CAN(keys.DELETE_CONNECTION.action, keys.DELETE_CONNECTION.subject) && (
		<Button id="delete-connection">Delete<Button/>
	)}
</React.Fragment> 
{% endcapture %}
{% include code.html code=code_content %}

 Once a user has logged in, the backend will send a response with the permissions that the user has. Those permissions will be used to create abilities on the frontend, CASL gets updated with those abilities. The UI maintains a constant file containing all allowed permissions, referred to as keys. With the help of these keys, CAN function evaluates the permissions during runtime and renders UI.

<div class="alert alert-dark" role="alert">
<h4 class="alert-heading">Note</h4>

It's important to understand not all pages uses CASL authorization, means even if you are not assigned with any role within organization you might access preferences page and Meshery UI dashboard.
</div>

## Authorization using Local Provider

Meshery's built-in identity provider, "Local" Provider, operates with a large set of predefined keys interspersed throughout Meshery UI and persisted in [Meshery Database](/concepts/architecture/database). These keys are used to evaluate the permissions of a given user and render the UI accordingly. The keys are grouped into three categories: `action`, `subject`, and `object`.

{% include discuss.html %}

