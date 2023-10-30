# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /coanda

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/item ./cmd/item
COPY api ./api
COPY internal/item ./internal/item

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o item ./cmd/item
EXPOSE 50051
EXPOSE 8081

# Run
CMD ["./item"]