package types

import (
	"encoding/json"
	"testing"
)

func TestMessageMarshaling(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Test message",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	if decoded.Role != msg.Role || decoded.Content != msg.Content {
		t.Errorf("Message mismatch: got %+v, want %+v", decoded, msg)
	}
}

func TestPerplexityRequestMarshaling(t *testing.T) {
	req := PerplexityRequest{
		Model: ModelSonarPro,
		Messages: []Message{
			{Role: "user", Content: "Test query"},
		},
		MaxTokens:           1024,
		Temperature:         0.7,
		ReturnCitations:     true,
		SearchDomainFilter:  []string{"example.com", "test.com"},
		SearchRecencyFilter: RecencyWeek,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	var decoded PerplexityRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if decoded.Model != req.Model {
		t.Errorf("Model mismatch: got %s, want %s", decoded.Model, req.Model)
	}
	if len(decoded.Messages) != len(req.Messages) {
		t.Errorf("Messages count mismatch: got %d, want %d", len(decoded.Messages), len(req.Messages))
	}
	if decoded.MaxTokens != req.MaxTokens {
		t.Errorf("MaxTokens mismatch: got %d, want %d", decoded.MaxTokens, req.MaxTokens)
	}
	if len(decoded.SearchDomainFilter) != len(req.SearchDomainFilter) {
		t.Errorf("SearchDomainFilter count mismatch: got %d, want %d", 
			len(decoded.SearchDomainFilter), len(req.SearchDomainFilter))
	}
}

func TestPerplexityResponseMarshaling(t *testing.T) {
	resp := PerplexityResponse{
		ID:      "test-id",
		Model:   ModelSonar,
		Object:  "chat.completion",
		Created: 1234567890,
		Choices: []Choice{
			{
				Index:        0,
				FinishReason: "stop",
				Message: Message{
					Role:    "assistant",
					Content: "Test response",
				},
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
			CitationTokens:   5,
		},
		Citations: []string{"https://example.com", "https://test.com"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var decoded PerplexityResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if decoded.ID != resp.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, resp.ID)
	}
	if decoded.Model != resp.Model {
		t.Errorf("Model mismatch: got %s, want %s", decoded.Model, resp.Model)
	}
	if len(decoded.Choices) != len(resp.Choices) {
		t.Errorf("Choices count mismatch: got %d, want %d", len(decoded.Choices), len(resp.Choices))
	}
	if decoded.Usage.TotalTokens != resp.Usage.TotalTokens {
		t.Errorf("TotalTokens mismatch: got %d, want %d", decoded.Usage.TotalTokens, resp.Usage.TotalTokens)
	}
	if len(decoded.Citations) != len(resp.Citations) {
		t.Errorf("Citations count mismatch: got %d, want %d", len(decoded.Citations), len(resp.Citations))
	}
}

func TestErrorResponseMarshaling(t *testing.T) {
	errResp := ErrorResponse{}
	errResp.Error.Type = "invalid_request_error"
	errResp.Error.Message = "Invalid API key"
	errResp.Error.Code = "401"

	data, err := json.Marshal(errResp)
	if err != nil {
		t.Fatalf("Failed to marshal error response: %v", err)
	}

	var decoded ErrorResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if decoded.Error.Type != errResp.Error.Type {
		t.Errorf("Error type mismatch: got %s, want %s", decoded.Error.Type, errResp.Error.Type)
	}
	if decoded.Error.Message != errResp.Error.Message {
		t.Errorf("Error message mismatch: got %s, want %s", decoded.Error.Message, errResp.Error.Message)
	}
	if decoded.Error.Code != errResp.Error.Code {
		t.Errorf("Error code mismatch: got %s, want %s", decoded.Error.Code, errResp.Error.Code)
	}
}

func TestSearchParametersWithPointers(t *testing.T) {
	boolTrue := true
	intVal := 512
	floatVal := 0.5

	params := SearchParameters{
		Query:           "test query",
		Model:           ModelSonarReasoning,
		ReturnCitations: &boolTrue,
		MaxTokens:       &intVal,
		Temperature:     &floatVal,
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal search parameters: %v", err)
	}

	var decoded SearchParameters
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal search parameters: %v", err)
	}

	if decoded.Query != params.Query {
		t.Errorf("Query mismatch: got %s, want %s", decoded.Query, params.Query)
	}
	if decoded.ReturnCitations == nil || *decoded.ReturnCitations != *params.ReturnCitations {
		t.Errorf("ReturnCitations mismatch")
	}
	if decoded.MaxTokens == nil || *decoded.MaxTokens != *params.MaxTokens {
		t.Errorf("MaxTokens mismatch")
	}
	if decoded.Temperature == nil || *decoded.Temperature != *params.Temperature {
		t.Errorf("Temperature mismatch")
	}
}

func TestModelConstants(t *testing.T) {
	models := []string{
		ModelSonar,
		ModelSonarPro,
		ModelSonarReasoning,
		ModelSonarFinance,
		ModelRelated,
	}

	for _, model := range models {
		if model == "" {
			t.Errorf("Model constant should not be empty")
		}
	}
}

func TestRecencyConstants(t *testing.T) {
	recencies := []string{
		RecencyHour,
		RecencyDay,
		RecencyWeek,
		RecencyMonth,
		RecencyYear,
	}

	for _, recency := range recencies {
		if recency == "" {
			t.Errorf("Recency constant should not be empty")
		}
	}
}

func TestDefaultConstants(t *testing.T) {
	if DefaultModel == "" {
		t.Error("DefaultModel should not be empty")
	}
	if DefaultMaxTokens <= 0 {
		t.Error("DefaultMaxTokens should be positive")
	}
	if DefaultTemperature < 0 || DefaultTemperature > 1 {
		t.Error("DefaultTemperature should be between 0 and 1")
	}
	if DefaultTopP < 0 || DefaultTopP > 1 {
		t.Error("DefaultTopP should be between 0 and 1")
	}
	if DefaultContextSize <= 0 {
		t.Error("DefaultContextSize should be positive")
	}
}