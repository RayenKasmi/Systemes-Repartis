# TP2/consumers/bo1/Dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bo2-consumer .

FROM alpine:latest
RUN ap# TP2/consumers/bo2/Dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bo2-consumer .

FROM alpine:latest
RUN ap