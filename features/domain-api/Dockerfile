# Use the official Golang image to build the Go application
FROM golang:1.23-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal image to run the Go application
FROM alpine:latest

# Install dig utility
RUN apk --no-cache add bind-tools

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Command to run the executable
CMD ["./main"]
