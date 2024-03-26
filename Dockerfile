# Start with a base image that includes Go for the build stage
FROM golang:1.21-bookworm as builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o quickwit-poc .

# Start a new stage for the runtime environment
FROM debian:bookworm-slim

# Install runtime dependencies, if any
RUN apt-get update && apt-get install -y librdkafka1 && rm -rf /var/lib/apt/lists/*

# Set the working directory in the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/quickwit-poc /app/

# Command to run the executable
CMD ["./quickwit-poc"]
