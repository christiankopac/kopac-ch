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

const discogsAPIBase = "https://api.discogs.com"

type ReleaseResult struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URI   string `json:"uri"`
}

type SearchResponse struct {
	Results []ReleaseResult `json:"results"`
}

type Artist struct {
	Name string `json:"name"`
}

type Label struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ResourceURL string `json:"resource_url"`
}

type Image struct {
	Type     string `json:"type"`
	URI      string `json:"uri"`
	ResourceURL string `json:"resource_url"`
	URI150    string `json:"uri150"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type ReleaseDetails struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Year    int      `json:"year"`
	Artists []Artist `json:"artists"`
	Labels  []Label  `json:"labels"`
	URI     string   `json:"uri"`
	Images  []Image  `json:"images"`
}

// AlbumData and AlbumInfo are now defined in markdown_helpers.go

func main() {
	var updatePages, skipExisting, includeDrafts bool
	flag.BoolVar(&updatePages, "update-pages", true, "Update markdown pages with fetched metadata (default: true)")
	flag.BoolVar(&skipExisting, "skip-existing", false, "Skip albums that already have all metadata")
	flag.BoolVar(&includeDrafts, "include-drafts", false, "Include draft albums when processing")
	flag.Parse()

	token := getUserToken()
	if token == "" {
		os.Exit(1)
	}

	baseDir := getBaseDir()
	imagesDir := filepath.Join(baseDir, "static", "images", "music")
	contentDir := filepath.Join(baseDir, "content")
	
	os.MkdirAll(imagesDir, 0755)

	// Get albums to process
	var albums []AlbumInfo
	if flag.NArg() > 0 {
		// If titles provided as args, find corresponding markdown files
		for _, title := range flag.Args() {
			// Try to find the file by title
			slug := slugify(title)
			filePath := filepath.Join(contentDir, "consumed", "music", slug+".md")
			if _, err := os.Stat(filePath); err == nil {
				albums = append(albums, AlbumInfo{Title: title, FilePath: filePath})
			} else {
				fmt.Printf("Warning: Could not find file for album: %s\n", title)
			}
		}
	} else {
		var err error
		albums, err = parseMarkdownMusicFiles(contentDir, includeDrafts)
		if err != nil {
			fmt.Printf("Error reading markdown files: %v\n", err)
			os.Exit(1)
		}
	}

	if len(albums) == 0 {
		fmt.Println("No albums found")
		os.Exit(0)
	}

	fmt.Printf("Found %d albums to process\n\n", len(albums))

	results := make(map[string]AlbumData)

	for _, album := range albums {
		if skipExisting && album.Artist != "" && album.Year != "" && album.Label != "" {
			fmt.Printf("\nSkipping %s (already has all metadata)\n", album.Title)
			continue
		}

		result := processAlbum(token, album.Title, album.Artist, album.Year)
		if result != nil {
			results[album.Title] = *result

			// Download cover image
			if result.CoverURL != "" {
				coverFile := getCoverFilename(album.Title)
				coverPath := filepath.Join(imagesDir, coverFile)
				if downloadCover(result.CoverURL, coverPath) {
					result.CoverPath = fmt.Sprintf("/images/music/%s", coverFile)
					fmt.Printf("  ✓ Downloaded cover: %s\n", coverFile)
				}
			}

			// Update markdown file
			if updatePages && album.FilePath != "" {
				if err := updateMarkdownMusicFrontmatter(album.FilePath, *result); err != nil {
					fmt.Printf("  ✗ Error updating markdown file: %v\n", err)
				} else {
					fmt.Printf("  ✓ Updated markdown file: %s\n", filepath.Base(album.FilePath))
				}
			}
		} else {
			fmt.Printf("  ⚠ Album not found\n")
		}

		// Rate limiting
		time.Sleep(1 * time.Second)
	}

	// Summary
	fmt.Printf("\n%s\n", strings.Repeat("=", 50))
	fmt.Printf("Summary:\n")
	fmt.Printf("  Processed: %d albums\n", len(results))
	artistCount := 0
	yearCount := 0
	labelCount := 0
	for _, r := range results {
		if r.Artist != "" {
			artistCount++
		}
		if r.Year != "" {
			yearCount++
		}
		if r.Label != "" {
			labelCount++
		}
	}
	fmt.Printf("  Artists found: %d\n", artistCount)
	fmt.Printf("  Years found: %d\n", yearCount)
	fmt.Printf("  Labels found: %d\n", labelCount)
}

func getUserToken() string {
	// Try .env file first (before checking environment)
	baseDir := getBaseDir()
	envFile := filepath.Join(baseDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err == nil {
			// .env loaded successfully
		}
	}

	// Now check environment (either from .env or system)
	token := os.Getenv("DISCOGS_USER_TOKEN")
	if token != "" {
		return token
	}

	fmt.Println("Error: DISCOGS_USER_TOKEN not found. Set it as environment variable or in .env file")
	fmt.Printf("  Looking for .env at: %s\n", envFile)
	fmt.Println("Get your token from: https://www.discogs.com/settings/developers")
	fmt.Println("Note: You need to create a personal access token")
	return ""
}

func getBaseDir() string {
	// Start from current working directory and walk up to find project root
	wd, _ := os.Getwd()
	startWd := wd
	for {
		// Check for multiple markers to identify project root
		if _, err := os.Stat(filepath.Join(wd, "data", "music", "music.toml")); err == nil {
			return wd
		}
		if _, err := os.Stat(filepath.Join(wd, "data", "movies", "movies.toml")); err == nil {
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

func searchAlbum(token, title, artist string) (*ReleaseResult, error) {
	u := fmt.Sprintf("%s/database/search", discogsAPIBase)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "HugoSite/1.0")
	req.Header.Set("Authorization", fmt.Sprintf("Discogs token=%s", token))

	q := req.URL.Query()
	query := title
	if artist != "" {
		query = artist + " " + title
	}
	q.Set("q", query)
	q.Set("type", "release")
	// Don't filter by format - some releases aren't tagged as "album"
	// q.Set("format", "album")
	q.Set("per_page", "25") // Increase results to find better matches
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
		// Try to find best match
		titleLower := strings.ToLower(strings.TrimSpace(title))
		artistLower := strings.ToLower(strings.TrimSpace(artist))
		
		// First, try exact title match
		for _, result := range searchResp.Results {
			resultTitleLower := strings.ToLower(result.Title)
			if strings.Contains(resultTitleLower, titleLower) {
				// If we have an artist, prefer results that also match artist
				if artistLower != "" {
					// Check if any artist in the result matches
					// Note: ReleaseResult doesn't have artists, so we'll check in details later
					return &result, nil
				}
				return &result, nil
			}
		}
		
		// If no title match, try partial match (for cases like "768" matching "768 / Can't Sleep")
		for _, result := range searchResp.Results {
			resultTitleLower := strings.ToLower(result.Title)
			// Check if title is contained in result or vice versa
			if strings.Contains(resultTitleLower, titleLower) || strings.Contains(titleLower, resultTitleLower) {
				return &result, nil
			}
		}
		
		// Return first result if no match found
		return &searchResp.Results[0], nil
	}
	return nil, nil
}

func getReleaseDetails(token string, releaseID int) (*ReleaseDetails, error) {
	u := fmt.Sprintf("%s/releases/%d", discogsAPIBase, releaseID)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "HugoSite/1.0")
	req.Header.Set("Authorization", fmt.Sprintf("Discogs token=%s", token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var details ReleaseDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return &details, nil
}

func processAlbum(token, title, existingArtist, existingYear string) *AlbumData {
	fmt.Printf("\nProcessing: %s\n", title)

	album, err := searchAlbum(token, title, existingArtist)
	if err != nil || album == nil {
		fmt.Printf("  ✗ Album not found on Discogs\n")
		return nil
	}

	fmt.Printf("  ✓ Found: %s\n", album.Title)

	details, err := getReleaseDetails(token, album.ID)
	if err != nil {
		fmt.Printf("  ✗ Could not fetch release details\n")
		return nil
	}

	artist := existingArtist
	if artist == "" && len(details.Artists) > 0 {
		artist = details.Artists[0].Name
		if artist != "" {
			fmt.Printf("  ✓ Artist: %s\n", artist)
		}
	}

	year := existingYear
	if year == "" && details.Year > 0 {
		year = fmt.Sprintf("%d", details.Year)
		if year != "" {
			fmt.Printf("  ✓ Year: %s\n", year)
		}
	}

	label := ""
	labelURL := ""
	if len(details.Labels) > 0 {
		label = details.Labels[0].Name
		if label != "" {
			fmt.Printf("  ✓ Label: %s\n", label)
		}
		// Construct label URL from label ID
		if details.Labels[0].ID > 0 {
			labelURL = fmt.Sprintf("https://www.discogs.com/label/%d", details.Labels[0].ID)
			fmt.Printf("  ✓ Label URL: %s\n", labelURL)
		} else if details.Labels[0].ResourceURL != "" {
			// Fallback: try to extract ID from resource_url
			// Resource URL format: https://api.discogs.com/labels/{id}
			if strings.Contains(details.Labels[0].ResourceURL, "/labels/") {
				parts := strings.Split(details.Labels[0].ResourceURL, "/labels/")
				if len(parts) > 1 {
					labelID := strings.TrimSuffix(parts[1], "/")
					labelURL = fmt.Sprintf("https://www.discogs.com/label/%s", labelID)
					fmt.Printf("  ✓ Label URL: %s\n", labelURL)
				}
			}
		}
	}

	discogsURL := details.URI
	if discogsURL == "" {
		discogsURL = fmt.Sprintf("https://www.discogs.com/release/%d", details.ID)
	}

	// Get cover image URL (prefer primary image, fallback to first available)
	coverURL := ""
	if len(details.Images) > 0 {
		// Look for primary image first
		for _, img := range details.Images {
			if img.Type == "primary" {
				coverURL = img.ResourceURL
				if coverURL == "" {
					coverURL = img.URI
				}
				break
			}
		}
		// If no primary, use first image
		if coverURL == "" {
			img := details.Images[0]
			coverURL = img.ResourceURL
			if coverURL == "" {
				coverURL = img.URI
			}
		}
	}

	return &AlbumData{
		Title:      album.Title,
		Artist:     artist,
		Year:       year,
		Label:      label,
		LabelURL:   labelURL,
		DiscogsURL: discogsURL,
		DiscogsID:  details.ID,
		CoverURL:   coverURL,
	}
}

// Old parseConsumedToml function removed - use parseMarkdownMusicFiles instead

func slugify(title string) string {
	// Convert to lowercase and replace spaces/special chars with underscores
	title = strings.ToLower(title)
	re := regexp.MustCompile(`[^\w\s-]`)
	title = re.ReplaceAllString(title, "")
	re = regexp.MustCompile(`[-\s]+`)
	title = re.ReplaceAllString(title, "_")
	return strings.Trim(title, "_")
}

func getCoverFilename(title string) string {
	return slugify(title) + "_cover.jpg"
}

func downloadCover(coverURL, outputPath string) bool {
	if coverURL == "" {
		return false
	}

	resp, err := http.Get(coverURL)
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

// Old updateConsumedToml function removed - use updateMarkdownMusicFrontmatter instead

