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

func getRandomKey(keyList []string) ([]string, string) {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(keyList))
	apiKey := keyList[index]
	keyList = append(keyList[:index], keyList[index+1:]...)
	return keyList, apiKey
}

func ChatCompletion(messages []openai.ChatCompletionMessage, model string) (openai.ChatCompletionResponse, error) {

	var resp openai.ChatCompletionResponse
	var err error
	var apiKey string
	var apiKeyList = strings.Split(os.Getenv("OPENAI_API_KEY"), "||")

	for len(apiKeyList) != 0 {
		apiKeyList, apiKey = getRandomKey(apiKeyList)

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

	var stream *openai.ChatCompletionStream
	var err error
	var apiKey string
	var apiKeyList = strings.Split(os.Getenv("OPENAI_API_KEY"), "||")

	for len(apiKeyList) != 0 {
		apiKeyList, apiKey = getRandomKey(apiKeyList)

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
