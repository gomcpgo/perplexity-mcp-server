# Perplexity MCP Server Developer Guide

This guide provides comprehensive documentation for developing and maintaining the Perplexity MCP server for Savant.

## Overview

The Perplexity MCP server integrates Perplexity's Sonar API to provide real-time web search capabilities with citations. This server focuses on search functionality rather than general chat, complementing Savant's existing LLM integrations.

## API Endpoint

All requests are made to:
```
https://api.perplexity.ai/chat/completions
```

## Authentication

All requests require a Bearer token in the Authorization header:
```
Authorization: Bearer YOUR_API_KEY
```

## Available Models

The model should be configurable in MCP settings. Available models include:
- `sonar` - Fast, efficient search (127k context)
- `sonar-pro` - Advanced search with larger context (200k context)
- `sonar-reasoning` - Search with reasoning capabilities (127k context)
- `sonar-reasoning-pro` - Advanced reasoning search (127k context)
- `sonar-deep-research` - Comprehensive multi-step research

## MCP Functions

### 1. perplexity_search

General web search with citations and real-time information.

#### Parameters
- `query` (string, required): The search query
- `search_domain_filter` (array, optional): Domains to include/exclude
  - Prefix with "-" to exclude (e.g., "-reddit.com")
  - Limited to 10 domains total
- `search_recency_filter` (string, optional): Time filter
  - Values: "hour", "day", "week", "month", "year"
- `return_citations` (boolean, optional): Include source URLs (default: true)
- `return_images` (boolean, optional): Include relevant images (default: false)
- `return_related_questions` (boolean, optional): Suggest follow-ups (default: false)
- `search_mode` (string, optional): "web" (default) or "academic"
- `search_context_size` (string, optional): "low", "medium", "high" (default: "medium")
- `max_tokens` (integer, optional): Maximum response length (default: 1024)
- `temperature` (float, optional): 0-2, controls randomness (default: 0.7)

#### cURL Example
```bash
curl --request POST \
  --url https://api.perplexity.ai/chat/completions \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
    "model": "sonar",
    "messages": [
      {
        "role": "user",
        "content": "What are the latest developments in quantum computing?"
      }
    ],
    "search_domain_filter": ["arxiv.org", "nature.com", "-reddit.com"],
    "search_recency_filter": "month",
    "return_citations": true,
    "return_images": false,
    "max_tokens": 1024,
    "temperature": 0.7,
    "web_search_options": {
      "search_context_size": "medium"
    }
  }'
```

#### Response Structure
```json
{
  "id": "uuid",
  "model": "sonar",
  "created": 1234567890,
  "usage": {
    "prompt_tokens": 50,
    "completion_tokens": 300,
    "total_tokens": 350,
    "citation_tokens": 150
  },
  "citations": [
    "https://example.com/article1",
    "https://example.com/article2"
  ],
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Response with citations..."
      },
      "finish_reason": "stop"
    }
  ]
}
```

### 2. perplexity_academic_search

Specialized search for academic and research content.

#### Parameters
- `query` (string, required): The academic search query
- `search_domain_filter` (array, optional): Academic domains
  - Suggested: ["arxiv.org", "pubmed.ncbi.nlm.nih.gov", "scholar.google.com", "researchgate.net"]
- `date_range` (object, optional):
  - `after`: "MM/DD/YYYY" format
  - `before`: "MM/DD/YYYY" format
- `return_citations` (boolean, optional): Include citations (default: true)
- `max_tokens` (integer, optional): Maximum response length (default: 2048)
- `temperature` (float, optional): 0-2 (default: 0.3 for factual accuracy)

#### cURL Example
```bash
curl --request POST \
  --url https://api.perplexity.ai/chat/completions \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
    "model": "sonar-pro",
    "messages": [
      {
        "role": "user",
        "content": "Recent advances in CRISPR gene editing techniques"
      }
    ],
    "search_mode": "academic",
    "search_domain_filter": ["arxiv.org", "nature.com", "pubmed.ncbi.nlm.nih.gov"],
    "search_recency_filter": "year",
    "return_citations": true,
    "max_tokens": 2048,
    "temperature": 0.3,
    "web_search_options": {
      "search_context_size": "high"
    }
  }'
```

### 3. perplexity_financial_search

Search SEC filings and financial documents.

#### Parameters
- `query` (string, required): The financial search query
- `company` (string, optional): Company name or ticker symbol
- `filing_type` (string, optional): Type of filing (10-K, 10-Q, 8-K, etc.)
- `date_range` (object, optional):
  - `after`: "MM/DD/YYYY" format
  - `before`: "MM/DD/YYYY" format
- `return_citations` (boolean, optional): Include citations (default: true)
- `max_tokens` (integer, optional): Maximum response length (default: 2048)
- `temperature` (float, optional): 0-2 (default: 0.2 for accuracy)

#### cURL Example
```bash
curl --request POST \
  --url https://api.perplexity.ai/chat/completions \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
    "model": "sonar-pro",
    "messages": [
      {
        "role": "user",
        "content": "Apple Inc latest quarterly earnings and revenue growth"
      }
    ],
    "search_domain_filter": ["sec.gov", "investor.apple.com"],
    "search_recency_filter": "month",
    "return_citations": true,
    "max_tokens": 2048,
    "temperature": 0.2,
    "web_search_options": {
      "search_context_size": "high"
    }
  }'
```

**Note**: For SEC-specific searches, include "sec" in the search_domain parameter:
```json
{
  "search_domain": "sec",
  "messages": [{"role": "user", "content": "AAPL 10-K filing analysis"}]
}
```

### 4. perplexity_filtered_search

Advanced search with multiple filtering options.

#### Parameters
- `query` (string, required): The search query
- `filters` (object, required):
  - `domains` (array, optional): Include/exclude domains
  - `date_range` (object, optional):
    - `after`: "MM/DD/YYYY"
    - `before`: "MM/DD/YYYY"
  - `location` (object, optional):
    - `latitude`: number
    - `longitude`: number
    - `country`: ISO country code
  - `recency` (string, optional): Time filter
  - `content_type` (string, optional): Filter by type
- `return_citations` (boolean, optional): Include citations (default: true)
- `return_images` (boolean, optional): Include images (default: false)
- `max_tokens` (integer, optional): Maximum response length (default: 1024)
- `temperature` (float, optional): 0-2 (default: 0.7)

#### cURL Example
```bash
curl --request POST \
  --url https://api.perplexity.ai/chat/completions \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
    "model": "sonar",
    "messages": [
      {
        "role": "user",
        "content": "Climate change policies in California"
      }
    ],
    "search_domain_filter": ["ca.gov", "climate.ca.gov", "-blog.com"],
    "search_after_date_filter": "01/01/2024",
    "search_before_date_filter": "12/31/2024",
    "return_citations": true,
    "max_tokens": 1024,
    "temperature": 0.7,
    "web_search_options": {
      "search_context_size": "medium",
      "user_location": {
        "latitude": 37.7749,
        "longitude": -122.4194,
        "country": "US"
      }
    }
  }'
```

## Default Values

When implementing the MCP server, use these fair defaults:

```json
{
  "model": "sonar",
  "temperature": 0.7,
  "max_tokens": 1024,
  "return_citations": true,
  "return_images": false,
  "return_related_questions": false,
  "search_context_size": "medium"
}
```

For specific function types:
- **Academic search**: Use temperature 0.3, max_tokens 2048, search_context_size "high"
- **Financial search**: Use temperature 0.2, max_tokens 2048, search_context_size "high"
- **General search**: Use defaults above

## Error Handling

Common error responses:

### Rate Limit Error (429)
```json
{
  "error": {
    "message": "Rate limit exceeded",
    "type": "rate_limit_error",
    "code": "rate_limit_exceeded"
  }
}
```

### Invalid API Key (401)
```json
{
  "error": {
    "message": "Invalid API key",
    "type": "authentication_error",
    "code": "invalid_api_key"
  }
}
```

### Bad Request (400)
```json
{
  "error": {
    "message": "Invalid request format",
    "type": "invalid_request_error",
    "code": "invalid_request"
  }
}
```

## Implementation Guidelines

### 1. Model Configuration
- Model should be configurable in MCP settings, not hardcoded
- Validate model selection against available models
- Consider different models for different search types

### 2. Response Processing
- Extract citations from the response
- Format citations as clickable links
- Handle cases where citations might be empty

### 3. Cost Optimization
- Use appropriate search_context_size based on query complexity
- Default to "sonar" model unless "sonar-pro" features needed
- Implement caching for repeated queries if appropriate

### 4. User Experience
- Show citation count and sources
- Provide clear error messages



## Testing

Test cases to implement:

1. **Basic Search**
   ```bash
   # Test with simple query
   curl -X POST ... --data '{"messages": [{"role": "user", "content": "What is quantum computing?"}]}'
   ```

2. **Domain Filtering**
   ```bash
   # Test include/exclude domains
   curl -X POST ... --data '{"search_domain_filter": ["wikipedia.org", "-reddit.com"]}'
   ```

3. **Academic Search**
   ```bash
   # Test academic mode
   curl -X POST ... --data '{"search_mode": "academic", "search_domain_filter": ["arxiv.org"]}'
   ```

4. **Date Filtering**
   ```bash
   # Test date range
   curl -X POST ... --data '{"search_after_date_filter": "01/01/2025"}'
   ```

5. **Error Handling**
   ```bash
   # Test with invalid API key
   curl -X POST ... --header 'Authorization: Bearer INVALID_KEY'
   ```


## Security Considerations

1. **API Key Storage**: Never hardcode API keys
2. **Input Validation**: Sanitize user queries
3. **Domain Validation**: Validate domain filter inputs
4. **Response Validation**: Verify response structure before processing

## Monitoring

Track these metrics:
- API usage and costs
- Average response time
- Citation count per query
- Error rates by type
- Model usage distribution

## Future Enhancements

Consider these potential improvements:
1. Implement response caching
2. Add support for structured output formats
3. Enable async operations for deep research
4. Add webhook support for long-running queries
5. Implement batch search operations

## Support

For API issues: api@perplexity.ai
For documentation: https://docs.perplexity.ai
For status updates: Check Perplexity's system status page
