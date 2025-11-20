# How to Use the Go Scripts

All scripts are written in Go for consistency with Hugo. You can either build them or run them directly with `go run`.

## Quick Start

### Option 1: Run Directly (No Build Required)

```bash
# Download movie metadata and posters
go run scripts/download_movie_metadata.go -update-toml

# Download music/album metadata
go run scripts/download_music_metadata.go -update-toml

# Download book metadata
go run scripts/download_book_metadata.go -update-toml

# Create missing review pages
go run scripts/create_missing_reviews.go
```

### Option 2: Build First (Faster for Repeated Use)

```bash
# Build all scripts
go build -o scripts/download_movie_metadata scripts/download_movie_metadata.go
go build -o scripts/download_music_metadata scripts/download_music_metadata.go
go build -o scripts/download_book_metadata scripts/download_book_metadata.go
go build -o scripts/create_missing_reviews scripts/create_missing_reviews.go

# Then run them
./scripts/download_movie_metadata -update-toml
./scripts/download_music_metadata -update-toml
./scripts/download_book_metadata -update-toml
./scripts/create_missing_reviews
```

## Setup

### 1. Install Go Dependencies

```bash
go mod tidy
```

This will download the required package (`github.com/joho/godotenv`).

### 2. Set API Keys

Create a `.env` file in the project root:

```bash
# .env
TMDB_API_KEY=your_tmdb_api_key_here
DISCOGS_USER_TOKEN=your_discogs_token_here
# Note: Open Library (for books) doesn't require an API key
```

Or set as environment variables:
```bash
export TMDB_API_KEY="your_key"
export DISCOGS_USER_TOKEN="your_token"
```

## Scripts

### download_movie_metadata

Fetches movie metadata and downloads posters from TMDB.

```bash
# Process all movies in consumed.toml
go run scripts/download_movie_metadata.go

# Process specific movies
go run scripts/download_movie_metadata.go "Movie Title" "Another Movie"

# Update consumed.toml with metadata
go run scripts/download_movie_metadata.go -update-toml

# Skip movies that already have posters and directors
go run scripts/download_movie_metadata.go -skip-existing

# Combine options
go run scripts/download_movie_metadata.go -update-toml -skip-existing
```

**What it does:**
- Searches TMDB for each movie
- Downloads poster images to `static/screenshots/`
- Fetches director, year, and TMDB URL
- Updates `consumed.toml` with metadata (if `-update-toml` flag is used)

### download_music_metadata

Fetches album metadata from Discogs.

```bash
# Process all albums in consumed.toml
go run scripts/download_music_metadata.go

# Process specific albums
go run scripts/download_music_metadata.go "Album Title" "Another Album"

# Update consumed.toml with metadata
go run scripts/download_music_metadata.go -update-toml

# Skip albums that already have all metadata
go run scripts/download_music_metadata.go -skip-existing
```

**What it does:**
- Searches Discogs for each album
- Gets first release information
- Extracts artist, year, label, and Discogs URL
- Updates `consumed.toml` with metadata (if `-update-toml` flag is used)

### download_book_metadata

Fetches book metadata from Open Library.

```bash
# Process all books in consumed.toml
go run scripts/download_book_metadata.go

# Process specific books
go run scripts/download_book_metadata.go "Book Title" "Another Book"

# Update consumed.toml with metadata
go run scripts/download_book_metadata.go -update-toml

# Skip books that already have all metadata
go run scripts/download_book_metadata.go -skip-existing
```

**What it does:**
- Searches Open Library for each book
- Gets publication details
- Extracts author, year, publisher, and Open Library URL
- Updates `consumed.toml` with metadata (if `-update-toml` flag is used)
- **No API key required** - Open Library is free and open!

### create_missing_reviews

Creates placeholder review pages for movies that don't have review pages yet.

```bash
go run scripts/create_missing_reviews.go
```

**What it does:**
- Reads `consumed.toml` for all movies with review links
- Creates minimal review pages in `content/reviews/`
- Review pages automatically use metadata from `consumed.toml`

## Typical Workflow

1. **Add a movie to consumed.toml:**
   ```toml
   [[collection]]
   title = "New Movie"
   link = "/reviews/new-movie"
   category = "movies"
   ```

2. **Create the review page:**
   ```bash
   go run scripts/create_missing_reviews.go
   ```

3. **Fetch metadata and poster:**
   ```bash
   go run scripts/download_movie_metadata.go -update-toml
   ```

4. **Write your review:**
   Edit `content/reviews/new-movie.md` - metadata is automatically pulled from `consumed.toml`!

## Notes

- All metadata is stored in `consumed.toml` (single source of truth)
- Review pages only contain the review content
- The template automatically displays metadata from `consumed.toml`
- No need to sync metadata between files!

