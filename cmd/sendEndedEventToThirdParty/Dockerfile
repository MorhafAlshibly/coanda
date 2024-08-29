# syntax=docker/dockerfile:1

FROM golang:1.23-alpine

WORKDIR /coanda

ARG PORT

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/sendEndedEventToThirdParty ./cmd/sendEndedEventToThirdParty
COPY api ./api
COPY internal/sendEndedEventToThirdParty ./internal/sendEndedEventToThirdParty

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o out.exe ./cmd/sendEndedEventToThirdParty

# Expose port
EXPOSE $PORT

# Run
CMD ["./out.exe"]