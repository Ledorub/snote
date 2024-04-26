MAIN_PACKAGE_PATH := "./cmd/app"
BINARY_NAME := "snote"

.PHONY: migrate-up
migrate-up:
	docker compose --profile migration up --abort-on-container-exit

.PHONY: migrate-down
migrate-down:
	docker compose --profile migration up --abort-on-container-exit "db", "5432", "/run/secrets/db_main", "up", "/usr/local/src/migrations"

.PHONY: build
build:
	docker compose build

.PHONY: run
run:
	docker compose --profile app up