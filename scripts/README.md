# Scripts

All scripts are written in Go for consistency with Hugo.

## Building

```bash
# Build all scripts
go build -o scripts/download_movie_metadata scripts/download_movie_metadata.go
go build -o scripts/download_music_metadata scripts/download_music_metadata.go
go build -o scripts/download_book_metadata scripts/download_book_metadata.go
go build -o scripts/create_missing_reviews scripts/create_missing_reviews.go

# Or run directly without building
go run scripts/download_movie_metadata.go -update-toml
go run scripts/download_music_metadata.go -update-toml
go run scripts/download_book_metadata.go -update-toml
go run scripts/create_missing_reviews.go
```

## download_movie_metadata

Downloads movie posters and fetches metadata (year, director, TMDB URL) from TMDB API.

### Setup

1. Get a TMDB API key from https://www.themoviedb.org/settings/api
2. Set it as an environment variable or in a `.env` file:
   ```bash
   export TMDB_API_KEY="your_api_key_here"
   ```
   
   Or create `.env` file in project root:
   ```
   TMDB_API_KEY=your_api_key_here
   ```

3. Install Go dependencies (if not already done):
   ```bash
   go mod tidy
   ```

### Usage

```bash
# Download posters and metadata for all movies in consumed.toml
go run scripts/download_movie_metadata.go

# Or if built:
./scripts/download_movie_metadata

# Download for specific movies
go run scripts/download_movie_metadata.go "Little Trouble Girls" "Bunny"

# Update consumed.toml with metadata
go run scripts/download_movie_metadata.go -update-toml

# Skip movies that already have posters and directors
go run scripts/download_movie_metadata.go -skip-existing

# Combine options
go run scripts/download_movie_metadata.go -update-toml -skip-existing
```

### What it does

1. Searches TMDB for each movie
2. Downloads poster images to `static/screenshots/`
3. Fetches director information and year
4. Gets TMDB URL
5. Optionally updates `consumed.toml` with year, director, and TMDB URL

### Output

- Posters saved as: `{movie_slug}_poster.jpg` in `static/screenshots/`
- Metadata can be written back to `consumed.toml` with `-update-toml` flag

## download_music_metadata

Downloads album metadata (artist, year, label, Discogs URL) from Discogs API for music in consumed.toml.

### Setup

1. Get a Discogs personal access token from https://www.discogs.com/settings/developers
2. Set it as an environment variable or in a `.env` file:
   ```bash
   export DISCOGS_USER_TOKEN="your_token_here"
   ```
   
   Or create `.env` file in project root:
   ```
   DISCOGS_USER_TOKEN=your_token_here
   ```

3. Install Go dependencies (if not already done):
   ```bash
   go mod tidy
   ```

### Usage

```bash
# Download metadata for all albums in consumed.toml
go run scripts/download_music_metadata.go

# Or if built:
./scripts/download_music_metadata

# Download for specific albums
go run scripts/download_music_metadata.go "Album Title" "Another Album"

# Update consumed.toml with fetched metadata
go run scripts/download_music_metadata.go -update-toml

# Skip albums that already have all metadata
go run scripts/download_music_metadata.go -skip-existing
```

### What it does

1. Searches Discogs for each album
2. Gets the first release information
3. Extracts artist, year, and label
4. Gets Discogs URL
5. Optionally updates `consumed.toml` with metadata

### Output

- Metadata saved to `consumed.toml`:
  - `artist` - Album artist
  - `year` - Release year (first release)
  - `label` - Record label
  - `discogs` - Discogs URL

## download_book_metadata

Downloads book metadata (author, year, publisher, Open Library URL) from Open Library API for books in consumed.toml.

### Setup

1. No API key required - Open Library is free and open!

2. Install Go dependencies (if not already done):
   ```bash
   go mod tidy
   ```

### Usage

```bash
# Download metadata for all books in consumed.toml
go run scripts/download_book_metadata.go

# Or if built:
./scripts/download_book_metadata

# Download for specific books
go run scripts/download_book_metadata.go "Book Title" "Another Book"

# Update consumed.toml with fetched metadata
go run scripts/download_book_metadata.go -update-toml

# Skip books that already have all metadata
go run scripts/download_book_metadata.go -skip-existing
```

### What it does

1. Searches Open Library for each book
2. Gets publication details
3. Extracts author, year, and publisher
4. Gets Open Library URL
5. Optionally updates `consumed.toml` with metadata

### Output

- Metadata saved to `consumed.toml`:
  - `author` - Book author
  - `year` - Publication year
  - `publisher` - Publisher name
  - `openlibrary` - Open Library URL

## create_missing_reviews

Creates placeholder review pages for movies in consumed.toml that don't have review pages yet.
Creates minimal review pages - all metadata comes from consumed.toml automatically.

### Usage

```bash
go run scripts/create_missing_reviews.go

# Or if built:
./scripts/create_missing_reviews
```

This ensures all movie links work even if reviews haven't been written yet.
The review pages will automatically display metadata (title, director, year, rating) from consumed.toml.

## Workflow: Single Source of Truth

**All metadata is stored in `consumed.toml`** - this is the single source of truth.

Review pages (`content/reviews/*.md`) contain **only the review content** - no duplicate metadata.

The template automatically pulls metadata (title, year, director, rating) from `consumed.toml` based on the review page slug.

### Benefits

- ✅ Edit metadata in one place (`consumed.toml`)
- ✅ Review pages are simple - just the review text
- ✅ No sync issues between files
- ✅ Title and links are automatically constructed from metadata

### Example

**consumed.toml:**
```toml
[[collection]]
title = "Lurker"
year = "2025"
director = "Alex Russell"
rating = 4
link = "/reviews/lurker"
```

**content/reviews/lurker.md:**
```markdown
+++
date = 2025-11-20
draft = false
+++

My review of Lurker goes here...
```

The template automatically:
- Uses title from `consumed.toml`
- Shows director, year, rating from `consumed.toml`
- No need to duplicate this in the review page!
