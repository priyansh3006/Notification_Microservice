FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /trial

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go application
RUN go build -o /trial/bin/app ./main.go

# Final stage
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /trial

# Copy the built binary from the builder stage
COPY --from=builder /trial/bin/app /trial/bin/app

#Explicitly exposing the port
EXPOSE 8080
# Set the entry point for the container
ENTRYPOINT ["/trial/bin/app"]

