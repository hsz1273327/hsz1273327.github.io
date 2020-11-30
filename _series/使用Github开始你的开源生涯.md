---
title: "使用Github开始你的开源生涯"
series_name: "get_along_well_with_github"
description: "Github是全球最流行的开源项目托管平台,也是全球最大的IT领域同性交友平台,本系列将介绍如何利用Github参与开源项目并介绍Github套件的开源可替代方案."
date: 2020-10-29
author: "Hsz"
category: recommend
tags:
    - github
    - Open Source
update: 2020-11-01
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
