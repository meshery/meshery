---
layout: default
title: Reference Permissions
permalink: reference/permissions
language: en
type: Reference
abstract: List of default permissions.
---


Permissions are represented as <b>keys</b>, each serving as a unique identifier for a specific permission. One or more keys can be grouped together and assigned to a <b>keychain</b>.<br>
<h6>Learn more about <a style="color: #00d3a9" onmouseover="this.style.color='#00b39f';" onmouseout="this.style.color='#00d3a9';" href="https://docs.layer5.io/cloud/security/keys/">Keys</a> and <a style="color: #00d3a9" onmouseover="this.style.color='#00b39f';" onmouseout="this.style.color='#00d3a9';" href="https://docs.layer5.io/cloud/security/keychains/">Keychains</a>.</h6>

{% include alert.html type="info" title="Customizable Permissions" content="Default permissions can be easily customized by simply creating your own <a href='https://docs.layer5.io/cloud/security/keychains/#keychains-management'>keychains</a> and <a href='https://docs.layer5.io/cloud/security/roles/'>roles</a>."
%}

<div style="overflow-x:auto;">
{% include_cached permissions.html %}
</div>
