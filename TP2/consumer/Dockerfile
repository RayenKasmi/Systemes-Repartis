FROM golang:1.24-alpine

WORKDIR /app

# Copy go.mod and go.sum from the TP2 root
COPY go.mod go.sum ./
RUN go mod download

# Copy consumer code into container
COPY ./consumer/ ./
COPY ./models/ ./models

RUN go build -o consumer .

CMD ["./consumer"]