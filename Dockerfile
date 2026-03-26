FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./

ENV GOPROXY=https://proxy.golang.org,direct

RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/main.go

CMD ["./main"]