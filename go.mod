module chatbot

go 1.20

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.7
	github.com/sashabaranov/go-openai v1.5.7
)

replace github.com/sashabaranov/go-openai v1.5.7 => github.com/guoqiangqi/go-openai v0.0.0-20230417064903-dc85f51edea7
