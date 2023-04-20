package chatbot

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func ChatCompletion(messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error) {
	apiKeyString := os.Getenv("OPENAI_API_KEY")
	apiKeyList := strings.Split(apiKeyString, "||")

	var resp openai.ChatCompletionResponse
	var err error

	for len(apiKeyList) != 0 {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(apiKeyList))
		apiKey := apiKeyList[index]
		apiKeyList = append(apiKeyList[:index], apiKeyList[index+1:]...)

		log.Println("Request with API token: ", apiKey)
		client := openai.NewClient(apiKey)

		resp, err = client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    model,
				Messages: messages,
			},
		)

		if err == nil {
			return resp, nil
		}
		log.Println("Request failed: ", err)
	}
	log.Println("ChatCompletion error: ", err)
	return resp, err
}

func ChatCompletionStream(messages []openai.ChatCompletionMessage, model string) (*openai.ChatCompletionStream, error) {
	apiKeyString := os.Getenv("OPENAI_API_KEY")
	apiKeyList := strings.Split(apiKeyString, "||")

	var stream *openai.ChatCompletionStream
	var err error

	for len(apiKeyList) != 0 {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(apiKeyList))
		apiKey := apiKeyList[index]
		apiKeyList = append(apiKeyList[:index], apiKeyList[index+1:]...)

		log.Println("Request with API token: ", apiKey)
		client := openai.NewClient(apiKey)

		stream, err = client.CreateChatCompletionStream(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    model,
				Messages: messages,
				Stream:   true,
			},
		)

		if err == nil {
			return stream, nil
		}
		log.Println("Request failed: ", err)
	}

	log.Println("ChatCompletionStream error: ", err)
	return stream, err
}
