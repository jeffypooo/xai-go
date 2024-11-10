package main

import (
	"bufio"
	"fmt"
	"github.com/jeffypooo/xai-go/pkg/xai"
	"log"
	"os"
)

func main() {
	if os.Getenv("X_API_KEY") == "" {
		log.Fatal("X_API_KEY environment variable is required")
	}
	client := xai.NewClient(os.Getenv("X_API_KEY"))

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
