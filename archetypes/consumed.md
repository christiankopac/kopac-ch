+++
title = "{{ replace .File.ContentBaseName "-" " " | title }}"
date = {{ .Date }}
draft = true
description = ""
+++

Log of consumed media (books, movies, podcasts, etc.)

Use the `collection` shortcode with category filter:

{{</* collection file="consumed.toml" style="poster-grid" category_filter="true" */>}}

