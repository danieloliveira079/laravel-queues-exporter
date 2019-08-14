FROM golang:1.12

WORKDIR /go/src/github.com/danieloliveira079/laravel-queues-exporter

COPY . .

RUN go build -o bin/laravel-queues-exporter cmd/laravel-queues-exporter/laravel-queues-exporter.go

ENTRYPOINT ["/go/src/github.com/danieloliveira079/laravel-queues-exporter/bin/laravel-queues-exporter"]