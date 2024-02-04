# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /coanda

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/tournament ./cmd/tournament
COPY api ./api
COPY internal/tournament ./internal/tournament

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o tournament ./cmd/tournament

# Run
CMD ["./tournament"]