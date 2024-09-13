# Use the official Golang Image
FROM golang:1.23.1-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /automatedgo

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy only the necessary source files
COPY cmd/automatedgo/ ./cmd/automatedgo/
COPY pkg/ ./pkg/

# Build the Go app
RUN go build -o automatedgo ./cmd/automatedgo

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /root/

# Copy the pre-built binary file from the builder stage
COPY --from=builder /automatedgo .

# Command to run the executable
CMD ["./automatedgo"]
