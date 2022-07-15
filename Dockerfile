FROM golang:1.17.6

LABEL Author="insaneEra" project="forum" version="1.0"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/main.go

CMD ["./main"]