---
layout: default
title: "Extensible Authorization: UI"
permalink: extensibility/extensible-authorization
type: Extensibility
abstract: "Meshery offers support for more adapters than any other project or product in the world. Meshery UI has a number of extension points that allow for users to customize their experience with third-party plugins."
language: en
list: include
---

Meshery now features an extensible authorization system that provides fine-tuned access control to Meshery's features. This system ensures users have secure and controlled access to Meshery's rich set of features and functionalities in context of your organization. Example of authorization: Ability to deploy design in cluster or Ability to delete any resource in cluster.

## What is permission

In Meshery, permissions are represented as keys, each serving as a unique identifier for a specific permission. For example, you can allow any users to view and manage connections but not to delete them. One or more keys can be grouped together and assigned to a keychain. Then this keychain can be assigned to a role and that role can be assigned to a user. This is the general flow of how keys are assigned to a user. You can learn more about keys and roles about here, [Learn more](https://docs.layer5.io/cloud/security/keys/)

Meshery utilizes CASL (JS-based permission framework) to enforce these permissions within its user interface.

<a href="/assets/img/permission-in-UI.png">
  <img style="width:min(100%,800px)" src="/assets/img/permission-in-UI.png" />
</a>

## Introduction to CASL.js

CASL is an isomorphic authorization JavaScript library which restricts what resources a given client is allowed to access. It's designed to be incrementally adoptable and can easily scale between a simple claim based and fully featured subject and attribute based authorization. It makes it easy to manage and share permissions/keys across UI components, API services, and database queries.

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


## Authorization for Local provider

While Meshery offers flexibility with various providers, the Local provider operates with predefined keys stored in a local database. Remember this limits customization of permissions.

## Summary

Extensible authorization in Meshery offers granular control over its features, empowering you to tailor your Meshery experience to your organization's needs. The system ensures secure access to Meshery's extensive capabilities, providing flexibility in managing permissions, fine-tuning access, and customizing functionalities. With Meshery's extensible authorization, you gain the ability to define and enforce permissions with precision, enhancing the overall security and control of your cloud-native environment.

<div class="alert alert-dark" role="alert">
<h4 class="alert-heading">Discussion Forum</h4>
Not finding what you're looking for? Ask on the <a href="http://discuss.meshery.io">Discussion Forum</a>.
</div>

# Additional Resources

- [CASLjs Documentation](https://casl.js.org/v4/en/guide/intro/)
