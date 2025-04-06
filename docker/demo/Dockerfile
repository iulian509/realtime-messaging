FROM golang:1.24.2-alpine

WORKDIR /app

COPY demo/ ./demo
COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o demo-bin -ldflags "-s -w" ./demo/
