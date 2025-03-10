# Use the official Golang image
FROM golang:1.24.1-alpine

# Set author metadata using labels
LABEL org.opencontainers.image.authors="https://github.com/Nicconike"

# Set the Current Working Directory inside the container
WORKDIR /automatedgo

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy only the necessary source files
COPY cmd/automatedgo/ ./cmd/automatedgo/
COPY pkg/ ./pkg/

# Build the Go app
RUN go build -o automatedgo ./cmd/automatedgo

# Command to run the executable
ENTRYPOINT ["./automatedgo"]
