# Single Source of Truth - Proposal

## Current Situation

Right now you have:
1. **Data files** (`data/movies/movies.toml`, etc.) - metadata, content, reviews
2. **Content pages** (`content/consumed/movie/*.md`) - mostly empty, just frontmatter
3. **Review pages** (`content/reviews/*.md`) - review content only

This means editing in multiple places.

## Proposed Solution: Data Files as Single Source

### Structure

**Everything in data files:**
```toml
[[collection]]
title = "Lurker"
year = "2025"
director = "Alex Russell"
rating = 4
content = "Description here..."
review = "Full review text here..."
img = "/images/movies/lurker_poster.jpg"
link = "/consumed/movie/lurker"
screenshots = "screenshots-lurker"  # Reference to gallery data
trailer = "https://youtube.com/..."
```

**Auto-generated content pages:**
- Minimal frontmatter only
- All content comes from data file
- Pages exist only for routing

### Benefits

✅ **One place to edit** - Everything in `data/movies/movies.toml`
✅ **No sync issues** - No duplicate data
✅ **Easy to maintain** - Edit metadata, content, reviews all in one place
✅ **Auto-generated pages** - Script creates pages from data

### Implementation

1. **Update template** to prioritize data file content over page content
2. **Update create_consumed_pages.go** to sync pages with data
3. **Store everything in data files:**
   - Metadata (title, year, director, etc.)
   - Content/description
   - Review text
   - Gallery references
   - Trailer links

### Workflow

1. Edit `data/movies/movies.toml` - add/edit everything
2. Run `go run scripts/create_consumed_pages.go` - syncs pages
3. Build site - pages auto-populate from data

### Optional: Keep Review Pages Separate

If reviews get long, you can:
- Keep review in data file for short reviews
- Or reference external markdown file from data
- Template checks data first, then page content

## Recommendation

**Use data files as single source of truth:**
- All metadata, content, reviews in data files
- Content pages are auto-generated and minimal
- Template prioritizes data file content
- Optional: Allow page content to override/extend data content

