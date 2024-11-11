package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jeffypooo/xai-go/pkg/xai"
	"log"
	"os"
)

const APIKeyEnvVarName = "XAI_KEY"

func main() {
	apiKey := os.Getenv(APIKeyEnvVarName)
	model := "grok-beta"

	flag.StringVar(&apiKey, "key", apiKey, "API key for xAI, defaults to XAI_KEY environment variable")
	flag.StringVar(&model, "model", model, "Model to use for chat, defaults to 'grok-beta'")
	flag.Parse()

	if apiKey == "" {
		panic("API key is required. Pass via -key flag or set XAI_KEY environment variable")
	}

	client := xai.NewClient(apiKey, model)

	// system prompt
	systemMessage := &xai.ChatMessage{
		Role:    "system",
		Content: "You are an unhinged lunatic, yet also extremely intelligent and kind. You cannot help but swear profusely when speaking.",
	}
	// keep running list of chat messages
	chatMessages := []*xai.ChatMessage{
		systemMessage,
	}

	// continuously loop and allow user to enter messages
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter your message: ")
		if !scanner.Scan() {
			break
		}
		userMessage := scanner.Text()
		if userMessage == "" {
			continue
		}
		chatMessages = append(
			chatMessages, &xai.ChatMessage{
				Role:    "user",
				Content: userMessage,
			},
		)
		chatCompletionResponse, err := client.GetChatCompletion(chatMessages)
		if err != nil {
			log.Fatal(err)
		}
		reply := chatCompletionResponse.Choices[0].Message
		chatMessages = append(chatMessages, reply)
		println("SYSTEM: " + reply.Content)
	}
}
