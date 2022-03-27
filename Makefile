ENV := local
ifdef $$APP_ENV
ENV := $$APP_ENV
endif

export PROJECT = github.com/LasseJacobs/go-starter-kit

build:
	env GOOS=linux GOARCH=amd64 go build -o bin/server $(PROJECT)/cmd
	chmod +x bin/server

build-mac:
	env GOOS=darwin GOARCH=arm64 go build -o bin/server $(PROJECT)/cmd
	chmod +x bin/server

run:
	go run ./cmd/main.go

start:
	./bin/server --config config/${ENV}.yaml serve

start-db:
	docker-compose up database

seed:
	./bin/server --config config/${ENV}.yaml admin seed test/testdata/seed.sql

test:
	go test ./... -count=1

tidy:
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache
