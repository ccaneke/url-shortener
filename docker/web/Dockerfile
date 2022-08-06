# note that all the commands here assume that docker compose is run from the project root directory. Which is
# why the command is sudo docker compose -f deployments/docker-compose.yml up. In other words, these commands
# are run in the build working directory and the build working directory is determined by where the docker
# compose command is run from.

# Start with a base image
FROM golang:alpine AS builder

# Update the image and add git
RUN apk add git

# Set environment variable in the container
ENV GO111MODULE=on

# Set the container working directory
WORKDIR /app

# Copy go mod and sum files into the container's current working directory
COPY go.mod go.sum ./

# Update the image and add dependencies
RUN go env && go mod download

# Copy content from the current directory into the container's working directory
COPY . .

# Set the working directory in the container to the go main package directory
WORKDIR /app/cmd/httpapi/server/

# Build the go main package and write the binary to the root directory
RUN go build -o /url-shortner

FROM alpine:latest

WORKDIR /app/

# Copy the binary built in the first stage into the current working directory in the second stage
COPY --from=builder /url-shortner .

# The container should listen on port 8080
EXPOSE 8080

# Command to run when the container starts
ENTRYPOINT ./url-shortner  
