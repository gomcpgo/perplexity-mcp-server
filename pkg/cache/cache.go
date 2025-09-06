package cache

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// QueryMetadata represents metadata for a cached query
type QueryMetadata struct {
	Query      string                 `yaml:"query"`
	SearchType string                 `yaml:"search_type"`
	Timestamp  time.Time              `yaml:"timestamp"`
	Model      string                 `yaml:"model"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

// QueryListItem represents an item in the previous queries list
type QueryListItem struct {
	Query      string    `json:"query"`
	UniqueID   string    `json:"unique_id"`
	DateTime   time.Time `json:"datetime"`
	SearchType string    `json:"search_type"`
}

const (
	metadataFile = "metadata.yaml"
	resultFile   = "result.md"
	idLength     = 10
	idCharset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateUniqueID generates a 10-character alphanumeric unique ID
func GenerateUniqueID(rootFolder string) (string, error) {
	maxAttempts := 100
	for attempt := 0; attempt < maxAttempts; attempt++ {
		id := generateRandomID()
		if !idExists(rootFolder, id) {
			return id, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique ID after %d attempts", maxAttempts)
}

// generateRandomID creates a random 10-character alphanumeric string
func generateRandomID() string {
	result := make([]byte, idLength)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(idCharset))))
		result[i] = idCharset[num.Int64()]
	}
	return string(result)
}

// idExists checks if a folder with the given ID already exists
func idExists(rootFolder, id string) bool {
	if rootFolder == "" {
		return false
	}
	folderPath := filepath.Join(rootFolder, id)
	_, err := os.Stat(folderPath)
	return err == nil
}

// SaveResult saves query result and metadata to the cache
func SaveResult(rootFolder, query, searchType, model, result string, parameters map[string]interface{}) (string, error) {
	if rootFolder == "" {
		return "", nil // No caching if root folder not set
	}

	// Generate unique ID
	uniqueID, err := GenerateUniqueID(rootFolder)
	if err != nil {
		return "", fmt.Errorf("failed to generate unique ID: %w", err)
	}

	// Create folder for this result
	resultFolder := filepath.Join(rootFolder, uniqueID)
	if err := os.MkdirAll(resultFolder, 0755); err != nil {
		return "", fmt.Errorf("failed to create result folder: %w", err)
	}

	// Save metadata
	metadata := QueryMetadata{
		Query:      query,
		SearchType: searchType,
		Timestamp:  time.Now(),
		Model:      model,
		Parameters: parameters,
	}

	metadataPath := filepath.Join(resultFolder, metadataFile)
	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := ioutil.WriteFile(metadataPath, metadataBytes, 0644); err != nil {
		return "", fmt.Errorf("failed to write metadata file: %w", err)
	}

	// Save result
	resultPath := filepath.Join(resultFolder, resultFile)
	if err := ioutil.WriteFile(resultPath, []byte(result), 0644); err != nil {
		return "", fmt.Errorf("failed to write result file: %w", err)
	}

	return uniqueID, nil
}

// ListPreviousQueries returns a list of previous queries sorted by recency
func ListPreviousQueries(rootFolder string) ([]QueryListItem, error) {
	if rootFolder == "" {
		return []QueryListItem{}, nil // Return empty list if no root folder set
	}

	// Check if root folder exists
	if _, err := os.Stat(rootFolder); os.IsNotExist(err) {
		return []QueryListItem{}, nil // Return empty list if folder doesn't exist
	}

	// Read all subdirectories
	entries, err := ioutil.ReadDir(rootFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read results directory: %w", err)
	}

	var queryItems []QueryListItem

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		uniqueID := entry.Name()
		metadataPath := filepath.Join(rootFolder, uniqueID, metadataFile)

		// Read metadata
		metadataBytes, err := ioutil.ReadFile(metadataPath)
		if err != nil {
			continue // Skip if metadata file doesn't exist or can't be read
		}

		var metadata QueryMetadata
		if err := yaml.Unmarshal(metadataBytes, &metadata); err != nil {
			continue // Skip if metadata can't be parsed
		}

		queryItems = append(queryItems, QueryListItem{
			Query:      metadata.Query,
			UniqueID:   uniqueID,
			DateTime:   metadata.Timestamp,
			SearchType: metadata.SearchType,
		})
	}

	// Sort by timestamp (most recent first)
	sort.Slice(queryItems, func(i, j int) bool {
		return queryItems[i].DateTime.After(queryItems[j].DateTime)
	})

	return queryItems, nil
}

// GetPreviousResult retrieves a cached result by unique ID
func GetPreviousResult(rootFolder, uniqueID string) (string, error) {
	if rootFolder == "" {
		return "", fmt.Errorf("results root folder not configured")
	}

	// Validate uniqueID format (10 alphanumeric characters)
	if len(uniqueID) != idLength || !isValidID(uniqueID) {
		return "", fmt.Errorf("invalid unique ID format: must be %d alphanumeric characters", idLength)
	}

	resultPath := filepath.Join(rootFolder, uniqueID, resultFile)

	// Check if result file exists
	if _, err := os.Stat(resultPath); os.IsNotExist(err) {
		return "", fmt.Errorf("result with ID '%s' not found", uniqueID)
	}

	// Read result file
	resultBytes, err := ioutil.ReadFile(resultPath)
	if err != nil {
		return "", fmt.Errorf("failed to read result file: %w", err)
	}

	return string(resultBytes), nil
}

// isValidID checks if the ID contains only valid characters
func isValidID(id string) bool {
	for _, char := range id {
		if !strings.ContainsRune(idCharset, char) {
			return false
		}
	}
	return true
}

// IsCachingEnabled returns true if caching is enabled (root folder is set)
func IsCachingEnabled(rootFolder string) bool {
	return rootFolder != ""
}