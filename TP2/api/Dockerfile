FROM golang:1.24-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies early
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . /app

# Set the working directory to the 'api' folder where your main.go file resides
WORKDIR /app/api

# Build the Go application
RUN go build -o /app/main .

# Expose the API port (optional)
EXPOSE 8081

# Run the API
CMD ["/app/main"]
