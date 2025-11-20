# Migration to One Page Per Item

## What Changed

We've migrated from data files to self-contained content pages where each item (movie, book, music) has all its metadata in the page frontmatter.

## Migration Steps

### 1. Run Migration Script

```bash
go run scripts/migrate_to_pages.go
```

This will:
- Read from `data/movies/movies.toml`, `data/music/music.toml`, `data/books/books.toml`
- Create pages in `content/consumed/movie/`, `content/consumed/music/`, `content/consumed/book/`
- Convert all metadata to frontmatter
- Skip existing pages (won't overwrite)

### 2. Review Generated Pages

Check the generated pages to ensure metadata was migrated correctly:

```bash
# Example
cat content/consumed/movie/lurker.md
```

### 3. Test the Site

```bash
hugo server -D
```

Verify:
- ✅ Consumed page gallery shows all items
- ✅ Individual consumed pages display correctly
- ✅ Screenshot galleries still work
- ✅ Picks page shows items correctly

## What's Updated

### Templates
- `themes/apacible/layouts/_default/single.html` - Reads from `.Params` instead of `.Site.Data`
- Screenshot galleries still work (they reference data files, which is fine)

### Shortcodes
- `themes/apacible/layouts/shortcodes/all-consumed.html` - Reads from pages instead of data files
- `themes/apacible/layouts/shortcodes/picks.html` - Reads from pages instead of data files

### Scripts
- `scripts/migrate_to_pages.go` - New migration script
- Download scripts still need to be updated (TODO)

## Page Structure

Each page now looks like:

```toml
+++
title = "Lurker"
date = 2025-11-20
draft = false
category = "movie"
year = "2025"
director = "Alex Russell"
rating = 4
img = "/images/movies/lurker_poster.jpg"
tmdb = "https://www.themoviedb.org/movie/1264573"
screenshots = "screenshots-lurker"
footer = "Watched Nov 2025"
processed = true
content = "Description..."
review = "Review text..."
+++
```

## Screenshot Galleries

Screenshot galleries still work! They:
- Reference data files in `data/movies/screenshots-*.toml`
- Can be specified in page frontmatter with `screenshots = "screenshots-lurker"`
- Fall back to auto-detection based on page slug

## Next Steps

1. ✅ Migration script created
2. ✅ Templates updated
3. ✅ Shortcodes updated
4. ⏳ Update download scripts to modify frontmatter
5. ⏳ Test everything works

## Benefits

- ✅ One file per item - easy to find and edit
- ✅ All metadata in one place
- ✅ Self-contained pages
- ✅ Better Git diffs (one file per change)
- ✅ More Hugo-native approach

