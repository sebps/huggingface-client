package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ListEndpoints - List all endpoints in a namespace
func (c *Client) ListEndpoints(namespace string, tags *string) ([]EndpointWithStatus, error) {
	url := fmt.Sprintf("%s/v2/endpoint/%s", c.Host, namespace)
	if tags != nil {
		url += fmt.Sprintf("?tags=%s", *tags)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []EndpointWithStatus `json:"items"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// CreateEndpoint - Create a new endpoint
func (c *Client) CreateEndpoint(namespace string, endpoint Endpoint) (*EndpointWithStatus, error) {
	rb, err := json.Marshal(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/endpoint/%s", c.Host, namespace), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var newEndpoint EndpointWithStatus
	err = json.Unmarshal(body, &newEndpoint)
	if err != nil {
		return nil, err
	}

	return &newEndpoint, nil
}

// GetEndpoint - Get endpoint information
func (c *Client) GetEndpoint(namespace, name string) (*EndpointWithStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/endpoint/%s/%s", c.Host, namespace, name), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var endpoint EndpointWithStatus
	err = json.Unmarshal(body, &endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

// UpdateEndpoint - Update an endpoint
func (c *Client) UpdateEndpoint(namespace, name string, endpointUpdate EndpointUpdate) (*EndpointWithStatus, error) {
	rb, err := json.Marshal(endpointUpdate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/endpoint/%s/%s", c.Host, namespace, name), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, &c.Token)
	if err != nil {
		return nil, err
	}

	var updatedEndpoint EndpointWithStatus
	err = json.Unmarshal(body, &updatedEndpoint)
	if err != nil {
		return nil, err
	}

	return &updatedEndpoint, nil
}

// DeleteEndpoint - Delete an endpoint
func (c *Client) DeleteEndpoint(namespace, name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/endpoint/%s/%s", c.Host, namespace, name), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, &c.Token)
	return err
}

// GetEndpointLogs - Get logs from an endpoint (optionally filter by replica ID)
func (c *Client) GetEndpointLogs(namespace, name string, replicaID *string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/endpoint/%s/%s/logs", c.Host, namespace, name)
	if replicaID != nil {
		url += fmt.Sprintf("?replica=%s", *replicaID)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, &c.Token)
}

// StreamEndpointLogs - Stream logs from an endpoint using SSE (optionally filter by replica ID)
func (c *Client) StreamEndpointLogs(namespace, name string, replicaID *string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/v2/endpoint/%s/%s/logs/sse", c.Host, namespace, name)
	if replicaID != nil {
		url += fmt.Sprintf("?replica=%s", *replicaID)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.doStreamRequest(req, &c.Token)
}

// GetEndpointMetrics - Get all metrics for an endpoint (plural version)
func (c *Client) GetEndpointMetrics(namespace, name string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/endpoint/%s/%s/metrics", c.Host, namespace, name), nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, &c.Token)
}

// GetEndpointMetric - Get metrics from an endpoint
func (c *Client) GetEndpointMetric(namespace, name, metricType string) ([]byte, error) {
	url := fmt.Sprintf("%s/v2/endpoint/%s/%s/metrics/%s", c.Host, namespace, name, metricType)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, &c.Token)
}

// PauseEndpoint - Pause a running endpoint
func (c *Client) PauseEndpoint(namespace, name string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/endpoint/%s/%s/pause", c.Host, namespace, name), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, &c.Token)
	return err
}

// GetEndpointReplicasStatuses - Get status of all endpoint replicas
func (c *Client) GetEndpointReplicasStatuses(namespace, name string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/endpoint/%s/%s/replicas", c.Host, namespace, name), nil)
	if err != nil {
		return nil, err
	}

	return c.doRequest(req, &c.Token)
}

// ResumeEndpoint - Resume a paused endpoint
func (c *Client) ResumeEndpoint(namespace, name string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/endpoint/%s/%s/resume", c.Host, namespace, name), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, &c.Token)
	return err
}

// ScaleEndpointToZero - Scale an endpoint down to zero replicas
func (c *Client) ScaleEndpointToZero(namespace, name string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/endpoint/%s/%s/scale-to-zero", c.Host, namespace, name), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, &c.Token)
	return err
}

// GetEndpointSSE - Stream server-sent events for endpoint
func (c *Client) GetEndpointSSE(namespace, name string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/endpoint/%s/%s/sse", c.Host, namespace, name), nil)
	if err != nil {
		return nil, err
	}

	return c.doStreamRequest(req, &c.Token)
}
