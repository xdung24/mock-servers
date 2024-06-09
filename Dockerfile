FROM golang:1.22 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
# -o mock-servers specifies the output name of the binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -o mock-servers

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
FROM alpine:3.20.0 as lean-production

# Create a folder to store the mock-servers binary
RUN mkdir -p /mock-servers

RUN apk add --no-cache ca-certificates

# Set workdir
WORKDIR /mock-servers

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/mock-servers /usr/local/bin/

# Run the web service on container startup with air.

# Production image
# Add air configuration and sample data
FROM lean-production as production

# Copy sample data
COPY data /mock-servers/data

# Copy configuration
COPY .env_sample /mock-servers/.env

# Set workdir
WORKDIR /mock-servers

# Set entrypoint
ENTRYPOINT ["mock-servers"]	
