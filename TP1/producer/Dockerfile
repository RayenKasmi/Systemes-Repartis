FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY producer.go .

# Build the producer binary
RUN go build -o producer .

# Command to run the producer
CMD ["./producer"]
