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

### Search Models
- `sonar` - Lightweight, cost-effective search model (127k context)
- `sonar-pro` - Advanced search offering supporting complex queries (200k context)

### Reasoning Models  
- `sonar-reasoning` - Fast, real-time reasoning model with search
- `sonar-reasoning-pro` - Precise reasoning powered by DeepSeek-R1 with search

### Research Model
- `sonar-deep-research` - Expert-level research model for comprehensive reports

### Offline Model (No Web Search)
- `r1-1776` - Conversational reasoning model without web search (based on DeepSeek-R1)

## Model Usage Guidelines

- **Quick factual queries**: Use `sonar` (fast, cost-effective)
- **General web search with more depth**: Use `sonar-pro`
- **Queries requiring analysis/reasoning**: Use `sonar-reasoning` or `sonar-reasoning-pro`
  - Complex problem-solving, comparisons, decision-making
  - Financial analysis, technical analysis
  - Multi-step logical reasoning
- **Comprehensive research reports**: Use `sonar-deep-research`
  - Exhaustive research with 20-50 searches
  - Generates detailed, structured reports
- **Academic search**: Use ANY model with `search_mode: "academic"`
  - Filters results to scholarly sources
  - Works with all search-enabled models
- **Creative/offline tasks**: Use `r1-1776` (Note: This model does NOT support web search)

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
    "model": "sonar-reasoning",
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
    "model": "sonar-reasoning-pro",
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

**Note**: For SEC-specific searches, you can also use:
```json
{
  "search_mode": "sec",
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

## Special Features

### Deep Research with sonar-deep-research

For comprehensive research tasks, use the `sonar-deep-research` model:

```bash
curl --request POST \
  --url https://api.perplexity.ai/chat/completions \
  --header 'Authorization: Bearer YOUR_API_KEY' \
  --header 'Content-Type: application/json' \
  --data '{
    "model": "sonar-deep-research",
    "messages": [
      {
        "role": "user",
        "content": "Provide an in-depth analysis of renewable energy adoption trends globally"
      }
    ],
    "max_tokens": 4000,
    "reasoning_effort": "high"
  }'
```

**Note**: Deep research queries can take 30+ seconds and may conduct 20-50 searches.

### Reasoning Models

When using reasoning models (`sonar-reasoning`, `sonar-reasoning-pro`), the response includes Chain of Thought (CoT) reasoning:

```json
{
  "usage": {
    "reasoning_tokens": 5000,
    "num_search_queries": 10
  }
}
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
- **Academic search**: Use any model with `search_mode: "academic"`, temperature 0.3, max_tokens 2048, search_context_size "high"
- **Financial search**: Use `sonar-reasoning-pro` for complex analysis, temperature 0.2, max_tokens 2048, search_context_size "high"
- **General search**: Use defaults above
- **Complex analysis**: Use `sonar-reasoning` or `sonar-reasoning-pro`, temperature 0.3-0.5
- **Deep research**: Use `sonar-deep-research`, max_tokens 4000, reasoning_effort "high"

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
- Note that `r1-1776` does NOT support web search

### 2. Response Processing
- Extract citations from the response
- Format citations as clickable links
- Handle cases where citations might be empty
- For reasoning models, handle the CoT output appropriately

### 3. Cost Optimization
- Use appropriate search_context_size based on query complexity
- Default to "sonar" model unless advanced features needed
- Use reasoning models for complex queries requiring step-by-step analysis
- Reserve `sonar-deep-research` for comprehensive research tasks

### 4. User Experience
- Stream responses when possible for better UX
- Show citation count and sources
- Provide clear error messages
- For deep research, warn users about longer processing time

### 5. Integration with Savant
- Store search results in conversation history
- Allow search results to be used as context for other LLMs
- Enable search results to be saved as project files
- Consider model recommendations based on query type

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
   # Test academic mode with reasoning model
   curl -X POST ... --data '{"model": "sonar-reasoning", "search_mode": "academic"}'
   ```

4. **Financial Search**
   ```bash
   # Test with reasoning pro model
   curl -X POST ... --data '{"model": "sonar-reasoning-pro", "search_domain_filter": ["sec.gov"]}'
   ```

5. **Deep Research**
   ```bash
   # Test comprehensive research
   curl -X POST ... --data '{"model": "sonar-deep-research", "reasoning_effort": "high"}'
   ```

6. **Error Handling**
   ```bash
   # Test with invalid API key
   curl -X POST ... --header 'Authorization: Bearer INVALID_KEY'
   ```

## Rate Limits

Be aware of Perplexity's rate limits:
- Requests per minute vary by usage tier (high, medium, low)
- Deep research queries count as multiple requests
- Implement exponential backoff for rate limit errors
- Consider queuing requests during high usage

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
- Deep research query frequency

## Future Enhancements

Consider these potential improvements:
1. Implement response caching
2. Add support for structured output formats (JSON Schema)
3. Enable async operations for deep research
4. Add webhook support for long-running queries
5. Implement batch search operations
6. Add search result ranking/filtering

## Additional Notes

- **Search Modes**: New search context sizes (low, medium, high) affect pricing
- **Citation Tokens**: No longer charged for citation tokens (except sonar-deep-research)
- **Reasoning Output**: Reasoning models include `<think>` sections in responses
- **Async API**: Available for sonar-deep-research at `https://api.perplexity.ai/async/chat/completions`

## Support

For API issues: api@perplexity.ai
For documentation: https://docs.perplexity.ai
For status updates: Check Perplexity's system status page
