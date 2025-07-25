package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/prasanthmj/perplexity/pkg/types"
)

// Config holds the configuration for the Perplexity MCP server
type Config struct {
	APIKey         string
	DefaultModel   string
	MaxTokens      int
	Temperature    float64
	TopP           float64
	TopK           int
	Timeout        time.Duration
	ReturnCitations bool
	ReturnImages    bool
	ReturnRelated   bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		// Set defaults
		DefaultModel:    types.DefaultModel,
		MaxTokens:       types.DefaultMaxTokens,
		Temperature:     types.DefaultTemperature,
		TopP:           types.DefaultTopP,
		TopK:           types.DefaultTopK,
		Timeout:        30 * time.Second,
		ReturnCitations: types.DefaultReturnCitations,
		ReturnImages:    types.DefaultReturnImages,
		ReturnRelated:   types.DefaultReturnRelated,
	}

	// API Key is required
	cfg.APIKey = os.Getenv("PERPLEXITY_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("PERPLEXITY_API_KEY environment variable is required")
	}

	// Override defaults with environment variables if set
	if model := os.Getenv("PERPLEXITY_DEFAULT_MODEL"); model != "" {
		if err := validateModel(model); err != nil {
			return nil, fmt.Errorf("invalid model: %w", err)
		}
		cfg.DefaultModel = model
	}

	if maxTokens := os.Getenv("PERPLEXITY_MAX_TOKENS"); maxTokens != "" {
		val, err := strconv.Atoi(maxTokens)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_MAX_TOKENS: %w", err)
		}
		if val <= 0 {
			return nil, fmt.Errorf("PERPLEXITY_MAX_TOKENS must be positive")
		}
		cfg.MaxTokens = val
	}

	if temp := os.Getenv("PERPLEXITY_TEMPERATURE"); temp != "" {
		val, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_TEMPERATURE: %w", err)
		}
		if val < 0 || val > 2 {
			return nil, fmt.Errorf("PERPLEXITY_TEMPERATURE must be between 0 and 2")
		}
		cfg.Temperature = val
	}

	if topP := os.Getenv("PERPLEXITY_TOP_P"); topP != "" {
		val, err := strconv.ParseFloat(topP, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_TOP_P: %w", err)
		}
		if val < 0 || val > 1 {
			return nil, fmt.Errorf("PERPLEXITY_TOP_P must be between 0 and 1")
		}
		cfg.TopP = val
	}

	if topK := os.Getenv("PERPLEXITY_TOP_K"); topK != "" {
		val, err := strconv.Atoi(topK)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_TOP_K: %w", err)
		}
		if val < 0 {
			return nil, fmt.Errorf("PERPLEXITY_TOP_K must be non-negative")
		}
		cfg.TopK = val
	}

	if timeout := os.Getenv("PERPLEXITY_TIMEOUT"); timeout != "" {
		val, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_TIMEOUT: %w", err)
		}
		if val <= 0 {
			return nil, fmt.Errorf("PERPLEXITY_TIMEOUT must be positive")
		}
		cfg.Timeout = val
	}

	if returnCitations := os.Getenv("PERPLEXITY_RETURN_CITATIONS"); returnCitations != "" {
		val, err := strconv.ParseBool(returnCitations)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_RETURN_CITATIONS: %w", err)
		}
		cfg.ReturnCitations = val
	}

	if returnImages := os.Getenv("PERPLEXITY_RETURN_IMAGES"); returnImages != "" {
		val, err := strconv.ParseBool(returnImages)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_RETURN_IMAGES: %w", err)
		}
		cfg.ReturnImages = val
	}

	if returnRelated := os.Getenv("PERPLEXITY_RETURN_RELATED"); returnRelated != "" {
		val, err := strconv.ParseBool(returnRelated)
		if err != nil {
			return nil, fmt.Errorf("invalid PERPLEXITY_RETURN_RELATED: %w", err)
		}
		cfg.ReturnRelated = val
	}

	return cfg, nil
}

// validateModel checks if the model is valid
func validateModel(model string) error {
	validModels := map[string]bool{
		types.ModelSonar:          true,
		types.ModelSonarPro:       true,
		types.ModelSonarReasoning: true,
		types.ModelSonarFinance:   true,
		types.ModelRelated:        true,
	}

	if !validModels[model] {
		return fmt.Errorf("model '%s' is not valid", model)
	}
	return nil
}

// GetAPIKey returns the API key (for testing purposes)
func (c *Config) GetAPIKey() string {
	return c.APIKey
}