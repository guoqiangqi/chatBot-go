package chatbot

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type MaxRetryError struct {
	message string
}

func (e *MaxRetryError) Error() string {
	return fmt.Sprintf("Error: %s", e.message)
}

type ChatWorkFunc interface {
	DoChatWork(messages []openai.ChatCompletionMessage, model string) (interface{}, interface{})
}

type ChatCompletionFunc struct{}
type ChatCompletionStreamFunc struct{}

func (ccf *ChatCompletionFunc) DoChatWork(messages []openai.ChatCompletionMessage, model string) (interface{}, interface{}) {
	return ChatCompletion(messages, model)
}

func (ccsf *ChatCompletionStreamFunc) DoChatWork(messages []openai.ChatCompletionMessage, model string) (interface{}, interface{}) {
	return ChatCompletionStream(messages, model)
}

// ExponentialBackOff is a function that implements exponential backoff algorithm.
// It returns an error when the maximum number of retries is reached.
func ExponentialBackOff(cwf ChatWorkFunc, initialDelay float64, exponentialBase float64, jitter float64, maxRetries int /*, errors []error*/) func(messages []openai.ChatCompletionMessage, model string) (interface{}, interface{}) {
	return func(messages []openai.ChatCompletionMessage, model string) (interface{}, interface{}) {

		res, err := cwf.DoChatWork(messages, model)
		if err == nil {
			return res, nil
		}

		delay := initialDelay
		for i := 0; i < maxRetries; i++ {

			log.Println("ExponentialBackOff Retry: ", i)
			rand.Seed(time.Now().UnixNano())
			delay *= exponentialBase * (1.0 + jitter*rand.Float64())
			log.Println("Failed, sleeping time: ", delay)
			time.Sleep(time.Duration(delay) * time.Second)

			res, err = cwf.DoChatWork(messages, model)
			if err == nil {
				return res, nil
			}

			/*
				// TODO: Optimize error checking here:
				// Due to the lack of unified error management for openai RESTful API in the go-openai,
				// error checking is temporarily disabled here.

				isProvidedError := false
				for _, value := range errors {
					if value == err {
						isProvidedError = true
						break
					}
				}
			*/
			isProvidedError := true
			if !isProvidedError {
				// return res, err
				log.Println("ExponentialBackOff faiiled: ", err)
				return res, errors.New("something wrong with the internal service")
			}
		}
		return res, &MaxRetryError{"ExponentialBackOff faiiled: Maximum number of retries exceeded."}
	}
}
