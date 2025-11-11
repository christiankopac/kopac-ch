# Archetypes Summary

All archetypes are now correct and using proper Hugo syntax.

## Available Archetypes

### 1. `default.md`
**Usage**: `hugo new content/page-name.md`

Basic page template with minimal frontmatter.

### 2. `posts.md` ✅
**Usage**: `hugo new content/posts/my-post.md`

Blog post template with:
- Title, date, draft status, description
- `[taxonomies]` section for tags and categories
- Optional features: math, mermaid, toc

**Frontmatter**:
```toml
+++
title = "Post Title"
date = 2024-01-01
draft = true
description = ""

[taxonomies]
tags = []
categories = []
+++
```

### 3. `reviews.md` ✅ (NEW)
**Usage**: `hugo new content/reviews/movie-name.md`

Movie review template with:
- Title, date, draft status, description
- `[params]` section for director, year, rating
- Spoiler gallery shortcode example

**Frontmatter**:
```toml
+++
title = "Movie Title"
date = 2024-01-01
draft = true
description = ""

[params]
director = "Director Name"
year = 2010
rating = 5
+++
```

### 4. `collections.md` ✅
**Usage**: `hugo new content/collections/collection-name.md`

Collection page template with:
- Instructions for using collection shortcode
- List of available styles (card, poster-grid, inline, etc.)
- Filter options (year, date, category)

### 5. `consumed.md` ✅
**Usage**: `hugo new content/consumed/_index.md`

Media consumption log template with:
- Example using `consumed.toml` with poster-grid style
- Category filter enabled

## Key Corrections Made

### ❌ Before (Zola Syntax):
```toml
template = "prose.html"  # Wrong

[extra]  # Wrong - This is Zola
director = "Name"
```

### ✅ After (Hugo Syntax):
```toml
# No template field needed

[params]  # Correct - This is Hugo
director = "Name"
```

## Data Files Structure

All collection data files should be in `data/`:
- `consumed.toml` - Combined media (movies, books, music)
- `projects.toml` - Projects collection
- `skills.toml` - Skills collection
- `*-screenshots.toml` - Screenshot galleries for reviews

## Creating New Content

```bash
# New blog post
hugo new content/posts/my-new-post.md

# New movie review
hugo new content/reviews/inception.md

# New collection
hugo new content/collections/bookmarks.md
```

All archetypes follow Hugo conventions and use correct TOML frontmatter syntax! ✨

