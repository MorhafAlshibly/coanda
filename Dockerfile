# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /coanda

ARG PORT
ARG SERVICE

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/$SERVICE ./cmd/$SERVICE
COPY api ./api
COPY internal/$SERVICE ./internal/$SERVICE

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o out.exe ./cmd/$SERVICE

# Expose port
EXPOSE $PORT

# Run
CMD ["./out.exe"]