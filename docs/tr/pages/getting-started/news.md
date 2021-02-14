---
title: News
permalink: tr/news/
language: tr
lang: tr
categories: tr
excluded_in_search: true
---

# Haberler

<p>İle abone ol <a href="{{ site.baseurl }}/feed.xml">RSS</a>en son haberleri takip etmek için.
Site değişiklikleri için bkz. <a href="https://github.com/{{ site.github_user }}/{{ site.github_repo }}/blob/master/CHANGELOG.md">
değişim günlüğü</a> kod tabanında tutulur.</p>

<br>

{% for post in site.posts limit:10 %}
   <div class="post-preview">
   <h2><a href="{{ site.baseurl }}{{ post.url }}">{{ post.title }}</a></h2>
   <span class="post-date">{{ post.date | date: "%B %d, %Y" }}</span><br>
   {% if post.badges %}{% for badge in post.badges %}<span class="badge badge-{{ badge.type }}">{{ badge.tag }}</span>{% endfor %}{% endif %}
   {{ post.content | split:'<!--more-->' | first }}
   {% if post.content contains '<!--more-->' %}
      <a href="{{ site.baseurl }}{{ post.url }}">read more</a>
   {% endif %}
   </div>
   <hr>
{% endfor %}

Daha fazlasını görmek ister misin? Bakın<a href="{{ site.baseurl }}/tr/archive/">Haber Arşivi</a>.
