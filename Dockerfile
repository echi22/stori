# syntax=docker/dockerfile:1

FROM golang:1.21
WORKDIR /app
COPY . .
RUN go mod tidy
RUN apt-get update && apt-get install -y gcc sqlite3
ENV CGO_ENABLED=1
RUN go build -o stori-app main.go
CMD ["./stori-app"] 