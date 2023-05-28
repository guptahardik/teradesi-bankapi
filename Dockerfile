# Use the official Go 1.17 base image
FROM golang:1.17

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Run the tests and exit
CMD ["go", "test", "./..."]

# Default command to run the application
ENTRYPOINT ["./main"]
