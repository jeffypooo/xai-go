package xai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client interface {
	GetEmbeddingModels() ([]*Model, error)
	GetLanguageModels() ([]*Model, error)
	GetModels() ([]*Model, error)
	GetChatCompletion(messages []*ChatMessage) (*ChatCompletionResponse, error)
}

type DefaultClient struct {
	ApiKey string
	Model  string
	client *http.Client
}

func NewClient(
	key string,
	model string,
) Client {
	return &DefaultClient{
		ApiKey: key,
		Model:  model,
		client: &http.Client{},
	}
}

func (x *DefaultClient) GetEmbeddingModels() ([]*Model, error) {
	resp, err := x.client.Do(get(x.ApiKey, "embedding-models"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	//prettyPrintResponse(data)
	if err != nil {
		return nil, err
	}
	var modelsResp ListSpecificModelsResponse
	err = json.Unmarshal(data, &modelsResp)
	if err != nil {
		return nil, err
	}
	return modelsResp.Models, nil
}

func (x *DefaultClient) GetLanguageModels() ([]*Model, error) {
	resp, err := x.client.Do(get(x.ApiKey, "language-models"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	//prettyPrintResponse(data)
	if err != nil {
		return nil, err
	}
	var modelsResp ListSpecificModelsResponse
	err = json.Unmarshal(data, &modelsResp)
	if err != nil {
		return nil, err
	}
	return modelsResp.Models, nil
}

func (x *DefaultClient) GetModels() ([]*Model, error) {
	resp, err := x.client.Do(get(x.ApiKey, "models"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	//prettyPrintResponse(data)
	if err != nil {
		return nil, err
	}
	var modelsResp ListModelsResponse
	err = json.Unmarshal(data, &modelsResp)
	if err != nil {
		return nil, err
	}
	return modelsResp.Data, nil
}

func (x *DefaultClient) GetChatCompletion(
	messages []*ChatMessage,
) (*ChatCompletionResponse, error) {
	completionRequest := NewChatCompletionRequest(x.Model, messages)
	completionRequestData, err := json.Marshal(completionRequest)
	if err != nil {
		return nil, err
	}
	//prettyPrintRequest(completionRequestData)
	resp, err := x.client.Do(post(x.ApiKey, "chat/completions", completionRequestData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	//prettyPrintResponse(data)
	if err != nil {
		return nil, err
	}
	var completionResp ChatCompletionResponse
	err = json.Unmarshal(data, &completionResp)
	if err != nil {
		return nil, err
	}
	return &completionResp, nil
}

func requestURL(endpoint string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "api.x.ai",
		Path:   "/v1/" + endpoint,
	}
}

func requestHeaders(apiKey string) http.Header {
	return http.Header{
		"Authorization": []string{"Bearer " + apiKey},
		"Content-Type":  []string{"application/json"},
	}
}

func apiRequest(
	apiKey string,
	method string,
	endpoint string,
	body *[]byte,
) *http.Request {
	if body == nil {
		return &http.Request{
			Method: method,
			URL:    requestURL(endpoint),
			Header: requestHeaders(apiKey),
		}
	} else {
		return &http.Request{
			Method: method,
			URL:    requestURL(endpoint),
			Header: requestHeaders(apiKey),
			Body:   io.NopCloser(bytes.NewReader(*body)),
		}
	}
}

func get(
	apiKey string,
	endpoint string,
) *http.Request {
	return apiRequest(apiKey, http.MethodGet, endpoint, nil)
}

func post(
	apiKey string,
	endpoint string,
	body []byte,
) *http.Request {
	return apiRequest(apiKey, http.MethodPost, endpoint, &body)
}

func prettyPrintResponse(data []byte) {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, data, "", "    ")
	if err != nil {
		log.Printf("Error during pretty-print: %v", err)
	}
	log.Printf("<< %v", pretty.String())
}

func prettyPrintRequest(data []byte) {
	var pretty bytes.Buffer
	err := json.Indent(&pretty, data, "", "    ")
	if err != nil {
		log.Printf("Error during pretty-print: %v", err)
	}
	log.Printf(">> %v", pretty.String())
}
