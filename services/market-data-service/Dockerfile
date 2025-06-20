FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build cmd/main.go into /app/main
RUN go build -o main ./cmd

# optional: copy env if your Go app uses godotenv
COPY .env ./

CMD ["./main"]