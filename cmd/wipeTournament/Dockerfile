# syntax=docker/dockerfile:1

FROM golang:1.21-alpine as builder

WORKDIR /coanda

# Install CA certificates
RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates


# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY pkg ./pkg
COPY cmd/wipeTournament ./cmd/wipeTournament
COPY api ./api
COPY internal/tournament ./internal/tournament

# Build with optional lambda.norpc tag
RUN CGO_ENABLED=0 GOOS=linux go build -tags lambda.norpc -o wipeTournament ./cmd/wipeTournament

# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /coanda/wipeTournament ./wipeTournament

# Run
ENTRYPOINT ["./wipeTournament"]