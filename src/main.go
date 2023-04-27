package main

import (
	chatbot "chatbot/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	openai "github.com/sashabaranov/go-openai"
)

var jwtSecreteKey = []byte("change your secret key here.")
var configFilePath = "config/config.json"
var configData map[string]string

func initConfig() {
	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&configData)
	if err != nil {
		panic(err)
	}
	return
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse := chatbot.ErrorResponse{
			ErrorMessage: "Authentication failed: cannot find username/password in request body.",
		}

		jsonData, _ := json.Marshal(errorResponse)
		w.Write(jsonData)
		return
	}

	user := chatbot.Authenticate(username, password)

	if user != nil {
		tokenString, _ := chatbot.GenerateToken(user.Id, jwtSecreteKey)
		log.Println("Authenticate successfully: ", user.Name)
		log.Println("Generate token successfully: ", tokenString)

		authResponse := chatbot.AuthResponse{
			AccessToken: tokenString,
		}
		jsonData, _ := json.Marshal(authResponse)
		w.Write(jsonData)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		errorResponse := chatbot.ErrorResponse{
			ErrorMessage: "Authentication failed: cannot authenticate with provided username and password.",
		}

		jsonData, _ := json.Marshal(errorResponse)
		w.Write(jsonData)
	}
}

func chatCompletionHandler(stream bool, source string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Verify the JWT token from http header.
		authHeader := r.Header.Get("Authorization")
		const bearerPrefix = "Bearer "
		if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := chatbot.ErrorResponse{
				ErrorMessage: "Invalid token: empty or not starts with 'Bearer '",
			}

			jsonData, _ := json.Marshal(errorResponse)
			w.Write(jsonData)
			return
		}

		tokenString := authHeader[len(bearerPrefix):]
		claims, err := chatbot.ParseToken(tokenString, jwtSecreteKey)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			errorResponse := chatbot.ErrorResponse{
				ErrorMessage: "Invalid token: expired or fake token.",
			}

			jsonData, _ := json.Marshal(errorResponse)
			w.Write(jsonData)
			return
		}

		// username := chatbot.IndexUserWithID(claims.Id).Name
		username := claims.StandardClaims.Subject
		log.Println("Hello, ", username)

		// handling chat and responsing answer of question.
		var chatPayload []openai.ChatCompletionMessage
		systemPayload := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: configData[source],
			},
		}
		err = json.NewDecoder(r.Body).Decode(&chatPayload)
		chatPayload = append(systemPayload, chatPayload...)
		// log.Println(chatPayload.Question)
		question := chatPayload[len(chatPayload)-1].Content
		answer := ""
		dbTableName := source + "_QA"

		if stream {
			// completionStream, err := chatbot.ChatCompletionStream(chatPayload, openai.GPT3Dot5Turbo)
			chatWorkFunc := chatbot.ExponentialBackOff(&chatbot.ChatCompletionStreamFunc{}, 1.0, 2.0, 1.0, 3)
			resp, err := chatWorkFunc(chatPayload, openai.GPT3Dot5Turbo)
			completionStream := resp.(*openai.ChatCompletionStream)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorResponse := chatbot.ErrorResponse{
					ErrorMessage: fmt.Sprintf("Failed with chatbot.ChatCompletionStream: %s", err),
				}

				jsonData, _ := json.Marshal(errorResponse)
				w.Write(jsonData)
				return
			}

			defer completionStream.Close()

			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")

			for {
				response, err := completionStream.Recv()

				if errors.Is(err, io.EOF) {
					log.Println("Stream finished")

					go chatbot.WriteQAToDB(question, answer, dbTableName)

					return
				}

				// log.Println(response.Choices[0].Delta.Content)
				answer += response.Choices[0].Delta.Content
				fmt.Fprintf(w, "event: message\ndata: %q\n\n", response.Choices[0].Delta.Content)
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}

		} else {
			// Rate limit for free account to use gpt-3.5-turbo is 20 per min,
			// set a exponential backoff here instead of original request to avoid reaching the limit:
			//
			// chatResponse, err := chatbot.ChatCompletion(chatPayload, openai.GPT3Dot5Turbo)
			//
			chatWorkFunc := chatbot.ExponentialBackOff(&chatbot.ChatCompletionFunc{}, 1.0, 2.0, 1.0, 3)
			chatResponse, err := chatWorkFunc(chatPayload, openai.GPT3Dot5Turbo)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				errorResponse := chatbot.ErrorResponse{
					ErrorMessage: fmt.Sprintf("Failed with chatbot.ChatCompletion: %s", err),
				}

				jsonData, _ := json.Marshal(errorResponse)
				w.Write(jsonData)
				return
			}

			res := chatResponse.(openai.ChatCompletionResponse)
			answer = res.Choices[0].Message.Content

			var answerPayload chatbot.Answer
			answerPayload.Message = answer
			jsonData, _ := json.Marshal(answerPayload)
			w.Write(jsonData)
		}

		go chatbot.WriteQAToDB(question, answer, dbTableName)
	}
}

func main() {
	initConfig()

	router := mux.NewRouter()

	router.HandleFunc("/auth", authHandler).Methods("POST")
	router.HandleFunc("/chatCompletion", chatCompletionHandler(false, "openEuler")).Methods("POST")
	router.HandleFunc("/chatCompletionStream", chatCompletionHandler(true, "openEuler")).Methods("POST")

	router.HandleFunc("/chatCompletion-compass", chatCompletionHandler(false, "Compass")).Methods("POST")
	router.HandleFunc("/chatCompletionStream-compass", chatCompletionHandler(true, "Compass")).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
