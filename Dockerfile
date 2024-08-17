FROM golang:1.23-alpine
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bitbird

EXPOSE 8080

# Run
CMD ["/bitbird"]
