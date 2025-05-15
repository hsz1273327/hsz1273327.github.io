---
title: "Ubuntu折腾记"
series_name: "ubuntu_experiment"
description: "记录折腾ubuntu的过程"
date: 2025-05-15
author: "Hsz"
category: experiment
tags:
    - Linux
update: 2025-05-15
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
