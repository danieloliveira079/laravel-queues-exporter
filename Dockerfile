FROM golang:1.12 AS builder

WORKDIR /go/src/github.com/danieloliveira079/laravel-queues-exporter

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/laravel-queues-exporter cmd/laravel-queues-exporter/main.go

ENTRYPOINT ["/go/src/github.com/danieloliveira079/laravel-queues-exporter/bin/laravel-queues-exporter"]

FROM alpine

WORKDIR /app

COPY --from=builder /go/src/github.com/danieloliveira079/laravel-queues-exporter/bin/laravel-queues-exporter /app/

RUN chmod +x laravel-queues-exporter

ENTRYPOINT ["./laravel-queues-exporter"]