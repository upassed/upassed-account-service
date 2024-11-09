FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o upassed-account-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir -p /upassed-account-service/config
RUN mkdir -p /upassed-account-service/migration/scripts
COPY --from=builder /app/upassed-account-service /upassed-account-service/upassed-account-service
COPY --from=builder /app/config/dev.yml /upassed-account-service/config/dev.yml
COPY --from=builder /app/migration/scripts/* /upassed-account-service/migration/scripts
RUN chmod +x /upassed-account-service/upassed-account-service
ENV APP_CONFIG_PATH="/upassed-account-service/config/local.yml"
EXPOSE 44044
CMD ["/upassed-account-service/upassed-account-service"]
