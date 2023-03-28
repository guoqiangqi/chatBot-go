package chatbot

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}

type ChatPayload struct {
	Question string `json:"question"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}
