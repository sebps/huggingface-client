package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/sebps/huggingface-client/client"
)

func BuildInferenceImage(
	imageType string,
	imageUrl string,
	imagePort int,
	modelPath string,
	username string,
	password string,
) (*client.EndpointModelImage, error) {
	image := &client.EndpointModelImage{}

	switch imageType {
	case "huggingface":
		image.HuggingFace = &client.HuggingFaceImage{}
	case "huggingfaceNeuron":
		image.HuggingFaceNeuron = &client.HuggingFaceNeuronImage{}
	case "tgi":
		if imageUrl == "" {
			return nil, fmt.Errorf("url is required for tgi image")
		}
		image.TGI = &client.TGIImage{
			URL:                  imageUrl,
			Port:                 imagePort,
			DisableCustomKernels: true,
		}
	case "tgiNeuron":
		if imageUrl == "" {
			return nil, fmt.Errorf("url is required for tgi image")
		}
		image.TGINeuron = &client.TGINeuronImage{
			URL:  imageUrl,
			Port: imagePort,
		}
	case "tei":
		if imageUrl == "" {
			return nil, fmt.Errorf("url is required for tei image")
		}
		image.TEI = &client.TEIImage{
			URL:  imageUrl,
			Port: imagePort,
		}
	case "llamacpp":
		if imageUrl == "" || modelPath == "" {
			return nil, fmt.Errorf("url and model path are required for llamacpp image")
		}
		image.LlamaCpp = &client.LlamaCppImage{
			URL:       imageUrl,
			Port:      imagePort,
			CtxSize:   4096,
			ModelPath: modelPath,
			NParallel: 1,
		}
	case "custom":
		if imageUrl == "" {
			return nil, fmt.Errorf("url is required for custom image")
		}
		image.Custom = &client.CustomImage{
			URL:  imageUrl,
			Port: imagePort,
			Credentials: &client.Credentials{
				Username: username,
				Password: &password,
			},
		}
	default:
		return nil, fmt.Errorf("unsupported image type: %s", imageType)
	}

	return image, nil
}

// parseSmartTime tries to parse a time string flexibly
func ParseTime(input string) (time.Time, error) {
	input = strings.TrimSpace(input)

	// First try known time formats
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.ANSIC,
		"2006-01-02 15:04:05", // common ISO8601 without timezone
		"2006-01-02",          // date only
	}

	for _, format := range formats {
		t, err := time.Parse(format, input)
		if err == nil {
			return t.UTC(), nil
		}
	}

	// Try parsing as a duration like "-1h", "-24h"
	if d, err := time.ParseDuration(input); err == nil {
		return time.Now().Add(d), nil
	}

	return time.Time{}, fmt.Errorf("could not parse time: %s", input)
}

func IsMetricValid(metric string) bool {
	switch metric {
	case "pending-requests", "request-count", "median-latency", "p95-latency", "success-throughput",
		"bad-request-throughput", "server-error-throughput", "cpu-usage", "memory-usage", "gpu-usage",
		"gpu-memory-usage", "neuron-usage", "neuron-memory-usage", "ready-replicas", "running-replicas",
		"target-replicas", "average-latency", "success-rate", "bad-request-rate", "server-error-rate":
		return true
	}

	return false
}
