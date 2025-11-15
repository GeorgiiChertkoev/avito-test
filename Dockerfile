# Build stage
FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download

RUN go build -o /pr-reviewer ./cmd/app

# Run stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /pr-reviewer .

EXPOSE $HTTP_PORT
CMD ["./pr-reviewer"]