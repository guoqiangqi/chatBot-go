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
	DoChatWork(messages []openai.ChatCompletionMessage, model string) (interface{}, error)
}

type ChatCompletionFunc struct{}
type ChatCompletionStreamFunc struct{}

func (ccf *ChatCompletionFunc) DoChatWork(messages []openai.ChatCompletionMessage, model string) (interface{}, error) {
	return ChatCompletion(messages, model)
}

func (ccsf *ChatCompletionStreamFunc) DoChatWork(messages []openai.ChatCompletionMessage, model string) (interface{}, error) {
	return ChatCompletionStream(messages, model)
}

// ExponentialBackOff is a function that implements exponential backoff algorithm.
// It returns an error when the maximum number of retries is reached.
func ExponentialBackOff(cwf ChatWorkFunc, initialDelay float64, exponentialBase float64, jitter float64, maxRetries int /*, errors []error*/) func(messages []openai.ChatCompletionMessage, model string) (interface{}, interface{}) {
	return func(messages []openai.ChatCompletionMessage, model string) (interface{}, interface{}) {
		delay := initialDelay
		var res interface{}
		for i := 0; i <= maxRetries; i++ {
			if i != 0 {
				log.Println("ExponentialBackOff Retry: ", i)
				rand.Seed(time.Now().UnixNano())
				delay *= exponentialBase * (1.0 + jitter*rand.Float64())
				log.Println("Sleeping time: ", delay)
				time.Sleep(time.Duration(delay) * time.Second)
			}
			// Using an intermediate variable triggers automatic type conversion.
			// Directly using 'res' to receive the return value of a function will cause a type mismatch.
			resp, err := cwf.DoChatWork(messages, model)
			res = resp
			if err == nil {
				return res, nil
			}

			e := &openai.APIError{}
			if errors.As(err, &e) {
				switch e.HTTPStatusCode {
				case 401:
					// invalid auth or key (do not retry)
					// do nothing, go to return
				case 429:
					// rate limiting or engine overload (wait and retry)
					continue
				case 500:
					// openai server error (retry)
					continue
				}
			}

			// return res, err
			log.Println("ExponentialBackOff faiiled: ", err)
			return res, errors.New("something wrong with the internal service")
		}
		return res, &MaxRetryError{"ExponentialBackOff faiiled: Maximum number of retries exceeded."}
	}
}
