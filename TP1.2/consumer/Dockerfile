FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the consumer code
COPY consumer.go .

# Build the consumer binary
RUN go build -o consumer .

# Command to run the consumer
CMD ["./consumer"]
