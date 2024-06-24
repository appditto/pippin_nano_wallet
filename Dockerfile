# Stage 1: Build the binary with OpenCL dependencies
FROM golang:1.22-bullseye as builder

ARG VERSION

# Set the working directory inside the container
WORKDIR /app

# Copy the workspace files
COPY . .

# Ensure dependencies are downloaded based on your workspace configuration
RUN go work sync

# Build the application
RUN CGO_ENABLED=1 go build -a -ldflags "-s -w -X main.Version=${VERSION}" -o pippin ./apps/cli

# Stage 2: Use a smaller base image
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/pippin .

# Add ca-certificates in case your app makes external network requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Command to run the executable
ENTRYPOINT ["./pippin", "-start-server"]
