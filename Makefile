.PHONY: help build up down clean test

help:
	@echo "Weather by CEP - Makefile Commands"
	@echo ""
	@echo "  make build          Build Docker images"
	@echo "  make up             Start services with docker-compose"
	@echo "  make down           Stop and remove services"
	@echo "  make test           Run all unit tests"
	@echo "  make test-valid     Test with valid CEP"
	@echo "  make test-invalid   Test with invalid CEP"
	@echo "  make test-notfound  Test with non-existent CEP"
	@echo ""

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

test:
	go test -v ./...

test-valid:
	curl -X POST http://localhost:8080/weather \
		-H "Content-Type: application/json" \
		-d '{"cep":"01310100"}' | jq .

test-invalid:
	curl -X POST http://localhost:8080/weather \
		-H "Content-Type: application/json" \
		-d '{"cep":"123"}' | jq .

test-notfound:
	curl -X POST http://localhost:8080/weather \
		-H "Content-Type: application/json" \
		-d '{"cep":"99999999"}' | jq .
