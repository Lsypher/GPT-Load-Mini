package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OpenAIChannel struct {
	upstreamURL string
	testModel   string
	httpClient  *http.Client
}

func NewOpenAIChannel(upstreamURL, testModel string, client *http.Client) *OpenAIChannel {
	return &OpenAIChannel{
		upstreamURL: strings.TrimSuffix(upstreamURL, "/"),
		testModel:   testModel,
		httpClient:  client,
	}
}

func (c *OpenAIChannel) BuildUpstreamURL(path string) string {
	return fmt.Sprintf("%s%s", c.upstreamURL, path)
}

func (c *OpenAIChannel) ModifyRequest(req *http.Request, apiKey string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
}

func (c *OpenAIChannel) IsStreamRequest(req *http.Request, body []byte) bool {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err == nil {
		if stream, ok := payload["stream"].(bool); ok && stream {
			return true
		}
	}
	return false
}

func (c *OpenAIChannel) ExtractModel(body []byte) string {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err == nil {
		if model, ok := payload["model"].(string); ok {
			return model
		}
	}
	return ""
}

func (c *OpenAIChannel) GetHTTPClient() *http.Client {
	return c.httpClient
}

func (c *OpenAIChannel) GetStreamClient() *http.Client {
	return c.httpClient
}

func (c *OpenAIChannel) HandleResponse(resp *http.Response, isStream bool) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
