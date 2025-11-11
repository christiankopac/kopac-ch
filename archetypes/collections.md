+++
title = "{{ replace .File.ContentBaseName "-" " " | title }}"
date = {{ .Date }}
draft = true
description = ""
+++

Use the `collection` shortcode to display structured data:

{{</* collection file="your-data.toml" style="card" */>}}

Available styles:
- `card` - Full cards with images
- `simple-card` - Simple cards without images
- `list` - Minimal list view
- `card-grid` - Grid of cards with images
- `card-horizontal` - Horizontal card layout
- `inline` - Inline list layout
- `poster-grid` - Compact poster grid (for media)

Optional filters:
- `filter="year"` - Enable year-based filtering
- `filter="date"` - Enable date-based filtering
- `category_filter="true"` - Enable category filtering

Data files go in the `data/` directory.

