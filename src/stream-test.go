package main

import (
	"bufio"
	"bytes"
	chatbot "chatbot/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

var baseURL = "http://localhost:8080/"

var headers = map[string]string{
	"Content-Type": "application/json",
	"Accept":       "text/event-stream",
}

func main() {

	// login with username and password, get token from response
	authURL := baseURL + "auth"

	authPayload := url.Values{
		"username": {"temporary_user"},
		"password": {"default_password"},
	}

	resp, err := http.PostForm(authURL, authPayload)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var authResponse chatbot.ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&authResponse)
		fmt.Println(authResponse.ErrorMessage)
		return
	}
	var authResponse chatbot.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	accessToken := authResponse.AccessToken
	fmt.Println(accessToken)

	// request to chatgpt with token
	chatCompletionStreamURL := baseURL + "chatCompletionStream"
	headers["Authorization"] = "Bearer " + accessToken

	// chatPayload type should be  []openai.ChatCompletionMessage
	chatPayload := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "What is your name.",
		},
	}
	chatPayloadBytes, _ := json.Marshal(chatPayload)

	req, err := http.NewRequest("POST", chatCompletionStreamURL, bytes.NewBuffer(chatPayloadBytes))
	if err != nil {
		fmt.Println(err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var chatResponse chatbot.ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&chatResponse)
		fmt.Println(chatResponse.ErrorMessage)
		return
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Finish stream.")
			return
		}
		line = strings.TrimSuffix(line, "\n")
		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			fmt.Println(data)
		}
	}

}
