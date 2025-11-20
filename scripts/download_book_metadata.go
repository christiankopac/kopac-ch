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

const openLibraryAPIBase = "https://openlibrary.org"
const googleBooksAPIBase = "https://www.googleapis.com/books/v1"

type BookSearchResult struct {
	Key      string   `json:"key"`
	Title    string   `json:"title"`
	Author   []string `json:"author_name"`
	Year     string   `json:"first_publish_year"`
	ISBN     []string `json:"isbn"`
	CoverKey string   `json:"cover_i"`
}

type BookSearchResponse struct {
	Docs []BookSearchResult `json:"docs"`
}

// Google Books API types
type GoogleBookItem struct {
	ID         string              `json:"id"`
	VolumeInfo GoogleVolumeInfo    `json:"volumeInfo"`
}

type GoogleVolumeInfo struct {
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	PublishedDate string   `json:"publishedDate"`
	Publisher     string   `json:"publisher"`
	IndustryIdentifiers []struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"industryIdentifiers"`
	ImageLinks struct {
		Thumbnail string `json:"thumbnail"`
		Small     string `json:"small"`
	} `json:"imageLinks"`
	InfoLink string `json:"infoLink"`
	PreviewLink string `json:"previewLink"`
}

type GoogleBooksResponse struct {
	Items []GoogleBookItem `json:"items"`
}

type BookDetails struct {
	Title    string   `json:"title"`
	Authors  []Author `json:"authors"`
	Publish  []string `json:"publish_dates"`
	ISBN10   []string `json:"isbn_10"`
	ISBN13   []string `json:"isbn_13"`
	Publishers []string `json:"publishers"`
}

type Author struct {
	Key string `json:"key"`
}

type AuthorDetails struct {
	Name string `json:"name"`
}

type BookData struct {
	Title    string
	Author   string
	Year     string
	Publisher string
	OpenLibraryURL string
	CoverURL string
	CoverPath string
}

func getBaseDir() string {
	wd, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(wd, "data", "books", "books.toml")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return wd
}

func searchBookGoogle(title string) (*GoogleBookItem, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	
	// Search for the book using Google Books API
	searchURL := fmt.Sprintf("%s/volumes?q=%s&maxResults=5", googleBooksAPIBase, strings.ReplaceAll(title, " ", "+"))
	
	// Retry up to 3 times for 503 errors
	maxRetries := 3
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry (exponential backoff)
			waitTime := time.Duration(attempt) * 2 * time.Second
			fmt.Printf("    Retrying in %v...\n", waitTime)
			time.Sleep(waitTime)
		}
		
		resp, err := client.Get(searchURL)
		if err != nil {
			lastErr = fmt.Errorf("failed to search Google Books: %w", err)
			continue
		}
		
		if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = fmt.Errorf("Google Books search failed with status: %d", resp.StatusCode)
			continue
		}
		
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("Google Books search failed with status: %d", resp.StatusCode)
		}
		
		var searchResp GoogleBooksResponse
		if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode Google Books response: %w", err)
		}
		resp.Body.Close()
		
		if len(searchResp.Items) == 0 {
			return nil, fmt.Errorf("no results found in Google Books")
		}
		
		// Return the first result
		return &searchResp.Items[0], nil
	}
	
	return nil, lastErr
}

func searchBook(title string) (*BookSearchResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Search for the book
	searchURL := fmt.Sprintf("%s/search.json?title=%s&limit=5", openLibraryAPIBase, strings.ReplaceAll(title, " ", "+"))
	
	resp, err := client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}
	
	var searchResp BookSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}
	
	if len(searchResp.Docs) == 0 {
		return nil, fmt.Errorf("no results found")
	}
	
	// Return the first result
	return &searchResp.Docs[0], nil
}

func getBookDetails(workKey string) (*BookDetails, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	// Get work details
	detailsURL := fmt.Sprintf("%s%s.json", openLibraryAPIBase, workKey)
	
	resp, err := client.Get(detailsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get details: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("details request failed with status: %d", resp.StatusCode)
	}
	
	var details BookDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("failed to decode details: %w", err)
	}
	
	return &details, nil
}

func getAuthorName(authorKey string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	authorURL := fmt.Sprintf("%s%s.json", openLibraryAPIBase, authorKey)
	
	resp, err := client.Get(authorURL)
	if err != nil {
		return "", fmt.Errorf("failed to get author: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("author request failed with status: %d", resp.StatusCode)
	}
	
	var author AuthorDetails
	if err := json.NewDecoder(resp.Body).Decode(&author); err != nil {
		return "", fmt.Errorf("failed to decode author: %w", err)
	}
	
	return author.Name, nil
}

func processBookGoogle(title string) (*BookData, error) {
	fmt.Printf("Searching Google Books for: %s\n", title)
	
	googleBook, err := searchBookGoogle(title)
	if err != nil {
		return nil, fmt.Errorf("Google Books search failed: %w", err)
	}
	
	volumeInfo := googleBook.VolumeInfo
	fmt.Printf("  Found: %s\n", volumeInfo.Title)
	
	// Get author name
	authorName := ""
	if len(volumeInfo.Authors) > 0 {
		authorName = volumeInfo.Authors[0]
	}
	
	// Extract year from publishedDate
	year := ""
	if volumeInfo.PublishedDate != "" {
		// PublishedDate can be "YYYY", "YYYY-MM", or "YYYY-MM-DD"
		yearMatch := regexp.MustCompile(`^\d{4}`).FindString(volumeInfo.PublishedDate)
		if yearMatch != "" {
			year = yearMatch
		}
	}
	
	// Get publisher
	publisher := volumeInfo.Publisher
	
	// Build Open Library URL (try to find ISBN and search Open Library)
	openLibraryURL := ""
	if len(volumeInfo.IndustryIdentifiers) > 0 {
		for _, id := range volumeInfo.IndustryIdentifiers {
			if id.Type == "ISBN_13" || id.Type == "ISBN_10" {
				// Try to find the book on Open Library using ISBN
				isbn := id.Identifier
				openLibraryURL = fmt.Sprintf("https://openlibrary.org/isbn/%s", isbn)
				break
			}
		}
	}
	
	// If no ISBN found, use Google Books link
	if openLibraryURL == "" {
		openLibraryURL = volumeInfo.InfoLink
		if openLibraryURL == "" {
			openLibraryURL = volumeInfo.PreviewLink
		}
	}
	
	// Get cover image URL (prefer thumbnail, fallback to small)
	coverURL := ""
	if volumeInfo.ImageLinks.Thumbnail != "" {
		coverURL = volumeInfo.ImageLinks.Thumbnail
		// Replace http:// with https:// and remove &edge=curl parameter if present
		coverURL = strings.ReplaceAll(coverURL, "http://", "https://")
		coverURL = strings.ReplaceAll(coverURL, "&edge=curl", "")
	} else if volumeInfo.ImageLinks.Small != "" {
		coverURL = volumeInfo.ImageLinks.Small
		coverURL = strings.ReplaceAll(coverURL, "http://", "https://")
		coverURL = strings.ReplaceAll(coverURL, "&edge=curl", "")
	}
	
	return &BookData{
		Title:         volumeInfo.Title,
		Author:        authorName,
		Year:          year,
		Publisher:     publisher,
		OpenLibraryURL: openLibraryURL,
		CoverURL:      coverURL,
	}, nil
}

func processBook(title string) (*BookData, error) {
	// Try Google Books first (more reliable)
	fmt.Printf("Searching for: %s\n", title)
	
	googleData, err := processBookGoogle(title)
	if err == nil {
		return googleData, nil
	}
	
	fmt.Printf("  Google Books failed: %v, trying Open Library...\n", err)
	
	// Fallback to Open Library
	searchResult, err := searchBook(title)
	if err != nil {
		return nil, fmt.Errorf("both Google Books and Open Library searches failed. Last error: %w", err)
	}
	
	fmt.Printf("  Found on Open Library: %s\n", searchResult.Title)
	
	// Get work details for more info
	details, err := getBookDetails(searchResult.Key)
	if err != nil {
		fmt.Printf("  Warning: Could not get details: %v\n", err)
		details = &BookDetails{}
	}
	
	// Get author name
	authorName := ""
	if len(searchResult.Author) > 0 {
		authorName = searchResult.Author[0]
	} else if len(details.Authors) > 0 {
		// Try to get author name from details
		authorName, err = getAuthorName(details.Authors[0].Key)
		if err != nil {
			fmt.Printf("  Warning: Could not get author name: %v\n", err)
		}
	}
	
	// Extract year
	year := ""
	if searchResult.Year != "" {
		year = searchResult.Year
	} else if len(details.Publish) > 0 {
		// Try to extract year from publish date
		yearMatch := regexp.MustCompile(`\d{4}`).FindString(details.Publish[0])
		if yearMatch != "" {
			year = yearMatch
		}
	}
	
	// Get publisher
	publisher := ""
	if len(details.Publishers) > 0 {
		publisher = details.Publishers[0]
	}
	
	// Build Open Library URL
	openLibraryURL := fmt.Sprintf("https://openlibrary.org%s", searchResult.Key)
	
	return &BookData{
		Title:         searchResult.Title,
		Author:        authorName,
		Year:          year,
		Publisher:     publisher,
		OpenLibraryURL: openLibraryURL,
	}, nil
}

func parseConsumedToml(filepath string) ([]map[string]string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	
	var books []map[string]string
	
	// Find all [[collection]] blocks - Go regexp doesn't support negative lookahead
	contentStr := string(content)
	re := regexp.MustCompile(`\[\[collection\]\]\s*\n`)
	indices := re.FindAllStringIndex(contentStr, -1)
	
	for i, idx := range indices {
		start := idx[1] // Start after [[collection]]\n
		var end int
		if i+1 < len(indices) {
			end = indices[i+1][0] // Next [[
		} else {
			end = len(contentStr) // End of file
		}
		
		block := contentStr[start:end]
		
		// Only process books
		if !strings.Contains(block, `category = "books"`) {
			continue
		}
		
		book := make(map[string]string)
		
		// Extract title
		titleMatch := regexp.MustCompile(`title\s*=\s*"([^"]+)"`).FindStringSubmatch(block)
		if len(titleMatch) > 1 {
			book["title"] = titleMatch[1]
		}
		
		// Extract existing fields
		authorMatch := regexp.MustCompile(`author\s*=\s*"([^"]*)"`).FindStringSubmatch(block)
		if len(authorMatch) > 1 {
			book["author"] = authorMatch[1]
		}
		
		yearMatch := regexp.MustCompile(`year\s*=\s*"([^"]*)"`).FindStringSubmatch(block)
		if len(yearMatch) > 1 {
			book["year"] = yearMatch[1]
		}
		
		publisherMatch := regexp.MustCompile(`publisher\s*=\s*"([^"]*)"`).FindStringSubmatch(block)
		if len(publisherMatch) > 1 {
			book["publisher"] = publisherMatch[1]
		}
		
		// Extract processed flag
		processedMatch := regexp.MustCompile(`processed\s*=\s*(true|yes)`).FindStringSubmatch(block)
		if len(processedMatch) > 1 {
			book["processed"] = "true"
		}
		
		if book["title"] != "" {
			books = append(books, book)
		}
	}
	
	return books, nil
}

func updateConsumedToml(bookUpdates map[string]*BookData) error {
	baseDir := getBaseDir()
	filepath := filepath.Join(baseDir, "data", "books", "books.toml")
	
	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	
	updated := false
	
	for originalTitle, data := range bookUpdates {
		// Find the collection block using index-based approach (no lookahead)
		contentStr := string(content)
		re := regexp.MustCompile(`\[\[collection\]\]\s*\n`)
		indices := re.FindAllStringIndex(contentStr, -1)
		
		var currentBlock string
		var blockStart, blockEnd int
		found := false
		
		for i, idx := range indices {
			start := idx[1] // Start after [[collection]]\n
			var end int
			if i+1 < len(indices) {
				end = indices[i+1][0] // Next [[
			} else {
				end = len(contentStr) // End of file
			}
			
			candidateBlock := contentStr[start:end]
			// Check if this block contains the title we're looking for
			titlePattern := regexp.MustCompile(fmt.Sprintf(`title\s*=\s*"%s"`, regexp.QuoteMeta(originalTitle)))
			if titlePattern.MatchString(candidateBlock) {
				currentBlock = candidateBlock
				blockStart = start
				blockEnd = end
				found = true
				break
			}
		}
		
		if !found {
			continue
		}
		
		modifiedBlock := currentBlock
		
		// Update author
		if data.Author != "" {
			if regexp.MustCompile(`author\s*=\s*"[^"]*"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(author\s*=\s*)"[^"]*"`).ReplaceAllString(modifiedBlock, fmt.Sprintf(`$1"%s"`, data.Author))
			} else {
				// Add after title, preserving newline
				if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nauthor = \"%s\"", data.Author))
				} else {
					modifiedBlock = strings.TrimRight(modifiedBlock, "\n") + fmt.Sprintf("\nauthor = \"%s\"", data.Author)
				}
			}
		}
		
		// Update year
		if data.Year != "" {
			if regexp.MustCompile(`year\s*=\s*"[^"]*"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(year\s*=\s*)"[^"]*"`).ReplaceAllString(modifiedBlock, fmt.Sprintf(`$1"%s"`, data.Year))
			} else {
				// Try to insert after author, then title
				if regexp.MustCompile(`author\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(author\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nyear = \"%s\"", data.Year))
				} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nyear = \"%s\"", data.Year))
				} else {
					modifiedBlock = strings.TrimRight(modifiedBlock, "\n") + fmt.Sprintf("\nyear = \"%s\"", data.Year)
				}
			}
		}
		
		// Update publisher
		if data.Publisher != "" {
			if regexp.MustCompile(`publisher\s*=\s*"[^"]*"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(publisher\s*=\s*)"[^"]*"`).ReplaceAllString(modifiedBlock, fmt.Sprintf(`$1"%s"`, data.Publisher))
			} else {
				// Try to insert after year, author, or title (in that order)
				if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\npublisher = \"%s\"", data.Publisher))
				} else if regexp.MustCompile(`author\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(author\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\npublisher = \"%s\"", data.Publisher))
				} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\npublisher = \"%s\"", data.Publisher))
				} else {
					modifiedBlock = strings.TrimRight(modifiedBlock, "\n") + fmt.Sprintf("\npublisher = \"%s\"", data.Publisher)
				}
			}
		}
		
		// Update or add Open Library URL
		if data.OpenLibraryURL != "" {
			if regexp.MustCompile(`openlibrary\s*=\s*"[^"]*"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(openlibrary\s*=\s*)"[^"]*"`).ReplaceAllString(modifiedBlock, fmt.Sprintf(`$1"%s"`, data.OpenLibraryURL))
			} else {
				// Try to insert after publisher, year, author, or title (in that order)
				if regexp.MustCompile(`publisher\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(publisher\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nopenlibrary = \"%s\"", data.OpenLibraryURL))
				} else if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nopenlibrary = \"%s\"", data.OpenLibraryURL))
				} else if regexp.MustCompile(`author\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(author\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nopenlibrary = \"%s\"", data.OpenLibraryURL))
				} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nopenlibrary = \"%s\"", data.OpenLibraryURL))
				} else {
					modifiedBlock = strings.TrimRight(modifiedBlock, "\n") + fmt.Sprintf("\nopenlibrary = \"%s\"", data.OpenLibraryURL)
				}
			}
		}
		
		// Update or add cover image path
		if data.CoverPath != "" {
			if regexp.MustCompile(`img\s*=\s*"[^"]*"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(img\s*=\s*)"[^"]*"`).ReplaceAllString(modifiedBlock, fmt.Sprintf(`$1"%s"`, data.CoverPath))
			} else {
				// Try to insert after title or at the beginning of the block
				if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
					modifiedBlock = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nimg = \"%s\"", data.CoverPath))
				} else {
					modifiedBlock = strings.TrimRight(modifiedBlock, "\n") + fmt.Sprintf("\nimg = \"%s\"", data.CoverPath)
				}
			}
		}
		
		// Add processed flag to mark book as processed
		if !regexp.MustCompile(`processed\s*=\s*(true|yes)`).MatchString(modifiedBlock) {
			// Insert after img, openlibrary, publisher, year, author, or title (in that order)
			if regexp.MustCompile(`img\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(img\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nprocessed = true"))
			} else if regexp.MustCompile(`openlibrary\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(openlibrary\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nprocessed = true"))
			} else if regexp.MustCompile(`publisher\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(publisher\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nprocessed = true"))
			} else if regexp.MustCompile(`year\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(year\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nprocessed = true"))
			} else if regexp.MustCompile(`author\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(author\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nprocessed = true"))
			} else if regexp.MustCompile(`title\s*=\s*"[^"]+"`).MatchString(modifiedBlock) {
				modifiedBlock = regexp.MustCompile(`(title\s*=\s*"[^"]+")`).ReplaceAllString(modifiedBlock, fmt.Sprintf("$1\nprocessed = true"))
			} else {
				modifiedBlock = strings.TrimRight(modifiedBlock, "\n") + fmt.Sprintf("\nprocessed = true")
			}
		}
		
		if modifiedBlock != currentBlock {
			// Replace the block in the original content
			contentStr = contentStr[:blockStart] + modifiedBlock + contentStr[blockEnd:]
			content = []byte(contentStr)
			updated = true
			fmt.Printf("  ✓ Updated %s in books.toml\n", originalTitle)
			// Recalculate indices after modification
			indices = re.FindAllStringIndex(contentStr, -1)
		}
	}
	
	if updated {
		return os.WriteFile(filepath, content, 0644)
	}
	
	return nil
}

func main() {
	updateToml := flag.Bool("update-toml", false, "Update consumed.toml with fetched metadata")
	skipExisting := flag.Bool("skip-existing", false, "Skip books that already have author and year")
	flag.Parse()
	
	// Load .env file if it exists (before checking environment)
	baseDir := getBaseDir()
	envFile := filepath.Join(baseDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err == nil {
			// .env loaded successfully
		}
	}
	
	booksFile := filepath.Join(baseDir, "data", "books", "books.toml")
	
	books, err := parseConsumedToml(booksFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading consumed.toml: %v\n", err)
		os.Exit(1)
	}
	
	if len(books) == 0 {
		fmt.Println("No books found in books.toml")
		return
	}
	
	// Filter by command-line arguments if provided
	args := flag.Args()
	if len(args) > 0 {
		var filtered []map[string]string
		for _, book := range books {
			for _, arg := range args {
				if strings.Contains(strings.ToLower(book["title"]), strings.ToLower(arg)) {
					filtered = append(filtered, book)
					break
				}
			}
		}
		books = filtered
	}
	
	fmt.Printf("Processing %d book(s)...\n\n", len(books))
	
	bookUpdates := make(map[string]*BookData)
	imagesDir := filepath.Join(baseDir, "static", "images", "books")
	os.MkdirAll(imagesDir, 0755)
	
	for _, book := range books {
		title := book["title"]
		
		// Skip if marked as processed
		if book["processed"] == "true" {
			fmt.Printf("Skipping %s (already processed)\n", title)
			continue
		}
		
		// Skip if already has all metadata (when using --skip-existing flag)
		if *skipExisting && book["author"] != "" && book["year"] != "" {
			fmt.Printf("Skipping %s (already has metadata)\n", title)
			continue
		}
		
		data, err := processBook(title)
		if err != nil {
			fmt.Printf("  ✗ Error: %v\n\n", err)
			continue
		}
		
		fmt.Printf("  Author: %s\n", data.Author)
		fmt.Printf("  Year: %s\n", data.Year)
		if data.Publisher != "" {
			fmt.Printf("  Publisher: %s\n", data.Publisher)
		}
		fmt.Printf("  Open Library: %s\n", data.OpenLibraryURL)
		
		// Download cover if available
		if data.CoverURL != "" {
			coverFilename := getCoverFilename(title)
			coverPath := filepath.Join(imagesDir, coverFilename)
			if downloadCover(data.CoverURL, coverPath) {
				data.CoverPath = fmt.Sprintf("/images/books/%s", coverFilename)
				fmt.Printf("  Cover: %s\n", data.CoverPath)
			} else {
				fmt.Printf("  Warning: Could not download cover\n")
			}
		}
		fmt.Println()
		
		bookUpdates[title] = data
		
		// Rate limiting
		time.Sleep(1 * time.Second)
	}
	
	if *updateToml && len(bookUpdates) > 0 {
		fmt.Println("Updating books.toml...")
		if err := updateConsumedToml(bookUpdates); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating books.toml: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ books.toml updated")
	}
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

func getCoverFilename(title string) string {
	return slugify(title) + "_cover.jpg"
}

func downloadCover(coverURL, outputPath string) bool {
	if coverURL == "" {
		return false
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(coverURL)
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

