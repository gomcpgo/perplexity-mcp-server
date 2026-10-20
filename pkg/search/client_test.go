package search

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prasanthmj/perplexity/pkg/types"
)

type captureTransport struct {
	req *http.Request
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"id":"1","choices":[]}`)),
		Request:    req,
	}, nil
}

func callWithCapture(t *testing.T, url string) *http.Request {
	t.Helper()
	transport := &captureTransport{}
	c := NewClient("test-key", time.Second)
	c.httpClient.Transport = transport
	c.baseURL = url

	if _, err := c.callAPI(context.Background(), &types.PerplexityRequest{Model: "sonar"}); err != nil {
		t.Fatalf("callAPI returned error: %v", err)
	}
	if transport.req == nil {
		t.Fatal("expected a request to be sent")
	}
	return transport.req
}

func TestCallAPISetsIntegrationHeaderForPerplexity(t *testing.T) {
	req := callWithCapture(t, baseURL)

	if got := req.Header.Values(integrationHeader); len(got) != 1 || got[0] != integrationHeaderValue {
		t.Errorf("expected %s=%q once, got %v", integrationHeader, integrationHeaderValue, got)
	}
}

func TestCallAPIOmitsIntegrationHeaderForOtherHosts(t *testing.T) {
	for _, url := range []string{
		"https://example.com/chat/completions",
		"https://api.perplexity.ai.example.com/chat/completions",
	} {
		req := callWithCapture(t, url)
		if got := req.Header.Get(integrationHeader); got != "" {
			t.Errorf("%s: expected no %s header, got %q", url, integrationHeader, got)
		}
	}
}
