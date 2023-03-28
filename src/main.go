package main

import (
	chatbot "chatbot/utils"
	"fmt"
	"log"
	"strings"
	"net/http"

	"github.com/gorilla/mux"
)

var jwtSecreteKey = []byte("change your secret key here.")

func authHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := chatbot.Authenticate(username, password)

	if user != nil {
		tokenString, _ := chatbot.GenerateToken(user.Id, jwtSecreteKey)
		fmt.Fprintln(w, tokenString)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Authentication failed")
	}
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	const bearerPrefix = "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Wrong token: empty or not start with Bearer.")
		return
	}

	tokenString := authHeader[len(bearerPrefix):]
	claims, err := chatbot.ParseToken(tokenString, jwtSecreteKey)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token")
		return
	}

	// username := chatbot.IndexUserWithID(claims.Id).Name
    username := claims.StandardClaims.Subject
	fmt.Fprintf(w, "Hello, %s!", username)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/auth", authHandler).Methods("POST")
	router.HandleFunc("/protected", protectedHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}