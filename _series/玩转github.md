---
title: "玩转Github"
series_name: "get_along_well_with_github"
description: "Github是全球最流行的开源项目托管平台,也是全球最大的IT领域同性交友平台(😁),本系列将介绍如何利用github参与开源项目和结交大佬."
date: 2019-03-14
author: "Hsz"
category: recommend
tags:
    - github
update: 2019-03-14
---

{% assign series_posts = "" | split: "" %}
{% for post in site.posts %}
    {% if post.series contains page.series_name %}
    {% assign series_posts=series_posts | push: post %}
    {% endif %}
{% endfor %}
{% assign indexed_posts = series_posts |sort: "series.get_along_well_with_github.index" %}

{% for post in indexed_posts %}
+ [{{ post.title }}]({{site.url}}{{ post.url }})
{% endfor %}
