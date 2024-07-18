FROM golang:1.22.3 as builder

# Set the Current Working Directory inside the container
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main ./features/domain-api

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
