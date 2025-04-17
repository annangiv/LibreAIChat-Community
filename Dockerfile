# Use Go base image
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -o libreai ./cmd/server

# Final minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/libreai .
COPY ./views ./views
COPY ./static ./static

ENV PORT=3000

EXPOSE 3000

CMD ["./libreai"]
