FROM golang:1.24.2-alpine

WORKDIR /app

COPY internal/ ./internal/
COPY auth-cli/ ./auth-cli
COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-cli-bin -ldflags "-s -w" ./auth-cli/
