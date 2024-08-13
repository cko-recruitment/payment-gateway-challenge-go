# Use a lightweight, minimal base image
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the executable
RUN go build -o main .

# Use a smaller base image for production
FROM alpine:latest

# Copy the built binary
COPY --from=builder /app/main /app/

# Set the working directory
WORKDIR /app

# Expose the port your application listens on
EXPOSE 8090

# Command to run when the container starts
CMD ["./main"]