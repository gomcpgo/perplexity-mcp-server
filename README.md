# Perplexity MCP Server

An MCP (Model Context Protocol) server that provides access to Perplexity AI's powerful search capabilities, including web search, academic research, financial data, and advanced filtering options.

## Features

The Perplexity MCP server offers **six functions** for comprehensive search and result management:

### Search Functions (4)
Each optimized for different use cases. **All functions automatically return source URLs** and save results locally if caching is enabled.

1. **`perplexity_search`**: General web search with real-time information. Best for current events, general knowledge, and quick facts.

2. **`perplexity_academic_search`**: Automatically filters to academic sources (arxiv.org, pubmed, journals). Best for research papers, scientific studies, and scholarly content.

3. **`perplexity_financial_search`**: Optimized for financial domains and recent data. Best for stock analysis, earnings reports, SEC filings, and market trends.

4. **`perplexity_filtered_search`**: Advanced search with multiple filtering options. Best when you need specific domain filtering, content types, or location-based results.

### Cache Management Functions (2)
Manage previously saved search results for easy reference and reuse.

5. **`list_previous`**: List all previous search queries with unique IDs, sorted by recency. Returns JSON array with query details.

6. **`get_previous_result`**: Retrieve a previously cached search result by its unique 10-character ID.

## Installation

1. Ensure you have Go 1.23 or later installed
2. Clone this repository
3. Build the server:
   ```bash
   ./run.sh build
   ```

## Configuration

The server requires a Perplexity API key and supports various configuration options through environment variables:

### Required
- `PERPLEXITY_API_KEY`: Your Perplexity AI API key

### Optional
- `PERPLEXITY_DEFAULT_MODEL`: Default model to use (default: "sonar")
  - `sonar`: Fast, cost-effective search for quick facts
  - `sonar-pro`: Comprehensive search with better depth and coverage
- `PERPLEXITY_MAX_TOKENS`: Maximum tokens in response (default: 1024)
- `PERPLEXITY_TEMPERATURE`: Response randomness 0-2 (default: 0.2)
- `PERPLEXITY_TOP_P`: Nucleus sampling parameter (default: 0.9)
- `PERPLEXITY_TOP_K`: Top-k sampling parameter (default: 0)
- `PERPLEXITY_TIMEOUT`: Request timeout duration (default: 30s)
- `PERPLEXITY_RETURN_IMAGES`: Include images by default (default: false)
- `PERPLEXITY_RETURN_RELATED`: Include related questions by default (default: false)
- `PERPLEXITY_RESULTS_ROOT_FOLDER`: Directory to store cached search results (default: empty/disabled)

## Usage

### MCP Server Mode

Run the server in MCP mode (default):
```bash
export PERPLEXITY_API_KEY="your-api-key"
./run.sh run
# or directly: ./perplexity
```

### Terminal Mode (CLI Testing)

Test individual functions directly from the command line:

```bash
export PERPLEXITY_API_KEY="your-api-key"

# Test different search types
./run.sh search "latest AI news" sonar-pro
./run.sh academic "quantum computing" sonar-pro
./run.sh financial "AAPL earnings" sonar-pro
./run.sh filtered "renewable energy" sonar-pro

# Cache management
./run.sh list                    # List previous queries
./run.sh get ABC123XYZ0         # Get cached result by ID
```

### Integration Tests

Run integration tests against the real Perplexity API:
```bash
export PERPLEXITY_API_KEY="your-api-key"
./run.sh integration-test
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

## Local Result Caching

The server automatically caches search results when `PERPLEXITY_RESULTS_ROOT_FOLDER` is configured:

- **Storage**: Each result is saved in `/unique_id/result.md` with metadata in `/unique_id/metadata.yaml`
- **Unique IDs**: 10-character alphanumeric identifiers (e.g., `A1B2C3D4E5`)
- **Result ID**: When caching is enabled, search responses include `**Result ID:** ABC123XYZ0`
- **No Reuse**: Each search creates a new cached entry, even for identical queries
- **LLM Integration**: Perfect for LLMs to reference previous searches in conversations

### Cache Management Examples

```bash
# List previous searches
echo '{"method": "tools/call", "params": {"name": "list_previous", "arguments": {}}}' | ./perplexity

# Get specific result
echo '{"method": "tools/call", "params": {"name": "get_previous_result", "arguments": {"unique_id": "A1B2C3D4E5"}}}' | ./perplexity
```

## Function Reference

### perplexity_search

Perform a general web search.

**Parameters:**
- `query` (required): The search query
- `model`: Choose 'sonar' for quick searches or 'sonar-pro' for comprehensive results (default: sonar)
- `search_domain_filter`: Array of domains to include
- `search_exclude_domains`: Array of domains to exclude
- `search_recency_filter`: Time filter (hour, day, week, month, year)
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
- `model`: Defaults to 'sonar-pro' for comprehensive academic results
- `search_domain_filter`: Array of academic domains
- `search_recency_filter`: Time filter
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
- `model`: Defaults to 'sonar-pro' for comprehensive financial data
- `search_recency_filter`: Time filter
- `date_range_start`: Report start date
- `date_range_end`: Report end date
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
- `model`: Choose based on needs (defaults to sonar-pro)
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

### list_previous

List all previous search queries with metadata.

**Parameters:** None

**Response:** JSON array with query history, sorted by recency (most recent first).

**Example:**
```json
[
  {
    "query": "latest AI developments",
    "unique_id": "A1B2C3D4E5",
    "datetime": "2025-01-15T10:30:45Z",
    "search_type": "general"
  },
  {
    "query": "quantum computing research",
    "unique_id": "X9Y8Z7W6V5",
    "datetime": "2025-01-15T09:15:30Z",
    "search_type": "academic"
  }
]
```

### get_previous_result

Retrieve a cached search result by unique ID.

**Parameters:**
- `unique_id` (required): The 10-character alphanumeric ID of the cached result

**Returns:** The complete markdown result from the cached search.

**Example:**
```json
{
  "unique_id": "A1B2C3D4E5"
}
```

## Response Format

All search functions return responses in the following format:

1. **Main Content**: The search results and answer
2. **Source URLs**: A list of source URLs that the LLM can fetch for more details
3. **Detailed Sources** (if available): Title, URL, and snippet for each source
4. **Related Questions** (if requested): Suggested follow-up questions
5. **Result ID** (if caching enabled): Unique 10-character ID for retrieving this result later

Example response structure:
```
[Main search results content...]

## Source URLs
1. https://example.com/article1
2. https://example.com/article2
3. https://example.com/article3

## Detailed Sources
1. **Article Title**
   URL: https://example.com/article1
   Snippet: Brief excerpt from the article...

## Related Questions
- What are the latest developments?
- How does this compare to...?

**Result ID:** A1B2C3D4E5
```

## Development

### Running Tests

Run unit tests:
```bash
./run.sh test
```

Run integration tests with real API:
```bash
./run.sh integration-test
```

### Project Structure

The server follows clean architecture principles with separation of concerns:

```
perplexity/
├── cmd/
│   └── main.go              # Thin entry point with terminal mode (~200 lines)
├── pkg/
│   ├── handler/             # MCP protocol layer
│   │   ├── handler.go       # Main MCP handler  
│   │   ├── tools.go         # Tool definitions
│   │   └── search_handlers.go # Parameter extraction
│   ├── search/              # Core business logic
│   │   ├── types.go         # Local search types
│   │   ├── search.go        # Strongly-typed search functions
│   │   └── client.go        # Perplexity API client
│   ├── cache/               # Result caching system
│   ├── config/              # Configuration management
│   └── types/               # Perplexity API types
├── test/
│   └── test.go             # Integration tests
└── README.md
```

### Architecture Benefits

- **Thin main.go**: Reduced from 360 to 197 lines (45% reduction)
- **Terminal mode**: Direct CLI testing without MCP protocol overhead
- **Separation of concerns**: MCP protocol handling separate from business logic
- **Strongly-typed**: Core functions use proper Go structs instead of `map[string]interface{}`
- **Local types**: Each package owns its types, preventing circular dependencies
- **Easy testing**: Business logic can be tested independently

## Error Handling

The server handles various error conditions:
- Invalid or missing API key (401)
- Rate limiting (429)
- Invalid parameters (400)
- Server errors (500)

Errors are returned with descriptive messages to help diagnose issues.

## License

MIT License - see LICENSE file for details.