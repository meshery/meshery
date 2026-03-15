---
layout: default
title: Permission Keys Reference 
permalink: reference/permissions
language: en
type: Reference
abstract: List of default permissions.
---


Permissions are represented as <b>keys</b>, each serving as a unique identifier for a specific permission. One or more keys can be grouped together and assigned to a <b>keychain</b>.<br>

{% include alert.html type="info" title="Customizable Permissions" content="Default permissions can be easily customized by simply creating your own keychains and roles."
%}

{% include alert.html type="info" title="Contributing to Permission Keys" content="To contribute permission keys, see the <a href='https://docs.meshery.io/extensibility/authorization'>Authorization Extensibility Documentation</a>. This guide explains Meshery’s extensible authorization system." %}

<div style="overflow-x:auto;">
{% include_cached permissions.html %}
</div>
