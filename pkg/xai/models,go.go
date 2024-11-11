package xai

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
