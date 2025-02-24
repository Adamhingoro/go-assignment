# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application (statically linked binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# Stage 2: Create a minimal runtime image
FROM alpine:3.18

# Install necessary certificates (optional, but recommended for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port your application listens on (e.g., 8080)
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]