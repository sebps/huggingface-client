package client

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type mockRoundTripper struct {
	mockRespond func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockRespond(req), nil
}

// newTestClient creates a new Client with a fake HTTP client
func newTestClient(mockRespond func(req *http.Request) *http.Response) *Client {
	httpClient := &http.Client{
		Transport: &mockRoundTripper{mockRespond},
	}
	return &Client{
		Host:   "https://fake.api",
		Token:  "fake-token",
		Client: httpClient,
	}
}

func TestListEndpoints(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"items":[{"name":"test-endpoint"}]}`)),
		}
	})

	endpoints, err := client.ListEndpoints("namespace", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(endpoints) != 1 || endpoints[0].Name != "test-endpoint" {
		t.Fatalf("unexpected result: %+v", endpoints)
	}
}

func TestCreateEndpoint(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"name":"created-endpoint"}`)),
		}
	})

	newEndpoint, err := client.CreateEndpoint("namespace", Endpoint{Name: "new"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newEndpoint.Name != "created-endpoint" {
		t.Fatalf("unexpected result: %+v", newEndpoint)
	}
}

func TestGetEndpoint(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"name":"fetched-endpoint"}`)),
		}
	})

	endpoint, err := client.GetEndpoint("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if endpoint.Name != "fetched-endpoint" {
		t.Fatalf("unexpected result: %+v", endpoint)
	}
}

func TestUpdateEndpoint(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"name":"updated-endpoint"}`)),
		}
	})

	updated, err := client.UpdateEndpoint("namespace", "endpoint", EndpointUpdate{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Name != "updated-endpoint" {
		t.Fatalf("unexpected result: %+v", updated)
	}
}

func TestDeleteEndpoint(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 204,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
	})

	err := client.DeleteEndpoint("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetEndpointLogs(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("logs-data")),
		}
	})

	data, err := client.GetEndpointLogs("namespace", "endpoint", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != "logs-data" {
		t.Fatalf("unexpected logs: %s", string(data))
	}
}

func TestStreamEndpointLogs(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("stream-data")),
		}
	})

	stream, err := client.StreamEndpointLogs("namespace", "endpoint", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer stream.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, stream)
	if err != nil {
		t.Fatalf("error while copying stream: %v", err)
	}

	if buf.String() != "stream-data" {
		t.Fatalf("unexpected stream: %s", buf.String())
	}
}

func TestGetEndpointMetrics(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("metrics-data")),
		}
	})

	data, err := client.GetEndpointMetrics("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != "metrics-data" {
		t.Fatalf("unexpected metrics: %s", string(data))
	}
}

func TestGetEndpointMetric(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("metric-data")),
		}
	})

	data, err := client.GetEndpointMetric("namespace", "endpoint", "hardwareUsage")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != "metric-data" {
		t.Fatalf("unexpected metric: %s", string(data))
	}
}

func TestPauseEndpoint(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
	})

	err := client.PauseEndpoint("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetEndpointReplicasStatuses(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("replicas-status")),
		}
	})

	data, err := client.GetEndpointReplicasStatuses("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != "replicas-status" {
		t.Fatalf("unexpected replicas: %s", string(data))
	}
}

func TestResumeEndpoint(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
	})

	err := client.ResumeEndpoint("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestScaleEndpointToZero(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
	})

	err := client.ScaleEndpointToZero("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetEndpointSSE(t *testing.T) {
	client := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("sse-data")),
		}
	})

	stream, err := client.GetEndpointSSE("namespace", "endpoint")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer stream.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, stream)
	if err != nil {
		t.Fatalf("error while copying stream: %v", err)
	}

	if buf.String() != "sse-data" {
		t.Fatalf("unexpected SSE: %s", buf.String())
	}
}
