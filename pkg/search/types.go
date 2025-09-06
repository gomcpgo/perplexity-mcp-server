package search

// SearchParams represents strongly-typed search parameters
type SearchParams struct {
	// Common parameters
	Query                    string             `json:"query"`
	SearchType               string             `json:"search_type"`
	Model                    string             `json:"model,omitempty"`
	SearchDomainFilter       []string           `json:"search_domain_filter,omitempty"`
	SearchExcludeDomains     []string           `json:"search_exclude_domains,omitempty"`
	SearchRecencyFilter      string             `json:"search_recency_filter,omitempty"`
	ReturnImages             *bool              `json:"return_images,omitempty"`
	ReturnRelatedQuestions   *bool              `json:"return_related_questions,omitempty"`
	MaxTokens                *int               `json:"max_tokens,omitempty"`
	Temperature              *float64           `json:"temperature,omitempty"`
	DateRangeStart           string             `json:"date_range_start,omitempty"`
	DateRangeEnd             string             `json:"date_range_end,omitempty"`
	Location                 string             `json:"location,omitempty"`

	// Academic-specific parameters
	SubjectArea              string             `json:"subject_area,omitempty"`

	// Financial-specific parameters
	Ticker                   string             `json:"ticker,omitempty"`
	CompanyName              string             `json:"company_name,omitempty"`
	ReportType               string             `json:"report_type,omitempty"`

	// Filtered search parameters
	ContentType              string             `json:"content_type,omitempty"`
	FileType                 string             `json:"file_type,omitempty"`
	Language                 string             `json:"language,omitempty"`
	Country                  string             `json:"country,omitempty"`
	CustomFilters            map[string]interface{} `json:"custom_filters,omitempty"`
}

// SearchResult represents a search operation result
type SearchResult struct {
	Content  string
	UniqueID string
	Error    error
}