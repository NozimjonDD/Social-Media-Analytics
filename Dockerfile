FROM golang:1.21.3-alpine

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN go build -o tgstat cmd/tgstat/main.go

CMD ["./tgstat"]
