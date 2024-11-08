FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o upassed-account-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/upassed-account-service /usr/local/bin/upassed-account-service
EXPOSE 44044
CMD ["upassed-account-service"]
