package chatbot

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net/http"
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
	var statusCode int

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

		//FIXME01: can not catch 'exceeded quota error' here, is a issue from upstream client.CreateChatCompletionStream api
		// but not with client.CreateChatCompletion.
		if err != nil {
			log.Println("Request failed: ", err)
			continue
		}
		// defer stream.Close()

		res, _ := stream.GetHttpResponse()
		statusCode = res.StatusCode
		// Rate limit and quota exceeded would get http.StatusTooManyRequests. Incorrect API key get http.StatusUnauthorized
		if statusCode == http.StatusTooManyRequests || statusCode == http.StatusRequestTimeout || statusCode >= http.StatusInternalServerError {
			err = errors.New("get bad response, please retry you http request")
			log.Println("Request failed: ", err)
			continue

		} else if statusCode == http.StatusUnauthorized {
			err = errors.New("get bad response, may cased by incorrect API key")
			log.Println("Request failed: ", err)
			continue
		}
		return stream, nil
	}

	log.Println("ChatCompletionStream error: ", err)
	return stream, err
}
