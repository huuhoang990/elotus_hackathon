FROM golang:1.23-alpine

WORKDIR /app
COPY . .
RUN go mod tidy
#RUN go build -o main ./cmd/api

EXPOSE 8080
ENTRYPOINT ["go", "run", "/app/cmd/api"]
#ENTRYPOINT ["go", "run", "/app/cmd/frontapi/main.go"]