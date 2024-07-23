# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /coanda

ARG PORT

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/matchmaking ./cmd/matchmaking
COPY api ./api
COPY internal/matchmaking ./internal/matchmaking

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o matchmaking ./cmd/matchmaking

# Expose port
EXPOSE $PORT

# Run
CMD ["./matchmaking"]