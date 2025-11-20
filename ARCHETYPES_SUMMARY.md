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

### 3. `reviews.md` ✅
**Usage**: `hugo new content/reviews/movie-name.md`

Movie review template with minimal frontmatter - **all metadata comes from `consumed.toml`**.

**Important**: The filename (slug) must match the link in `consumed.toml`.
- If `consumed.toml` has `link = "/reviews/lurker"`, the file should be `lurker.md`
- Metadata (title, year, director, rating, tmdb) is automatically pulled from `consumed.toml`
- You only write the review content in the markdown file

**Frontmatter**:
```toml
+++
# Minimal frontmatter - all metadata comes from consumed.toml
date = 2024-01-01
draft = false
+++

Your review content goes here...
```

**Workflow**:
1. Add movie to `consumed.toml` with `link = "/reviews/movie-slug"`
2. Run `go run scripts/create_missing_reviews.go` to create placeholder review page
3. Write review content in `content/reviews/movie-slug.md`
4. Metadata (title, director, year, rating) is automatically displayed from `consumed.toml`

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

## Key Workflow: Single Source of Truth

**All metadata is stored in `consumed.toml`** - this is the single source of truth.

- **Movies**: Edit metadata in `consumed.toml`, write reviews in `content/reviews/*.md`
- **Music**: Edit metadata in `consumed.toml` (artist, year, label)
- **Books**: Edit metadata in `consumed.toml` (author, year)

Review pages contain **only the review content** - no duplicate metadata needed!

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
# Metadata comes from consumed.toml, not frontmatter
+++
date = 2024-01-01
draft = false
+++
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

## Helper Scripts (Go)

All scripts are written in Go for consistency with Hugo:

```bash
# Build scripts
go build -o scripts/download_movie_metadata scripts/download_movie_metadata.go
go build -o scripts/download_music_metadata scripts/download_music_metadata.go
go build -o scripts/download_book_metadata scripts/download_book_metadata.go
go build -o scripts/create_missing_reviews scripts/create_missing_reviews.go

# Or run directly
go run scripts/download_movie_metadata.go -update-toml
go run scripts/download_music_metadata.go -update-toml
go run scripts/download_book_metadata.go -update-toml
go run scripts/create_missing_reviews.go
```

**Scripts**:
- `download_movie_metadata.go` - Fetch movie metadata and posters from TMDB
- `download_music_metadata.go` - Fetch album metadata from Discogs
- `download_book_metadata.go` - Fetch book metadata from Open Library
- `create_missing_reviews.go` - Create placeholder review pages from consumed.toml

All archetypes follow Hugo conventions and use correct TOML frontmatter syntax! ✨

