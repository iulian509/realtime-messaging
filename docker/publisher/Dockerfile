FROM golang:1.24.2-alpine

WORKDIR /app

COPY config/ ./config/
COPY internal/ ./internal
COPY publisher/ ./publisher/
COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o publisher-bin -ldflags "-s -w" ./publisher/cmd/publisher

EXPOSE 3000

CMD ["/app/publisher-bin"]
