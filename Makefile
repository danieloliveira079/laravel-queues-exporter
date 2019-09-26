APP_BIN=laravel-queues-exporter
DOCKER_REPO=danieloliv
DOCKER_TAG=${DOCKER_REPO}/laravel-queues-exporter:latest

test:
	go test -v -cover ./...

build-docker:
	@echo "=============building Local Exporter============="
	docker build -t ${DOCKER_TAG} .

push: build-docker
	docker push ${DOCKER_TAG}

build:
	rm -f bin/${APP_BIN}
	go build -o bin/${APP_BIN} cmd/laravel-queues-exporter/main.go

run: build
	./bin/${APP_BIN}

release:
	@echo "==========building release candidate============="
	GOOS=darwin GOARCH=amd64 go build -o bin/${APP_BIN}-darwin-amd64 cmd/laravel-queues-exporter/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/${APP_BIN}-linux-amd64 cmd/laravel-queues-exporter/main.go

up:
	@echo "=============starting exporter locally============="
	docker-compose up --build

services:
	@echo "=============starting local services============="
	 docker-compose up redis statsd redis-commander

logs:
	docker-compose logs -f

down:
	docker-compose down