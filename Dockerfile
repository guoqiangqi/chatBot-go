FROM golang:alpine
COPY . /chatbot-go
WORKDIR /chatbot-go
RUN go mod download
RUN go build -o main src/main.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /chatbot-go/main ./
CMD ["./main"]

# NOTICE: set OPENAI_API_KEY env variable when run whith docker, just like:
# docker run -dit -p 8080:8080 -e OPENAI_API_KEY="xx" image_name