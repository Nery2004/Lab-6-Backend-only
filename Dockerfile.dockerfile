# Imagen base de Go
FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main

EXPOSE 8080

CMD ["./main"]
