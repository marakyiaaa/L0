include .env

build:
	docker compose build

run:
	docker compose up -d

reboot:
	docker compose down

topic:
	docker exec kafka kafka-topics --bootstrap-server kafka:9092 --create --topic orders

write_model:
	go run cmd/app/main.go --write-data

go:
	go run cmd/app/main.go

test:
	go test -v internal/service/service_test.go internal/service/service.go

cover:
	go test -cover internal/service/service_test.go internal/service/service.go


style: install-deps
	${LOCAL_BIN}/golangci-lint run

.PHONY: build run reboot topic go style write_model go test cover