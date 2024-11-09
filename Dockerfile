FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY config/dev.yml config/dev.yml
RUN go build -o upassed-account-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir -p /upassed-account-service/config
COPY --from=builder /app/upassed-account-service /upassed-account-service/upassed-account-service
COPY --from=builder /app/config/dev.yml /upassed-account-service/config/dev.yml
RUN chmod +x /upassed-account-service/upassed-account-service
ENV APP_CONFIG_PATH="/upassed-account-service/config/dev.yml"
EXPOSE 44044
CMD ["/upassed-account-service/upassed-account-service"]
