# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /coanda

ARG PORT

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/wipeTournament ./cmd/wipeTournament
COPY api ./api
COPY internal/wipeTournament ./internal/wipeTournament
COPY internal/tournament ./internal/tournament

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o out.exe ./cmd/wipeTournament

# Expose port
EXPOSE $PORT

# Run
CMD ["./out.exe"]