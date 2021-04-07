---
title: 
sitemap: false
permalink: /search/
not_editable: true
excluded_in_search: true
---

<input class="form-control td-search-input fa fa-search" type="search" name="q" id="search-input" placeholder="&#xf002 Search this site…"  style="margin-top:5px" autofocus>
<i style="color:white; margin-right:8px; margin-left:5px" class="fa fa-search"></i>

<p><span id="search-process">Loading</span> results <span id="search-query-container" style="display: none;">for "<strong id="search-query"></strong>"</span></p>

<ul id="search-results"></ul>

<script>
	window.data = {
		{% for post in site.docs %}
				{% unless post.excluded_in_search %}
					{% if added %},{% endif %}
					{% assign added = false %}
					"{{ post.url | slugify }}": {
						"id": "{{ post.url | slugify }}",
						"title": "{{ post.title | xml_escape }}",
						"categories": "{{ post.categories | join: ", " | xml_escape }}",
						"url": " {{ post.url | xml_escape }}",
						"content": {{ post.content | strip_html | replace_regex: "[\s/\n]+"," " | strip | jsonify }}
					}
					{% assign added = true %}
				{% endunless %}
		{% endfor %}
		{% for post in site.pages %}
				{% unless post.excluded_in_search %}
					{% if added %},{% endif %}
					{% assign added = false %}
					"{{ post.url | slugify }}": {
						"id": "{{ post.url | slugify }}",
						"title": "{{ post.title | xml_escape }}",
						"categories": "{{ post.categories | join: ", " | xml_escape }}",
						"url": " {{ post.url | xml_escape }}",
						"content": {{ post.content | strip_html | replace_regex: "[\s/\n]+"," " | strip | jsonify }}
					}
					{% assign added = true %}
				{% endunless %}
		{% endfor %}
	};
</script>
<script src="{{ site.baseurl }}/assets/js/lunr.min.js"></script>
<script src="{{ site.baseurl }}/assets/js/search.js"></script>
