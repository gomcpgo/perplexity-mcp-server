# Perplexity MCP Server

An MCP (Model Context Protocol) server that provides access to Perplexity AI's powerful search capabilities, including web search, academic research, financial data, and advanced filtering options.

## Features

The Perplexity MCP server offers four specialized search functions:

### 1. `perplexity_search`
General web search with comprehensive parameters for filtering and customization.

### 2. `perplexity_academic_search`
Specialized search for academic papers and scholarly content with subject area filtering.

### 3. `perplexity_financial_search`
Search financial data, SEC filings, and market information with ticker and report type filtering.

### 4. `perplexity_filtered_search`
Advanced search with comprehensive filtering options including content type, language, location, and custom filters.

## Installation

1. Ensure you have Go 1.23 or later installed
2. Clone this repository
3. Build the server:
   ```bash
   ./run.sh
   ```

## Configuration

The server requires a Perplexity API key and supports various configuration options through environment variables:

### Required
- `PERPLEXITY_API_KEY`: Your Perplexity AI API key

### Optional
- `PERPLEXITY_DEFAULT_MODEL`: Default model to use (default: "sonar")
  - Available models: `sonar`, `sonar-pro`, `sonar-reasoning`, `sonar-reasoning-pro`, `sonar-deep-research`, `r1-1776`
- `PERPLEXITY_MAX_TOKENS`: Maximum tokens in response (default: 1024)
- `PERPLEXITY_TEMPERATURE`: Response randomness 0-2 (default: 0.2)
- `PERPLEXITY_TOP_P`: Nucleus sampling parameter (default: 0.9)
- `PERPLEXITY_TOP_K`: Top-k sampling parameter (default: 0)
- `PERPLEXITY_TIMEOUT`: Request timeout duration (default: 30s)
- `PERPLEXITY_RETURN_CITATIONS`: Include citations by default (default: true)
- `PERPLEXITY_RETURN_IMAGES`: Include images by default (default: false)
- `PERPLEXITY_RETURN_RELATED`: Include related questions by default (default: false)

## Usage

### MCP Server Mode

Run the server in MCP mode (default):
```bash
export PERPLEXITY_API_KEY="your-api-key"
./perplexity
```

### Test Mode

Run integration tests against the real Perplexity API:
```bash
export PERPLEXITY_API_KEY="your-api-key"
./perplexity -test
```

### MCP Client Configuration

To use this server with an MCP client, add it to your client configuration:

```json
{
  "servers": {
    "perplexity": {
      "command": "path/to/perplexity",
      "env": {
        "PERPLEXITY_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Function Reference

### perplexity_search

Perform a general web search.

**Parameters:**
- `query` (required): The search query
- `model`: Model to use (sonar, sonar-pro, sonar-reasoning, sonar-reasoning-pro, sonar-deep-research, r1-1776)
- `search_domain_filter`: Array of domains to include
- `search_exclude_domains`: Array of domains to exclude
- `search_recency_filter`: Time filter (hour, day, week, month, year)
- `return_citations`: Include citations
- `return_images`: Include images
- `return_related_questions`: Include related questions
- `max_tokens`: Maximum response tokens
- `temperature`: Response randomness (0-2)
- `date_range_start`: Start date (YYYY-MM-DD)
- `date_range_end`: End date (YYYY-MM-DD)
- `location`: Geo-specific search location

**Example:**
```json
{
  "query": "latest AI developments",
  "model": "sonar-pro",
  "search_recency_filter": "week",
  "return_citations": true
}
```

### perplexity_academic_search

Search academic papers and scholarly content.

**Parameters:**
- `query` (required): The academic search query
- `subject_area`: Academic subject (e.g., "Physics", "Computer Science")
- `model`: Model to use (defaults to sonar-reasoning)
- `search_domain_filter`: Array of academic domains
- `search_recency_filter`: Time filter
- `return_citations`: Include citations (default: true)
- `max_tokens`: Maximum response tokens
- `temperature`: Response randomness

**Example:**
```json
{
  "query": "quantum computing applications",
  "subject_area": "Physics",
  "search_recency_filter": "year"
}
```

### perplexity_financial_search

Search financial data and SEC filings.

**Parameters:**
- `query` (required): The financial search query
- `ticker`: Stock ticker symbol (e.g., "AAPL")
- `company_name`: Company name
- `report_type`: Financial report type (e.g., "10-K", "10-Q", "8-K")
- `model`: Model to use (defaults to sonar-reasoning-pro)
- `search_recency_filter`: Time filter
- `date_range_start`: Report start date
- `date_range_end`: Report end date
- `return_citations`: Include citations (default: true)
- `max_tokens`: Maximum response tokens

**Example:**
```json
{
  "query": "quarterly earnings",
  "ticker": "MSFT",
  "report_type": "10-Q",
  "search_recency_filter": "month"
}
```

### perplexity_filtered_search

Advanced search with comprehensive filtering.

**Parameters:**
- `query` (required): The search query
- `model`: Model to use (defaults to sonar-pro)
- `search_domain_filter`: Array of domains to include
- `search_exclude_domains`: Array of domains to exclude
- `search_recency_filter`: Time filter
- `content_type`: Type of content (news, academic, blog, etc.)
- `file_type`: File type filter (pdf, doc, html, etc.)
- `language`: Language filter
- `country`: Country for geo-specific search
- `date_range_start`: Start date
- `date_range_end`: End date
- `return_citations`: Include citations
- `return_images`: Include images
- `return_related_questions`: Include related questions
- `max_tokens`: Maximum response tokens
- `temperature`: Response randomness
- `custom_filters`: Object with additional key-value filters

**Example:**
```json
{
  "query": "renewable energy innovations",
  "content_type": "news",
  "language": "English",
  "country": "Germany",
  "search_recency_filter": "month",
  "custom_filters": {
    "industry": "energy",
    "technology": "solar"
  }
}
```

## Development

### Running Tests

Run unit tests:
```bash
go test ./pkg/...
```

Run integration tests with real API:
```bash
go run cmd/main.go -test
```

### Project Structure
```
perplexity/
├── cmd/
│   └── main.go              # MCP server entry point
├── pkg/
│   ├── types/               # API types and constants
│   ├── perplexity/          # Perplexity client and search functions
│   └── config/              # Configuration management
├── test/
│   └── integration_test.go  # Integration tests
└── README.md
```

## Error Handling

The server handles various error conditions:
- Invalid or missing API key (401)
- Rate limiting (429)
- Invalid parameters (400)
- Server errors (500)

Errors are returned with descriptive messages to help diagnose issues.

## License

MIT License - see LICENSE file for details.