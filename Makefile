.PHONY: build test test-examples regen start psql mysqlsh

build:
	go build ./...

test:
	go test ./...

test-examples:
	go test --tags=examples ./...

regen: sqlc-dev
	go run ./scripts/regenerate/

sqlc-dev:
	go build -o ~/bin/sqlc-dev ./cmd/sqlc/

sqlc-pg-gen:
	go build -o ~/bin/sqlc-pg-gen ./internal/tools/sqlc-pg-gen

start:
	docker-compose up -d

psql:
	PGPASSWORD=mysecretpassword psql --host=127.0.0.1 --port=5432 --username=postgres dinotest

mysqlsh:
	mysqlsh --sql --user root --password mysecretpassword --database dinotest 127.0.0.1:3306
