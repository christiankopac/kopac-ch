package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// MovieInfo represents a movie from markdown frontmatter
type MovieInfo struct {
	Title     string
	Year      string
	Director  string
	Processed bool
	Draft     bool
	FilePath  string
}

// MovieData represents fetched movie metadata to be written to frontmatter
type MovieData struct {
	Title      string
	Year       string
	Director   string
	PosterPath string
	PosterURL  string
	TMDBID     int
	TMDBURL    string
	ImagePath  string // Path for frontmatter img field
	TrailerURL string // YouTube trailer URL
	Draft      bool   // Mark as draft if no poster found
}

// parseMarkdownFiles reads all markdown files in content/consumed/movie/ and extracts movie info
func parseMarkdownFiles(contentDir string, includeDrafts bool) ([]MovieInfo, error) {
	movieDir := filepath.Join(contentDir, "consumed", "movie")
	
	// Check if directory exists
	if _, err := os.Stat(movieDir); err != nil {
		return nil, fmt.Errorf("movie directory not found: %s", movieDir)
	}

	var movies []MovieInfo
	
	files, err := os.ReadDir(movieDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(movieDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		contentStr := string(content)
		
		// Find frontmatter
		frontmatterStart := strings.Index(contentStr, "+++")
		if frontmatterStart == -1 {
			continue
		}
		
		frontmatterEnd := strings.Index(contentStr[frontmatterStart+3:], "+++")
		if frontmatterEnd == -1 {
			continue
		}
		
		frontmatter := contentStr[frontmatterStart+3 : frontmatterStart+3+frontmatterEnd]

		// Extract title
		titleMatch := regexp.MustCompile(`title\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(titleMatch) < 2 {
			continue
		}
		title := titleMatch[1]

		// Extract category - must be "movie"
		if !regexp.MustCompile(`category\s*=\s*"movie"`).MatchString(frontmatter) {
			continue
		}

		// Check draft flag
		draft := false
		if regexp.MustCompile(`draft\s*=\s*true`).MatchString(frontmatter) {
			draft = true
		}

		// Skip drafts unless explicitly including them
		if draft && !includeDrafts {
			continue
		}

		// Check processed flag
		processed := false
		if regexp.MustCompile(`processed\s*=\s*(true|yes)`).MatchString(frontmatter) {
			processed = true
		}

		// Extract year
		year := ""
		yearMatch := regexp.MustCompile(`year\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(yearMatch) >= 2 {
			year = yearMatch[1]
		}

		// Extract director
		director := ""
		directorMatch := regexp.MustCompile(`director\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(directorMatch) >= 2 {
			director = directorMatch[1]
		}

		// Check for missing metadata
		hasDirector := director != ""
		hasYear := year != ""
		hasTMDB := regexp.MustCompile(`tmdb\s*=\s*"[^"]+"`).MatchString(frontmatter)
		hasImg := regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(frontmatter)
		
		// Process if:
		// 1. processed is missing or false, OR
		// 2. metadata is incomplete (missing director, year, tmdb, or img)
		needsProcessing := !processed || !hasDirector || !hasYear || !hasTMDB || !hasImg
		
		if !needsProcessing {
			continue
		}

		movies = append(movies, MovieInfo{
			Title:     title,
			Year:      year,
			Director:  director,
			Processed: processed,
			Draft:     draft,
			FilePath:  filePath,
		})
	}

	return movies, nil
}

// updateMarkdownFrontmatter updates a markdown file's frontmatter with movie data
func updateMarkdownFrontmatter(filePath string, data MovieData) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	
	// Find frontmatter boundaries
	frontmatterStart := strings.Index(contentStr, "+++")
	if frontmatterStart == -1 {
		return fmt.Errorf("no frontmatter found")
	}
	
	frontmatterEnd := strings.Index(contentStr[frontmatterStart+3:], "+++")
	if frontmatterEnd == -1 {
		return fmt.Errorf("no closing frontmatter found")
	}
	frontmatterEnd += frontmatterStart + 3

	frontmatter := contentStr[frontmatterStart+3 : frontmatterEnd]
	body := contentStr[frontmatterEnd+3:]

	updated := false

	// Update year
	if data.Year != "" {
		if regexp.MustCompile(`year\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`year\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`year = "%s"`, data.Year))
			updated = true
		} else {
			// Insert after title
			if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nyear = \"%s\"", data.Year))
				updated = true
			}
		}
	}

	// Update director
	if data.Director != "" {
		if regexp.MustCompile(`director\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`director\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`director = "%s"`, data.Director))
			updated = true
		} else {
			// Insert after year or title
			if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ndirector = \"%s\"", data.Director))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ndirector = \"%s\"", data.Director))
				updated = true
			}
		}
	}

	// Update or add TMDB URL
	if data.TMDBURL != "" {
		if regexp.MustCompile(`tmdb\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`tmdb\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`tmdb = "%s"`, data.TMDBURL))
			updated = true
		} else {
			// Insert after director, rating, or title
			if regexp.MustCompile(`director\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(director\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ntmdb = \"%s\"", data.TMDBURL))
				updated = true
			} else if regexp.MustCompile(`rating\s*=\s*\d+`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(rating\s*=\s*\d+)`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ntmdb = \"%s\"", data.TMDBURL))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ntmdb = \"%s\"", data.TMDBURL))
				updated = true
			}
		}
	}

	// Update or add image path
	if data.ImagePath != "" {
		if regexp.MustCompile(`img\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`img\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`img = "%s"`, data.ImagePath))
			updated = true
		} else {
			// Insert after category or title
			if regexp.MustCompile(`category\s*=\s*"movie"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(category\s*=\s*"movie")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nimg = \"%s\"", data.ImagePath))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nimg = \"%s\"", data.ImagePath))
				updated = true
			}
		}
	}

	// Update or add trailer URL
	if data.TrailerURL != "" {
		if regexp.MustCompile(`trailer\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`trailer\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`trailer = "%s"`, data.TrailerURL))
			updated = true
		} else {
			// Insert after tmdb, img, director, or title
			if regexp.MustCompile(`tmdb\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(tmdb\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ntrailer = \"%s\"", data.TrailerURL))
				updated = true
			} else if regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(img\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ntrailer = \"%s\"", data.TrailerURL))
				updated = true
			} else if regexp.MustCompile(`director\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(director\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ntrailer = \"%s\"", data.TrailerURL))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ntrailer = \"%s\"", data.TrailerURL))
				updated = true
			}
		}
	}

	// Update draft status
	if data.Draft {
		if regexp.MustCompile(`draft\s*=\s*(true|false)`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`draft\s*=\s*(true|false)`).ReplaceAllString(frontmatter, "draft = true")
			updated = true
		} else {
			// Insert after category or title
			if regexp.MustCompile(`category\s*=\s*"movie"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(category\s*=\s*"movie")`).ReplaceAllString(frontmatter, "$1\ndraft = true")
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\ndraft = true")
				updated = true
			}
		}
	} else {
		// Remove draft if it exists and we have a poster
		if regexp.MustCompile(`draft\s*=\s*true`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`\ndraft\s*=\s*true\s*\n?`).ReplaceAllString(frontmatter, "\n")
			updated = true
		}
	}

	// Add or update processed flag (set to true after successful processing)
	if !regexp.MustCompile(`processed\s*=\s*(true|false)`).MatchString(frontmatter) {
		// Insert after img, director, year, or title
		if regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(img\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`director\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(director\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		}
	} else {
		// Update existing processed flag to true
		frontmatter = regexp.MustCompile(`processed\s*=\s*(true|false)`).ReplaceAllString(frontmatter, "processed = true")
		updated = true
	}

	if updated {
		newContent := "+++" + frontmatter + "+++" + body
		return os.WriteFile(filePath, []byte(newContent), 0644)
	}

	return nil
}

// AlbumInfo represents a music album from markdown frontmatter
type AlbumInfo struct {
	Title     string
	Artist    string
	Year      string
	Label     string
	Processed bool
	Draft     bool
	FilePath  string
}

// AlbumData represents fetched album metadata to be written to frontmatter
type AlbumData struct {
	Title      string
	Artist     string
	Year       string
	Label      string
	LabelURL   string
	DiscogsURL string
	DiscogsID  int
	CoverURL   string
	CoverPath  string
}

// parseMarkdownMusicFiles reads all markdown files in content/consumed/music/ and extracts album info
func parseMarkdownMusicFiles(contentDir string, includeDrafts bool) ([]AlbumInfo, error) {
	musicDir := filepath.Join(contentDir, "consumed", "music")
	
	// Check if directory exists
	if _, err := os.Stat(musicDir); err != nil {
		return nil, fmt.Errorf("music directory not found: %s", musicDir)
	}

	var albums []AlbumInfo
	
	files, err := os.ReadDir(musicDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(musicDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		contentStr := string(content)
		
		// Find frontmatter
		frontmatterStart := strings.Index(contentStr, "+++")
		if frontmatterStart == -1 {
			continue
		}
		
		frontmatterEnd := strings.Index(contentStr[frontmatterStart+3:], "+++")
		if frontmatterEnd == -1 {
			continue
		}
		
		frontmatter := contentStr[frontmatterStart+3 : frontmatterStart+3+frontmatterEnd]

		// Extract title
		titleMatch := regexp.MustCompile(`title\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(titleMatch) < 2 {
			continue
		}
		title := titleMatch[1]

		// Extract category - must be "music"
		if !regexp.MustCompile(`category\s*=\s*"music"`).MatchString(frontmatter) {
			continue
		}

		// Check draft flag
		draft := false
		if regexp.MustCompile(`draft\s*=\s*true`).MatchString(frontmatter) {
			draft = true
		}

		// Skip drafts unless explicitly including them
		if draft && !includeDrafts {
			continue
		}

		// Check processed flag
		processed := false
		if regexp.MustCompile(`processed\s*=\s*(true|yes)`).MatchString(frontmatter) {
			processed = true
		}

		// Extract fields
		artist := ""
		artistMatch := regexp.MustCompile(`artist\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(artistMatch) >= 2 {
			artist = artistMatch[1]
		}

		year := ""
		yearMatch := regexp.MustCompile(`year\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(yearMatch) >= 2 {
			year = yearMatch[1]
		}

		label := ""
		labelMatch := regexp.MustCompile(`label\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(labelMatch) >= 2 {
			label = labelMatch[1]
		}

		// Check for missing metadata
		hasArtist := artist != ""
		hasYear := year != ""
		hasLabel := label != ""
		hasDiscogs := regexp.MustCompile(`discogs\s*=\s*"[^"]+"`).MatchString(frontmatter)
		hasImg := regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(frontmatter)
		
		// Process if:
		// 1. processed is missing or false, OR
		// 2. metadata is incomplete (missing artist, year, label, discogs, or img)
		needsProcessing := !processed || !hasArtist || !hasYear || !hasLabel || !hasDiscogs || !hasImg
		
		if !needsProcessing {
			continue
		}

		albums = append(albums, AlbumInfo{
			Title:     title,
			Artist:    artist,
			Year:      year,
			Label:     label,
			Processed: processed,
			Draft:     draft,
			FilePath:  filePath,
		})
	}

	return albums, nil
}

// updateMarkdownMusicFrontmatter updates a markdown file's frontmatter with album data
func updateMarkdownMusicFrontmatter(filePath string, data AlbumData) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	
	// Find frontmatter boundaries
	frontmatterStart := strings.Index(contentStr, "+++")
	if frontmatterStart == -1 {
		return fmt.Errorf("no frontmatter found")
	}
	
	frontmatterEnd := strings.Index(contentStr[frontmatterStart+3:], "+++")
	if frontmatterEnd == -1 {
		return fmt.Errorf("no closing frontmatter found")
	}
	frontmatterEnd += frontmatterStart + 3

	frontmatter := contentStr[frontmatterStart+3 : frontmatterEnd]
	body := contentStr[frontmatterEnd+3:]

	updated := false

	// Update artist
	if data.Artist != "" {
		if regexp.MustCompile(`artist\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`artist\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`artist = "%s"`, data.Artist))
			updated = true
		} else {
			if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nartist = \"%s\"", data.Artist))
				updated = true
			}
		}
	}

	// Update year
	if data.Year != "" {
		if regexp.MustCompile(`year\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`year\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`year = "%s"`, data.Year))
			updated = true
		} else {
			if regexp.MustCompile(`artist\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(artist\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nyear = \"%s\"", data.Year))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nyear = \"%s\"", data.Year))
				updated = true
			}
		}
	}

	// Update label
	if data.Label != "" {
		if regexp.MustCompile(`label\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`label\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`label = "%s"`, data.Label))
			updated = true
		} else {
			if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nlabel = \"%s\"", data.Label))
				updated = true
			} else if regexp.MustCompile(`artist\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(artist\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nlabel = \"%s\"", data.Label))
				updated = true
			}
		}
	}

	// Update or add Discogs URL
	if data.DiscogsURL != "" {
		if regexp.MustCompile(`discogs\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`discogs\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`discogs = "%s"`, data.DiscogsURL))
			updated = true
		} else {
			if regexp.MustCompile(`label\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(label\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ndiscogs = \"%s\"", data.DiscogsURL))
				updated = true
			} else if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ndiscogs = \"%s\"", data.DiscogsURL))
				updated = true
			}
		}
	}

	// Update or add label URL
	if data.LabelURL != "" {
		if regexp.MustCompile(`discogsLabel\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`discogsLabel\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`discogsLabel = "%s"`, data.LabelURL))
			updated = true
		} else {
			if regexp.MustCompile(`discogs\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(discogs\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\ndiscogsLabel = \"%s\"", data.LabelURL))
				updated = true
			}
		}
	}

	// Update or add image path
	if data.CoverPath != "" {
		if regexp.MustCompile(`img\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`img\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`img = "%s"`, data.CoverPath))
			updated = true
		} else {
			if regexp.MustCompile(`category\s*=\s*"music"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(category\s*=\s*"music")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nimg = \"%s\"", data.CoverPath))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nimg = \"%s\"", data.CoverPath))
				updated = true
			}
		}
	}

	// Add or update processed flag (set to true after successful processing)
	if !regexp.MustCompile(`processed\s*=\s*(true|false)`).MatchString(frontmatter) {
		// Insert after img, discogs, label, or title
		if regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(img\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`discogs\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(discogs\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`label\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(label\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		}
	} else {
		// Update existing processed flag to true
		frontmatter = regexp.MustCompile(`processed\s*=\s*(true|false)`).ReplaceAllString(frontmatter, "processed = true")
		updated = true
	}

	if updated {
		newContent := "+++" + frontmatter + "+++" + body
		return os.WriteFile(filePath, []byte(newContent), 0644)
	}

	return nil
}

// BookInfo represents a book from markdown frontmatter
type BookInfo struct {
	Title     string
	Author    string
	Year      string
	Publisher string
	Processed bool
	Draft     bool
	FilePath  string
}

// BookData represents fetched book metadata to be written to frontmatter
type BookData struct {
	Title         string
	Author        string
	Year          string
	Publisher     string
	OpenLibraryURL string
	CoverURL      string
	CoverPath     string
}

// parseMarkdownBookFiles reads all markdown files in content/consumed/book/ and extracts book info
func parseMarkdownBookFiles(contentDir string, includeDrafts bool) ([]BookInfo, error) {
	bookDir := filepath.Join(contentDir, "consumed", "book")
	
	// Check if directory exists
	if _, err := os.Stat(bookDir); err != nil {
		return nil, fmt.Errorf("book directory not found: %s", bookDir)
	}

	var books []BookInfo
	
	files, err := os.ReadDir(bookDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(bookDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		contentStr := string(content)
		
		// Find frontmatter
		frontmatterStart := strings.Index(contentStr, "+++")
		if frontmatterStart == -1 {
			continue
		}
		
		frontmatterEnd := strings.Index(contentStr[frontmatterStart+3:], "+++")
		if frontmatterEnd == -1 {
			continue
		}
		
		frontmatter := contentStr[frontmatterStart+3 : frontmatterStart+3+frontmatterEnd]

		// Extract title
		titleMatch := regexp.MustCompile(`title\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(titleMatch) < 2 {
			continue
		}
		title := titleMatch[1]

		// Extract category - must be "book"
		if !regexp.MustCompile(`category\s*=\s*"book"`).MatchString(frontmatter) {
			continue
		}

		// Check draft flag
		draft := false
		if regexp.MustCompile(`draft\s*=\s*true`).MatchString(frontmatter) {
			draft = true
		}

		// Skip drafts unless explicitly including them
		if draft && !includeDrafts {
			continue
		}

		// Check processed flag
		processed := false
		if regexp.MustCompile(`processed\s*=\s*(true|yes)`).MatchString(frontmatter) {
			processed = true
		}

		// Extract fields
		author := ""
		authorMatch := regexp.MustCompile(`author\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(authorMatch) >= 2 {
			author = authorMatch[1]
		}

		year := ""
		yearMatch := regexp.MustCompile(`year\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(yearMatch) >= 2 {
			year = yearMatch[1]
		}

		publisher := ""
		publisherMatch := regexp.MustCompile(`publisher\s*=\s*"([^"]+)"`).FindStringSubmatch(frontmatter)
		if len(publisherMatch) >= 2 {
			publisher = publisherMatch[1]
		}

		// Check for missing metadata
		hasAuthor := author != ""
		hasYear := year != ""
		hasPublisher := publisher != ""
		hasOpenLibrary := regexp.MustCompile(`openlibrary\s*=\s*"[^"]+"`).MatchString(frontmatter)
		hasImg := regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(frontmatter)
		
		// Process if:
		// 1. processed is missing or false, OR
		// 2. metadata is incomplete (missing author, year, publisher, openlibrary, or img)
		needsProcessing := !processed || !hasAuthor || !hasYear || !hasPublisher || !hasOpenLibrary || !hasImg
		
		if !needsProcessing {
			continue
		}

		books = append(books, BookInfo{
			Title:     title,
			Author:    author,
			Year:      year,
			Publisher: publisher,
			Processed: processed,
			Draft:     draft,
			FilePath:  filePath,
		})
	}

	return books, nil
}

// updateMarkdownBookFrontmatter updates a markdown file's frontmatter with book data
func updateMarkdownBookFrontmatter(filePath string, data BookData) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	
	// Find frontmatter boundaries
	frontmatterStart := strings.Index(contentStr, "+++")
	if frontmatterStart == -1 {
		return fmt.Errorf("no frontmatter found")
	}
	
	frontmatterEnd := strings.Index(contentStr[frontmatterStart+3:], "+++")
	if frontmatterEnd == -1 {
		return fmt.Errorf("no closing frontmatter found")
	}
	frontmatterEnd += frontmatterStart + 3

	frontmatter := contentStr[frontmatterStart+3 : frontmatterEnd]
	body := contentStr[frontmatterEnd+3:]

	updated := false

	// Update author
	if data.Author != "" {
		if regexp.MustCompile(`author\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`author\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`author = "%s"`, data.Author))
			updated = true
		} else {
			if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nauthor = \"%s\"", data.Author))
				updated = true
			}
		}
	}

	// Update year
	if data.Year != "" {
		if regexp.MustCompile(`year\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`year\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`year = "%s"`, data.Year))
			updated = true
		} else {
			if regexp.MustCompile(`author\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(author\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nyear = \"%s\"", data.Year))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nyear = \"%s\"", data.Year))
				updated = true
			}
		}
	}

	// Update publisher
	if data.Publisher != "" {
		if regexp.MustCompile(`publisher\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`publisher\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`publisher = "%s"`, data.Publisher))
			updated = true
		} else {
			if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\npublisher = \"%s\"", data.Publisher))
				updated = true
			} else if regexp.MustCompile(`author\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(author\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\npublisher = \"%s\"", data.Publisher))
				updated = true
			}
		}
	}

	// Update or add OpenLibrary URL
	if data.OpenLibraryURL != "" {
		if regexp.MustCompile(`openlibrary\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`openlibrary\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`openlibrary = "%s"`, data.OpenLibraryURL))
			updated = true
		} else {
			if regexp.MustCompile(`publisher\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(publisher\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nopenlibrary = \"%s\"", data.OpenLibraryURL))
				updated = true
			} else if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nopenlibrary = \"%s\"", data.OpenLibraryURL))
				updated = true
			}
		}
	}

	// Update or add image path
	if data.CoverPath != "" {
		if regexp.MustCompile(`img\s*=\s*"[^"]*"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`img\s*=\s*"[^"]*"`).ReplaceAllString(frontmatter, fmt.Sprintf(`img = "%s"`, data.CoverPath))
			updated = true
		} else {
			if regexp.MustCompile(`category\s*=\s*"book"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(category\s*=\s*"book")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nimg = \"%s\"", data.CoverPath))
				updated = true
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
				frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, fmt.Sprintf("$1\nimg = \"%s\"", data.CoverPath))
				updated = true
			}
		}
	}

	// Add or update processed flag (set to true after successful processing)
	if !regexp.MustCompile(`processed\s*=\s*(true|false)`).MatchString(frontmatter) {
		// Insert after img, openlibrary, publisher, or title
		if regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(img\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`openlibrary\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(openlibrary\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`publisher\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(publisher\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(frontmatter) {
			frontmatter = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(frontmatter, "$1\nprocessed = true")
			updated = true
		}
	} else {
		// Update existing processed flag to true
		frontmatter = regexp.MustCompile(`processed\s*=\s*(true|false)`).ReplaceAllString(frontmatter, "processed = true")
		updated = true
	}

	if updated {
		newContent := "+++" + frontmatter + "+++" + body
		return os.WriteFile(filePath, []byte(newContent), 0644)
	}

	return nil
}

