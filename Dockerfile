# Use a Go base image for building the application
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o bin/mini_agent ./cmd/mini_agent

# Use a minimal Alpine base image for the final stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/bin/mini_agent .

# Create the agent_workspace directory for file tool if it doesn't exist
# and ensure it's writable.
RUN mkdir -p agent_workspace && chmod -R 777 agent_workspace

# Create the data.db file
# TODO: mount a volume for persistent data.
RUN touch data.db

# Command to run the application
CMD ["./mini_agent"]
