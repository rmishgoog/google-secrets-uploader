# Use the official Golang image to create a build artifact.
# This is the "builder" stage.
FROM golang:1.21-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# -o /app/google-secrets-uploader: output the binary to /app/google-secrets-uploader
# CGO_ENABLED=0: disable CGO for a statically linked binary
# GOOS=linux: target Linux OS
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/google-secrets-uploader .

# Start a new, smaller image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/google-secrets-uploader .

# This will be the default command to run when the container starts
ENTRYPOINT ["/app/google-secrets-uploader"]

# You can pass arguments to the entrypoint like this:
# docker run my-app --project-id my-project --secrets-file secrets.csv
