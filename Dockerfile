# Stage 1: Build the binary with OpenCL dependencies
FROM golang:1.22-bullseye as builder

ARG VERSION

# Set the working directory inside the container
WORKDIR /app

# Copy the workspace files
COPY . .

# Ensure dependencies are downloaded based on your workspace configuration
RUN go work sync

# Build the application statically
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w -X main.Version=${VERSION}" -o pippin ./apps/cli

# Stage 2: Use a smaller base image
FROM scratch

# Set the working directory inside the container
WORKDIR /root/

ENV PIPPIN_HOME=/root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/pippin .

# Copy SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Command to run the executable
ENTRYPOINT ["./pippin", "-start-server"]
