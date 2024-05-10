FROM golang:1.22.3
# Add Maintainer Info
LABEL maintainer="ly.lester@rayprosoft.com"

# ENV GOPROXY=http://192.168.50.245:3000,direct
# ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY . .
# Copy the source from the current directory to the Working Directory inside the container
RUN go mod download
RUN go build -o main ./features/daily-first-deposit
RUN ls -al .

# Expose port 8080 to the outside world
EXPOSE 8080
# Command to run the executable
CMD ["./main"]

