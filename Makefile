# автовыбор shell под платформу
ifeq ($(OS),Windows_NT)
	SHELL := powershell.exe
else
	SHELL := /bin/bash
endif

APP_NAME := url-shortener
CMD_DIR := ./cmd/url-shortener

# загрузка .env
ifneq (,$(wildcard .env))
    include .env
    export
endif

# указывает make, что это команды, а не файлы
.PHONY: run build migrate-up migrate-down migrate-new docker-up docker-down

run:
	$(SHELL) -Command "setx CONFIG_PATH '$(CONFIG_PATH)'"
	CONFIG_PATH=$(CONFIG_PATH) go run $(CMD_DIR)

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-new:
	migrate create -ext sql -dir migrations -seq $(name)

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
