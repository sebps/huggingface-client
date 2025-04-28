package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Default host url is huggingface API url
const HostURL string = "https://api.endpoints.huggingface.cloud"

type Client struct {
	Host   string
	Token  string
	Client *http.Client
}

func NewClient(host, token *string) (*Client, error) {
	c := Client{
		Client: &http.Client{Timeout: 10 * time.Second},
		// Default Huggingface Endpoints URL
		Host: HostURL,
	}

	if host != nil && *host != "" {
		c.Host = *host
	}

	// If token not provided, return empty client
	if token == nil {
		return &c, nil
	}

	c.Token = *token

	return &c, nil
}

// doRequest - for normal HTTP APIs (returns whole response body)
func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	if c.Client == nil {
		c.Client = &http.Client{Timeout: 30 * time.Second}
	}

	if authToken != nil {
		req.Header.Set("Authorization", "Bearer "+*authToken)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP errors (status >= 400)
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// doStreamRequest - for streaming APIs (returns a reader directly)
func (c *Client) doStreamRequest(req *http.Request, authToken *string) (io.ReadCloser, error) {
	if c.Client == nil {
		c.Client = &http.Client{Timeout: 0} // No timeout for streams
	}

	if authToken != nil {
		req.Header.Set("Authorization", "Bearer "+*authToken)
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Important: caller must close resp.Body when done!
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp.Body, nil
}
