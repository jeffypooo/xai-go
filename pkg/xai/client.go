package xai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	ApiKey string
	client *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
		client: &http.Client{},
	}
}

func (x *Client) GetEmbeddingModels() ([]*Model, error) {
	resp, err := x.client.Do(xAIGet(x.ApiKey, "embedding-models"))
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

func (x *Client) GetLanguageModels() ([]*Model, error) {
	resp, err := x.client.Do(xAIGet(x.ApiKey, "language-models"))
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

func (x *Client) GetModels() ([]*Model, error) {
	resp, err := x.client.Do(xAIGet(x.ApiKey, "models"))
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

func (x *Client) GetChatCompletion(
	messages []*ChatMessage,
) (*ChatCompletionResponse, error) {
	completionRequest := NewChatCompletionRequest(
		"grok-beta",
		messages,
	)
	completionRequestData, err := json.Marshal(completionRequest)
	if err != nil {
		return nil, err
	}
	//prettyPrintRequest(completionRequestData)
	resp, err := x.client.Do(xAIPost(x.ApiKey, "chat/completions", completionRequestData))
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

type ListSpecificModelsResponse struct {
	Models []*Model `json:"models"`
}

type ListModelsResponse struct {
	Object string   `json:"object"`
	Data   []*Model `json:"data"`
}

type Model struct {
	Created               int64    `json:"created"`
	ID                    string   `json:"id"`
	InputModalities       []string `json:"input_modalities"`
	Object                string   `json:"object"`
	OwnedBy               string   `json:"owned_by"`
	PromptImageTokenPrice int      `json:"prompt_image_token_price"`
	PromptTextTokenPrice  int      `json:"prompt_text_token_price"`
	Version               string   `json:"version"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	FrequencyPenalty int                `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	Logprobs         bool               `json:"logprobs,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	Messages         []*ChatMessage     `json:"messages"`
	Model            string             `json:"model"`
	N                int                `json:"n,omitempty"`
	PresencePenalty  int                `json:"presence_penalty,omitempty"`
	ResponseFormat   string             `json:"response_format,omitempty"`
	Seed             int                `json:"seed,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	Stream           bool               `json:"stream,omitempty"`
	StreamOptions    string             `json:"stream_options,omitempty"`
}

func NewChatCompletionRequest(
	model string,
	messages []*ChatMessage,
) *ChatCompletionRequest {
	return &ChatCompletionRequest{
		MaxTokens: 500,
		Messages:  messages,
		Model:     model,
	}
}

type ChatChoice struct {
	FinishReason string       `json:"finish_reason"`
	Index        int          `json:"index"`
	Message      *ChatMessage `json:"message"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	Choices []*ChatChoice `json:"choices"`
	Created int64         `json:"created"`
	ID      string        `json:"id"`
	Model   string        `json:"model"`
	Object  string        `json:"object"`
	Usage   *Usage        `json:"usage"`
}

func xAIGet(
	apiKey string,
	endpoint string,
) *http.Request {
	return &http.Request{
		Method: http.MethodGet,
		URL:    xAIRequestURL(endpoint),
		Header: xAIRequestHeaders(apiKey),
	}
}

func xAIPost(
	apiKey string,
	endpoint string,
	body []byte,
) *http.Request {
	return &http.Request{
		Method: http.MethodPost,
		URL:    xAIRequestURL(endpoint),
		Header: xAIRequestHeaders(apiKey),
		Body:   io.NopCloser(bytes.NewReader(body)),
	}
}

func xAIRequestURL(endpoint string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "api.x.ai",
		Path:   "/v1/" + endpoint,
	}
}

func xAIRequestHeaders(apiKey string) http.Header {
	return http.Header{
		"Authorization": []string{"Bearer " + apiKey},
		"Content-Type":  []string{"application/json"},
	}
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
