---
title: Releases
description: List of released Meshery versions and their release notes.
categories: [project]
---

See also Meshery's [Compatibility Matrix](/installation) and [Build and Release](/project/contributing/build-and-release) process.

{% assign sorted_release = site.releases | sort: 'date' | reverse %}

<table>
<tr><th> Version </th><th> Date </th></tr>

{% for releases in sorted_release %}
    {% assign releasedate = releases.date | date: "%B %d, %Y" %}
    <tr>
        <td style="text-align:center"> <a href="{{site.baseurl}}/{{page.permalink}}/{{ releases.tag }}">{{ releases.tag }}</a> </td>
        <td style="text-align:center"> {{releasedate}} </td>
    </tr>
{% endfor %}

</table>
