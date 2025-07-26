package config

import (
	"os"
	"testing"
	"time"

	"github.com/prasanthmj/perplexity/pkg/types"
)

func TestLoadConfigWithDefaults(t *testing.T) {
	// Set only required env var
	os.Setenv("PERPLEXITY_API_KEY", "test-api-key")
	defer os.Unsetenv("PERPLEXITY_API_KEY")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.APIKey != "test-api-key" {
		t.Errorf("APIKey mismatch: got %s, want test-api-key", cfg.APIKey)
	}
	if cfg.DefaultModel != types.DefaultModel {
		t.Errorf("DefaultModel mismatch: got %s, want %s", cfg.DefaultModel, types.DefaultModel)
	}
	if cfg.MaxTokens != types.DefaultMaxTokens {
		t.Errorf("MaxTokens mismatch: got %d, want %d", cfg.MaxTokens, types.DefaultMaxTokens)
	}
	if cfg.Temperature != types.DefaultTemperature {
		t.Errorf("Temperature mismatch: got %f, want %f", cfg.Temperature, types.DefaultTemperature)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout mismatch: got %v, want %v", cfg.Timeout, 30*time.Second)
	}
}

func TestLoadConfigMissingAPIKey(t *testing.T) {
	// Ensure API key is not set
	os.Unsetenv("PERPLEXITY_API_KEY")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing API key, got nil")
	}
	if err.Error() != "PERPLEXITY_API_KEY environment variable is required" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestLoadConfigWithCustomValues(t *testing.T) {
	// Set all environment variables
	envVars := map[string]string{
		"PERPLEXITY_API_KEY":          "custom-api-key",
		"PERPLEXITY_DEFAULT_MODEL":    types.ModelSonarPro,
		"PERPLEXITY_MAX_TOKENS":       "2048",
		"PERPLEXITY_TEMPERATURE":      "0.8",
		"PERPLEXITY_TOP_P":            "0.95",
		"PERPLEXITY_TOP_K":            "10",
		"PERPLEXITY_TIMEOUT":          "60s",
		"PERPLEXITY_RETURN_CITATIONS": "false",
		"PERPLEXITY_RETURN_IMAGES":    "true",
		"PERPLEXITY_RETURN_RELATED":   "true",
	}

	for k, v := range envVars {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.APIKey != "custom-api-key" {
		t.Errorf("APIKey mismatch: got %s, want custom-api-key", cfg.APIKey)
	}
	if cfg.DefaultModel != types.ModelSonarPro {
		t.Errorf("DefaultModel mismatch: got %s, want %s", cfg.DefaultModel, types.ModelSonarPro)
	}
	if cfg.MaxTokens != 2048 {
		t.Errorf("MaxTokens mismatch: got %d, want 2048", cfg.MaxTokens)
	}
	if cfg.Temperature != 0.8 {
		t.Errorf("Temperature mismatch: got %f, want 0.8", cfg.Temperature)
	}
	if cfg.TopP != 0.95 {
		t.Errorf("TopP mismatch: got %f, want 0.95", cfg.TopP)
	}
	if cfg.TopK != 10 {
		t.Errorf("TopK mismatch: got %d, want 10", cfg.TopK)
	}
	if cfg.Timeout != 60*time.Second {
		t.Errorf("Timeout mismatch: got %v, want %v", cfg.Timeout, 60*time.Second)
	}
	if cfg.ReturnCitations != false {
		t.Errorf("ReturnCitations mismatch: got %v, want false", cfg.ReturnCitations)
	}
	if cfg.ReturnImages != true {
		t.Errorf("ReturnImages mismatch: got %v, want true", cfg.ReturnImages)
	}
	if cfg.ReturnRelated != true {
		t.Errorf("ReturnRelated mismatch: got %v, want true", cfg.ReturnRelated)
	}
}

func TestLoadConfigInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr string
	}{
		{
			name: "invalid model",
			envVars: map[string]string{
				"PERPLEXITY_API_KEY":       "test-key",
				"PERPLEXITY_DEFAULT_MODEL": "invalid-model",
			},
			wantErr: "invalid model:",
		},
		{
			name: "invalid max tokens",
			envVars: map[string]string{
				"PERPLEXITY_API_KEY":    "test-key",
				"PERPLEXITY_MAX_TOKENS": "not-a-number",
			},
			wantErr: "invalid PERPLEXITY_MAX_TOKENS:",
		},
		{
			name: "negative max tokens",
			envVars: map[string]string{
				"PERPLEXITY_API_KEY":    "test-key",
				"PERPLEXITY_MAX_TOKENS": "-100",
			},
			wantErr: "PERPLEXITY_MAX_TOKENS must be positive",
		},
		{
			name: "invalid temperature",
			envVars: map[string]string{
				"PERPLEXITY_API_KEY":     "test-key",
				"PERPLEXITY_TEMPERATURE": "not-a-float",
			},
			wantErr: "invalid PERPLEXITY_TEMPERATURE:",
		},
		{
			name: "temperature too high",
			envVars: map[string]string{
				"PERPLEXITY_API_KEY":     "test-key",
				"PERPLEXITY_TEMPERATURE": "2.5",
			},
			wantErr: "PERPLEXITY_TEMPERATURE must be between 0 and 2",
		},
		{
			name: "invalid timeout",
			envVars: map[string]string{
				"PERPLEXITY_API_KEY": "test-key",
				"PERPLEXITY_TIMEOUT": "invalid",
			},
			wantErr: "invalid PERPLEXITY_TIMEOUT:",
		},
		{
			name: "invalid boolean",
			envVars: map[string]string{
				"PERPLEXITY_API_KEY":          "test-key",
				"PERPLEXITY_RETURN_CITATIONS": "not-a-bool",
			},
			wantErr: "invalid PERPLEXITY_RETURN_CITATIONS:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			for k := range tt.envVars {
				os.Unsetenv(k)
			}

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			_, err := LoadConfig()
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if !containsString(err.Error(), tt.wantErr) {
				t.Errorf("Error message mismatch: got %v, want to contain %s", err, tt.wantErr)
			}
		})
	}
}

func TestValidateModel(t *testing.T) {
	validModels := []string{
		types.ModelSonar,
		types.ModelSonarPro,
	}

	for _, model := range validModels {
		if err := validateModel(model); err != nil {
			t.Errorf("validateModel(%s) failed: %v", model, err)
		}
	}

	invalidModels := []string{"gpt-4", "claude", "invalid", ""}
	for _, model := range invalidModels {
		if err := validateModel(model); err == nil {
			t.Errorf("validateModel(%s) should have failed", model)
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}