# This is a master dockerfile used for fly io to point entrypoint to the correct file
# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /coanda

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build  -o . ./...

EXPOSE 8080

# Run
ENTRYPOINT [ "./bff" ]