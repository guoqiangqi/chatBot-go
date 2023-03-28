package main

import (
	chatbot "chatbot/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var jwtSecreteKey = []byte("change your secret key here.")

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
		fmt.Println("Authenticate successfully: ", user.Name)
		fmt.Println("Generate token successfully: ", tokenString)

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

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	fmt.Println("Hello, ", username)

	chatResponse := chatbot.ChatResponse{
		Answer: "Receive your question, handling...",
	}

	jsonData, _ := json.Marshal(chatResponse)
	w.Write(jsonData)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/auth", authHandler).Methods("POST")
	router.HandleFunc("/protected", protectedHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
