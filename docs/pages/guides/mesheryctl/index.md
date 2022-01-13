---
layout: default
title: Meshery CLI Guides
permalink: guides/mesheryctl
redirect_from: guides/mesheryctl/
language: en
type: Guides
category: mesheryctl
list: exclude
---

Guides to using Meshery's various features and components.

{% capture tag %}
<li><a href="{{ site.baseurl }}/guides/mesheryctl/working-with-mesheryctl">Using mesheryctl</a></li>
<li><a href="{{ site.baseurl }}/guides/mesheryctl/configuring-autocompletion-for-mesheryctl">Configure autocompletion for mesheryctl</a></li>
<li><a href="{{ site.baseurl }}/guides/mesheryctl/authenticate-with-meshery-via-cli">Authenticate via Meshery CLI</a></li>
<li><a href="{{ site.baseurl }}/guides/upgrade#upgrading-meshery-cli">Upgrading Meshery CLI</a></li>
<li><a href="{{ site.baseurl }}/guides/mesheryctl/running-system-checks-using-mesheryctl">Running system checks using Meshery CLI</a></li>
{% endcapture %}

{% include suggested-reading.html diffName="true" isDiffTag="true" diffTag=tag %}

<!-- {% include toc.html page=Guides %} -->

{:toc}

<!-- {% comment %}
#
#  Change date order by adding '| reversed'
#  To sort by title or other variables use {% assign sorted_posts = category[1] | sort: 'title' %}
#
{% endcomment %}

{% for guide in site.adapter %}
<h2 id="{{guide[0] | uri_escape | downcase }}">{{guide[0] | capitalize}}1</h2>

{% endfor %}

{% assign sorted_guides = site.guides | sort %}
{% for guide in sorted_guides %}
<h2 id="{{guide[0] | uri_escape | downcase }}">{{guide[0] | capitalize}}</h2>

{% endfor %} -->
