# Variables
APP_NAME = url-shortener
SRC = ./cmd/main.go

build:
	@echo "docker-compose build - Building $(APP_NAME)..."
	docker-compose build

run:
	@echo "docker-compose up - Running $(APP_NAME)..."
	docker-compose up

generate-swagger:
	@echo "Generating swagger documentation..."
	swag init -g $(SRC)