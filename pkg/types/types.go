package types

// Model constants
const (
	ModelSonar    = "sonar"
	ModelSonarPro = "sonar-pro"
)

// Recency filter constants
const (
	RecencyHour  = "hour"
	RecencyDay   = "day"
	RecencyWeek  = "week"
	RecencyMonth = "month"
	RecencyYear  = "year"
)

// Default values
const (
	DefaultModel           = ModelSonar
	DefaultMaxTokens       = 1024
	DefaultTemperature     = 0.2
	DefaultTopP            = 0.9
	DefaultTopK            = 0
	DefaultReturnImages    = false
	DefaultReturnRelated   = false
	DefaultSearchMode      = "web"
	DefaultContextSize     = 5
)

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// PerplexityRequest represents the request to Perplexity API
type PerplexityRequest struct {
	Model                    string   `json:"model"`
	Messages                 []Message `json:"messages"`
	MaxTokens                int      `json:"max_tokens,omitempty"`
	Temperature              float64  `json:"temperature,omitempty"`
	TopP                     float64  `json:"top_p,omitempty"`
	TopK                     int      `json:"top_k,omitempty"`
	Stream                   bool     `json:"stream,omitempty"`
	PresencePenalty          float64  `json:"presence_penalty,omitempty"`
	FrequencyPenalty         float64  `json:"frequency_penalty,omitempty"`
	SearchDomainFilter       []string `json:"search_domain_filter,omitempty"`
	SearchExcludeDomains     []string `json:"search_exclude_domains,omitempty"`
	ReturnImages             bool     `json:"return_images,omitempty"`
	ReturnRelatedQuestions   bool     `json:"return_related_questions,omitempty"`
	SearchRecencyFilter      string   `json:"search_recency_filter,omitempty"`
	ReturnCitations          bool     `json:"return_citations"`
	CitationQuality          string   `json:"citation_quality,omitempty"`
	SearchMode               string   `json:"search_mode,omitempty"`
	DateRangeStart           string   `json:"date_range_start,omitempty"`
	DateRangeEnd             string   `json:"date_range_end,omitempty"`
	Location                 string   `json:"location,omitempty"`
	SearchContextSize        int      `json:"search_context_size,omitempty"`
}

// PerplexityResponse represents the response from Perplexity API
type PerplexityResponse struct {
	ID                string     `json:"id"`
	Model             string     `json:"model"`
	Object            string     `json:"object"`
	Created           int64      `json:"created"`
	Choices           []Choice   `json:"choices"`
	Usage             Usage      `json:"usage"`
	Citations         []string   `json:"citations,omitempty"`
	SearchResults     []SearchResult `json:"search_results,omitempty"`
	RelatedQuestions  []string   `json:"related_questions,omitempty"`
}

// Choice represents a response choice
type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
	Delta        *Message `json:"delta,omitempty"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens      int `json:"prompt_tokens"`
	CompletionTokens  int `json:"completion_tokens"`
	TotalTokens       int `json:"total_tokens"`
	CitationTokens    int `json:"citation_tokens,omitempty"`
}

// SearchResult represents a search result with citation
type SearchResult struct {
	URL     string `json:"url"`
	Title   string `json:"title,omitempty"`
	Snippet string `json:"snippet,omitempty"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

// SearchParameters contains common parameters for search functions
type SearchParameters struct {
	Query                    string   `json:"query"`
	Model                    string   `json:"model,omitempty"`
	SearchDomainFilter       []string `json:"search_domain_filter,omitempty"`
	SearchExcludeDomains     []string `json:"search_exclude_domains,omitempty"`
	SearchRecencyFilter      string   `json:"search_recency_filter,omitempty"`
	ReturnCitations          *bool    `json:"return_citations,omitempty"`
	ReturnImages             *bool    `json:"return_images,omitempty"`
	ReturnRelatedQuestions   *bool    `json:"return_related_questions,omitempty"`
	MaxTokens                *int     `json:"max_tokens,omitempty"`
	Temperature              *float64 `json:"temperature,omitempty"`
	TopP                     *float64 `json:"top_p,omitempty"`
	TopK                     *int     `json:"top_k,omitempty"`
	SearchMode               string   `json:"search_mode,omitempty"`
	CitationQuality          string   `json:"citation_quality,omitempty"`
	DateRangeStart           string   `json:"date_range_start,omitempty"`
	DateRangeEnd             string   `json:"date_range_end,omitempty"`
	Location                 string   `json:"location,omitempty"`
	SearchContextSize        *int     `json:"search_context_size,omitempty"`
}

// AcademicSearchParameters contains parameters specific to academic search
type AcademicSearchParameters struct {
	SearchParameters
	SubjectArea string `json:"subject_area,omitempty"`
}

// FinancialSearchParameters contains parameters specific to financial search
type FinancialSearchParameters struct {
	SearchParameters
	Ticker       string `json:"ticker,omitempty"`
	CompanyName  string `json:"company_name,omitempty"`
	ReportType   string `json:"report_type,omitempty"`
}

// FilteredSearchParameters contains all advanced filtering options
type FilteredSearchParameters struct {
	SearchParameters
	ContentType      string   `json:"content_type,omitempty"`
	FileType         string   `json:"file_type,omitempty"`
	Language         string   `json:"language,omitempty"`
	Country          string   `json:"country,omitempty"`
	CustomFilters    map[string]string `json:"custom_filters,omitempty"`
}