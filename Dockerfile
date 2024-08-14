# Start from the latest golang base image
FROM docker.io/golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main main.go

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]
CMD ["web"]
