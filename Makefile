include .env
export

ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: run
run:
	cd src/app && \
	go run cmd/api/main.go

.PHONY: devup
devup:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d

.PHONY: devdown
devdown:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

.PHONY: mgc
mgc:
	docker run --rm -v $(ROOT_DIR)data/migrations:/migrations migrate/migrate -verbose create -ext sql -dir /migrations $(filename)

.PHONY: mgu
mgu:
	docker run --rm --network host -v $(ROOT_DIR)data/migrations:/migrations migrate/migrate -verbose -path=/migrations/ -database "$(DB_DSN)" up

.PHONY: mgd
mgd:
	docker run --rm --network host -v $(ROOT_DIR)/data/migrations:/migrations migrate/migrate -verbose -path=/migrations/ -database $(DB_DSN) down 1