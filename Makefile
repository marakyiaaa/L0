include .env

build:
	docker compose build

run:
	docker compose up -d

write_model:
	go run cmd/app/main.go --write-data

go:
	go run cmd/app/main.go

#style: install-deps
#	${LOCAL_BIN}/golangci-lint run


.PHONY: build run go style write_model go