# One Page Per Item - Proposal

## Current Structure (Data Files + Pages)

```
data/movies/movies.toml  →  All metadata
content/consumed/movie/lurker.md  →  Minimal frontmatter only
```

**Problems:**
- Two places to maintain
- Data file is separate from content
- Harder to find/edit specific items

## Proposed Structure (Self-Contained Pages)

```
content/consumed/movie/lurker.md  →  ALL metadata in frontmatter + content
```

**Benefits:**
- ✅ One file per item - easy to find
- ✅ All metadata visible in one place
- ✅ Self-contained - no external data file
- ✅ Scripts update frontmatter directly
- ✅ Git diffs are clearer (one file per change)
- ✅ Can use Hugo's built-in content management

## Example Structure

### Movie Page (`content/consumed/movie/lurker.md`)

```toml
+++
title = "Lurker"
date = 2025-11-20
draft = false

# Movie metadata
year = "2025"
director = "Alex Russell"
rating = 4
tmdb = "https://www.themoviedb.org/movie/1264573"
img = "/images/movies/lurker_poster.jpg"
trailer = "https://youtube.com/..."
screenshots = "screenshots-lurker"  # Reference to gallery data
processed = true
footer = "Watched Nov 2025"

# Content
content = "Movie description here..."
review = "Full review text here..."
+++

# Optional: Additional markdown content below
```

### Book Page (`content/consumed/book/how-to-take-smart-notes.md`)

```toml
+++
title = "How to Take Smart Notes"
date = 2025-10-01
draft = false

# Book metadata
author = "Sönke Ahrens"
year = "2022"
publisher = "Sönke Ahrens"
openlibrary = "https://openlibrary.org/isbn/9783982438818"
img = "/images/books/how_to_take_smart_notes_cover.jpg"
date_read = "2025-10"
processed = true

# Content
content = "Book description..."
review = "My review..."
+++
```

## Script Updates

Scripts would update frontmatter directly:

```go
// Update movie page frontmatter
func updateMoviePage(filepath string, data MovieData) error {
    // Read file
    // Parse frontmatter
    // Update fields
    // Write back
}
```

## Gallery Listings

Gallery shortcodes would read from content pages instead of data files:

```go
// Read all pages in content/consumed/movie/
// Extract frontmatter
// Build collection from pages
```

## Migration Path

1. **Create script** to migrate data files → content pages
2. **Update scripts** to modify frontmatter instead of TOML
3. **Update templates** to read from `.Params` instead of `.Site.Data`
4. **Update shortcodes** to read from pages

## Comparison

| Aspect | Current (Data Files) | Proposed (Pages) |
|--------|---------------------|------------------|
| **Edit location** | `data/movies/movies.toml` | `content/consumed/movie/lurker.md` |
| **Find item** | Search in large TOML file | One file per item |
| **Git diffs** | Large diffs in one file | One file per change |
| **Script updates** | Update TOML blocks | Update frontmatter |
| **Self-contained** | No | Yes |
| **Hugo native** | Uses data files | Uses content pages |

## Recommendation

**Yes, this is better!** One page per item is:
- More intuitive
- Easier to maintain
- Better for version control
- More aligned with Hugo's content model

