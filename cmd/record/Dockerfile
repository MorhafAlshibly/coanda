# syntax=docker/dockerfile:1

FROM golang:1.23-alpine

WORKDIR /coanda

ARG PORT

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/record ./cmd/record
COPY api ./api
COPY internal/record ./internal/record

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o record ./cmd/record

# Expose port
EXPOSE $PORT

# Run
CMD ["./record"]