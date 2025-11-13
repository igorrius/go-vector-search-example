# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go-vector-search ./cmd/server

# Stage 2: Create the final image
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /go-vector-search /go-vector-search

# Expose the application port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/go-vector-search"]