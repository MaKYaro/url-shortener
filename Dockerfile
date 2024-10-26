FROM golang:1.22-alpine

WORKDIR /app

RUN go version
ENV GOPATH=/

COPY . .

RUN go mod download
RUN go build -o main ./cmd/url-shortener/main.go

CMD ["./main", "--config-path=./config/local.json"]