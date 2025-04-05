FROM golang:1.24.2-alpine

WORKDIR /app

COPY config/ ./config/
COPY internal/ ./internal
COPY subscriber/ ./subscriber/
COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o subscriber-bin -ldflags "-s -w" ./subscriber/cmd/subscriber

EXPOSE 3001

CMD ["/app/subscriber-bin"]
