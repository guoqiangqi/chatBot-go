package main

import (
	"bytes"
	chatbot "chatbot/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var baseURL = "http://localhost:8080/"

var headers = map[string]string{
	"Content-Type": "application/json",
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
	chatURL := baseURL + "protected"
	headers["Authorization"] = "Bearer " + accessToken
	chatPayload := chatbot.ChatPayload{
		Question: "who are you?",
	}

	chatPayloadBytes, err := json.Marshal(chatPayload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", chatURL,bytes.NewBuffer(chatPayloadBytes))
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

	var chatResponse chatbot.ChatResponse
	err = json.NewDecoder(resp.Body).Decode(&chatResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(chatResponse.Answer)
}
