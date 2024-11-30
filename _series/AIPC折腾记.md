---
title: "AIPC折腾记"
series_name: "aipc_experiment"
description: ""
date: 2024-11-06
author: "Hsz"
category: experiment
tags:
    - AIPC
    - AI
update: 2024-11-27
---

{% assign series_posts = "" | split: "" %}
{% for post in site.posts %}
    {% if post.series contains page.series_name %}
    {% assign series_posts=series_posts | push: post %}
    {% endif %}
{% endfor %}
{% assign indexed_posts = series_posts |sort: series.aipc_experiment.index %}

{% for post in indexed_posts %}
+ [{{ post.title }}]({{site.url}}{{ post.url }})
{% endfor %}
