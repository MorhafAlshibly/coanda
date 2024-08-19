# syntax=docker/dockerfile:1

FROM golang:1.23-alpine

WORKDIR /coanda

ARG PORT

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/event ./cmd/event
COPY api ./api
COPY internal/event ./internal/event

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o event ./cmd/event

# Expose port
EXPOSE $PORT

# Run
CMD ["./event"]