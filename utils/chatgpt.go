package chatbot

import (
	"context"
	"errors"
	"log"
	"net/http"
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
	//FIXME01: can not catch 'exceeded quota error' here, is a issue from upstream client.CreateChatCompletionStream api
	// but not with client.CreateChatCompletion.
	if err != nil {
		log.Println("ChatCompletionStream error: ", err)
		return stream, err
	}
	// defer stream.Close()

	res, _ := stream.GetHttpResponse()
	statusCode := res.StatusCode
	// Rate limit and quota exceeded would get http.StatusTooManyRequests. Incorrect API key get http.StatusUnauthorized
	if statusCode == http.StatusTooManyRequests || statusCode == http.StatusRequestTimeout || statusCode >= http.StatusInternalServerError {
		err = errors.New("get bad response, please retry you http request")

	} else if statusCode == http.StatusUnauthorized {
		err = errors.New("get bad response, may cased by incorrected API key")
	}

	log.Println("ChatCompletionStream error with status code: ", statusCode)
	return stream, err
}
