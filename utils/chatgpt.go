package chatbot

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func ChatCompletion(messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)

	if err != nil {
		log.Println("ChatCompletion error: ", err)
		return resp, err
	}

	return resp, nil
}

func ChatCompletionStream(messages []openai.ChatCompletionMessage, model string) (*openai.ChatCompletionStream, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
			Stream:   true,
		},
	)
	if err != nil {
		log.Println("ChatCompletionStream error: ", err)
		return stream, err
	}
	// defer stream.Close()

	return stream, nil
}
