# syntax=docker/dockerfile:1

FROM golang:1.23-alpine

WORKDIR /coanda

ARG PORT

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/bff ./cmd/bff
COPY api ./api
COPY internal/bff ./internal/bff

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bff ./cmd/bff

# Expose port
EXPOSE $PORT

# Run
CMD ["./bff"]