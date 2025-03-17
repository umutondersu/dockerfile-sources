FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o dockerfile-sources ./cmd/dockerfile-sources

FROM alpine:latest
# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/dockerfile-sources /dockerfile-sources

ENTRYPOINT ["/dockerfile-sources"]
