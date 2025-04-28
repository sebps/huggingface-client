package utils

import (
	"fmt"

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
