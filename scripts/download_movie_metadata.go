package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	tmdbAPIBase   = "https://api.themoviedb.org/3"
	tmdbImageBase = "https://image.tmdb.org/t/p/w500"
)

type MovieResult struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	PosterPath  string `json:"poster_path"`
}

type SearchResponse struct {
	Results []MovieResult `json:"results"`
}

type Credits struct {
	Crew []CrewMember `json:"crew"`
}

type CrewMember struct {
	Job  string `json:"job"`
	Name string `json:"name"`
}

type MovieDetails struct {
	ID      int     `json:"id"`
	Title   string  `json:"title"`
	Credits Credits `json:"credits"`
}

// MovieData is now defined in markdown_helpers.go

func main() {
	var updatePages, skipExisting, includeDrafts bool
	flag.BoolVar(&updatePages, "update-pages", true, "Update markdown pages with fetched metadata (default: true)")
	flag.BoolVar(&skipExisting, "skip-existing", false, "Skip movies that already have posters and directors")
	flag.BoolVar(&includeDrafts, "include-drafts", false, "Include draft movies when processing")
	flag.Parse()

	apiKey := getAPIKey()
	if apiKey == "" {
		os.Exit(1)
	}

	baseDir := getBaseDir()
	imagesDir := filepath.Join(baseDir, "static", "images", "movies")
	contentDir := filepath.Join(baseDir, "content")

	os.MkdirAll(imagesDir, 0755)

	// Get movies to process
	var movies []MovieInfo
	if flag.NArg() > 0 {
		// If titles provided as args, find corresponding markdown files
		for _, title := range flag.Args() {
			// Try to find the file by title
			slug := slugify(title)
			filePath := filepath.Join(contentDir, "consumed", "movie", slug+".md")
			if _, err := os.Stat(filePath); err == nil {
				movies = append(movies, MovieInfo{Title: title, FilePath: filePath})
			} else {
				fmt.Printf("Warning: Could not find file for movie: %s\n", title)
			}
		}
	} else {
		var err error
		movies, err = parseMarkdownFiles(contentDir, includeDrafts)
		if err != nil {
			fmt.Printf("Error reading markdown files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(movies) == 0 {
		fmt.Println("No movies found")
		os.Exit(0)
	}

	fmt.Printf("Found %d movies to process\n\n", len(movies))

	results := make(map[string]MovieData)

	for _, movie := range movies {
		// Skip if marked as processed and has all metadata (safety check)
		// This should already be filtered in parseMarkdownFiles, but check again for safety
		if skipExisting && movie.Director != "" {
			posterFile := getPosterFilename(movie.Title)
			posterPath := filepath.Join(imagesDir, posterFile)
			if _, err := os.Stat(posterPath); err == nil {
				fmt.Printf("\nSkipping %s (already has poster and director)\n", movie.Title)
				continue
			}
		}
		
		// Show what's missing
		missing := []string{}
		if movie.Year == "" {
			missing = append(missing, "year")
		}
		if movie.Director == "" {
			missing = append(missing, "director")
		}
		if len(missing) > 0 {
			fmt.Printf("\nProcessing: %s (missing: %s)\n", movie.Title, strings.Join(missing, ", "))
		} else {
			fmt.Printf("\nProcessing: %s\n", movie.Title)
		}

		result := processMovie(apiKey, movie.Title, movie.Year, movie.Director)
		if result != nil {
			results[movie.Title] = *result

			// Download poster
			posterDownloaded := false
			if result.PosterPath != "" {
				posterFile := getPosterFilename(movie.Title)
				posterPath := filepath.Join(imagesDir, posterFile)
				if downloadPoster(result.PosterURL, posterPath) {
					result.ImagePath = fmt.Sprintf("/images/movies/%s", posterFile)
					posterDownloaded = true
					fmt.Printf("  ✓ Downloaded poster: %s\n", posterFile)
				} else {
					fmt.Printf("  ⚠ Failed to download poster\n")
				}
			} else {
				fmt.Printf("  ⚠ No poster available\n")
			}

			// Mark as draft if no poster was downloaded
			if !posterDownloaded {
				result.Draft = true
			}

			// Update markdown file
			if updatePages && movie.FilePath != "" {
				if err := updateMarkdownFrontmatter(movie.FilePath, *result); err != nil {
					fmt.Printf("  ✗ Error updating markdown file: %v\n", err)
				} else {
					fmt.Printf("  ✓ Updated markdown file: %s\n", filepath.Base(movie.FilePath))
				}
			}
		} else {
			// Movie not found - mark as draft
			fmt.Printf("  ⚠ Movie not found, marking as draft\n")
			if updatePages && movie.FilePath != "" {
				draftData := MovieData{
					Title: movie.Title,
					Draft: true,
				}
				if err := updateMarkdownFrontmatter(movie.FilePath, draftData); err != nil {
					fmt.Printf("  ✗ Error updating markdown file: %v\n", err)
				} else {
					fmt.Printf("  ✓ Updated markdown file (marked as draft)\n")
				}
			}
		}
	}

	// Summary
	fmt.Printf("\n%s\n", strings.Repeat("=", 50))
	fmt.Printf("Summary:\n")
	fmt.Printf("  Processed: %d movies\n", len(results))
	posterCount := 0
	directorCount := 0
	for _, r := range results {
		if r.PosterPath != "" {
			posterCount++
		}
		if r.Director != "" {
			directorCount++
		}
	}
	fmt.Printf("  Posters downloaded: %d\n", posterCount)
	fmt.Printf("  Directors found: %d\n", directorCount)
}

func getAPIKey() string {
	// Try .env file first (before checking environment)
	baseDir := getBaseDir()
	envFile := filepath.Join(baseDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err == nil {
			// .env loaded successfully
		}
	}

	// Now check environment (either from .env or system)
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	fmt.Println("Error: TMDB_API_KEY not found. Set it as environment variable or in .env file")
	fmt.Printf("  Looking for .env at: %s\n", envFile)
	fmt.Println("Get your API key from: https://www.themoviedb.org/settings/api")
	return ""
}

func getBaseDir() string {
	// Start from current working directory and walk up to find project root
	wd, _ := os.Getwd()
	startWd := wd
	for {
		// Check for multiple markers to identify project root
		if _, err := os.Stat(filepath.Join(wd, "data", "movies", "movies.toml")); err == nil {
			return wd
		}
		if _, err := os.Stat(filepath.Join(wd, "data", "music", "music.toml")); err == nil {
			return wd
		}
		if _, err := os.Stat(filepath.Join(wd, ".env")); err == nil {
			return wd
		}
		if _, err := os.Stat(filepath.Join(wd, "content")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	// Fallback: return the directory where the script was run from
	return startWd
}

func searchMovie(apiKey, title, year string) (*MovieResult, error) {
	u := fmt.Sprintf("%s/search/movie", tmdbAPIBase)
	req, _ := http.NewRequest("GET", u, nil)
	q := req.URL.Query()
	q.Set("api_key", apiKey)
	q.Set("query", title)
	q.Set("language", "en-US")
	if year != "" {
		q.Set("year", year)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	if len(searchResp.Results) > 0 {
		return &searchResp.Results[0], nil
	}
	return nil, nil
}

func getMovieDetails(apiKey string, movieID int) (*MovieDetails, error) {
	u := fmt.Sprintf("%s/movie/%d", tmdbAPIBase, movieID)
	req, _ := http.NewRequest("GET", u, nil)
	q := req.URL.Query()
	q.Set("api_key", apiKey)
	q.Set("language", "en-US")
	q.Set("append_to_response", "credits")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var details MovieDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

func getDirector(credits Credits) string {
	for _, member := range credits.Crew {
		if member.Job == "Director" {
			return member.Name
		}
	}
	return ""
}

func downloadPoster(posterURL, outputPath string) bool {
	if posterURL == "" {
		return false
	}

	resp, err := http.Get(posterURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err == nil
}

func slugify(title string) string {
	// Convert to lowercase and replace spaces/special chars with underscores
	title = strings.ToLower(title)
	re := regexp.MustCompile(`[^\w\s-]`)
	title = re.ReplaceAllString(title, "")
	re = regexp.MustCompile(`[-\s]+`)
	title = re.ReplaceAllString(title, "_")
	return strings.Trim(title, "_")
}

func getPosterFilename(title string) string {
	return slugify(title) + "_poster.jpg"
}

// MovieInfo is now defined in markdown_helpers.go
// Old parseConsumedToml function removed - use parseMarkdownFiles instead

func processMovie(apiKey, title, year, existingDirector string) *MovieData {
	fmt.Printf("\nProcessing: %s", title)
	if year != "" {
		fmt.Printf(" (%s)", year)
	}
	fmt.Println()

	movie, err := searchMovie(apiKey, title, year)
	if err != nil || movie == nil {
		fmt.Printf("  ✗ Movie not found on TMDB\n")
		return nil
	}

	fmt.Printf("  ✓ Found: %s", movie.Title)
	if movie.ReleaseDate != "" && len(movie.ReleaseDate) >= 4 {
		releaseYear := movie.ReleaseDate[:4]
		fmt.Printf(" (%s)", releaseYear)
	}
	fmt.Println()

	details, err := getMovieDetails(apiKey, movie.ID)
	if err != nil {
		fmt.Printf("  ✗ Could not fetch movie details\n")
		return nil
	}

	director := existingDirector
	if director == "" {
		director = getDirector(details.Credits)
		if director != "" {
			fmt.Printf("  ✓ Director: %s\n", director)
		}
	}

	posterURL := ""
	if movie.PosterPath != "" {
		posterURL = tmdbImageBase + movie.PosterPath
	}

	releaseYear := year
	if releaseYear == "" && movie.ReleaseDate != "" && len(movie.ReleaseDate) >= 4 {
		releaseYear = movie.ReleaseDate[:4]
	}

	tmdbURL := fmt.Sprintf("https://www.themoviedb.org/movie/%d", movie.ID)

	return &MovieData{
		Title:      movie.Title,
		Year:       releaseYear,
		Director:   director,
		PosterPath: movie.PosterPath,
		PosterURL:  posterURL,
		TMDBID:     movie.ID,
		TMDBURL:    tmdbURL,
	}
}

// Old updateConsumedToml function removed - use updateMarkdownFrontmatter instead

