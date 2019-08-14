build:
	@echo "=============building Local Exporter============="
	docker build -t laravel-queues-exporter .

up:
	@echo "=============starting exporter locally============="
	docker-compose up --build

logs:
	docker-compose logs -f

down:
	docker-compose down

test:
	go test -v -cover ./...

clean: down
	@echo "=============cleaning up============="
	rm -f bin
	#docker system prune -f
	#docker volume prune -f