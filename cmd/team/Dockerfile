# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /coanda

ARG PORT

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/team ./cmd/team
COPY api ./api
COPY internal/team ./internal/team

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o team ./cmd/team

# Expose port
EXPOSE $PORT

# Run
CMD ["./team"]